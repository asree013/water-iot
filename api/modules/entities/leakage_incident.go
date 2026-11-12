package entities

import "time"

// TLeakageIncident ประวัติรับแจ้งซ่อมท่อแตกรั่วจากระบบ OCS
type TLeakageIncident struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TicketID            string     `gorm:"type:varchar(50);uniqueIndex" json:"ticket_id"` // รหัสใบแจ้งซ่อม เช่น OCS-2026-001
	DMACode             string     `gorm:"type:varchar(50);index" json:"dma_code"`        // โซนที่เกิดเหตุ
	PipeID              string     `gorm:"type:varchar(50)" json:"pipe_id"`               // รหัสท่อที่เกิดเหตุ
	IncidentType        string     `gorm:"type:varchar(50)" json:"incident_type"`         // ประเภท: Main Burst (ท่อประธานแตก), Leak (ท่อรั่วซึม)
	Severity            string     `gorm:"type:varchar(20)" json:"severity"`              // ความรุนแรง: LOW, MEDIUM, HIGH, CRITICAL
	LocationDescription string     `gorm:"type:text" json:"location_description"`         // สถานที่เกิดเหตุ / จุดสังเกต
	Latitude            float64    `json:"latitude"`                                      // พิกัดที่เกิดเหตุ
	Longitude           float64    `json:"longitude"`
	WaterLossEstimate   float64    `json:"water_loss_estimate"`                   // ปริมาณน้ำสูญเสียประเมิน (ลบ.ม.)
	ReportedAt          time.Time  `json:"reported_at"`                           // วันเวลาที่รับแจ้ง
	RepairedAt          *time.Time `json:"repaired_at"`                           // วันเวลาที่ซ่อมเสร็จ (เป็น pointer เผื่อยังซ่อมไม่เสร็จ)
	RepairStatus        string     `gorm:"type:varchar(30)" json:"repair_status"` // PENDING, IN_PROGRESS, RESOLVED, CLOSED
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
