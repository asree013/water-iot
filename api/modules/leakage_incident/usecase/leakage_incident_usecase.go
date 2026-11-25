package usecase

import "api/modules/entities"

type LeakageIncidentUseCase struct {
	incidentRepo entities.TLeakageIncidentRepo
}

// AddData implements [entities.TLeakageIncidentUseCase].
func (u LeakageIncidentUseCase) AddData(data entities.TLeakageIncidentReq) (entities.TLeakageIncident, error) {
	return u.incidentRepo.Create(data)
}

// DeleteData implements [entities.TLeakageIncidentUseCase].
func (u LeakageIncidentUseCase) DeleteData(id uint) (entities.TLeakageIncident, error) {
	return u.incidentRepo.DeleteById(id)
}

// FindDataById implements [entities.TLeakageIncidentUseCase].
func (u LeakageIncidentUseCase) FindDataById(id uint) (entities.TLeakageIncident, error) {
	return u.incidentRepo.GetById(id)
}

// FindDatas implements [entities.TLeakageIncidentUseCase].
func (u LeakageIncidentUseCase) FindDatas(page *int, limit *int) ([]entities.TLeakageIncident, error) {
	return u.incidentRepo.Gets(page, limit)
}

// UpdateData implements [entities.TLeakageIncidentUseCase].
func (u LeakageIncidentUseCase) UpdateData(id uint, data entities.TLeakageIncidentReq) (entities.TLeakageIncident, error) {
	return u.incidentRepo.UpdateById(id, data)
}

func NewLeakageIncidentUseCase(incidentRepo entities.TLeakageIncidentRepo) entities.TLeakageIncidentUseCase {
	return LeakageIncidentUseCase{
		incidentRepo: incidentRepo,
	}
}
