package repository

import (
	"api/modules/entities"
	"errors"

	"gorm.io/gorm"
)

type PumpRepository struct {
	db *gorm.DB
}

func NewPumpRepository(db *gorm.DB) entities.TPumpRepo {
	return &PumpRepository{db: db}
}

func (r *PumpRepository) GetAllPumps() ([]entities.TPumpAsset, error) {
	var pumps []entities.TPumpAsset
	if err := r.db.Find(&pumps).Error; err != nil {
		return nil, err
	}
	return pumps, nil
}

func (r *PumpRepository) GetPumpByID(pumpID string) (entities.TPumpAsset, error) {
	var pump entities.TPumpAsset
	if err := r.db.Where("pump_id = ?", pumpID).First(&pump).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pump, errors.New("pump not found")
		}
		return pump, err
	}
	return pump, nil
}

func (r *PumpRepository) GetAllPumpHealth() ([]entities.TPumpHealthScore, error) {
	var list []entities.TPumpHealthScore
	if err := r.db.Order("health_score ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PumpRepository) GetPumpHealthByID(pumpID string) (entities.TPumpHealthScore, error) {
	var score entities.TPumpHealthScore
	if err := r.db.Where("pump_id = ?", pumpID).Order("evaluated_at DESC").First(&score).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return score, errors.New("pump health record not found")
		}
		return score, err
	}
	return score, nil
}

func (r *PumpRepository) SavePumpHealthBatch(items []entities.TPumpHealthScore) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}
