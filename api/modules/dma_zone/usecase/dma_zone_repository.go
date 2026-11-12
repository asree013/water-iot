package usecase

import "api/modules/entities"

type DMAZoneUseCase struct {
	TDMAZoneRepo entities.TDMAZoneRepo
}

// AddData implements [entities.TDMAZoneUseCase].
func (d DMAZoneUseCase) AddData(data entities.TDMAZone) (entities.TDMAZone, error) {
	return d.TDMAZoneRepo.Create(data)
}

// DeleteData implements [entities.TDMAZoneUseCase].
func (d DMAZoneUseCase) DeleteData(id uint) error {
	_, err := d.TDMAZoneRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
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
func (d DMAZoneUseCase) UpdateData(id uint, data entities.TDMAZone) (entities.TDMAZone, error) {
	return d.TDMAZoneRepo.Update(id, data)
}

func NewDMAZoneUseCase(TDMAZoneRepo entities.TDMAZoneRepo) entities.TDMAZoneUseCase {
	return DMAZoneUseCase{
		TDMAZoneRepo: TDMAZoneRepo,
	}
}
