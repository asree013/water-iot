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
	panic("unimplemented")
}

// Delete implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Delete(id int) (entities.TDMAZone, error) {
	panic("unimplemented")
}

// GetById implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) GetById(id int) (entities.TDMAZone, error) {
	panic("unimplemented")
}

// Gets implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Gets(page *int, limit *int) ([]entities.TDMAZone, error) {
	panic("unimplemented")
}

// Update implements [entities.TDMAZoneRepo].
func (d DMAZoneRepository) Update(data entities.TDMAZone) (entities.TDMAZone, error) {
	panic("unimplemented")
}

func NewDMAZoneRepository(db *gorm.DB) entities.TDMAZoneRepo {
	return DMAZoneRepository{db: db}
}
