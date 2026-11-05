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
	data, err := w.FindById(id)
	if err != nil {
		return entities.WaterData{}, nil
	}
	return data, nil
}

// FindAll implements [entities.WaterDataRepository].
func (w *waterDataRepository) FindAll() ([]entities.WaterData, error) {
	var list []entities.WaterData
	if err := w.db.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindById implements [entities.WaterDataRepository].
func (w *waterDataRepository) FindById(id uint) (entities.WaterData, error) {
	var data entities.WaterData
	if err := w.db.First(&data, id).Error; err != nil {
		return entities.WaterData{}, err
	}
	return data, nil
}

// Save implements [entities.WaterDataRepository].
func (w *waterDataRepository) Save(data entities.WaterData) (entities.WaterData, error) {
	if err := w.db.Create(&data).Error; err != nil {
		return entities.WaterData{}, err
	}
	return data, nil
}

// Update implements [entities.WaterDataRepository].
func (w *waterDataRepository) Update(data entities.WaterData) (entities.WaterData, error) {
	if err := w.db.Save(&data).Error; err != nil {
		return entities.WaterData{}, err
	}
	return data, nil
}
