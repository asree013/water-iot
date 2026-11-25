---
name: dma-zone-management
description: จัดการข้อมูลพื้นที่จ่ายน้ำย่อย District Metered Area (DMA Zone) ตามสถาปัตยกรรม Clean Architecture
---

# DMA Zone Management Module (`dma_zone`)

โมดูลจัดการข้อมูลพื้นที่จ่ายน้ำย่อย (District Metered Area: DMA) สำหรับระบบบริหารจัดการน้ำประปาอัจฉริยะ (Water IoT) ตามข้อกำหนด TOR สอดคล้องกับ Tier 2 (API & GIS Tier)

---

## 1. วัตถุประสงค์และบทบาทในระบบ
- จัดการข้อมูลเขตพื้นที่จ่ายน้ำย่อย (DMA) ข้อมูลความดันเป้าหมาย ปริมาณน้ำจ่ายเข้า (Total Inflow) ปริมาณน้ำจำหน่าย (Billed Consumption)
- คำนวณอัตราน้ำสูญเสีย (Non-Revenue Water: NRW %)
- รองรับการเชื่อมต่อกับระบบแผนที่ GIS (Polygon พื้นที่ DMA) และเป็น Key อ้างอิงหลักสำหรับโมเดล ML ใน Tier 3 (Demand Forecasting & MNF Leak Detection)

---

## 2. โครงสร้างสถาปัตยกรรม (Clean Architecture)

```
api/modules/dma_zone/
├── controller/
│   ├── dma_zone_controller.go     # Fiber HTTP Controller (CRUD, Swagger Annotations)
│   └── dma_zone_test.go           # Unit Test Suite ครอบคลุม Mock Data 4 จังหวัด
├── entities/                      # นิยาม Entity ใน api/modules/entities/dma_zone.go
├── repository/
│   └── dma_zone_repository.go     # GORM MySQL Data Access Layer
└── usecase/
    └── dma_zone_usecase.go        # Business Logic Layer (IBaseUseCase Implementation)
```

---

## 3. โครงสร้างฐานข้อมูล (Database Schema: `tdma_zones`)

| Column | Type | Nullable | Description |
| :--- | :--- | :--- | :--- |
| `id` | `BIGINT UNSIGNED` (PK, Auto Inc) | NO | รหัสอ้างอิงลำดับ |
| `dma_code` | `VARCHAR(50)` (Unique Index) | NO | รหัสโค้ด DMA เช่น `DMA-SK-01`, `DMA-BKK-01` |
| `dma_name` | `VARCHAR(255)` | NO | ชื่อพื้นที่เขตจ่ายน้ำ |
| `branch_name` | `VARCHAR(100)` | YES | สาขาการประปาที่รับผิดชอบ |
| `region_code` | `VARCHAR(50)` | YES | เขตพื้นที่ (เช่น เขต 5, เขต 3) |
| `target_pressure` | `DOUBLE` | YES | ความดันเป้าหมายในโซน (bar) |
| `total_inflow` | `DOUBLE` | YES | ปริมาณน้ำจ่ายเข้าโซน ($m^3/day$) |
| `billed_consumption` | `DOUBLE` | YES | ปริมาณน้ำที่จำหน่ายได้ ($m^3/day$) |
| `nrw_ratio` | `DOUBLE` | YES | อัตราน้ำสูญเสีย (% NRW) |
| `customer_count` | `INT` | YES | จำนวนผู้ใช้น้ำในโซน |
| `status` | `VARCHAR(50)` | YES | สถานะโซน (`NORMAL`, `WARNING`, `CRITICAL`) |
| `created_at` | `DATETIME(3)` | YES | เวลาที่สร้างข้อมูล |
| `updated_at` | `DATETIME(3)` | YES | เวลาที่แก้ไขข้อมูลล่าสุด |

---

## 4. รายการ API Endpoints (RESTful API)

- **`GET /api/v1/dma-zone`**: ดึงรายการพื้นที่ DMA ทั้งหมด (รองรับ Pagination: `?page=1&limit=10`)
- **`GET /api/v1/dma-zone/:id`**: ดึงข้อมูล DMA ตาม ID
- **`POST /api/v1/dma-zone`**: สร้างพื้นที่ DMA ใหม่
- **`PUT /api/v1/dma-zone/:id`**: ปรับปรุงข้อมูลพื้นที่ DMA
- **`DELETE /api/v1/dma-zone/:id`**: ลบข้อมูลพื้นที่ DMA

### ตัวอย่าง Response (`GET /api/v1/dma-zone/1`):
```json
{
  "status": 200,
  "message": "Fetch DMA Zone Successfully",
  "data": {
    "id": 1,
    "dma_code": "DMA-SK-01",
    "dma_name": "โซนเทศบาลนครสงขลา",
    "branch_name": "สาขาสงขลา",
    "region_code": "เขต 5",
    "target_pressure": 2.2,
    "total_inflow": 3500,
    "billed_consumption": 2800,
    "nrw_ratio": 20,
    "customer_count": 12500,
    "status": "NORMAL",
    "created_at": "2026-11-25T13:49:04+07:00",
    "updated_at": "2026-11-25T13:49:04+07:00"
  }
}
```

---

## 5. การทดสอบ (Verification)
- **Unit Test File:** `api/modules/dma_zone/controller/dma_zone_test.go`
- **สถานะ:** ทดสอบผ่าน 100% พร้อมชุดข้อมูลจำลองครอบคลุม 4 จังหวัดนำร่อง (สงขลา, เพชรบุรี, กาญจนบุรี, กทม.)
