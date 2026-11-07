package usecases

import (
	"api/modules/entities"
	"time"
)

type waterDataUseCase struct {
	waterDataRepo entities.WaterDataRepository
}

// Constructor
func NewWaterDataUseCase(waterDataRepo entities.WaterDataRepository) entities.WaterDataUseCase {
	return &waterDataUseCase{
		waterDataRepo: waterDataRepo,
	}
}

// GetLatestTelemetry implements [entities.WaterDataUseCase].
func (w *waterDataUseCase) GetLatestTelemetry() ([]entities.WaterData, error) {
	return w.waterDataRepo.FindAll()
}

// SaveTelemetry implements [entities.WaterDataUseCase].
func (w *waterDataUseCase) SaveTelemetry(req *entities.WaterDataReq) error {
	data := entities.WaterData{
		Pressure:         req.Pressure,
		FlowRate:         req.FlowRate,
		Temperature:      req.Temperature,
		Vibration:        req.Vibration,
		RPM:              req.RPM,
		OperationalHours: req.OperationalHours,
		Zone:             req.Zone,
		Block:            req.Block,
		Pipe:             req.Pipe,
		LocationCode:     req.LocationCode,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		LeakageFlag:      req.LeakageFlag,
		Timestamp:        req.Timestamp,
		CreatedAt:        time.Now(),
	}
	_, err := w.waterDataRepo.Save(data)
	if err != nil {
		return err
	}
	return nil
}
