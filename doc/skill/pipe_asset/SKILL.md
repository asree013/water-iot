---
name: pipe-asset-management
description: จัดการข้อมูลสินทรัพย์เส้นท่อประปา (Pipe Assets) และคุณลักษณะทางกายภาพตามสถาปัตยกรรม Clean Architecture
---

# Pipe Asset Management Module (`pipe_asset`)

โมดูลจัดการข้อมูลสินทรัพย์เส้นท่อประปา (Pipe Assets) สำหรับระบบ Water IoT เพื่อบันทึก ติดตามสภาพ และเชื่อมโยงพิกัดเส้นท่อ (GIS LineString) เข้ากับโซนจ่ายน้ำ DMA ตามข้อกำหนด TOR

---

## 1. วัตถุประสงค์และบทบาทในระบบ
- จัดการทะเบียนสินทรัพย์เส้นท่อประปา เช่น ชนิดท่อ (HDPE, Steel, AC, PVC), ขนาดเส้นผ่านศูนย์กลาง (Diameter mm), ความยาว (Length m), ปีที่ติดตั้ง (Install Year)
- ติดตามสภาพความสมบูรณ์ของท่อ (Condition Rating 1-5 หรือ GOOD/FAIR/POOR)
- รองรับการคำนวณอายุการใช้งานและเป็นฐานข้อมูลสำหรับโมเดลการประเมินความเสี่ยงท่อแตก (Pipe Burst Risk Analysis)

---

## 2. โครงสร้างสถาปัตยกรรม (Clean Architecture)

```
api/modules/pipe_asset/
├── controller/
│   ├── pipe_asset_controller.go   # Fiber HTTP Controller (CRUD, Swagger Annotations)
│   └── pipe_asset_test.go         # Unit Test Suite ครอบคลุมเส้นท่อ 4 จังหวัด
├── entities/                      # นิยาม Entity ใน api/modules/entities/pipe_asset.go
├── repository/
│   └── pipe_asset_repository.go   # GORM MySQL Data Access Layer
└── usecase/
    └── pipe_asset_usecase.go      # Business Logic Layer
```

---

## 3. โครงสร้างฐานข้อมูล (Database Schema: `tpipe_assets`)

| Column | Type | Nullable | Description |
| :--- | :--- | :--- | :--- |
| `id` | `BIGINT UNSIGNED` (PK, Auto Inc) | NO | รหัสอ้างอิงลำดับ |
| `pipe_id` | `VARCHAR(50)` (Unique Index) | NO | รหัสระบุเส้นท่อ เช่น `PIPE-SK-01`, `PIPE-BKK-01` |
| `dma_code` | `VARCHAR(50)` (Index) | NO | รหัสโซน DMA ที่เส้นท่อพาดผ่าน |
| `pipe_material` | `VARCHAR(50)` | YES | ชนิดวัสดุท่อ (HDPE, Steel, AC, PVC, DI) |
| `diameter_mm` | `DOUBLE` | YES | ขนาดเส้นผ่านศูนย์กลางท่อ (มม.) |
| `length_meter` | `DOUBLE` | YES | ความยาวของช่วงท่อ (เมตร) |
| `install_year` | `INT` | YES | ปีที่เริ่มติดตั้งท่อ (ค.ศ.) |
| `pressure_rating_bar` | `DOUBLE` | YES | พิกัดแรงดันทนได้ของท่อ (bar) |
| `condition_rating` | `VARCHAR(50)` | YES | ระดับสภาพท่อ (`GOOD`, `FAIR`, `POOR`, `CRITICAL`) |
| `status` | `VARCHAR(50)` | YES | สถานะการใช้งาน (`ACTIVE`, `MAINTENANCE`, `RETIRED`) |
| `created_at` | `DATETIME(3)` | YES | เวลาที่สร้างข้อมูล |
| `updated_at` | `DATETIME(3)` | YES | เวลาที่แก้ไขข้อมูลล่าสุด |

---

## 4. รายการ API Endpoints (RESTful API)

- **`GET /api/v1/pipe-asset`**: ดึงรายการท่อประปาทั้งหมด (รองรับ Pagination: `?page=1&limit=10`)
- **`GET /api/v1/pipe-asset/:id`**: ดึงข้อมูลเส้นท่อตาม ID
- **`POST /api/v1/pipe-asset`**: เพิ่มข้อมูลเส้นท่อใหม่
- **`PUT /api/v1/pipe-asset/:id`**: ปรับปรุงข้อมูลเส้นท่อ
- **`DELETE /api/v1/pipe-asset/:id`**: ลบข้อมูลเส้นท่อ

### ตัวอย่าง Response (`GET /api/v1/pipe-asset/1`):
```json
{
  "status": 200,
  "message": "Fetch Pipe Asset Successfully",
  "data": {
    "id": 1,
    "pipe_id": "PIPE-SK-01",
    "dma_code": "DMA-SK-01",
    "pipe_material": "HDPE",
    "diameter_mm": 300,
    "length_meter": 1250,
    "install_year": 2018,
    "pressure_rating_bar": 10,
    "condition_rating": "GOOD",
    "status": "ACTIVE",
    "created_at": "2026-11-25T13:49:04+07:00",
    "updated_at": "2026-11-25T13:49:04+07:00"
  }
}
```

---

## 5. การทดสอบ (Verification)
- **Unit Test File:** `api/modules/pipe_asset/controller/pipe_asset_test.go`
- **สถานะ:** ทดสอบผ่าน 100% พร้อมชุดข้อมูลท่อประปาครอบคลุม 4 จังหวัดนำร่อง
