package database

import (
	"api/modules/entities"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

// GetMockDMAZones คืนค่า Mock Data ของพื้นที่ DMA ครอบคลุม 4 จังหวัด (สงขลา, เพชรบุรี, กาญจนบุรี, กทม)
func GetMockDMAZones() []entities.TDMAZoneReq {
	return []entities.TDMAZoneReq{
		// 1. สงขลา (3 โซน)
		{
			DMACode:           "DMA-SK-01",
			DMAName:           "โซนเทศบาลนครสงขลา",
			BranchName:        "สาขาสงขลา",
			RegionCode:        "เขต 5",
			TargetPressure:    2.2,
			TotalInflow:       3500.0,
			BilledConsumption: 2800.0,
			NRWRatio:          20.0,
			CustomerCount:     12500,
			Status:            "NORMAL",
		},
		{
			DMACode:           "DMA-SK-02",
			DMAName:           "โซนหาดใหญ่-คอหงส์",
			BranchName:        "สาขาหาดใหญ่",
			RegionCode:        "เขต 5",
			TargetPressure:    2.5,
			TotalInflow:       6200.0,
			BilledConsumption: 4700.0,
			NRWRatio:          24.2,
			CustomerCount:     28000,
			Status:            "WARNING",
		},
		{
			DMACode:           "DMA-SK-03",
			DMAName:           "โซนสะเดา-ด่านนอก",
			BranchName:        "สาขาสะเดา",
			RegionCode:        "เขต 5",
			TargetPressure:    2.0,
			TotalInflow:       2100.0,
			BilledConsumption: 1750.0,
			NRWRatio:          16.7,
			CustomerCount:     7600,
			Status:            "NORMAL",
		},

		// 2. เพชรบุรี (3 โซน)
		{
			DMACode:           "DMA-PB-01",
			DMAName:           "โซนเทศบาลเมืองเพชรบุรี",
			BranchName:        "สาขาเพชรบุรี",
			RegionCode:        "เขต 3",
			TargetPressure:    2.5,
			TotalInflow:       2800.0,
			BilledConsumption: 2100.0,
			NRWRatio:          25.0,
			CustomerCount:     8900,
			Status:            "WARNING",
		},
		{
			DMACode:           "DMA-PB-02",
			DMAName:           "โซนชายหาดชะอำ",
			BranchName:        "สาขาชะอำ",
			RegionCode:        "เขต 3",
			TargetPressure:    2.3,
			TotalInflow:       4500.0,
			BilledConsumption: 3300.0,
			NRWRatio:          26.7,
			CustomerCount:     14200,
			Status:            "CRITICAL",
		},
		{
			DMACode:           "DMA-PB-03",
			DMAName:           "โซนท่ายาง-เขื่อนเพชร",
			BranchName:        "สาขาท่ายาง",
			RegionCode:        "เขต 3",
			TargetPressure:    2.0,
			TotalInflow:       1900.0,
			BilledConsumption: 1600.0,
			NRWRatio:          15.8,
			CustomerCount:     6100,
			Status:            "NORMAL",
		},

		// 3. กาญจนบุรี (3 โซน)
		{
			DMACode:           "DMA-KN-01",
			DMAName:           "โซนเทศบาลเมืองกาญจนบุรี",
			BranchName:        "สาขากาญจนบุรี",
			RegionCode:        "เขต 7",
			TargetPressure:    2.8,
			TotalInflow:       4200.0,
			BilledConsumption: 3100.0,
			NRWRatio:          26.2,
			CustomerCount:     15300,
			Status:            "CRITICAL",
		},
		{
			DMACode:           "DMA-KN-02",
			DMAName:           "โซนท่าม่วง-ท่าเรือ",
			BranchName:        "สาขาท่าม่วง",
			RegionCode:        "เขต 7",
			TargetPressure:    2.2,
			TotalInflow:       2600.0,
			BilledConsumption: 2150.0,
			NRWRatio:          17.3,
			CustomerCount:     9400,
			Status:            "NORMAL",
		},
		{
			DMACode:           "DMA-KN-03",
			DMAName:           "โซนไทรโยค",
			BranchName:        "สาขาไทรโยค",
			RegionCode:        "เขต 7",
			TargetPressure:    2.4,
			TotalInflow:       1800.0,
			BilledConsumption: 1400.0,
			NRWRatio:          22.2,
			CustomerCount:     5200,
			Status:            "WARNING",
		},

		// 4. กรุงเทพมหานคร (4 โซน)
		{
			DMACode:           "DMA-BKK-01",
			DMAName:           "โซนดินแดง-ห้วยขวาง",
			BranchName:        "สาขาพญาไท",
			RegionCode:        "เขต 1",
			TargetPressure:    2.0,
			TotalInflow:       8500.0,
			BilledConsumption: 7200.0,
			NRWRatio:          15.3,
			CustomerCount:     32000,
			Status:            "NORMAL",
		},
		{
			DMACode:           "DMA-BKK-02",
			DMAName:           "โซนบางเขน-จตุจักร",
			BranchName:        "สาขาบางเขน",
			RegionCode:        "เขต 1",
			TargetPressure:    2.1,
			TotalInflow:       9200.0,
			BilledConsumption: 7600.0,
			NRWRatio:          17.4,
			CustomerCount:     35400,
			Status:            "NORMAL",
		},
		{
			DMACode:           "DMA-BKK-03",
			DMAName:           "โซนบางกะปิ-รามคำแหง",
			BranchName:        "สาขาหัวหมาก",
			RegionCode:        "เขต 1",
			TargetPressure:    2.2,
			TotalInflow:       7800.0,
			BilledConsumption: 5900.0,
			NRWRatio:          24.3,
			CustomerCount:     29100,
			Status:            "WARNING",
		},
		{
			DMACode:           "DMA-BKK-04",
			DMAName:           "โซนพระโขนง-คลองเตย",
			BranchName:        "สาขาสุขุมวิท",
			RegionCode:        "เขต 1",
			TargetPressure:    1.9,
			TotalInflow:       8100.0,
			BilledConsumption: 6100.0,
			NRWRatio:          24.7,
			CustomerCount:     31200,
			Status:            "CRITICAL",
		},
	}
}

// GetMockPipeAssets คืนค่า Mock Data ท่อประปา
func GetMockPipeAssets() []entities.TPipeAssetReq {
	return []entities.TPipeAssetReq{
		// สงขลา
		{
			AssetID: "AST-SK-001", PipeID: "PIPE-SK-01", DMACode: "DMA-SK-01",
			Material: "HDPE", Diameter: 200.0, Length: 1200.0, InstallationYear: 2018,
			Latitude: 7.1898, Longitude: 100.5954, Status: "ACTIVE",
		},
		{
			AssetID: "AST-SK-002", PipeID: "PIPE-SK-02", DMACode: "DMA-SK-01",
			Material: "PVC", Diameter: 150.0, Length: 800.0, InstallationYear: 2016,
			Latitude: 7.1920, Longitude: 100.5970, Status: "ACTIVE",
		},
		{
			AssetID: "AST-SK-003", PipeID: "PIPE-SK-03", DMACode: "DMA-SK-02",
			Material: "เหล็กเหนียว", Diameter: 300.0, Length: 1800.0, InstallationYear: 2014,
			Latitude: 7.0084, Longitude: 100.4767, Status: "ACTIVE",
		},
		{
			AssetID: "AST-SK-004", PipeID: "PIPE-SK-04", DMACode: "DMA-SK-03",
			Material: "HDPE", Diameter: 150.0, Length: 950.0, InstallationYear: 2021,
			Latitude: 6.6382, Longitude: 100.4211, Status: "ACTIVE",
		},

		// เพชรบุรี
		{
			AssetID: "AST-PB-001", PipeID: "PIPE-PB-01", DMACode: "DMA-PB-01",
			Material: "PVC", Diameter: 150.0, Length: 850.0, InstallationYear: 2015,
			Latitude: 13.1110, Longitude: 99.9392, Status: "ACTIVE",
		},
		{
			AssetID: "AST-PB-002", PipeID: "PIPE-PB-02", DMACode: "DMA-PB-01",
			Material: "HDPE", Diameter: 200.0, Length: 1100.0, InstallationYear: 2019,
			Latitude: 13.1145, Longitude: 99.9450, Status: "ACTIVE",
		},
		{
			AssetID: "AST-PB-003", PipeID: "PIPE-PB-03", DMACode: "DMA-PB-02",
			Material: "AC", Diameter: 200.0, Length: 1400.0, InstallationYear: 2008,
			Latitude: 12.7984, Longitude: 99.9678, Status: "REPLACED",
		},
		{
			AssetID: "AST-PB-004", PipeID: "PIPE-PB-04", DMACode: "DMA-PB-03",
			Material: "PVC", Diameter: 100.0, Length: 600.0, InstallationYear: 2017,
			Latitude: 12.9712, Longitude: 99.9189, Status: "ACTIVE",
		},

		// กาญจนบุรี
		{
			AssetID: "AST-KN-001", PipeID: "PIPE-KN-01", DMACode: "DMA-KN-01",
			Material: "เหล็กเหนียว", Diameter: 300.0, Length: 1500.0, InstallationYear: 2012,
			Latitude: 14.0228, Longitude: 99.5328, Status: "ACTIVE",
		},
		{
			AssetID: "AST-KN-002", PipeID: "PIPE-KN-02", DMACode: "DMA-KN-01",
			Material: "HDPE", Diameter: 250.0, Length: 1250.0, InstallationYear: 2020,
			Latitude: 14.0350, Longitude: 99.5200, Status: "ACTIVE",
		},
		{
			AssetID: "AST-KN-003", PipeID: "PIPE-KN-03", DMACode: "DMA-KN-02",
			Material: "PVC", Diameter: 150.0, Length: 900.0, InstallationYear: 2016,
			Latitude: 13.9620, Longitude: 99.6410, Status: "ACTIVE",
		},
		{
			AssetID: "AST-KN-004", PipeID: "PIPE-KN-04", DMACode: "DMA-KN-03",
			Material: "HDPE", Diameter: 100.0, Length: 750.0, InstallationYear: 2019,
			Latitude: 14.1200, Longitude: 99.1500, Status: "ACTIVE",
		},

		// กทม
		{
			AssetID: "AST-BKK-001", PipeID: "PIPE-BKK-01", DMACode: "DMA-BKK-01",
			Material: "เหล็กเหนียว", Diameter: 400.0, Length: 2200.0, InstallationYear: 2020,
			Latitude: 13.7630, Longitude: 100.5600, Status: "ACTIVE",
		},
		{
			AssetID: "AST-BKK-002", PipeID: "PIPE-BKK-02", DMACode: "DMA-BKK-01",
			Material: "HDPE", Diameter: 300.0, Length: 1400.0, InstallationYear: 2022,
			Latitude: 13.7710, Longitude: 100.5680, Status: "ACTIVE",
		},
		{
			AssetID: "AST-BKK-003", PipeID: "PIPE-BKK-03", DMACode: "DMA-BKK-02",
			Material: "เหล็กเหนียว", Diameter: 500.0, Length: 3100.0, InstallationYear: 2015,
			Latitude: 13.8590, Longitude: 100.5970, Status: "ACTIVE",
		},
		{
			AssetID: "AST-BKK-004", PipeID: "PIPE-BKK-04", DMACode: "DMA-BKK-03",
			Material: "PVC", Diameter: 200.0, Length: 1100.0, InstallationYear: 2018,
			Latitude: 13.7620, Longitude: 100.6300, Status: "ACTIVE",
		},
		{
			AssetID: "AST-BKK-005", PipeID: "PIPE-BKK-05", DMACode: "DMA-BKK-04",
			Material: "เหล็กเหนียว", Diameter: 350.0, Length: 1900.0, InstallationYear: 2011,
			Latitude: 13.7120, Longitude: 100.5850, Status: "ACTIVE",
		},
	}
}

// GetMockLeakageIncidents คืนค่า Mock Data รายการแจ้งซ่อมท่อแตกรั่ว 16 เคส ครบทั้ง 4 จังหวัด
func GetMockLeakageIncidents() []entities.TLeakageIncidentReq {
	tRepaired1 := time.Date(2026, 11, 21, 14, 30, 0, 0, time.UTC)
	tRepaired2 := time.Date(2026, 11, 22, 17, 0, 0, 0, time.UTC)
	tRepaired3 := time.Date(2026, 11, 23, 11, 0, 0, 0, time.UTC)
	tRepaired4 := time.Date(2026, 11, 24, 16, 45, 0, 0, time.UTC)

	return []entities.TLeakageIncidentReq{
		// ---------------- สงขลา ----------------
		{
			TicketID:            "OCS-2026-SK01",
			DMACode:             "DMA-SK-01",
			PipeID:              "PIPE-SK-01",
			IncidentType:        "Main Burst",
			Severity:            "HIGH",
			LocationDescription: "ถนนชลาทัศน์ ใกล้หาดสมิหลา อ.เมือง จ.สงขลา",
			Latitude:            7.1950,
			Longitude:           100.5980,
			WaterLossEstimate:   150.0,
			ReportedAt:          time.Date(2026, 11, 20, 8, 30, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "IN_PROGRESS",
		},
		{
			TicketID:            "OCS-2026-SK02",
			DMACode:             "DMA-SK-01",
			PipeID:              "PIPE-SK-02",
			IncidentType:        "Leak",
			Severity:            "LOW",
			LocationDescription: "ซอยรามคำแหง ถนนไทรบุรี อ.เมือง จ.สงขลา",
			Latitude:            7.1915,
			Longitude:           100.5960,
			WaterLossEstimate:   15.0,
			ReportedAt:          time.Date(2026, 11, 21, 9, 10, 0, 0, time.UTC),
			RepairedAt:          &tRepaired1,
			RepairStatus:        "CLOSED",
		},
		{
			TicketID:            "OCS-2026-SK03",
			DMACode:             "DMA-SK-02",
			PipeID:              "PIPE-SK-03",
			IncidentType:        "Main Burst",
			Severity:            "CRITICAL",
			LocationDescription: "ถนนกาญจนวนิช หน้า ม.อ.หาดใหญ่ อ.หาดใหญ่ จ.สงขลา",
			Latitude:            7.0090,
			Longitude:           100.4990,
			WaterLossEstimate:   280.0,
			ReportedAt:          time.Date(2026, 11, 23, 14, 0, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "PENDING",
		},
		{
			TicketID:            "OCS-2026-SK04",
			DMACode:             "DMA-SK-03",
			PipeID:              "PIPE-SK-04",
			IncidentType:        "Leak",
			Severity:            "MEDIUM",
			LocationDescription: "ถนนกาญจนะรุจี ซอย 3 ด่านนอก อ.สะเดา จ.สงขลา",
			Latitude:            6.6410,
			Longitude:           100.4230,
			WaterLossEstimate:   35.0,
			ReportedAt:          time.Date(2026, 11, 24, 11, 20, 0, 0, time.UTC),
			RepairedAt:          &tRepaired4,
			RepairStatus:        "RESOLVED",
		},

		// ---------------- เพชรบุรี ----------------
		{
			TicketID:            "OCS-2026-PB01",
			DMACode:             "DMA-PB-01",
			PipeID:              "PIPE-PB-01",
			IncidentType:        "Leak",
			Severity:            "MEDIUM",
			LocationDescription: "ถนนเพชรเกษม หน้าตลาดเทศบาลเมืองเพชรบุรี จ.เพชรบุรี",
			Latitude:            13.1125,
			Longitude:           99.9410,
			WaterLossEstimate:   45.0,
			ReportedAt:          time.Date(2026, 11, 21, 10, 15, 0, 0, time.UTC),
			RepairedAt:          &tRepaired1,
			RepairStatus:        "RESOLVED",
		},
		{
			TicketID:            "OCS-2026-PB02",
			DMACode:             "DMA-PB-01",
			PipeID:              "PIPE-PB-02",
			IncidentType:        "Main Burst",
			Severity:            "HIGH",
			LocationDescription: "ถนนราชดำเนิน ใกล้วังบ้านปืน อ.เมือง จ.เพชรบุรี",
			Latitude:            13.1020,
			Longitude:           99.9480,
			WaterLossEstimate:   180.0,
			ReportedAt:          time.Date(2026, 11, 22, 13, 40, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "IN_PROGRESS",
		},
		{
			TicketID:            "OCS-2026-PB03",
			DMACode:             "DMA-PB-02",
			PipeID:              "PIPE-PB-03",
			IncidentType:        "Main Burst",
			Severity:            "CRITICAL",
			LocationDescription: "ถนนเลียบชายหาดชะอำ ใกล้สะพานหิน อ.ชะอำ จ.เพชรบุรี",
			Latitude:            12.8010,
			Longitude:           99.9720,
			WaterLossEstimate:   310.0,
			ReportedAt:          time.Date(2026, 11, 23, 7, 30, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "PENDING",
		},
		{
			TicketID:            "OCS-2026-PB04",
			DMACode:             "DMA-PB-03",
			PipeID:              "PIPE-PB-04",
			IncidentType:        "Leak",
			Severity:            "LOW",
			LocationDescription: "ถนนท่ายาง-ชะอำ หมู่ 4 ต.ท่ายาง อ.ท่ายาง จ.เพชรบุรี",
			Latitude:            12.9690,
			Longitude:           99.9210,
			WaterLossEstimate:   22.0,
			ReportedAt:          time.Date(2026, 11, 24, 15, 10, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "IN_PROGRESS",
		},

		// ---------------- กาญจนบุรี ----------------
		{
			TicketID:            "OCS-2026-KN01",
			DMACode:             "DMA-KN-01",
			PipeID:              "PIPE-KN-01",
			IncidentType:        "Main Burst",
			Severity:            "CRITICAL",
			LocationDescription: "ถนนแสงชูโต ใกล้สะพานข้ามแม่น้ำแคว อ.เมือง จ.กาญจนบุรี",
			Latitude:            14.0410,
			Longitude:           99.5040,
			WaterLossEstimate:   320.0,
			ReportedAt:          time.Date(2026, 11, 22, 6, 0, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "PENDING",
		},
		{
			TicketID:            "OCS-2026-KN02",
			DMACode:             "DMA-KN-01",
			PipeID:              "PIPE-KN-02",
			IncidentType:        "Leak",
			Severity:            "MEDIUM",
			LocationDescription: "ถนนแม่กลอง ปากทางเข้าสุสานสัมพันธมิตร อ.เมือง จ.กาญจนบุรี",
			Latitude:            14.0250,
			Longitude:           99.5290,
			WaterLossEstimate:   50.0,
			ReportedAt:          time.Date(2026, 11, 22, 11, 30, 0, 0, time.UTC),
			RepairedAt:          &tRepaired2,
			RepairStatus:        "CLOSED",
		},
		{
			TicketID:            "OCS-2026-KN03",
			DMACode:             "DMA-KN-02",
			PipeID:              "PIPE-KN-03",
			IncidentType:        "Leak",
			Severity:            "LOW",
			LocationDescription: "ถนนแสงชูโตสายเก่า หน้าตลาดท่าม่วง อ.ท่าม่วง จ.กาญจนบุรี",
			Latitude:            13.9650,
			Longitude:           99.6380,
			WaterLossEstimate:   18.0,
			ReportedAt:          time.Date(2026, 11, 23, 16, 20, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "IN_PROGRESS",
		},
		{
			TicketID:            "OCS-2026-KN04",
			DMACode:             "DMA-KN-03",
			PipeID:              "PIPE-KN-04",
			IncidentType:        "Main Burst",
			Severity:            "HIGH",
			LocationDescription: "ทางหลวงชนบท กจ.4032 ใกล้ทางเข้าน้ำตกไทรโยคน้อย อ.ไทรโยค จ.กาญจนบุรี",
			Latitude:            14.2380,
			Longitude:           99.0600,
			WaterLossEstimate:   210.0,
			ReportedAt:          time.Date(2026, 11, 24, 8, 45, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "PENDING",
		},

		// ---------------- กรุงเทพมหานคร ----------------
		{
			TicketID:            "OCS-2026-BKK01",
			DMACode:             "DMA-BKK-01",
			PipeID:              "PIPE-BKK-01",
			IncidentType:        "Leak",
			Severity:            "LOW",
			LocationDescription: "ถนนประชาสงเคราะห์ เขตดินแดง กทม",
			Latitude:            13.7680,
			Longitude:           100.5645,
			WaterLossEstimate:   20.0,
			ReportedAt:          time.Date(2026, 11, 23, 9, 0, 0, 0, time.UTC),
			RepairedAt:          &tRepaired3,
			RepairStatus:        "CLOSED",
		},
		{
			TicketID:            "OCS-2026-BKK02",
			DMACode:             "DMA-BKK-01",
			PipeID:              "PIPE-BKK-02",
			IncidentType:        "Main Burst",
			Severity:            "HIGH",
			LocationDescription: "ถนนรัชดาภิเษก ใกล้แยกห้วยขวาง เขตห้วยขวาง กทม",
			Latitude:            13.7780,
			Longitude:           100.5730,
			WaterLossEstimate:   220.0,
			ReportedAt:          time.Date(2026, 11, 24, 2, 15, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "IN_PROGRESS",
		},
		{
			TicketID:            "OCS-2026-BKK03",
			DMACode:             "DMA-BKK-02",
			PipeID:              "PIPE-BKK-03",
			IncidentType:        "Main Burst",
			Severity:            "CRITICAL",
			LocationDescription: "ถนนพหลโยธิน หน้าสวนจตุจักร เขตจตุจักร กทม",
			Latitude:            13.8050,
			Longitude:           100.5540,
			WaterLossEstimate:   450.0,
			ReportedAt:          time.Date(2026, 11, 24, 18, 30, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "PENDING",
		},
		{
			TicketID:            "OCS-2026-BKK04",
			DMACode:             "DMA-BKK-04",
			PipeID:              "PIPE-BKK-05",
			IncidentType:        "Leak",
			Severity:            "MEDIUM",
			LocationDescription: "ถนนสุขุมวิท 71 (ปรีดี พนมยงค์) เขตวัฒนา กทม",
			Latitude:            13.7250,
			Longitude:           100.5930,
			WaterLossEstimate:   60.0,
			ReportedAt:          time.Date(2026, 11, 25, 8, 0, 0, 0, time.UTC),
			RepairedAt:          nil,
			RepairStatus:        "IN_PROGRESS",
		},
	}
}

// SeedAll นำข้อมูล Mock ลงฐานข้อมูล (Upsert / Check Exists)
func SeedAll(db *gorm.DB) error {
	log.Println("[Seeder] Starting to seed mock data...")

	// 1. Seed DMA Zones
	for _, zoneReq := range GetMockDMAZones() {
		var count int64
		db.Model(&entities.TDMAZone{}).Where("dma_code = ?", zoneReq.DMACode).Count(&count)
		if count == 0 {
			zone := entities.TDMAZone{
				DMACode:           zoneReq.DMACode,
				DMAName:           zoneReq.DMAName,
				BranchName:        zoneReq.BranchName,
				RegionCode:        zoneReq.RegionCode,
				TargetPressure:    zoneReq.TargetPressure,
				TotalInflow:       zoneReq.TotalInflow,
				BilledConsumption: zoneReq.BilledConsumption,
				NRWRatio:          zoneReq.NRWRatio,
				CustomerCount:     zoneReq.CustomerCount,
				Status:            zoneReq.Status,
			}
			if err := db.Create(&zone).Error; err != nil {
				return fmt.Errorf("failed to seed DMA zone %s: %w", zoneReq.DMACode, err)
			}
			log.Printf("[Seeder] Seeded DMA Zone: %s (%s)", zone.DMACode, zone.DMAName)
		}
	}

	// 2. Seed Pipe Assets
	for _, pipeReq := range GetMockPipeAssets() {
		var count int64
		db.Model(&entities.TPipeAssets{}).Where("asset_id = ?", pipeReq.AssetID).Count(&count)
		if count == 0 {
			pipe := entities.TPipeAssets{
				AssetID:          pipeReq.AssetID,
				PipeID:           pipeReq.PipeID,
				DMACode:          pipeReq.DMACode,
				Material:         pipeReq.Material,
				Diameter:         pipeReq.Diameter,
				Length:           pipeReq.Length,
				InstallationYear: pipeReq.InstallationYear,
				Latitude:         pipeReq.Latitude,
				Longitude:        pipeReq.Longitude,
				Status:           pipeReq.Status,
			}
			if err := db.Create(&pipe).Error; err != nil {
				return fmt.Errorf("failed to seed Pipe Asset %s: %w", pipeReq.AssetID, err)
			}
			log.Printf("[Seeder] Seeded Pipe Asset: %s (%s)", pipe.AssetID, pipe.PipeID)
		}
	}

	// 3. Seed Leakage Incidents
	for _, incidentReq := range GetMockLeakageIncidents() {
		var count int64
		db.Model(&entities.TLeakageIncident{}).Where("ticket_id = ?", incidentReq.TicketID).Count(&count)
		if count == 0 {
			incident := entities.TLeakageIncident{
				TicketID:            incidentReq.TicketID,
				DMACode:             incidentReq.DMACode,
				PipeID:              incidentReq.PipeID,
				IncidentType:        incidentReq.IncidentType,
				Severity:            incidentReq.Severity,
				LocationDescription: incidentReq.LocationDescription,
				Latitude:            incidentReq.Latitude,
				Longitude:           incidentReq.Longitude,
				WaterLossEstimate:   incidentReq.WaterLossEstimate,
				ReportedAt:          incidentReq.ReportedAt,
				RepairedAt:          incidentReq.RepairedAt,
				RepairStatus:        incidentReq.RepairStatus,
			}
			if err := db.Create(&incident).Error; err != nil {
				return fmt.Errorf("failed to seed Leakage Incident %s: %w", incidentReq.TicketID, err)
			}
			log.Printf("[Seeder] Seeded Leakage Incident: %s (%s)", incident.TicketID, incident.LocationDescription)
		}
	}

	// 4. Seed Demand Forecasts (จากผลลัพธ์ ML Model Pipeline)
	if err := SeedDemandForecasts(db); err != nil {
		log.Printf("[Seeder] Warning: could not seed demand forecasts from JSON: %v", err)
	}

	log.Println("[Seeder] Mock data seeding completed successfully!")
	return nil
}

// SeedDemandForecasts โหลดผลการพยากรณ์ 7 วันจาก JSON ของ ML Pipeline ลง MySQL
func SeedDemandForecasts(db *gorm.DB) error {
	possiblePaths := []string{
		"../python-etl/data/demand_forecast_7d.json",
		"python-etl/data/demand_forecast_7d.json",
		"./data/demand_forecast_7d.json",
	}

	var content []byte
	var err error
	for _, p := range possiblePaths {
		content, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}

	type jsonForecast struct {
		DMACode                  string  `json:"dma_code"`
		ForecastTimestamp        string  `json:"forecast_timestamp"`
		StepAhead                int     `json:"step_ahead"`
		PredictedDemand          float64 `json:"predicted_demand"`
		LowerBound80             float64 `json:"lower_bound_80"`
		UpperBound80             float64 `json:"upper_bound_80"`
		CIMargin                 float64 `json:"ci_margin"`
		Hour                     int     `json:"hour"`
		IsPeakPeriod             bool    `json:"is_peak_period"`
		PumpOperationRecommended string  `json:"pump_operation_recommended"`
	}

	var items []jsonForecast
	if err := json.Unmarshal(content, &items); err != nil {
		return err
	}

	var existingDMAs []string
	db.Model(&entities.TDemandForecast{}).Distinct().Pluck("dma_code", &existingDMAs)
	existingMap := make(map[string]bool)
	for _, d := range existingDMAs {
		existingMap[d] = true
	}

	var batch []entities.TDemandForecast
	for _, it := range items {
		if !existingMap[it.DMACode] {
			ts, _ := time.Parse("2006-01-02 15:04:05", it.ForecastTimestamp)
			batch = append(batch, entities.TDemandForecast{
				DMACode:                  it.DMACode,
				ForecastTimestamp:        ts,
				StepAhead:                it.StepAhead,
				PredictedDemand:          it.PredictedDemand,
				LowerBound80:             it.LowerBound80,
				UpperBound80:             it.UpperBound80,
				CIMargin:                 it.CIMargin,
				Hour:                     it.Hour,
				IsPeakPeriod:             it.IsPeakPeriod,
				PumpOperationRecommended: it.PumpOperationRecommended,
				CreatedAt:                time.Now(),
			})
		}
	}

	if len(batch) > 0 {
		if err := db.CreateInBatches(batch, 100).Error; err != nil {
			return err
		}
		log.Printf("[Seeder] Successfully seeded %d Demand Forecast records into database!", len(batch))
	}
	return nil
}
