package repository

import (
	"api/modules/entities"

	"gorm.io/gorm"
)

type PipeAssetRepository struct {
	db *gorm.DB
}

// Create implements [entities.TPipeAssetRepo].
func (p PipeAssetRepository) Create(data entities.TPipeAssetReq) (entities.TPipeAssets, error) {
	pipe := entities.TPipeAssets{
		AssetID:          data.AssetID,
		PipeID:           data.PipeID,
		DMACode:          data.DMACode,
		Material:         data.Material,
		Diameter:         data.Diameter,
		Length:           data.Length,
		InstallationYear: data.InstallationYear,
		Latitude:         data.Latitude,
		Longitude:        data.Longitude,
		Status:           data.Status,
	}
	if err := p.db.Create(&pipe).Error; err != nil {
		return entities.TPipeAssets{}, err
	}
	return pipe, nil
}

// DeleteById implements [entities.TPipeAssetRepo].
func (p PipeAssetRepository) DeleteById(id uint) (entities.TPipeAssets, error) {
	data, err := p.GetById(id)
	if err != nil {
		return entities.TPipeAssets{}, err
	}
	if err := p.db.Delete(&data).Error; err != nil {
		return entities.TPipeAssets{}, err
	}
	return data, nil
}

// GetById implements [entities.TPipeAssetRepo].
func (p PipeAssetRepository) GetById(id uint) (entities.TPipeAssets, error) {
	var data entities.TPipeAssets
	if err := p.db.First(&data, id).Error; err != nil {
		return entities.TPipeAssets{}, err
	}
	return data, nil
}

// Gets implements [entities.TPipeAssetRepo].
func (p PipeAssetRepository) Gets(page *int, limit *int) ([]entities.TPipeAssets, error) {
	var data []entities.TPipeAssets
	query := p.db.Model(&entities.TPipeAssets{})
	if page != nil && limit != nil && *page > 0 && *limit > 0 {
		query = query.Offset((*page - 1) * *limit).Limit(*limit)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateById implements [entities.TPipeAssetRepo].
func (p PipeAssetRepository) UpdateById(id uint, data entities.TPipeAssetReq) (entities.TPipeAssets, error) {
	pipe, err := p.GetById(id)
	if err != nil {
		return entities.TPipeAssets{}, err
	}

	pipe.AssetID = data.AssetID
	pipe.PipeID = data.PipeID
	pipe.DMACode = data.DMACode
	pipe.Material = data.Material
	pipe.Diameter = data.Diameter
	pipe.Length = data.Length
	pipe.InstallationYear = data.InstallationYear
	pipe.Latitude = data.Latitude
	pipe.Longitude = data.Longitude
	pipe.Status = data.Status

	if err := p.db.Save(&pipe).Error; err != nil {
		return entities.TPipeAssets{}, err
	}
	return pipe, nil
}

func NewPipeAssetRepository(db *gorm.DB) entities.TPipeAssetRepo {
	return PipeAssetRepository{db: db}
}
