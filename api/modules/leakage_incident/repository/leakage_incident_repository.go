package repository

import (
	"api/modules/entities"

	"gorm.io/gorm"
)

type LeakageIncidentRepository struct {
	db *gorm.DB
}

// Create implements [entities.TLeakageIncidentRepo].
func (r LeakageIncidentRepository) Create(data entities.TLeakageIncidentReq) (entities.TLeakageIncident, error) {
	incident := entities.TLeakageIncident{
		TicketID:            data.TicketID,
		DMACode:             data.DMACode,
		PipeID:              data.PipeID,
		IncidentType:        data.IncidentType,
		Severity:            data.Severity,
		LocationDescription: data.LocationDescription,
		Latitude:            data.Latitude,
		Longitude:           data.Longitude,
		WaterLossEstimate:   data.WaterLossEstimate,
		ReportedAt:          data.ReportedAt,
		RepairedAt:          data.RepairedAt,
		RepairStatus:        data.RepairStatus,
	}
	if err := r.db.Create(&incident).Error; err != nil {
		return entities.TLeakageIncident{}, err
	}
	return incident, nil
}

// DeleteById implements [entities.TLeakageIncidentRepo].
func (r LeakageIncidentRepository) DeleteById(id uint) (entities.TLeakageIncident, error) {
	data, err := r.GetById(id)
	if err != nil {
		return entities.TLeakageIncident{}, err
	}
	if err := r.db.Delete(&data).Error; err != nil {
		return entities.TLeakageIncident{}, err
	}
	return data, nil
}

// GetById implements [entities.TLeakageIncidentRepo].
func (r LeakageIncidentRepository) GetById(id uint) (entities.TLeakageIncident, error) {
	var data entities.TLeakageIncident
	if err := r.db.First(&data, id).Error; err != nil {
		return entities.TLeakageIncident{}, err
	}
	return data, nil
}

// Gets implements [entities.TLeakageIncidentRepo].
func (r LeakageIncidentRepository) Gets(page *int, limit *int) ([]entities.TLeakageIncident, error) {
	var data []entities.TLeakageIncident
	query := r.db.Model(&entities.TLeakageIncident{})
	if page != nil && limit != nil && *page > 0 && *limit > 0 {
		query = query.Offset((*page - 1) * *limit).Limit(*limit)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateById implements [entities.TLeakageIncidentRepo].
func (r LeakageIncidentRepository) UpdateById(id uint, data entities.TLeakageIncidentReq) (entities.TLeakageIncident, error) {
	incident, err := r.GetById(id)
	if err != nil {
		return entities.TLeakageIncident{}, err
	}

	incident.TicketID = data.TicketID
	incident.DMACode = data.DMACode
	incident.PipeID = data.PipeID
	incident.IncidentType = data.IncidentType
	incident.Severity = data.Severity
	incident.LocationDescription = data.LocationDescription
	incident.Latitude = data.Latitude
	incident.Longitude = data.Longitude
	incident.WaterLossEstimate = data.WaterLossEstimate
	incident.ReportedAt = data.ReportedAt
	incident.RepairedAt = data.RepairedAt
	incident.RepairStatus = data.RepairStatus

	if err := r.db.Save(&incident).Error; err != nil {
		return entities.TLeakageIncident{}, err
	}
	return incident, nil
}

func NewLeakageIncidentRepository(db *gorm.DB) entities.TLeakageIncidentRepo {
	return LeakageIncidentRepository{db: db}
}
