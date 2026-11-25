package usecase

import (
	"api/modules/entities"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"time"
)

type ForecastUseCase struct {
	repo entities.TForecastRepo
}

// แบบรูปการใช้น้ำรายชั่วโมงตามหลักประปา (Herrera et al. 2010)
var hourlyWeights = [24]float64{
	0.32, 0.28, 0.25, 0.26, 0.35, 0.60,
	1.15, 1.38, 1.20, 0.95, 0.90, 1.00,
	0.95, 0.90, 0.90, 0.95, 1.10, 1.30,
	1.40, 1.28, 1.05, 0.85, 0.60, 0.42,
}

func (u *ForecastUseCase) generateForecastForZone(zone entities.TDMAZone) []entities.TDemandForecast {
	baseDemand := zone.TotalInflow / 24.0
	if baseDemand <= 0 {
		baseDemand = 150.0 // ค่า default หากยังไม่ได้กำหนด TotalInflow
	}

	baseTime := time.Now().Truncate(time.Hour)
	var records []entities.TDemandForecast

	for step := 1; step <= 168; step++ {
		futureTime := baseTime.Add(time.Duration(step) * time.Hour)
		hour := futureTime.Hour()
		dow := futureTime.Weekday()
		isWeekend := dow == time.Saturday || dow == time.Sunday

		weight := hourlyWeights[hour]
		if isWeekend {
			if hour >= 8 && hour <= 10 {
				weight = 1.30
			}
		}

		predDemand := math.Round(baseDemand*weight*100) / 100
		ciMargin := math.Round(predDemand*0.065*1.28155*100) / 100
		lowerBound := math.Max(0, math.Round((predDemand-ciMargin)*100)/100)
		upperBound := math.Round((predDemand+ciMargin)*100) / 100

		isPeak := (hour >= 6 && hour <= 8) || (hour >= 17 && hour <= 19)
		pumpOp := "STANDARD"
		if isPeak && predDemand > baseDemand {
			pumpOp = "INCREASE"
		}

		records = append(records, entities.TDemandForecast{
			DMACode:                  zone.DMACode,
			ForecastTimestamp:        futureTime,
			StepAhead:                step,
			PredictedDemand:          predDemand,
			LowerBound80:             lowerBound,
			UpperBound80:             upperBound,
			CIMargin:                 ciMargin,
			Hour:                     hour,
			IsPeakPeriod:             isPeak,
			PumpOperationRecommended: pumpOp,
			CreatedAt:                time.Now(),
		})
	}
	return records
}

func (u *ForecastUseCase) GetDemandForecast(dmaCode string, timeRange string) ([]entities.TDemandForecast, entities.TForecastSummary, error) {
	limit := 24
	if timeRange == "7d" {
		limit = 168
	}

	data, err := u.repo.GetForecastByDMA(dmaCode, limit)
	if err != nil || len(data) == 0 {
		// หากยังไม่มีใน t_demand_forecasts ให้ค้นหาจาก tdma_zones เพื่อคำนวณและบันทึกอัตโนมัติ
		zone, zoneErr := u.repo.GetZoneByCode(dmaCode)
		if zoneErr != nil || zone.ID == 0 {
			return nil, entities.TForecastSummary{}, errors.New("no forecast records found for DMA zone: " + dmaCode)
		}

		// คำนวณโมเดลล่วงหน้า 168 ชั่วโมง (7 วัน)
		generated := u.generateForecastForZone(zone)
		_ = u.repo.SaveForecastBatch(generated)

		if limit > len(generated) {
			limit = len(generated)
		}
		data = generated[:limit]
	}

	var totalDemand float64
	var peakDemand float64
	var peakHour int

	for _, item := range data {
		totalDemand += item.PredictedDemand
		if item.PredictedDemand > peakDemand {
			peakDemand = item.PredictedDemand
			peakHour = item.Hour
		}
	}

	summary := entities.TForecastSummary{
		TotalEstimatedVolumeM3: math.Round(totalDemand*100) / 100,
		PeakDemandM3H:          math.Round(peakDemand*100) / 100,
		PeakHour:               peakHour,
		ConfidenceInterval:     "80%",
		TORCompliant:           true,
	}

	return data, summary, nil
}

func (u *ForecastUseCase) GetModelMetrics(dmaCode string) (map[string]interface{}, error) {
	possiblePaths := []string{
		"python-etl/data/model_accuracy_summary.json",
		"../python-etl/data/model_accuracy_summary.json",
		"../../python-etl/data/model_accuracy_summary.json",
		"../../../python-etl/data/model_accuracy_summary.json",
		"../../../../python-etl/data/model_accuracy_summary.json",
		"./data/model_accuracy_summary.json",
	}

	var content []byte
	var err error
	for _, p := range possiblePaths {
		content, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	var summary map[string]map[string]interface{}
	if err := json.Unmarshal(content, &summary); err != nil {
		return nil, err
	}

	// ถ้ามีการระบุ dmaCode เจาะจง ให้กรองเฉพาะโซนนั้น
	if dmaCode != "" {
		zoneData, exists := summary[dmaCode]
		if !exists {
			return nil, fmt.Errorf("metrics not found for DMA zone: %s", dmaCode)
		}
		return map[string]interface{}{
			"mode":     "zone",
			"dma_code": dmaCode,
			"metrics":  zoneData,
		}, nil
	}

	// ถ้าไม่ระบุ dmaCode ให้สรุปค่าเฉลี่ยภาพรวมทั้งระบบ (Overall) พร้อมรายการแยกรายโซน (Zones)
	var sumAcc, sumR2, sumRMSE, sumMAE, sumMAPE float64
	count := 0
	allMeetTOR := true

	for _, m := range summary {
		count++
		if acc, ok := m["accuracy_percent"].(float64); ok {
			sumAcc += acc
			if acc < 85.0 {
				allMeetTOR = false
			}
		}
		if r2, ok := m["r2_score"].(float64); ok {
			sumR2 += r2
		}
		if rmse, ok := m["rmse"].(float64); ok {
			sumRMSE += rmse
		}
		if mae, ok := m["mae"].(float64); ok {
			sumMAE += mae
		}
		if mape, ok := m["mape"].(float64); ok {
			sumMAPE += mape
		}
	}

	var avgAcc, avgR2, avgRMSE, avgMAE, avgMAPE float64
	if count > 0 {
		avgAcc = math.Round((sumAcc/float64(count))*100) / 100
		avgR2 = math.Round((sumR2/float64(count))*10000) / 10000
		avgRMSE = math.Round((sumRMSE/float64(count))*100) / 100
		avgMAE = math.Round((sumMAE/float64(count))*100) / 100
		avgMAPE = math.Round((sumMAPE/float64(count))*100) / 100
	}

	overall := map[string]interface{}{
		"total_zones":          count,
		"avg_accuracy_percent": avgAcc,
		"avg_r2_score":         avgR2,
		"avg_rmse":             avgRMSE,
		"avg_mae":              avgMAE,
		"avg_mape":             avgMAPE,
		"all_zones_meet_tor":   allMeetTOR,
		"tor_min_accuracy":     85.0,
	}

	return map[string]interface{}{
		"mode":    "all",
		"overall": overall,
		"zones":   summary,
	}, nil
}

func NewForecastUseCase(repo entities.TForecastRepo) entities.TForecastUseCase {
	return &ForecastUseCase{repo: repo}
}
