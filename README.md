# 🌊 PWA Water-IoT & Big Data Management System

ระบบบริหารจัดการข้อมูลขนาดใหญ่และการตรวจวัด IoT สำหรับติดตาม วิเคราะห์น้ำสูญเสีย (Non-Revenue Water: NRW) และเฝ้าระวังแรงดันน้ำประปา ตามข้อกำหนดของ**การประปาส่วนภูมิภาค (กปภ.)**

---

## 📌 ภาพรวมโครงการ (Project Overview)

โครงการนี้พัฒนาขึ้นเพื่อยกระดับการบริหารจัดการข้อมูลของ กปภ. ให้เป็นระบบรวมศูนย์ (Single Source of Truth) รองรับทั้งการประมวลผลข้อมูลสตรีมมิ่งแบบเรียลไทม์ (**Real-Time IoT Streaming**) ผ่าน **Mosquitto MQTT Broker** และการรวบรวมข้อมูลจากระบบเดิม (**Legacy Systems**) เช่น **SWMAP, DMAMA, GIS, OCS, SAP** มาจัดเก็บในฐานข้อมูลกลาง เพื่อนำไปแสดงผลบน **Executive Dashboard** และสนับสนุนการตัดสินใจเชิงกลยุทธ์

- **Backend API:** Go (Golang) + Fiber v2 + GORM + Clean Architecture
- **Database:** MySQL 8.0
- **Message Broker:** Eclipse Mosquitto MQTT (Port 1883)
- **API Documentation:** Scalar API Reference (`http://localhost:8080/docs`)
- **IoT Streaming Simulator:** โมดูล `stream-data` จำลองส่งข้อมูลเซ็นเซอร์ทุก 1 วินาที
- **Batch Data Pipeline:** โมดูล `python-etl` ทำ Data Validation, Cleansing และ Dead Letter Queue (DLQ)

---

## 📖 พจนานุกรมข้อมูล (Data Dictionary)

---

### 1. ตาราง `water_data` (Real-time Telemetry Data)
> **วัตถุประสงค์:** จัดเก็บข้อมูลการตรวจวัดสดจากเซ็นเซอร์ IoT ตามสถานีสูบน้ำและแนวท่อส่งน้ำผ่าน MQTT

| ชื่อฟิลด์ (Field) | ประเภทข้อมูล (Type) | คีย์ (Key) | ว่างได้ (Null) | คำอธิบาย (Description) | ตัวอย่างข้อมูล (Example) |
| :--- | :--- | :---: | :---: | :--- | :--- |
| `id` | `BIGINT UNSIGNED` | **PK** | No | ลำดับไอดีของข้อมูล (Auto Increment) | `1` |
| `zone` | `VARCHAR(255)` | - | Yes | โซนพื้นที่ที่ติดตั้งเซ็นเซอร์ (ผูกกับ DMA) | `"Zone_A"` |
| `pipe` | `VARCHAR(255)` | - | Yes | รหัสเส้นท่อที่ติดตั้งเซ็นเซอร์ | `"Pipe_01"` |
| `pressure` | `DOUBLE` | - | Yes | ค่าแรงดันน้ำในท่อ (bar) | `2.85` |
| `flow_rate` | `DOUBLE` | - | Yes | อัตราการไหลของน้ำ (m³/h) | `45.20` |
| `temperature` | `DOUBLE` | - | Yes | อุณหภูมิของน้ำ (°C) | `28.50` |
| `vibration` | `DOUBLE` | - | Yes | ค่าความสั่นสะเทือนของอุปกรณ์ | `0.12` |
| `rpm` | `DOUBLE` | - | Yes | รอบการหมุนของเครื่องสูบน้ำ | `1450.0` |
| `operational_hours` | `INT` | - | Yes | ชั่วโมงการทำงานสะสมของอุปกรณ์ | `320` |
| `location_code` | `VARCHAR(255)` | - | Yes | รหัสระบุตำแหน่งติดตั้งอุปกรณ์ | `"LOC-HY-01"` |
| `latitude` | `DOUBLE` | - | Yes | พิกัดภูมิศาสตร์ละติจูด | `13.7563` |
| `longitude` | `DOUBLE` | - | Yes | พิกัดภูมิศาสตร์ลองจิจูด | `100.5018` |
| `leakage_flag` | `INT` | - | Yes | ธงเตือนท่อรั่ว (0 = ปกติ, 1 = เสี่ยงรั่วซึม) | `0` |
| `timestamp` | `BIGINT` | - | Yes | เวลาที่เซ็นเซอร์ตรวจวัด (Unix Timestamp) | `1731427200` |
| `created_at` | `DATETIME(3)` | - | Yes | วันเวลาที่บันทึกข้อมูลเข้าฐานข้อมูล | `2026-11-12 21:00:00` |

