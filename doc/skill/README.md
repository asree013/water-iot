# Water IoT System Modules Documentation (`doc/skill/`)

เอกสารบันทึกรายละเอียดทางเทคนิค แผนการพัฒนา สถาปัตยกรรม Schema ฐานข้อมูล และรายการ API Endpoint ของแต่ละโมดูลในระบบ Water IoT ตามข้อกำหนด TOR

---

## สารบัญโมดูลทั้งหมด (Module Directory)

| โมดูล | โฟลเดอร์เอกสาร | รายละเอียดและขอบเขตงาน | สถานะการพัฒนา |
| :--- | :--- | :--- | :--- |
| **DMA Zone** | [`dma_zone/SKILL.md`](./dma_zone/SKILL.md) | จัดการข้อมูลพื้นที่จ่ายน้ำย่อย District Metered Area, NRW %, แรงดันเป้าหมาย และปริมาณน้ำจ่ายเข้า | ✅ เสร็จสมบูรณ์ (100%) |
| **Pipe Asset** | [`pipe_asset/SKILL.md`](./pipe_asset/SKILL.md) | ทะเบียนสินทรัพย์เส้นท่อประปา (วัสดุ, ขนาดเส้นผ่านศูนย์กลาง, ความยาว, ปีติดตั้ง, สภาพท่อ) | ✅ เสร็จสมบูรณ์ (100%) |
| **Leakage Incident** | [`leakage_incident/SKILL.md`](./leakage_incident/SKILL.md) | บันทึกและติดตามเหตุการณ์ท่อแตก/รั่ว, ประเมินปริมาณน้ำสูญเสีย, พิกัดจุดเกิดเหตุ และสถานะงานซ่อม | ✅ เสร็จสมบูรณ์ (100%) |
| **Demand Forecasting** | [`demand_forecast/SKILL.md`](./demand_forecast/SKILL.md) | พยากรณ์ความต้องการใช้น้ำล่วงหน้า 24 ชม. และ 7 วัน, 80% CI, แนะนำการเดินปั๊ม Peak/Off-Peak และประเมินความแม่นยำ ($\ge 85\%$) | ✅ เสร็จสมบูรณ์ (100%) |
| **Pump Health & PTTF** | [`pump_health_pttf/SKILL.md`](./pump_health_pttf/SKILL.md) | ดัชนีสุขภาพเครื่องสูบน้ำ (Health Score 0–100%), เวลาคาดการณ์ก่อนชำรุด (PTTF), วิเคราะห์ความเสื่อมสภาพตาม ISO 10816-3 | ✅ เสร็จสมบูรณ์ (100%) |

---

## โครงสร้างสถาปัตยกรรมภาพรวมของระบบ (5 Tiers Architecture)

```
1. MQTT STREAM DATA TIER
   IoT Meters / SCADA RTU ──> MQTT Broker ──> TimescaleDB / Database
2. API & GIS TIER (Go Fiber + GORM + PostGIS)
   CRUD REST API, RBAC 4 ระดับ, ท่อประปา, พื้นที่ DMA, จุดมิเตอร์ผู้ใช้น้ำ
3. ML TIER (Python + Scikit-learn + Statsmodels)
   • Demand Forecasting (24h / 7d)
   • Pump Health Score & PTTF
   • MNF Leak & Burst Anomaly Detection
4. RAG TIER (Vector DB + Local Thai LLM)
   คู่มือ SOP, กฎเกณฑ์, สรุปสถานการณ์รายวัน, Chatbot ผู้ช่วย กปภ.
5. PRESENTATION TIER (Single Pane of Glass)
   Interactive Map, SCADA Canvas 2D, Real-time Charts, Executive Cards
```
