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

	"api/modules/entities"
	"api/modules/leakage_incident/controller"
	"api/modules/leakage_incident/usecase"
	"api/pkg/database"

	"github.com/gofiber/fiber/v2"
)

// MockLeakageIncidentRepo จำลอง Repository ใน Memory สำหรับทดสอบ CRUD
type MockLeakageIncidentRepo struct {
	mu     sync.RWMutex
	data   map[uint]entities.TLeakageIncident
	nextID uint
}

func NewMockLeakageIncidentRepo() entities.TLeakageIncidentRepo {
	return &MockLeakageIncidentRepo{
		data:   make(map[uint]entities.TLeakageIncident),
		nextID: 1,
	}
}

func (m *MockLeakageIncidentRepo) Create(data entities.TLeakageIncidentReq) (entities.TLeakageIncident, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := entities.TLeakageIncident{
		ID:                  m.nextID,
		TicketID:            data.TicketID,
		DMACode:             data.DMACode,
		PipeID:              data.PipeID,
		IncidentType:        data.IncidentType,
		Severity:            data.Severity,
		LocationDescription: data.LocationDescription,
		Latitude:            data.Latitude,
		Longitude:           data.Longitude,
		WaterLossEstimate:   data.WaterLossEstimate,
		ReportedAt:          data.ReportedAt,
		RepairedAt:          data.RepairedAt,
		RepairStatus:        data.RepairStatus,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	m.data[m.nextID] = item
	m.nextID++
	return item, nil
}

func (m *MockLeakageIncidentRepo) Gets(page *int, limit *int) ([]entities.TLeakageIncident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []entities.TLeakageIncident
	for _, v := range m.data {
		list = append(list, v)
	}
	return list, nil
}

func (m *MockLeakageIncidentRepo) GetById(id uint) (entities.TLeakageIncident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TLeakageIncident{}, errors.New("record not found")
	}
	return item, nil
}

func (m *MockLeakageIncidentRepo) UpdateById(id uint, data entities.TLeakageIncidentReq) (entities.TLeakageIncident, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TLeakageIncident{}, errors.New("record not found")
	}
	item.TicketID = data.TicketID
	item.DMACode = data.DMACode
	item.PipeID = data.PipeID
	item.IncidentType = data.IncidentType
	item.Severity = data.Severity
	item.LocationDescription = data.LocationDescription
	item.Latitude = data.Latitude
	item.Longitude = data.Longitude
	item.WaterLossEstimate = data.WaterLossEstimate
	item.ReportedAt = data.ReportedAt
	item.RepairedAt = data.RepairedAt
	item.RepairStatus = data.RepairStatus
	item.UpdatedAt = time.Now()

	m.data[id] = item
	return item, nil
}

func (m *MockLeakageIncidentRepo) DeleteById(id uint) (entities.TLeakageIncident, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TLeakageIncident{}, errors.New("record not found")
	}
	delete(m.data, id)
	return item, nil
}

func setupLeakageIncidentApp() (*fiber.App, entities.TLeakageIncidentRepo) {
	app := fiber.New()
	mockRepo := NewMockLeakageIncidentRepo()
	uc := usecase.NewLeakageIncidentUseCase(mockRepo)
	group := app.Group("/api/v1/leakage-incident")
	controller.NewLeakageIncidentController(group, uc)
	return app, mockRepo
}

func TestLeakageIncidentCRUD_MockDataProvinces(t *testing.T) {
	app, _ := setupLeakageIncidentApp()
	mockIncidents := database.GetMockLeakageIncidents() // เหตุการณ์ใน สงขลา, เพชรบุรี, กาญจนบุรี, กทม

	// 1. CREATE: ทดสอบสร้างเหตุการณ์ท่อแตกรั่วทั้ง 4 จังหวัด
	for _, incReq := range mockIncidents {
		body, _ := json.Marshal(incReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/leakage-incident", bytes.NewReader(body))
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
		if resBody["message"] != "Leakage Incident Created Successfully" {
			t.Fatalf("Unexpected message: %v", resBody["message"])
		}
	}

	// 2. READ ALL (GET): ทดสอบดึงข้อมูลทั้งหมด ควรได้ 4 รายการ
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/leakage-incident", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		var resBody struct {
			Data []entities.TLeakageIncident `json:"data"`
		}
		_ = json.Unmarshal(bodyBytes, &resBody)
		if len(resBody.Data) != len(mockIncidents) {
			t.Fatalf("Expected %d incidents, got %d", len(mockIncidents), len(resBody.Data))
		}
	}

	// 3. READ BY ID (GET /:id): ทดสอบดึงข้อมูล ID 1 (สงขลา)
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/leakage-incident/1", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var resBody struct {
			Data entities.TLeakageIncident `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&resBody)
		if resBody.Data.TicketID != "OCS-2026-SK01" {
			t.Fatalf("Expected OCS-2026-SK01, got %s", resBody.Data.TicketID)
		}
		if resBody.Data.IncidentType != "Main Burst" {
			t.Fatalf("Expected Main Burst, got %s", resBody.Data.IncidentType)
		}
		if resBody.Data.Severity != "HIGH" {
			t.Fatalf("Expected HIGH, got %s", resBody.Data.Severity)
		}
		if resBody.Data.DMACode != "DMA-SK-01" {
			t.Fatalf("Expected DMA-SK-01, got %s", resBody.Data.DMACode)
		}
	}

	// 4. UPDATE (PUT /:id): ทดสอบอัปเดตสถานะการซ่อมของ ID 1
	{
		updatedData := mockIncidents[0]
		now := time.Now()
		updatedData.RepairedAt = &now
		updatedData.RepairStatus = "CLOSED"

		body, _ := json.Marshal(updatedData)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/leakage-incident/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var resBody struct {
			Data entities.TLeakageIncident `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&resBody)
		if resBody.Data.RepairStatus != "CLOSED" {
			t.Fatalf("Expected status CLOSED, got %s", resBody.Data.RepairStatus)
		}
		if resBody.Data.RepairedAt == nil {
			t.Fatalf("Expected non-nil RepairedAt")
		}
	}

	// 5. DELETE (DELETE /:id): ทดสอบลบรายการ ID 1
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/leakage-incident/1", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// ตรวจสอบว่าดึงอีกครั้งต้องเป็น 404 Not Found
		reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/leakage-incident/1", nil)
		respGet, _ := app.Test(reqGet)
		if respGet.StatusCode != fiber.StatusNotFound {
			t.Fatalf("Expected status 404 after deletion, got %d", respGet.StatusCode)
		}
	}
}
