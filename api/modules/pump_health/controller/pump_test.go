package controller_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api/modules/entities"
	"api/modules/pump_health/controller"
	"api/modules/pump_health/usecase"

	"github.com/gofiber/fiber/v2"
)

type MockPumpRepo struct {
	pumps  map[string]entities.TPumpAsset
	scores map[string]entities.TPumpHealthScore
}

func (m *MockPumpRepo) GetAllPumps() ([]entities.TPumpAsset, error) {
	var list []entities.TPumpAsset
	for _, p := range m.pumps {
		list = append(list, p)
	}
	return list, nil
}

func (m *MockPumpRepo) GetPumpByID(pumpID string) (entities.TPumpAsset, error) {
	p, ok := m.pumps[pumpID]
	if !ok {
		return p, errors.New("pump not found")
	}
	return p, nil
}

func (m *MockPumpRepo) GetAllPumpHealth() ([]entities.TPumpHealthScore, error) {
	var list []entities.TPumpHealthScore
	for _, s := range m.scores {
		list = append(list, s)
	}
	return list, nil
}

func (m *MockPumpRepo) GetPumpHealthByID(pumpID string) (entities.TPumpHealthScore, error) {
	s, ok := m.scores[pumpID]
	if !ok {
		return s, errors.New("health not found")
	}
	return s, nil
}

func (m *MockPumpRepo) SavePumpHealthBatch(items []entities.TPumpHealthScore) error {
	for _, item := range items {
		m.scores[item.PumpID] = item
	}
	return nil
}

func setupPumpApp() *fiber.App {
	app := fiber.New()

	pumpsMap := map[string]entities.TPumpAsset{
		"PUMP-SK-01": {
			ID:           1,
			PumpID:       "PUMP-SK-01",
			PumpName:     "เครื่องสูบน้ำแรงดันสูง หมายเลข 1 (สงขลา)",
			StationName:  "สถานีผลิตน้ำสงขลา",
			DMACode:      "DMA-SK-01",
			RatedFlowM3H: 500,
			RatedHeadBar: 5.0,
			RatedPowerKW: 90,
			InstallYear:  2021,
			Status:       "ACTIVE",
		},
		"PUMP-BKK-02": {
			ID:           2,
			PumpID:       "PUMP-BKK-02",
			PumpName:     "เครื่องสูบน้ำจ่ายน้ำแรงดันสูง หมายเลข 2 (กทม)",
			StationName:  "สถานีจ่ายน้ำสุทธิสาร",
			DMACode:      "DMA-BKK-01",
			RatedFlowM3H: 750,
			RatedHeadBar: 5.2,
			RatedPowerKW: 150,
			InstallYear:  2017,
			Status:       "CRITICAL",
		},
	}

	scoresMap := map[string]entities.TPumpHealthScore{
		"PUMP-SK-01": {
			ID:                  1,
			PumpID:              "PUMP-SK-01",
			PumpName:            "เครื่องสูบน้ำแรงดันสูง หมายเลข 1 (สงขลา)",
			StationName:         "สถานีผลิตน้ำสงขลา",
			DMACode:             "DMA-SK-01",
			HealthScore:         88.3,
			HealthStatus:        "GOOD",
			PTTFDays:            166,
			PTTFHours:           3984,
			FailureRisk:         "LOW",
			PrimaryStressFactor: "NORMAL_OPERATION",
			Recommendation:      "สภาพเครื่องสูบน้ำสมบูรณ์ดี",
			ISO10816Zone:        "Zone A",
			Telemetry: entities.TPumpTelemetry{
				VibrationRMSMMs:   2.15,
				BearingTempC:      57.8,
				MotorTempC:        63.4,
				FlowRateM3H:       485.0,
				HeadBar:           4.9,
				PowerKW:           88.5,
				EfficiencyPercent: 82.5,
				OperatingHours:    3420,
				CavitationIndex:   0.42,
			},
			EvaluatedAt: time.Now(),
		},
		"PUMP-BKK-02": {
			ID:                  2,
			PumpID:              "PUMP-BKK-02",
			PumpName:            "เครื่องสูบน้ำจ่ายน้ำแรงดันสูง หมายเลข 2 (กทม)",
			StationName:         "สถานีจ่ายน้ำสุทธิสาร",
			DMACode:             "DMA-BKK-01",
			HealthScore:         41.6,
			HealthStatus:        "CRITICAL",
			PTTFDays:            10,
			PTTFHours:           240,
			FailureRisk:         "CRITICAL",
			PrimaryStressFactor: "BEARING_AND_VIBRATION_CRITICAL",
			Recommendation:      "วิกฤต: ตรวจพบการสั่นสะเทือนรุนแรง",
			ISO10816Zone:        "Zone D",
			Telemetry: entities.TPumpTelemetry{
				VibrationRMSMMs:   7.65,
				BearingTempC:      84.5,
				MotorTempC:        89.2,
				FlowRateM3H:       580.0,
				HeadBar:           4.1,
				PowerKW:           158.0,
				EfficiencyPercent: 55.0,
				OperatingHours:    12800,
				CavitationIndex:   2.15,
			},
			EvaluatedAt: time.Now(),
		},
	}

	repo := &MockPumpRepo{pumps: pumpsMap, scores: scoresMap}
	uc := usecase.NewPumpUseCase(repo)
	group := app.Group("/api/v1/pump")
	controller.NewPumpController(group, uc)
	return app
}

