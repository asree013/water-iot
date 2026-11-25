---
name: pump-health-and-pttf
description: แผนพัฒนาและข้อกำหนดสถาปัตยกรรม Pump Health Score & PTTF (Predictive Maintenance สำหรับเครื่องสูบน้ำ)
---

# Pump Health Score & PTTF Module (`pump_health_pttf`)

โมดูลปัญญาประดิษฐ์และวิเคราะห์ความเสื่อมสภาพเพื่อการบำรุงรักษาเชิงคาดการณ์ (Predictive Maintenance: PdM) สำหรับเครื่องสูบน้ำ (Distribution & Booster Pumps) ในระบบประปา ตามข้อกำหนด TOR Tier 3 (ML Tier)

---

## 1. วัตถุประสงค์และบทบาทในระบบ
- **ประเมินดัชนีสุขภาพของเครื่องสูบน้ำ (Pump Health Score: 0–100%):** เพื่อให้วิศวกรผู้ควบคุมระบบทราบสถานะความสมบูรณ์ของเครื่องจักรในแบบ Real-time / Near Real-time
- **คาดการณ์เวลาล่วงหน้าก่อนเกิดความเสียหาย (Predicted Time To Failure: PTTF):** คำนวณ Remaining Useful Life (RUL) เป็นจำนวนวัน/ชั่วโมงก่อนที่อุปกรณ์จะขัดข้องจนเกิด Downtime
- **จำแนกสาเหตุความเสื่อมสภาพหลัก (Root Cause Identification):** ระบุว่าความเสื่อมเกิดจาก ตลับลูกปืน (Bearing), การสั่นสะเทือน (Vibration), อุณหภูมิมอเตอร์ (Motor Temperature), หรือการเกิดฟองอากาศกัดกร่อนใบพัด (Cavitation)
- **แนะนำมาตรการซ่อมบำรุงเชิงรุก (Prescriptive Maintenance Action):** เช่น แนะนำการอัดจาระบี, การตั้งศูนย์เพลา (Shaft Alignment), หรือการเปลี่ยนตลับลูกปืนล่วงหน้า

---

## 2. พารามิเตอร์เซนเซอร์และเกณฑ์มาตรฐานสากล (Engineering Standards)

| พารามิเตอร์ | หน่วย | มาตรฐาน / เกณฑ์ประเมิน | ช่วงปกติ (Normal) | ช่วงเตือน (Warning) | ช่วงวิกฤต (Critical) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Vibration (RMS)** | $mm/s$ | **ISO 10816-3** (Zone A-D) | $< 2.8$ | $2.8 - 4.5$ | $> 7.1$ |
| **Bearing Temp** | $^{\circ}C$ | อุณหภูมิตลับลูกปืน | $< 65$ | $65 - 80$ | $> 80$ |
| **Motor Winding Temp** | $^{\circ}C$ | อุณหภูมิขดลวดมอเตอร์ | $< 75$ | $75 - 85$ | $> 90$ |
| **Hydraulic Efficiency** | $\%$ | $Q \times H / P$ (การกินไฟต่อแรงดัน) | $> 75\%$ | $60 - 75\%$ | $< 60\%$ |
| **Cavitation Index** | Index | Suction vs Discharge Delta | Normal ($< 1.0$) | High ($1.0 - 2.0$) | Severe ($> 2.0$) |
| **Operating Hours** | ชม. | ชั่วโมงทำงานตั้งแต่ Overhaul ล่าสุด | $< 4,000$ | $4,000 - 8,000$ | $> 8,000$ |

---

## 3. ระดับคะแนนสุขภาพและมาตรการตอบสนอง (Health Score Tiers)

- **90 – 100% (EXCELLENT):** เครื่องสูบน้ำสมบูรณ์ เดินเครื่องเต็มประสิทธิภาพ
- **75 – 89% (GOOD):** สภาพปกติ ไม่พบสัญญาณความเสื่อมสภาพที่มีนัยสำคัญ
- **50 – 74% (WARNING / MONITOR):** เริ่มตรวจพบความเสื่อมสภาพ แนะนำเพิ่มความถี่ในการตรวจเช็คและเตรียมอะไหล่
- **< 50% (CRITICAL):** ปั๊มอยู่ในสภาวะเสี่ยงชำรุดเสียหายร้ายแรง ต้องออกใบงานสั่งซ่อมบำรุง (Work Order) ทันที

---

## 4. แผนงานสถาปัตยกรรม (Implementation Plan)

