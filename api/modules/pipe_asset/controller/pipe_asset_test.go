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
	"api/modules/pipe_asset/controller"
	"api/modules/pipe_asset/usecase"
	"api/pkg/database"

	"github.com/gofiber/fiber/v2"
)

// MockPipeAssetRepo จำลอง Repository ใน Memory สำหรับทดสอบ CRUD
type MockPipeAssetRepo struct {
	mu     sync.RWMutex
	data   map[uint]entities.TPipeAssets
	nextID uint
}

func NewMockPipeAssetRepo() entities.TPipeAssetRepo {
	return &MockPipeAssetRepo{
		data:   make(map[uint]entities.TPipeAssets),
		nextID: 1,
	}
}

func (m *MockPipeAssetRepo) Create(data entities.TPipeAssetReq) (entities.TPipeAssets, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := entities.TPipeAssets{
		ID:               m.nextID,
		AssetID:          data.AssetID,
		PipeID:           data.PipeID,
		DMACode:          data.DMACode,
		Material:         data.Material,
		Diameter:         data.Diameter,
		Length:           data.Length,
		InstallationYear: data.InstallationYear,
		Latitude:         data.Latitude,
		Longitude:        data.Longitude,
		Status:           data.Status,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	m.data[m.nextID] = item
	m.nextID++
	return item, nil
}

func (m *MockPipeAssetRepo) Gets(page *int, limit *int) ([]entities.TPipeAssets, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []entities.TPipeAssets
	for _, v := range m.data {
		list = append(list, v)
	}
	return list, nil
}

func (m *MockPipeAssetRepo) GetById(id uint) (entities.TPipeAssets, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TPipeAssets{}, errors.New("record not found")
	}
	return item, nil
}

func (m *MockPipeAssetRepo) UpdateById(id uint, data entities.TPipeAssetReq) (entities.TPipeAssets, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TPipeAssets{}, errors.New("record not found")
	}
	item.AssetID = data.AssetID
	item.PipeID = data.PipeID
	item.DMACode = data.DMACode
	item.Material = data.Material
	item.Diameter = data.Diameter
	item.Length = data.Length
	item.InstallationYear = data.InstallationYear
	item.Latitude = data.Latitude
	item.Longitude = data.Longitude
	item.Status = data.Status
	item.UpdatedAt = time.Now()

	m.data[id] = item
	return item, nil
}

func (m *MockPipeAssetRepo) DeleteById(id uint) (entities.TPipeAssets, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[id]
	if !ok {
		return entities.TPipeAssets{}, errors.New("record not found")
	}
	delete(m.data, id)
	return item, nil
}

func setupPipeAssetApp() (*fiber.App, entities.TPipeAssetRepo) {
	app := fiber.New()
	mockRepo := NewMockPipeAssetRepo()
	uc := usecase.NewPipeAssetUseCase(mockRepo)
	group := app.Group("/api/v1/pipe-asset")
	controller.NewPipeAssetController(group, uc)
	return app, mockRepo
}

func TestPipeAssetCRUD_MockDataProvinces(t *testing.T) {
	app, _ := setupPipeAssetApp()
	mockPipes := database.GetMockPipeAssets() // ท่อใน สงขลา, เพชรบุรี, กาญจนบุรี, กทม

	// 1. CREATE: ทดสอบสร้าง Pipe Asset ทั้ง 4 จังหวัด
	for _, pipeReq := range mockPipes {
		body, _ := json.Marshal(pipeReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/pipe-asset", bytes.NewReader(body))
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
		if resBody["message"] != "Pipe Asset Created Successfully" {
			t.Fatalf("Unexpected message: %v", resBody["message"])
		}
	}

	// 2. READ ALL (GET): ทดสอบดึงข้อมูลทั้งหมด ควรได้ 4 ท่อ
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pipe-asset", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		var resBody struct {
			Data []entities.TPipeAssets `json:"data"`
		}
		_ = json.Unmarshal(bodyBytes, &resBody)
		if len(resBody.Data) != len(mockPipes) {
			t.Fatalf("Expected %d pipes, got %d", len(mockPipes), len(resBody.Data))
		}
	}

	// 3. READ BY ID (GET /:id): ทดสอบดึงท่อ ID 1 (สงขลา)
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pipe-asset/1", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var resBody struct {
			Data entities.TPipeAssets `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&resBody)
		if resBody.Data.AssetID != "AST-SK-001" {
			t.Fatalf("Expected AST-SK-001, got %s", resBody.Data.AssetID)
		}
		if resBody.Data.PipeID != "PIPE-SK-01" {
			t.Fatalf("Expected PIPE-SK-01, got %s", resBody.Data.PipeID)
		}
		if resBody.Data.DMACode != "DMA-SK-01" {
			t.Fatalf("Expected DMA-SK-01, got %s", resBody.Data.DMACode)
		}
	}

	// 4. UPDATE (PUT /:id): ทดสอบอัปเดตข้อมูลท่อ ID 1
	{
		updatedData := mockPipes[0]
		updatedData.Diameter = 250.0
		updatedData.Status = "REPLACED"

		body, _ := json.Marshal(updatedData)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/pipe-asset/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var resBody struct {
			Data entities.TPipeAssets `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&resBody)
		if resBody.Data.Diameter != 250.0 {
			t.Fatalf("Expected diameter 250.0, got %f", resBody.Data.Diameter)
		}
		if resBody.Data.Status != "REPLACED" {
			t.Fatalf("Expected status REPLACED, got %s", resBody.Data.Status)
		}
	}

	// 5. DELETE (DELETE /:id): ทดสอบลบท่อ ID 1
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/pipe-asset/1", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// ตรวจสอบว่าดึงอีกครั้งต้องเป็น 404 Not Found
		reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/pipe-asset/1", nil)
		respGet, _ := app.Test(reqGet)
		if respGet.StatusCode != fiber.StatusNotFound {
			t.Fatalf("Expected status 404 after deletion, got %d", respGet.StatusCode)
		}
	}
}
