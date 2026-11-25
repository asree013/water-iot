package usecase

import "api/modules/entities"

type PipeAssetUseCase struct {
	pipeAssetRepo entities.TPipeAssetRepo
}

// AddData implements [entities.TPipeAssetUseCase].
func (u PipeAssetUseCase) AddData(data entities.TPipeAssetReq) (entities.TPipeAssets, error) {
	return u.pipeAssetRepo.Create(data)
}

// DeleteData implements [entities.TPipeAssetUseCase].
func (u PipeAssetUseCase) DeleteData(id uint) (entities.TPipeAssets, error) {
	return u.pipeAssetRepo.DeleteById(id)
}

// FindDataById implements [entities.TPipeAssetUseCase].
func (u PipeAssetUseCase) FindDataById(id uint) (entities.TPipeAssets, error) {
	return u.pipeAssetRepo.GetById(id)
}

// FindDatas implements [entities.TPipeAssetUseCase].
func (u PipeAssetUseCase) FindDatas(page *int, limit *int) ([]entities.TPipeAssets, error) {
	return u.pipeAssetRepo.Gets(page, limit)
}

// UpdateData implements [entities.TPipeAssetUseCase].
func (u PipeAssetUseCase) UpdateData(id uint, data entities.TPipeAssetReq) (entities.TPipeAssets, error) {
	return u.pipeAssetRepo.UpdateById(id, data)
}

func NewPipeAssetUseCase(pipeAssetRepo entities.TPipeAssetRepo) entities.TPipeAssetUseCase {
	return PipeAssetUseCase{
		pipeAssetRepo: pipeAssetRepo,
	}
}