---

### 2. ตาราง `tdma_zones` (DMA / Zone Master Data)
> **วัตถุประสงค์:** จัดเก็บข้อมูลพื้นที่จ่ายน้ำย่อยจากระบบ **DMAMA** ใช้บริหารจัดการแรงดันและคำนวณน้ำสูญเสีย (NRW)

| ชื่อฟิลด์ (Field) | ประเภทข้อมูล (Type) | คีย์ (Key) | ว่างได้ (Null) | คำอธิบาย (Description) | ตัวอย่างข้อมูล (Example) |
| :--- | :--- | :---: | :---: | :--- | :--- |
| `id` | `BIGINT UNSIGNED` | **PK** | No | ลำดับไอดีของโซน (Auto Increment) | `1` |
| `dma_code` | `VARCHAR(50)` | **UK** | No | รหัสประจำพื้นที่จ่ายน้ำย่อย (Unique Key) | `"DMA-01"` |
| `dma_name` | `VARCHAR(100)` | - | Yes | ชื่อพื้นที่จ่ายน้ำย่อย | `"โซนเทศบาลนครหาดใหญ่ 1"` |
| `branch_name` | `VARCHAR(100)` | - | Yes | สาขา กปภ. ที่รับผิดชอบ (นำร่อง 4 สาขา) | `"สาขาหาดใหญ่"` |
| `region_code` | `VARCHAR(20)` | - | Yes | รหัสเขตการประปา (10 เขตทั่วประเทศ) | `"เขต 5"` |
| `target_pressure` | `DOUBLE` | - | Yes | แรงดันน้ำเป้าหมายในโซน (bar) | `2.5` |
| `total_inflow` | `DOUBLE` | - | Yes | ปริมาณน้ำรวมที่จ่ายเข้าโซน (ลบ.ม./วัน) | `10000.0` |
| `billed_consumption`| `DOUBLE` | - | Yes | ปริมาณน้ำที่จำหน่ายได้จริง (ลบ.ม./วัน) | `8200.0` |
| `nrw_ratio` | `DOUBLE` | - | Yes | อัตราน้ำสูญเสีย % `((Inflow - Billed)/Inflow)*100` | `18.0` |
| `customer_count` | `INT` | - | Yes | จำนวนผู้ใช้น้ำในโซน (ราย) | `4200` |
| `status` | `VARCHAR(20)` | - | Yes | สถานะความเสี่ยง (`NORMAL`, `WARNING`, `CRITICAL`) | `"NORMAL"` |
| `created_at` | `DATETIME(3)` | - | Yes | วันเวลาที่สร้างข้อมูล | `2026-11-12 21:00:00` |
| `updated_at` | `DATETIME(3)` | - | Yes | วันเวลาที่แก้ไขข้อมูลล่าสุด | `2026-11-12 21:00:00` |

---

### 3. ตาราง `t_pipe_assets` (Pipe & Asset Master Data)
> **วัตถุประสงค์:** จัดเก็บข้อมูลทะเบียนสินทรัพย์ท่อส่งน้ำจากระบบ **SAP** ผสานกับพิกัดแนวท่อจากระบบ **SWMAP/GIS**

| ชื่อฟิลด์ (Field) | ประเภทข้อมูล (Type) | คีย์ (Key) | ว่างได้ (Null) | คำอธิบาย (Description) | ตัวอย่างข้อมูล (Example) |
| :--- | :--- | :---: | :---: | :--- | :--- |
| `id` | `BIGINT UNSIGNED` | **PK** | No | ลำดับไอดีสินทรัพย์ (Auto Increment) | `1` |
| `asset_id` | `VARCHAR(50)` | **UK** | No | รหัสสินทรัพย์จาก SAP (ใช้ตรวจ DLQ) | `"AST-001"` |
| `pipe_id` | `VARCHAR(50)` | **IDX** | Yes | รหัสประจำเส้นท่อ | `"Pipe_01"` |
| `dma_code` | `VARCHAR(50)` | **IDX** | Yes | รหัสโซน DMA ที่ท่อนี้สังกัดอยู่ | `"DMA-01"` |
| `material` | `VARCHAR(50)` | - | Yes | วัสดุท่อ เช่น `HDPE`, `PVC`, `Ductile Iron` | `"HDPE"` |
| `diameter` | `DOUBLE` | - | Yes | ขนาดเส้นผ่านศูนย์กลางท่อ (มม.) | `150.0` |
| `length` | `DOUBLE` | - | Yes | ความยาวของเส้นท่อ (เมตร) | `500.0` |
| `installation_year`| `INT` | - | Yes | ปี ค.ศ. ที่ติดตั้งท่อ (ใช้คำนวณอายุท่อ) | `2018` |
| `latitude` | `DOUBLE` | - | Yes | พิกัด GIS ละติจูด | `13.7563` |
| `longitude` | `DOUBLE` | - | Yes | พิกัด GIS ลองจิจูด | `100.5018` |
| `status` | `VARCHAR(20)` | - | Yes | สถานะสินทรัพย์ (`ACTIVE`, `INACTIVE`, `REPLACED`) | `"ACTIVE"` |
| `created_at` | `DATETIME(3)` | - | Yes | วันเวลาที่สร้างข้อมูล | `2026-11-12 21:00:00` |
| `updated_at` | `DATETIME(3)` | - | Yes | วันเวลาที่แก้ไขข้อมูลล่าสุด | `2026-11-12 21:00:00` |

