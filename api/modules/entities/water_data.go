package entities

import "time"

// Domain Entity
type WaterData struct {
	ID               uint      `gorm:"primarykey;autoIncrement" json:"id"`
	Pressure         float64   `json:"pressure"`
	FlowRate         float64   `json:"flow_rate"`
	Temperature      float64   `json:"temperature"`
	Vibration        float64   `json:"vibration"`
	RPM              float64   `json:"rpm"`
	OperationalHours int       `json:"operational_hours"`
	Zone             string    `json:"zone"`
	Block            string    `json:"block"`
	Pipe             string    `json:"pipe"`
	LocationCode     string    `json:"location_code"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	LeakageFlag      int       `json:"leakage_flag"`
	Timestamp        int64     `json:"timestamp"`
	CreatedAt        time.Time `json:"created_at"`
}

// DTO

type WaterDataReq struct {
	Pressure         float64 `json:"pressure"`
	FlowRate         float64 `json:"flow_rate"`
	Temperature      float64 `json:"temperature"`
	Vibration        float64 `json:"vibration"`
	RPM              float64 `json:"rpm"`
	OperationalHours int     `json:"operational_hours"`
	Zone             string  `json:"zone"`
	Block            string  `json:"block"`
	Pipe             string  `json:"pipe"`
	LocationCode     string  `json:"location_code"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	LeakageFlag      int     `json:"leakage_flag"`
	Timestamp        int64   `json:"timestamp"`
}

type WaterDataRes struct {
	WaterDataReq
}

// Repository Interface

type WaterDataRepository interface {
	Save(data WaterData) (WaterData, error)
	FindById(id uint) (WaterData, error)
	FindAll() ([]WaterData, error)
	Update(data WaterData) (WaterData, error)
	Delete(id uint) (WaterData, error)
}

// UseCase Interface

type WaterDataUseCase interface {
	SaveTelemetry(req *WaterDataReq) error
	GetLatestTelemetry() ([]WaterData, error)
}
