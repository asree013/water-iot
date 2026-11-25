package repository

import (
	"api/modules/entities"

	"gorm.io/gorm"
)

type DMAZoneRepository struct {
	db *gorm.DB
}

// Create implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Create(data entities.TDMAZoneReq) (entities.TDMAZone, error) {
	zone := entities.TDMAZone{
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
	}
	if err := d.db.Create(&zone).Error; err != nil {
		return entities.TDMAZone{}, err
	}
	return zone, nil
}

// DeleteById implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) DeleteById(id uint) (entities.TDMAZone, error) {
	data, err := d.GetById(id)
	if err != nil {
		return entities.TDMAZone{}, err
	}
	if err := d.db.Delete(&data).Error; err != nil {
		return entities.TDMAZone{}, err
	}
	return data, nil
}

// GetById implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) GetById(id uint) (entities.TDMAZone, error) {
	var data entities.TDMAZone
	if err := d.db.First(&data, id).Error; err != nil {
		return entities.TDMAZone{}, err
	}
	return data, nil
}

// Gets implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Gets(page *int, limit *int) ([]entities.TDMAZone, error) {
	var data []entities.TDMAZone
	query := d.db.Model(&entities.TDMAZone{})
	if page != nil && limit != nil && *page > 0 && *limit > 0 {
		query = query.Offset((*page - 1) * *limit).Limit(*limit)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateById implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) UpdateById(id uint, data entities.TDMAZoneReq) (entities.TDMAZone, error) {
	zone, err := d.GetById(id)
	if err != nil {
		return entities.TDMAZone{}, err
	}

	zone.DMACode = data.DMACode
	zone.DMAName = data.DMAName
	zone.BranchName = data.BranchName
	zone.RegionCode = data.RegionCode
	zone.TargetPressure = data.TargetPressure
	zone.TotalInflow = data.TotalInflow
	zone.BilledConsumption = data.BilledConsumption
	zone.NRWRatio = data.NRWRatio
	zone.CustomerCount = data.CustomerCount
	zone.Status = data.Status

	if err := d.db.Save(&zone).Error; err != nil {
		return entities.TDMAZone{}, err
	}
	return zone, nil
}

func NewDMAZoneRepository(db *gorm.DB) entities.TDMAZoneRepo {
	return DMAZoneRepository{db: db}
}
