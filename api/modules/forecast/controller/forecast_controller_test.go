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
	"api/modules/forecast/controller"
	"api/modules/forecast/usecase"

	"github.com/gofiber/fiber/v2"
)

type MockForecastRepo struct {
	data map[string][]entities.TDemandForecast
}

func (m *MockForecastRepo) GetForecastByDMA(dmaCode string, limit int) ([]entities.TDemandForecast, error) {
	items, ok := m.data[dmaCode]
	if !ok || len(items) == 0 {
		return nil, errors.New("no forecast found")
	}
	if limit > len(items) {
		limit = len(items)
	}
	return items[:limit], nil
}

func (m *MockForecastRepo) GetZoneByCode(dmaCode string) (entities.TDMAZone, error) {
	if dmaCode == "DMA-03" {
		return entities.TDMAZone{
			ID:          4,
			DMACode:     "DMA-03",
			DMAName:     "โซนนิคมอุตสาหกรรมภาคใต้",
			TotalInflow: 3600.0,
		}, nil
	}
	return entities.TDMAZone{}, errors.New("zone not found")
}

func (m *MockForecastRepo) SaveForecastBatch(items []entities.TDemandForecast) error {
	if len(items) > 0 {
		m.data[items[0].DMACode] = items
	}
	return nil
}

func setupForecastApp() *fiber.App {
	app := fiber.New()
	mockData := make(map[string][]entities.TDemandForecast)

	// สร้าง mock forecast 24 ชั่วโมง
	var skItems []entities.TDemandForecast
	baseTime := time.Date(2026, 11, 25, 13, 0, 0, 0, time.UTC)
	for i := 1; i <= 24; i++ {
		skItems = append(skItems, entities.TDemandForecast{
			ID:                       uint(i),
			DMACode:                  "DMA-SK-01",
			ForecastTimestamp:        baseTime.Add(time.Duration(i) * time.Hour),
			StepAhead:                i,
			PredictedDemand:          130.0 + float64(i%5)*5.0,
			LowerBound80:             122.0,
			UpperBound80:             138.0,
			CIMargin:                 8.0,
			Hour:                     (13 + i) % 24,
			IsPeakPeriod:             i == 5 || i == 6,
			PumpOperationRecommended: "STANDARD",
		})
	}
	mockData["DMA-SK-01"] = skItems

	repo := &MockForecastRepo{data: mockData}
	uc := usecase.NewForecastUseCase(repo)
	group := app.Group("/api/v1/forecast")
	controller.NewForecastController(group, uc)
	return app
}

