package repository

import (
	"api/modules/entities"

	"gorm.io/gorm"
)

// connection go
type waterDataRepository struct {
	db *gorm.DB
}

// Constructor
func NewWaterDataRepository(db *gorm.DB) entities.WaterDataRepository {
	return &waterDataRepository{
		db: db,
	}
}

// Delete implements [entities.WaterDataRepository].
func (w *waterDataRepository) Delete(id uint) (entities.WaterData, error) {
	panic("unimplemented")
}

// FindAll implements [entities.WaterDataRepository].
func (w *waterDataRepository) FindAll() ([]entities.WaterData, error) {
	panic("unimplemented")
}

// FindById implements [entities.WaterDataRepository].
func (w *waterDataRepository) FindById(id uint) (entities.WaterData, error) {
	panic("unimplemented")
}

// Save implements [entities.WaterDataRepository].
func (w *waterDataRepository) Save(data entities.WaterData) (entities.WaterData, error) {
	panic("unimplemented")
}

// Update implements [entities.WaterDataRepository].
func (w *waterDataRepository) Update(data entities.WaterData) (entities.WaterData, error) {
	panic("unimplemented")
}

// type WaterData struct {
// 	Pressure         float64 `json:"pressure"`
// 	FlowRate         float64 `json:"flow_rate"`
// 	Temperature      float64 `json:"temperature"`
// 	Vibration        float64 `json:"vibration"`
// 	RPM              float64 `json:"rpm"`
// 	OperationalHours int     `json:"operational_hours"`
// 	Zone             string  `json:"zone"`
// 	Block            string  `json:"block"`
// 	Pipe             string  `json:"pipe"`
// 	LocationCode     string  `json:"location_code"`
// 	Latitude         float64 `json:"latitude"`
// 	Longitude        float64 `json:"longitude"`
// 	LeakageFlag      int     `json:"leakage_flag"`
// 	Timestamp        int64   `json:"timestamp"` // Unix timestamp เพิ่มความสมจริงของ IoT
// }