func TestGetFleetHealth_Success(t *testing.T) {
	app := setupPumpApp()

	// 1. ทดสอบ GET /api/v1/pump/health (ดึงสุขภาพปั๊มทั้งหมด)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pump/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var resBody struct {
		Status  int                         `json:"status"`
		Summary entities.TPumpFleetSummary  `json:"summary"`
		Data    []entities.TPumpHealthScore `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resBody.Data) != 2 {
		t.Fatalf("Expected 2 pumps, got %d", len(resBody.Data))
	}
	if resBody.Summary.TotalPumps != 2 {
		t.Fatalf("Expected total_pumps = 2, got %d", resBody.Summary.TotalPumps)
	}
	if !resBody.Summary.UrgentActionRequired {
		t.Fatalf("Expected UrgentActionRequired to be true because PUMP-BKK-02 is CRITICAL")
	}

	// 2. ทดสอบ GET /api/v1/pump/health/PUMP-SK-01 (ดึงเจาะจงตัวปั๊ม)
	reqSK01 := httptest.NewRequest(http.MethodGet, "/api/v1/pump/health/PUMP-SK-01", nil)
	respSK01, err := app.Test(reqSK01)
	if err != nil {
		t.Fatalf("Failed to execute request for SK-01: %v", err)
	}
	if respSK01.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200 for SK-01, got %d", respSK01.StatusCode)
	}

	var resSK01 struct {
		Status int                     `json:"status"`
		Data   entities.TPumpHealthScore `json:"data"`
	}
	bodySK01Bytes, _ := io.ReadAll(respSK01.Body)
	if err := json.Unmarshal(bodySK01Bytes, &resSK01); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resSK01.Data.PumpID != "PUMP-SK-01" {
		t.Fatalf("Expected PUMP-SK-01, got %s", resSK01.Data.PumpID)
	}
	if resSK01.Data.ISO10816Zone != "Zone A" {
		t.Fatalf("Expected Zone A, got %s", resSK01.Data.ISO10816Zone)
	}
	if resSK01.Data.PTTFDays <= 0 {
		t.Fatalf("Expected PTTFDays > 0, got %d", resSK01.Data.PTTFDays)
	}

	// 3. ทดสอบ GET /api/v1/pump/summary (ดึงภาพรวมสำหรับ Executive Dashboard)
	reqSummary := httptest.NewRequest(http.MethodGet, "/api/v1/pump/summary", nil)
	respSummary, err := app.Test(reqSummary)
	if err != nil {
		t.Fatalf("Failed to execute summary request: %v", err)
	}
	if respSummary.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200 for summary, got %d", respSummary.StatusCode)
	}

	// 4. ทดสอบ GET /api/v1/pump/health/UNKNOWN-PUMP (404 Not Found)
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/v1/pump/health/UNKNOWN-PUMP", nil)
	respNotFound, _ := app.Test(reqNotFound)
	if respNotFound.StatusCode != fiber.StatusNotFound {
		t.Fatalf("Expected status 404 for unknown pump, got %d", respNotFound.StatusCode)
	}
}