func TestGetDemandForecast_Success(t *testing.T) {
	app := setupForecastApp()

	// 1. ทดสอบเรียกพยากรณ์ 24 ชั่วโมง: GET /api/v1/forecast/demand/DMA-SK-01?range=24h
	req := httptest.NewRequest(http.MethodGet, "/api/v1/forecast/demand/DMA-SK-01?range=24h", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var resBody struct {
		Status  int                        `json:"status"`
		DMACode string                     `json:"dma_code"`
		Range   string                     `json:"range"`
		Summary entities.TForecastSummary  `json:"summary"`
		Data    []entities.TDemandForecast `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resBody.DMACode != "DMA-SK-01" {
		t.Fatalf("Expected DMA-SK-01, got %s", resBody.DMACode)
	}
	if len(resBody.Data) != 24 {
		t.Fatalf("Expected 24 forecast records, got %d", len(resBody.Data))
	}
	if !resBody.Summary.TORCompliant {
		t.Fatalf("Expected TORCompliant to be true")
	}
	if resBody.Summary.ConfidenceInterval != "80%" {
		t.Fatalf("Expected ConfidenceInterval 80%%, got %s", resBody.Summary.ConfidenceInterval)
	}

	// 2. ทดสอบโซน DMA-03 (ที่เพิ่งมีในระบบแต่ยังไม่มี forecast ล่วงหน้า) ต้อง auto-generate สำเร็จ 200 OK
	reqDMA03 := httptest.NewRequest(http.MethodGet, "/api/v1/forecast/demand/DMA-03?range=24h", nil)
	respDMA03, errDMA03 := app.Test(reqDMA03)
	if errDMA03 != nil {
		t.Fatalf("Failed to execute request for DMA-03: %v", errDMA03)
	}
	if respDMA03.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200 for DMA-03, got %d", respDMA03.StatusCode)
	}

	bodyDMA03Bytes, _ := io.ReadAll(respDMA03.Body)
	var resBodyDMA03 struct {
		Status  int                        `json:"status"`
		DMACode string                     `json:"dma_code"`
		Data    []entities.TDemandForecast `json:"data"`
	}
	if err := json.Unmarshal(bodyDMA03Bytes, &resBodyDMA03); err != nil {
		t.Fatalf("Failed to decode response for DMA-03: %v", err)
	}
	if resBodyDMA03.DMACode != "DMA-03" {
		t.Fatalf("Expected DMA-03, got %s", resBodyDMA03.DMACode)
	}
	if len(resBodyDMA03.Data) != 24 {
		t.Fatalf("Expected 24 forecast records for DMA-03, got %d", len(resBodyDMA03.Data))
	}

	// 3. ทดสอบกรณีไม่พบโซนเลย: GET /api/v1/forecast/demand/UNKNOWN-ZONE
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/v1/forecast/demand/UNKNOWN-ZONE", nil)
	respNotFound, _ := app.Test(reqNotFound)
	if respNotFound.StatusCode != fiber.StatusNotFound {
		t.Fatalf("Expected status 404 for unknown zone, got %d", respNotFound.StatusCode)
	}

	// 4. ทดสอบ GET /api/v1/forecast/metrics (ภาพรวม Overall + ทุกโซน)
	reqMetricsAll := httptest.NewRequest(http.MethodGet, "/api/v1/forecast/metrics", nil)
	respMetricsAll, err := app.Test(reqMetricsAll)
	if err != nil {
		t.Fatalf("Failed to execute metrics all request: %v", err)
	}
	if respMetricsAll.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200 for metrics all, got %d", respMetricsAll.StatusCode)
	}

	var resMetricsAll struct {
		Status           int `json:"status"`
		ModelPerformance struct {
			Mode    string                 `json:"mode"`
			Overall map[string]interface{} `json:"overall"`
			Zones   map[string]interface{} `json:"zones"`
		} `json:"model_performance"`
	}
	bodyMetricsAll, _ := io.ReadAll(respMetricsAll.Body)
	if err := json.Unmarshal(bodyMetricsAll, &resMetricsAll); err != nil {
		t.Fatalf("Failed to decode metrics all response: %v", err)
	}
	if resMetricsAll.ModelPerformance.Mode != "all" {
		t.Fatalf("Expected mode 'all', got %s", resMetricsAll.ModelPerformance.Mode)
	}
	if resMetricsAll.ModelPerformance.Overall["total_zones"] == 0 {
		t.Fatalf("Expected total_zones > 0")
	}

	// 5. ทดสอบ GET /api/v1/forecast/metrics?dma_code=DMA-03 (กรองเฉพาะโซน)
	reqMetricsZone := httptest.NewRequest(http.MethodGet, "/api/v1/forecast/metrics?dma_code=DMA-03", nil)
	respMetricsZone, err := app.Test(reqMetricsZone)
	if err != nil {
		t.Fatalf("Failed to execute metrics zone request: %v", err)
	}
	if respMetricsZone.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200 for metrics DMA-03, got %d", respMetricsZone.StatusCode)
	}

	var resMetricsZone struct {
		Status           int `json:"status"`
		ModelPerformance struct {
			Mode    string                 `json:"mode"`
			DMACode string                 `json:"dma_code"`
			Metrics map[string]interface{} `json:"metrics"`
		} `json:"model_performance"`
	}
	bodyMetricsZone, _ := io.ReadAll(respMetricsZone.Body)
	if err := json.Unmarshal(bodyMetricsZone, &resMetricsZone); err != nil {
		t.Fatalf("Failed to decode metrics zone response: %v", err)
	}
	if resMetricsZone.ModelPerformance.Mode != "zone" || resMetricsZone.ModelPerformance.DMACode != "DMA-03" {
		t.Fatalf("Expected mode 'zone' and dma_code 'DMA-03', got %s, %s",
			resMetricsZone.ModelPerformance.Mode, resMetricsZone.ModelPerformance.DMACode)
	}

	// 6. ทดสอบ GET /api/v1/forecast/metrics?dma_code=INVALID-ZONE (ต้องได้ 404)
	reqMetricsInvalid := httptest.NewRequest(http.MethodGet, "/api/v1/forecast/metrics?dma_code=INVALID-ZONE", nil)
	respMetricsInvalid, _ := app.Test(reqMetricsInvalid)
	if respMetricsInvalid.StatusCode != fiber.StatusNotFound {
		t.Fatalf("Expected status 404 for invalid zone metrics, got %d", respMetricsInvalid.StatusCode)
	}
}