---

### 4. ตาราง `t_leakage_incidents` (Leakage & Maintenance Records)
> **วัตถุประสงค์:** จัดเก็บประวัติการรับแจ้งเหตุและบันทึกการซ่อมแซมท่อแตกรั่วจากระบบ **OCS**

| ชื่อฟิลด์ (Field) | ประเภทข้อมูล (Type) | คีย์ (Key) | ว่างได้ (Null) | คำอธิบาย (Description) | ตัวอย่างข้อมูล (Example) |
| :--- | :--- | :---: | :---: | :--- | :--- |
| `id` | `BIGINT UNSIGNED` | **PK** | No | ลำดับไอดีของบันทึกเหตุการณ์ (Auto Increment) | `1` |
| `ticket_id` | `VARCHAR(50)` | **UK** | No | รหัสใบแจ้งซ่อมจากระบบ OCS | `"OCS-2026-001"` |
| `dma_code` | `VARCHAR(50)` | **IDX** | Yes | โซน DMA ที่เกิดเหตุ | `"DMA-01"` |
| `pipe_id` | `VARCHAR(50)` | - | Yes | รหัสเส้นท่อที่เกิดเหตุ | `"Pipe_01"` |
| `incident_type` | `VARCHAR(50)` | - | Yes | ประเภทเหตุการณ์ เช่น `Main Burst`, `Service Leak` | `"Main Burst"` |
| `severity` | `VARCHAR(20)` | - | Yes | ระดับความรุนแรง (`LOW`, `MEDIUM`, `HIGH`, `CRITICAL`) | `"HIGH"` |
| `location_description` | `TEXT` | - | Yes | คำอธิบายจุดเกิดเหตุหรือสถานที่สังเกต | `"หน้าตลาดสดเทศบาล ซอย 2"` |
| `latitude` | `DOUBLE` | - | Yes | พิกัดจุดเกิดเหตุ ละติจูด | `13.7580` |
| `longitude` | `DOUBLE` | - | Yes | พิกัดจุดเกิดเหตุ ลองจิจูด | `100.5030` |
| `water_loss_estimate` | `DOUBLE` | - | Yes | ประมาณการปริมาณน้ำสูญเสีย (ลบ.ม.) | `120.5` |
| `reported_at` | `DATETIME(3)` | - | Yes | วันเวลาที่ได้รับแจ้งเหตุ | `2026-11-12 10:30:00` |
| `repaired_at` | `DATETIME(3)` | - | Yes | วันเวลาที่ซ่อมแซมเสร็จสิ้น (Null ได้ถ้ายังไม่เสร็จ) | `2026-11-12 14:00:00` |
| `repair_status` | `VARCHAR(30)` | - | Yes | สถานะ (`PENDING`, `IN_PROGRESS`, `RESOLVED`, `CLOSED`) | `"RESOLVED"` |
| `created_at` | `DATETIME(3)` | - | Yes | วันเวลาที่สร้างข้อมูล | `2026-11-12 10:35:00` |
| `updated_at` | `DATETIME(3)` | - | Yes | วันเวลาที่แก้ไขข้อมูลล่าสุด | `2026-11-12 14:05:00` |

---

## 🚀 การเริ่มต้นใช้งาน (Getting Started)

1. **รัน Database และ MQTT Broker:**
   ```bash
   docker compose up -d mysql_db mqtt_broker
   ```
2. **รันเซิร์ฟเวอร์ Backend API:**
   ```bash
   cd api
   go run apps/main.go
   ```
   - เข้าดูเอกสาร API เอกสารสไตล์ **Scalar UI**: `http://localhost:8080/docs`
   - ตรวจสอบสถานะเซิร์ฟเวอร์: `http://localhost:8080/health`
3. **รันตัวจำลอง IoT Sensor Streaming:**
   ```bash
   cd stream-data
   go run server/main.go
   ```
