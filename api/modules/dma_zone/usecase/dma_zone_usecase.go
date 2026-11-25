package usecase

import "api/modules/entities"

type DMAZoneUseCase struct {
	TDMAZoneRepo entities.TDMAZoneRepo
}

// AddData implements [entities.TDMAZoneUseCase].
func (d DMAZoneUseCase) AddData(data entities.TDMAZoneReq) (entities.TDMAZone, error) {
	return d.TDMAZoneRepo.Create(data)
}

// DeleteData implements [entities.TDMAZoneUseCase].
func (d DMAZoneUseCase) DeleteData(id uint) (entities.TDMAZone, error) {
	return d.TDMAZoneRepo.DeleteById(id)
}

// FindDataById implements [entities.TDMAZoneUseCase].
func (d DMAZoneUseCase) FindDataById(id uint) (entities.TDMAZone, error) {
	return d.TDMAZoneRepo.GetById(id)
}

// FindDatas implements [entities.TDMAZoneUseCase].
func (d DMAZoneUseCase) FindDatas(page *int, limit *int) ([]entities.TDMAZone, error) {
	return d.TDMAZoneRepo.Gets(page, limit)
}

// UpdateData implements [entities.TDMAZoneUseCase].
func (d DMAZoneUseCase) UpdateData(id uint, data entities.TDMAZoneReq) (entities.TDMAZone, error) {
	return d.TDMAZoneRepo.UpdateById(id, data)
}

func NewDMAZoneUseCase(TDMAZoneRepo entities.TDMAZoneRepo) entities.TDMAZoneUseCase {
	return DMAZoneUseCase{
		TDMAZoneRepo: TDMAZoneRepo,
	}
}
