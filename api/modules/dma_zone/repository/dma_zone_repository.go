package repository

import (
	"api/modules/entities"

	"gorm.io/gorm"
)

type DMAZoneRepository struct {
	db *gorm.DB
}

// Create implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Create(data entities.TDMAZone) (entities.TDMAZone, error) {
	if err := d.db.Create(&data).Error; err != nil {
		return entities.TDMAZone{}, err
	}
	return data, nil
}

// Delete implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Delete(id uint) (entities.TDMAZone, error) {
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
	if page != nil && limit != nil {
		query = query.Offset((*page - 1) * *limit).Limit(*limit)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Update(id uint, data entities.TDMAZone) (entities.TDMAZone, error) {
	_, err := d.GetById(id)
	if err != nil {
		return entities.TDMAZone{}, nil
	}
	if err := d.db.Save(&data).Error; err != nil {
		return entities.TDMAZone{}, err
	}
	return data, nil
}

func NewDMAZoneRepository(db *gorm.DB) entities.TDMAZoneRepo {
	return DMAZoneRepository{db: db}
}
