package entities

import "time"

// TPipeAsset ข้อมูลสินทรัพย์ท่อประปาจาก SAP + SWMAP/GIS
type TPipeAsset struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AssetID          string    `gorm:"type:varchar(50);uniqueIndex" json:"asset_id"` // รหัสสินทรัพย์จาก SAP เช่น AST-001 (ใช้ตรวจ DLQ หน้า 33)
	PipeID           string    `gorm:"type:varchar(50);index" json:"pipe_id"`        // รหัสท่อ เช่น Pipe_01
	DMACode          string    `gorm:"type:varchar(50);index" json:"dma_code"`       // อยู่ในโซน DMA ใด
	Material         string    `gorm:"type:varchar(50)" json:"material"`             // วัสดุท่อ: HDPE, PVC, เหล็กเหนียว, AC
	Diameter         float64   `json:"diameter"`                                     // ขนาดท่อ (มม.) เช่น 100, 150, 300
	Length           float64   `json:"length"`                                       // ความยาวท่อ (เมตร)
	InstallationYear int       `json:"installation_year"`                            // ปีที่ติดตั้ง (ใช้คำนวณอายุท่อเพื่อทำนายความเสี่ยง)
	Latitude         float64   `json:"latitude"`                                     // พิกัด GIS ละติจูด
	Longitude        float64   `json:"longitude"`                                    // พิกัด GIS ลองจิจูด
	Status           string    `gorm:"type:varchar(20)" json:"status"`               // สถานะสินทรัพย์: ACTIVE, INACTIVE, REPLACED
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
