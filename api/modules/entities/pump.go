package entities

import "time"

// TPumpAsset ข้อมูลทะเบียนสินทรัพย์เครื่องสูบน้ำ (Water Pump Asset)
type TPumpAsset struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PumpID        string    `gorm:"type:varchar(50);uniqueIndex" json:"pump_id"`
	PumpName      string    `gorm:"type:varchar(255)" json:"pump_name"`
	StationName   string    `gorm:"type:varchar(255)" json:"station_name"`
	DMACode       string    `gorm:"type:varchar(50);index" json:"dma_code"`
	RatedFlowM3H  float64   `json:"rated_flow_m3_h"`
	RatedHeadBar  float64   `json:"rated_head_bar"`
	RatedPowerKW  float64   `json:"rated_power_kw"`
	InstallYear   int       `json:"install_year"`
	Status        string    `gorm:"type:varchar(50)" json:"status"` // ACTIVE, MAINTENANCE, RETIRED
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TPumpTelemetry ข้อมูลเซนเซอร์ตรวจวัดเครื่องสูบน้ำ (SCADA Telemetry)
type TPumpTelemetry struct {
	VibrationRMSMMs   float64 `json:"vibration_rms_mm_s"`
	BearingTempC      float64 `json:"bearing_temp_c"`
	MotorTempC        float64 `json:"motor_temp_c"`
	FlowRateM3H       float64 `json:"flow_rate_m3_h"`
	HeadBar           float64 `json:"head_bar"`
	PowerKW           float64 `json:"power_kw"`
	EfficiencyPercent float64 `json:"efficiency_percent"`
	OperatingHours    int     `json:"operating_hours"`
	CavitationIndex   float64 `json:"cavitation_index"`
}

// TPumpHealthScore ผลการประเมินคะแนนสุขภาพและทำนาย PTTF ตามข้อกำหนด TOR ML Tier
type TPumpHealthScore struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	PumpID              string         `gorm:"type:varchar(50);index" json:"pump_id"`
	PumpName            string         `gorm:"type:varchar(255)" json:"pump_name"`
	StationName         string         `gorm:"type:varchar(255)" json:"station_name"`
	DMACode             string         `gorm:"type:varchar(50);index" json:"dma_code"`
	HealthScore         float64        `json:"health_score"`                                  // 0.0 - 100.0%
	HealthStatus        string         `gorm:"type:varchar(50)" json:"health_status"`         // EXCELLENT, GOOD, WARNING, CRITICAL
	PTTFDays            int            `json:"pttf_days"`                                     // Predicted Time To Failure (วัน)
	PTTFHours           int            `json:"pttf_hours"`                                    // Predicted Time To Failure (ชั่วโมง)
	FailureRisk         string         `gorm:"type:varchar(50)" json:"failure_risk"`          // LOW, MEDIUM, HIGH, CRITICAL
	PrimaryStressFactor string         `gorm:"type:varchar(100)" json:"primary_stress_factor"`// สาเหตุความเสื่อมสภาพ
	Recommendation      string         `gorm:"type:text" json:"recommendation"`               // คำแนะนำซ่อมบำรุงเชิงรุก
	ISO10816Zone        string         `gorm:"type:varchar(50)" json:"iso_10816_zone"`        // Zone A, B, C, D
	Telemetry           TPumpTelemetry `gorm:"embedded" json:"telemetry"`
	EvaluatedAt         time.Time      `json:"evaluated_at"`
}

// TPumpFleetSummary ภาพรวมสุขภาพเครื่องสูบน้ำทั้งระบบ (Fleet Overview) สำหรับ Executive Dashboard
type TPumpFleetSummary struct {
	TotalPumps           int            `json:"total_pumps"`
	AvgHealthScore       float64        `json:"avg_health_score"`
	StatusDistribution   map[string]int `json:"status_distribution"`
	UrgentActionRequired bool           `json:"urgent_action_required"`
	CriticalPumpsCount   int            `json:"critical_pumps_count"`
	WarningPumpsCount    int            `json:"warning_pumps_count"`
	HealthyPumpsCount    int            `json:"healthy_pumps_count"`
	TORStandard          string         `json:"tor_standard"`
}

type TPumpRepo interface {
	GetAllPumps() ([]TPumpAsset, error)
	GetPumpByID(pumpID string) (TPumpAsset, error)
	GetAllPumpHealth() ([]TPumpHealthScore, error)
	GetPumpHealthByID(pumpID string) (TPumpHealthScore, error)
	SavePumpHealthBatch(items []TPumpHealthScore) error
}

type TPumpUseCase interface {
	GetFleetHealth() ([]TPumpHealthScore, TPumpFleetSummary, error)
	GetPumpHealthByID(pumpID string) (TPumpHealthScore, error)
	GetFleetSummary() (TPumpFleetSummary, error)
}
