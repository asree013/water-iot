---
name: water-demand-forecasting
description: พยากรณ์ความต้องการใช้น้ำล่วงหน้า 24 ชม. และ 7 วัน พร้อมช่วงความเชื่อมั่น 80% CI และคำแนะนำการสูบน้ำตามเกณฑ์ TOR
---

# Demand Forecasting Module (`forecast` & `python-etl`)

โมดูลปัญญาประดิษฐ์พยากรณ์ความต้องการใช้น้ำล่วงหน้ารายชั่วโมง 24 ชั่วโมง (24-Hour Ahead) และรายวัน 7 วัน (7-Day Ahead) พร้อมช่วงความเชื่อมั่น 80% (80% Confidence Interval) และคำแนะนำการสั่งเดินเครื่องสูบน้ำ (Peak/Off-Peak Pump Scheduling Optimization) ตามข้อกำหนด TOR ภาคผนวก ๔

---

## 1. วัตถุประสงค์และบทบาทในระบบ
- **พยากรณ์ระยะสั้น 24 ชั่วโมง (24h Ahead):** สำหรับวางแผนรอบการเดินเครื่องสูบน้ำแบบ Peak / Off-Peak เพื่อลดค่าไฟฟ้าตามอัตรา Time of Use (TOU) และควบคุมแรงดันในท่อไม่ให้สูงเกินไปในช่วงกลางคืน
- **พยากรณ์ระยะกลาง 7 วัน (7d Ahead):** สำหรับวางแผนจัดสรรกำลังการผลิตน้ำประปา และการบริหารจัดการน้ำดิบของสถานีผลิตน้ำ
- **ช่วงความเชื่อมั่น 80% Confidence Interval:** คำนวณขอบเขตล่าง (`lower_bound_80`) และขอบเขตบน (`upper_bound_80`) ตามหลักการสถิติ ($Z = 1.28155$) เพื่อให้วิศวกรประปาทราบความผันผวน
- **เกณฑ์มาตรฐาน TOR:** ความแม่นยำของโมเดลต้อง $\ge 85\%$ (ระบบทำได้จริงเฉลี่ย $\approx 96.60\%$, $R^2 \approx 0.989$)

---

## 2. สถาปัตยกรรมและการเชื่อมต่อ (Data Flow)

```
[ Historical Inflow & Weather (180 Days) ]
                  │
                  ▼
   [ python-etl/data_generator.py ] ──> data/dma_hourly_timeseries.parquet (~30,240 records)
                  │
                  ▼
  [ python-etl/demand_forecaster.py ] ──> Feature Engineering (Lags, Cyclical, Rolling)
                  │
                  ├──> models/demand_model_{DMA}.joblib (GradientBoostingRegressor)
                  ├──> data/model_accuracy_summary.json (Per-Zone & Overall Metrics)
                  └──> data/demand_forecast_7d.json (1,176 records, 7 Zones × 168h)
                  │
                  ▼ (Seeder / REST API)
    [ MySQL: t_demand_forecasts ]
                  │
                  ▼
     [ api/modules/forecast/ ] ──> RESTful Endpoints & Dynamic Fallback
                  │
                  ▼
[ Presentation Tier / SCADA 2D Canvas ]
```

---

## 3. โครงสร้างฐานข้อมูล (Database Schema: `t_demand_forecasts`)

| Column | Type | Nullable | Description |
| :--- | :--- | :--- | :--- |
| `id` | `BIGINT UNSIGNED` (PK, Auto Inc) | NO | รหัสอ้างอิงลำดับ |
| `dma_code` | `VARCHAR(50)` (Index) | NO | รหัสโซน DMA |
| `forecast_timestamp` | `DATETIME(3)` (Index) | NO | เวลาที่คาดการณ์ล่วงหน้า |
| `step_ahead` | `INT` | NO | จำนวนชั่วโมงล่วงหน้า (1–168) |
| `predicted_demand` | `DOUBLE` | NO | ปริมาณน้ำที่พยากรณ์ ($m^3/h$) |
| `lower_bound80` | `DOUBLE` | NO | ขอบเขตล่างช่วงความเชื่อมั่น 80% ($m^3/h$) |
| `upper_bound80` | `DOUBLE` | NO | ขอบเขตบนช่วงความเชื่อมั่น 80% ($m^3/h$) |
| `ci_margin` | `DOUBLE` | NO | ค่าเผื่อความคลาดเคลื่อน ($\pm m^3/h$) |
| `hour` | `INT` | NO | ชั่วโมงของวัน (0–23) |
| `is_peak_period` | `TINYINT(1)` | NO | เป็นช่วงเวลาใช้น้ำสูงสุด (Peak) หรือไม่ |
| `pump_operation_recommended` | `VARCHAR(50)` | NO | คำแนะนำการเดินปั๊ม: `INCREASE`, `STANDARD`, `REDUCE` |
| `created_at` | `DATETIME(3)` | YES | เวลาที่สร้างข้อมูล |

