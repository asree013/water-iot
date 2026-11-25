---
name: leakage-incident-management
description: จัดการเหตุการณ์ท่อรั่วและท่อแตก (Leakage & Burst Incidents) ติดตามสถานะงานซ่อม และประเมินปริมาณน้ำสูญเสีย
---

# Leakage Incident Management Module (`leakage_incident`)

โมดูลจัดการข้อมูลเหตุการณ์ท่อรั่วและท่อแตก (Leakage Incident) เพื่อรองรับกระบวนการแจ้งเตือน บันทึกพิกัดจุดเกิดเหตุ ติดตามสถานะการซ่อมบำรุง และประเมินปริมาณน้ำสูญเสียสะสม ตามข้อกำหนด TOR

---

## 1. วัตถุประสงค์และบทบาทในระบบ
- บันทึกประวัติและ Ticket ของเหตุการณ์ท่อแตก (Burst) หรือท่อรั่วซึม (Leak)
- บันทึกพิกัดภูมิศาสตร์ (Latitude, Longitude) และคำอธิบายสถานที่เกิดเหตุ เพื่อแสดงผลบนแผนที่ Interactive Map (Tier 5)
- ติดตามวงจรสถานะงานซ่อม (`REPORTED`, `IN_PROGRESS`, `REPAIRED`, `RESOLVED`, `CLOSED`)
- คำนวณปริมาณน้ำสูญเสียโดยประมาณ (Water Loss Estimate: $m^3$) และระยะเวลาที่ใช้ในการซ่อมบำรุง (MTTR - Mean Time To Repair)

---

## 2. โครงสร้างสถาปัตยกรรม (Clean Architecture)

```
api/modules/leakage_incident/
├── controller/
│   ├── leakage_incident_controller.go   # Fiber HTTP Controller (CRUD, Swagger Annotations)
│   └── leakage_incident_test.go         # Unit Test Suite
├── entities/                            # นิยาม Entity ใน api/modules/entities/leakage_incident.go
├── repository/
│   └── leakage_incident_repository.go   # GORM MySQL Data Access Layer
└── usecase/
    └── leakage_incident_usecase.go      # Business Logic Layer
```

---

## 3. โครงสร้างฐานข้อมูล (Database Schema: `tleakage_incidents`)

| Column | Type | Nullable | Description |
| :--- | :--- | :--- | :--- |
| `id` | `BIGINT UNSIGNED` (PK, Auto Inc) | NO | รหัสอ้างอิงลำดับ |
| `ticket_id` | `VARCHAR(50)` (Unique Index) | NO | รหัส Ticket แจ้งเหตุ เช่น `OCS-2026-SK01` |
| `dma_code` | `VARCHAR(50)` (Index) | NO | รหัสโซน DMA ที่เกิดเหตุ |
| `pipe_id` | `VARCHAR(50)` (Index) | YES | รหัสเส้นท่อที่เกิดเหตุ |
| `incident_type` | `VARCHAR(50)` | NO | ประเภทเหตุการณ์ (`Main Burst`, `Sub-main Leak`, `Joint Leak`, `Valve Leak`) |
| `severity` | `VARCHAR(50)` | NO | ระดับความรุนแรง (`CRITICAL`, `HIGH`, `MEDIUM`, `LOW`) |
| `location_description` | `TEXT` | YES | รายละเอียดสถานที่เกิดเหตุ |
| `latitude` | `DOUBLE` | YES | พิกัดละติจูดจุดเกิดเหตุ |
| `longitude` | `DOUBLE` | YES | พิกัดลองจิจูดจุดเกิดเหตุ |
| `water_loss_estimate` | `DOUBLE` | YES | ปริมาณน้ำสูญเสียโดยประมาณ ($m^3$) |
| `reported_at` | `DATETIME(3)` | YES | วันเวลาที่ได้รับแจ้งเหตุ |
| `repaired_at` | `DATETIME(3)` | YES | วันเวลาที่ดำเนินการซ่อมแล้วเสร็จ |
| `repair_status` | `VARCHAR(50)` | YES | สถานะงานซ่อม (`REPORTED`, `IN_PROGRESS`, `REPAIRED`, `RESOLVED`, `CLOSED`) |
| `created_at` | `DATETIME(3)` | YES | เวลาที่สร้างข้อมูล |
| `updated_at` | `DATETIME(3)` | YES | เวลาที่แก้ไขข้อมูลล่าสุด |

---

## 4. รายการ API Endpoints (RESTful API)

- **`GET /api/v1/leakage-incident`**: ดึงรายการเหตุการณ์ทั้งหมด (รองรับ Pagination: `?page=1&limit=10`)
- **`GET /api/v1/leakage-incident/:id`**: ดึงข้อมูลเหตุการณ์ตาม ID
- **`POST /api/v1/leakage-incident`**: สร้างรายงานเหตุการณ์ท่อแตก/รั่วใหม่
- **`PUT /api/v1/leakage-incident/:id`**: ปรับปรุงสถานะหรือรายละเอียดงานซ่อม
- **`DELETE /api/v1/leakage-incident/:id`**: ลบข้อมูลเหตุการณ์

### ตัวอย่าง Response (`GET /api/v1/leakage-incident/1`):
```json
{
  "status": 200,
  "message": "Fetch Leakage Incident Successfully",
  "data": {
    "id": 1,
    "ticket_id": "OCS-2026-SK01",
    "dma_code": "DMA-SK-01",
    "pipe_id": "PIPE-SK-01",
    "incident_type": "Main Burst",
    "severity": "HIGH",
    "location_description": "ถนนชลาทัศน์ ใกล้หาดสมิหลา อ.เมือง จ.สงขลา",
    "latitude": 7.195,
    "longitude": 100.598,
    "water_loss_estimate": 150,
    "reported_at": "2026-11-20T15:30:00+07:00",
    "repaired_at": null,
    "repair_status": "IN_PROGRESS",
    "created_at": "2026-11-25T13:49:04+07:00",
    "updated_at": "2026-11-25T13:49:04+07:00"
  }
}
```

---

## 5. การทดสอบ (Verification)
- **Unit Test File:** `api/modules/leakage_incident/controller/leakage_incident_test.go`
- **สถานะ:** ทดสอบผ่าน 100% พร้อมข้อมูล Mock Incident ครอบคลุมทั้ง 4 จังหวัด
