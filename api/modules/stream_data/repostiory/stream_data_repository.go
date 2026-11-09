package repostiory

import (
	"api/modules/entities"
	waterRepo "api/modules/water_data/repository"

	"gorm.io/gorm"
)

func NewStreamDataRepository(db *gorm.DB) entities.IStreamDataRepo {
	return waterRepo.NewWaterDataRepository(db)
}