### 4.1 ส่วน Machine Learning Pipeline (`python-etl/`)
- ไฟล์: `python-etl/pump_health_model.py`
- การสร้างข้อมูลจำลอง Telemetry รายชั่วโมงของเครื่องสูบน้ำประจำสถานีสูบน้ำใน 4 จังหวัดนำร่อง และโซนนิคมอุตสาหกรรม
- โมเดล Health Score: Multi-sensor Degradation Fusion
- โมเดล PTTF: Remaining Useful Life (RUL) Regression ผสานกับ Weibull Degradation Curve
- ผลลัพธ์: บันทึกลง `python-etl/data/pump_health_summary.json`

### 4.2 ส่วน Backend API Module (`api/modules/pump_health/`)
- **Entities (`api/modules/entities/pump.go`):**
  - `TPumpAsset`: ข้อมูลคุณลักษณะของปั๊มน้ำ (พิกัดแรงดัน, อัตราการไหล, กำลังมอเตอร์, สถานีสูบน้ำ)
  - `TPumpHealthScore`: ข้อมูลคะแนนสุขภาพ, ค่า PTTF, สถานะความเสี่ยง, และคำแนะนำ
  - `TPumpTelemetry`: ข้อมูลตรวจวัด Sensor ล่าสุด (Vibration, Bearing Temp, Motor Temp, Current, Flow)
- **Repository & Usecase:**
  - จัดการข้อมูลใน Database พร้อม Dynamic Fallback เพื่อคำนวณคะแนนให้อัตโนมัติหากยังไม่มีข้อมูล
- **Controller & Endpoints:**
  - `GET /api/v1/pump/health`: สรุปภาพรวม Fleet Health ของปั๊มน้ำทุกตัว
  - `GET /api/v1/pump/health/:pump_id`: รายละเอียดสุขภาพ Sensor และ PTTF เจาะจงตัวปั๊ม
  - `GET /api/v1/pump/health/summary`: สถิติสรุปสำหรับ Executive Dashboard (จำนวนปั๊มปกติ/เตือน/วิกฤต)
- **Database Seeder & Swagger Documentation:**
  - สร้าง Seed Data และอัปเดต Swagger `/docs`
- **Unit Tests:**
  - ทดสอบการคำนวณและ Endpoint ทั้งหมด

---

## 5. รายการ API Endpoints (RESTful API)

- **`GET /api/v1/pump/health`**: สรุปภาพรวม Fleet Health ของปั๊มน้ำทุกตัวในระบบ พร้อม Summary
- **`GET /api/v1/pump/health/:pump_id`**: รายละเอียดสุขภาพ ค่าเซนเซอร์ SCADA Telemetry และ PTTF เจาะจงตัวปั๊ม
- **`GET /api/v1/pump/summary`**: สถิติภาพรวมสถานะเครื่องสูบน้ำทั้งระบบสำหรับ Executive Dashboard

### ตัวอย่าง Response (`GET /api/v1/pump/health/PUMP-SK-01`):
```json
{
  "status": 200,
  "message": "Fetch Pump Health Successfully",
  "data": {
    "id": 1,
    "pump_id": "PUMP-SK-01",
    "pump_name": "เครื่องสูบน้ำแรงดันสูง หมายเลข 1 (สงขลา)",
    "station_name": "สถานีผลิตน้ำสงขลา (หัวเขาแดง)",
    "dma_code": "DMA-SK-01",
    "health_score": 88.3,
    "health_status": "GOOD",
    "pttf_days": 166,
    "pttf_hours": 3984,
    "failure_risk": "LOW",
    "primary_stress_factor": "NORMAL_OPERATION",
    "recommendation": "สภาพเครื่องสูบน้ำสมบูรณ์ดีมาก ทำงานในเกณฑ์มาตรฐาน ISO Zone A ตรวจเช็คตามรอบปกติในอีก 60 วัน",
    "iso_10816_zone": "Zone A",
    "telemetry": {
      "vibration_rms_mm_s": 2.15,
      "bearing_temp_c": 57.8,
      "motor_temp_c": 63.4,
      "flow_rate_m3_h": 485,
      "head_bar": 4.9,
      "power_kw": 88.5,
      "efficiency_percent": 74.6,
      "operating_hours": 3420,
      "cavitation_index": 0.42
    },
    "evaluated_at": "2026-11-25T14:22:12Z"
  }
}
```

---

## 6. การทดสอบและการรับรอง (Verification)
- **Unit Test File:** `api/modules/pump_health/controller/pump_test.go`
- **ครอบคลุม:**
  - Fleet Health List และการคำนวณ `UrgentActionRequired`
  - ข้อมูล Telemetry และ PTTF เจาะจงตัวปั๊ม (`Zone A`, `Zone D`, etc.)
  - Executive Fleet Summary
  - กรณีไม่พบปั๊มในระบบ (404 Not Found)
- **สถานะ:** ทดสอบผ่าน 100% (`PASS: TestGetFleetHealth_Success`)

