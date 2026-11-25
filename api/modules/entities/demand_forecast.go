package entities

import "time"

// TDemandForecast ผลการพยากรณ์ความต้องการใช้น้ำ (Water Demand Forecasting) ตาม TOR ภาคผนวก ๔
type TDemandForecast struct {
	ID                       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DMACode                  string    `gorm:"type:varchar(50);index" json:"dma_code"`
	ForecastTimestamp        time.Time `gorm:"index" json:"forecast_timestamp"`
	StepAhead                int       `json:"step_ahead"`                            // ชั่วโมงล่วงหน้า (1-24 หรือ 1-168)
	PredictedDemand          float64   `json:"predicted_demand"`                      // ปริมาณน้ำที่พยากรณ์ (ลบ.ม./ชม.)
	LowerBound80             float64   `json:"lower_bound_80"`                        // ขอบเขตล่างช่วงความเชื่อมั่น 80%
	UpperBound80             float64   `json:"upper_bound_80"`                        // ขอบเขตบนช่วงความเชื่อมั่น 80%
	CIMargin                 float64   `json:"ci_margin"`                             // ค่าเผื่อความเชื่อมั่น (+/-)
	Hour                     int       `json:"hour"`                                  // ชั่วโมง (0-23)
	IsPeakPeriod             bool      `json:"is_peak_period"`                        // ช่วงเวลา Peak การใช้น้ำ
	PumpOperationRecommended string    `gorm:"type:varchar(50)" json:"pump_operation_recommended"` // คำแนะนำเดินเครื่องสูบ: INCREASE / STANDARD / REDUCE
	CreatedAt                time.Time `json:"created_at"`
}

type TDemandForecastReq struct {
	DMACode                  string    `json:"dma_code"`
	ForecastTimestamp        time.Time `json:"forecast_timestamp"`
	StepAhead                int       `json:"step_ahead"`
	PredictedDemand          float64   `json:"predicted_demand"`
	LowerBound80             float64   `json:"lower_bound_80"`
	UpperBound80             float64   `json:"upper_bound_80"`
	CIMargin                 float64   `json:"ci_margin"`
	Hour                     int       `json:"hour"`
	IsPeakPeriod             bool      `json:"is_peak_period"`
	PumpOperationRecommended string    `json:"pump_operation_recommended"`
}

type TForecastSummary struct {
	TotalEstimatedVolumeM3 float64 `json:"total_estimated_volume_m3"`
	PeakDemandM3H          float64 `json:"peak_demand_m3_h"`
	PeakHour               int     `json:"peak_hour"`
	ConfidenceInterval     string  `json:"confidence_interval"`
	TORCompliant           bool    `json:"tor_compliant"`
}

type TForecastRepo interface {
	GetForecastByDMA(dmaCode string, limit int) ([]TDemandForecast, error)
	GetZoneByCode(dmaCode string) (TDMAZone, error)
	SaveForecastBatch(items []TDemandForecast) error
}

type TForecastUseCase interface {
	GetDemandForecast(dmaCode string, timeRange string) ([]TDemandForecast, TForecastSummary, error)
	GetModelMetrics(dmaCode string) (map[string]interface{}, error)
}
