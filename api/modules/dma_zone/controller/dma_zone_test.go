package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"api/modules/dma_zone/controller"
	"api/modules/dma_zone/usecase"
	"api/modules/entities"
	"api/pkg/database"

	"github.com/gofiber/fiber/v2"
)

// MockDMAZoneRepo จำลอง Repository ใน Memory สำหรับทดสอบ CRUD
type MockDMAZoneRepo struct {
	mu     sync.RWMutex
	data   map[uint]entities.TDMAZone
	nextID uint
}

func NewMockDMAZoneRepo() entities.TDMAZoneRepo {
	return &MockDMAZoneRepo{
		data:   make(map[uint]entities.TDMAZone),
		nextID: 1,
	}
}

func (m *MockDMAZoneRepo) Create(data entities.TDMAZoneReq) (entities.TDMAZone, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := entities.TDMAZone{
		ID:                m.nextID,
		DMACode:           data.DMACode,
		DMAName:           data.DMAName,
		BranchName:        data.BranchName,
		RegionCode:        data.RegionCode,
		TargetPressure:    data.TargetPressure,
		TotalInflow:       data.TotalInflow,
		BilledConsumption: data.BilledConsumption,
		NRWRatio:          data.NRWRatio,
		CustomerCount:     data.CustomerCount,
		Status:            data.Status,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	m.data[m.nextID] = item
	m.nextID++
	return item, nil
}

func (m *MockDMAZoneRepo) Gets(page *int, limit *int) ([]entities.TDMAZone, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []entities.TDMAZone
	for _, v := range m.data {
		list = append(list, v)
	}
	return list, nil
}

func (m *MockDMAZoneRepo) GetById(id uint) (entities.TDMAZone, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TDMAZone{}, errors.New("record not found")
	}
	return item, nil
}

func (m *MockDMAZoneRepo) UpdateById(id uint, data entities.TDMAZoneReq) (entities.TDMAZone, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TDMAZone{}, errors.New("record not found")
	}
	item.DMACode = data.DMACode
	item.DMAName = data.DMAName
	item.BranchName = data.BranchName
	item.RegionCode = data.RegionCode
	item.TargetPressure = data.TargetPressure
	item.TotalInflow = data.TotalInflow
	item.BilledConsumption = data.BilledConsumption
	item.NRWRatio = data.NRWRatio
	item.CustomerCount = data.CustomerCount
	item.Status = data.Status
	item.UpdatedAt = time.Now()

	m.data[id] = item
	return item, nil
}

func (m *MockDMAZoneRepo) DeleteById(id uint) (entities.TDMAZone, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TDMAZone{}, errors.New("record not found")
	}
	delete(m.data, id)
	return item, nil
}

func setupDMAZoneApp() (*fiber.App, entities.TDMAZoneRepo) {
	app := fiber.New()
	mockRepo := NewMockDMAZoneRepo()
	uc := usecase.NewDMAZoneUseCase(mockRepo)
	group := app.Group("/api/v1/dma-zone")
	controller.NewDMAZoneController(group, uc)
	return app, mockRepo
}

func TestDMAZoneCRUD_MockDataProvinces(t *testing.T) {
	app, _ := setupDMAZoneApp()
	mockZones := database.GetMockDMAZones() // สงขลา, เพชรบุรี, กาญจนบุรี, กทม

	// 1. CREATE: ทดสอบสร้าง DMA Zone ทั้ง 4 จังหวัด
	for _, zoneReq := range mockZones {
		body, _ := json.Marshal(zoneReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/dma-zone", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusCreated {
			t.Fatalf("Expected status 201, got %d", resp.StatusCode)
		}

		var resBody map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&resBody)
		if resBody["message"] != "DMA Zone Created Successfully" {
			t.Fatalf("Unexpected message: %v", resBody["message"])
		}
	}

	// 2. READ ALL (GET): ทดสอบดึงข้อมูลทั้งหมด ควรได้ 4 โซน
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/dma-zone", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		var resBody struct {
			Data []entities.TDMAZone `json:"data"`
		}
		_ = json.Unmarshal(bodyBytes, &resBody)
		if len(resBody.Data) != len(mockZones) {
			t.Fatalf("Expected %d zones, got %d", len(mockZones), len(resBody.Data))
		}
	}

	// 3. READ BY ID (GET /:id): ทดสอบดึงข้อมูล ID 1 (สงขลา)
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/dma-zone/1", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var resBody struct {
			Data entities.TDMAZone `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&resBody)
		if resBody.Data.DMACode != "DMA-SK-01" {
			t.Fatalf("Expected DMA-SK-01, got %s", resBody.Data.DMACode)
		}
		if resBody.Data.DMAName != "โซนเทศบาลนครสงขลา" {
			t.Fatalf("Expected โซนเทศบาลนครสงขลา, got %s", resBody.Data.DMAName)
		}
		if resBody.Data.BranchName != "สาขาสงขลา" {
			t.Fatalf("Expected สาขาสงขลา, got %s", resBody.Data.BranchName)
		}
	}

	// 4. UPDATE (PUT /:id): ทดสอบอัปเดตสถานะของโซน ID 1
	{
		updatedData := mockZones[0]
		updatedData.Status = "WARNING"
		updatedData.TargetPressure = 2.4

		body, _ := json.Marshal(updatedData)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/dma-zone/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var resBody struct {
			Data entities.TDMAZone `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&resBody)
		if resBody.Data.Status != "WARNING" {
			t.Fatalf("Expected status WARNING, got %s", resBody.Data.Status)
		}
		if resBody.Data.TargetPressure != 2.4 {
			t.Fatalf("Expected pressure 2.4, got %f", resBody.Data.TargetPressure)
		}
	}

	// 5. DELETE (DELETE /:id): ทดสอบลบโซน ID 1
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/dma-zone/1", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// ตรวจสอบว่าดึงอีกครั้งต้องเป็น 404 Not Found
		reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/dma-zone/1", nil)
		respGet, _ := app.Test(reqGet)
		if respGet.StatusCode != fiber.StatusNotFound {
			t.Fatalf("Expected status 404 after deletion, got %d", respGet.StatusCode)
		}
	}
}
