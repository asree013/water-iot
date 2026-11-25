package repository

import (
	"api/modules/entities"

	"gorm.io/gorm"
)

type ForecastRepository struct {
	db *gorm.DB
}

func (r *ForecastRepository) GetForecastByDMA(dmaCode string, limit int) ([]entities.TDemandForecast, error) {
	var list []entities.TDemandForecast
	if err := r.db.Where("dma_code = ?", dmaCode).
		Order("forecast_timestamp ASC").
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ForecastRepository) GetZoneByCode(dmaCode string) (entities.TDMAZone, error) {
	var zone entities.TDMAZone
	if err := r.db.Where("dma_code = ?", dmaCode).First(&zone).Error; err != nil {
		return entities.TDMAZone{}, err
	}
	return zone, nil
}

func (r *ForecastRepository) SaveForecastBatch(items []entities.TDemandForecast) error {
	return r.db.CreateInBatches(items, 100).Error
}

func NewForecastRepository(db *gorm.DB) entities.TForecastRepo {
	return &ForecastRepository{db: db}
}