---

## 4. รายการ API Endpoints (RESTful API)

### 4.1 ดึงผลการพยากรณ์ความต้องการใช้น้ำ: `GET /api/v1/forecast/demand/:dma_code`
- **Path Param:** `dma_code` (เช่น `DMA-SK-01`, `DMA-03`)
- **Query Param:** `range` (`24h` หรือ `7d`, ค่าเริ่มต้น `24h`)
- **คุณลักษณะพิเศษ (Dynamic Fallback):** หากเป็นโซนใหม่ที่ยังไม่มีข้อมูลใน `t_demand_forecasts` แต่มีอยู่ใน `tdma_zones` ตัว Usecase จะคำนวณและบันทึกผลพยากรณ์ 168 ชั่วโมงลง Database ให้อัตโนมัติทันที

#### ตัวอย่าง Response (`GET /api/v1/forecast/demand/DMA-03?range=24h`):
```json
{
  "status": 200,
  "message": "Fetch Demand Forecast Successfully",
  "dma_code": "DMA-03",
  "range": "24h",
  "summary": {
    "total_estimated_volume_m3": 5842.15,
    "peak_demand_m3_h": 382.26,
    "peak_hour": 18,
    "confidence_interval": "80%",
    "tor_compliant": true
  },
  "data": [
    {
      "step_ahead": 1,
      "forecast_timestamp": "2026-11-25T14:00:00Z",
      "hour": 7,
      "predicted_demand": 364.5,
      "lower_bound_80": 349.07,
      "upper_bound_80": 379.93,
      "ci_margin": 15.43,
      "is_peak_period": true,
      "pump_operation_recommended": "INCREASE"
    }
  ]
}
```

### 4.2 ดึงผลการประเมินความแม่นยำของโมเดล: `GET /api/v1/forecast/metrics`
- **Query Param:** `dma_code` (Optional: หากระบุจะกรองเฉพาะโซนนั้น หากไม่ระบุจะแสดงภาพรวม `overall` และ `zones`)

#### ตัวอย่าง Response ภาพรวม (`GET /api/v1/forecast/metrics`):
```json
{
  "status": 200,
  "message": "Fetch Model Accuracy Summary Successfully",
  "tor_standard": "Accuracy >= 85%, Confidence Interval 80%, 24h / 7d Horizons",
  "model_performance": {
    "mode": "all",
    "overall": {
      "total_zones": 7,
      "avg_accuracy_percent": 96.6,
      "avg_r2_score": 0.9892,
      "avg_rmse": 7.96,
      "avg_mae": 5.83,
      "avg_mape": 3.4,
      "all_zones_meet_tor": true,
      "tor_min_accuracy": 85
    },
    "zones": {
      "DMA-03": {
        "accuracy_percent": 96.49,
        "r2_score": 0.9891,
        "rmse": 10.75,
        "mae": 7.98,
        "mape": 3.51,
        "tor_standard_met": true
      }
    }
  }
}
```

---

## 5. การทดสอบ (Verification)
- **Unit Test File:** `api/modules/forecast/controller/forecast_controller_test.go`
- **ครอบคลุม:**
  - พยากรณ์ 24 ชั่วโมง และ 7 วัน
  - Dynamic Fallback ของโซนใหม่ (`DMA-03`)
  - โซนที่ไม่มีในระบบ (404 Not Found)
  - Metrics ภาพรวมทั้งระบบ และ Metrics กรองรายโซน
