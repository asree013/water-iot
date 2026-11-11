package entities

import (
	baserepository "api/modules/base_interface"
	"time"
)

// domain entity

type TDMAZone struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DMACode           string    `gorm:"type:varchar(50);uniqueIndex" json:"dma_code"` // รหัส DMA เช่น DMA-01
	DMAName           string    `gorm:"type:varchar(100)" json:"dma_name"`            // ชื่อโซน เช่น โซนเทศบาล 1
	BranchName        string    `gorm:"type:varchar(100)" json:"branch_name"`         // สาขา (ตาม TOR: นำร่อง 4 สาขา)
	RegionCode        string    `gorm:"type:varchar(20)" json:"region_code"`          // เขต (ตาม TOR: 10 เขต)
	TargetPressure    float64   `json:"target_pressure"`                              // แรงดันเป้าหมาย (bar) เช่น 2.5 (แนะนำ float64 เพื่อใช้เทียบกับ IoT)
	TotalInflow       float64   `json:"total_inflow"`                                 // ปริมาณน้ำจ่ายเข้าโซน (ลบ.ม./วัน)
	BilledConsumption float64   `json:"billed_consumption"`                           // ปริมาณน้ำที่ขายได้ (ลบ.ม./วัน)
	NRWRatio          float64   `json:"nrw_ratio"`                                    // อัตราน้ำสูญเสีย % (Non-Revenue Water)
	CustomerCount     int       `json:"customer_count"`                               // จำนวนผู้ใช้น้ำในโซน
	Status            string    `gorm:"type:varchar(20)" json:"status"`               // NORMAL, WARNING, CRITICAL (ตามเกณฑ์แจ้งเตือน)
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// DTO
type TDMAZoneReq struct {
	DMACode           string  `gorm:"type:varchar(50);uniqueIndex" json:"dma_code"`
	DMAName           string  `gorm:"type:varchar(100)" json:"dma_name"`
	BranchName        string  `gorm:"type:varchar(100)" json:"branch_name"`
	RegionCode        string  `gorm:"type:varchar(20)" json:"region_code"`
	TargetPressure    float64 `json:"target_pressure"`
	TotalInflow       float64 `json:"total_inflow"`
	BilledConsumption float64 `json:"billed_consumption"`
	NRWRatio          float64 `json:"nrw_ratio"`
	CustomerCount     int     `json:"customer_count"`
	Status            string  `gorm:"type:varchar(20)" json:"status"`
}

type TDMAZoneRes struct {
	TDMAZoneReq
}

// repository interface
type TDMAZoneRepo interface {
	baserepository.IBaseRepository[TDMAZone]
}

// useCase Interface
type TDMAZoneUseCase interface {
	baserepository.IBaseUseCase[TDMAZone]
}
