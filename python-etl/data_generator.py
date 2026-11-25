"""
water_demand_data_generator.py
สร้างชุดข้อมูลอนุกรมเวลาย้อนหลัง (Historical Time-Series Data) สำหรับการพยากรณ์ความต้องการใช้น้ำ
ครอบคลุม 4 โซนนำร่องตาม TOR: สงขลา, เพชรบุรี, กาญจนบุรี, กทม
อ้างอิงพฤติกรรมการใช้น้ำตามงานวิจัย Herrera et al. (2010)
"""

import os
import math
import numpy as np
import pandas as pd
from datetime import datetime, timedelta

def get_hourly_pattern(hour, is_weekend):
    """
    แบบรูปการใช้น้ำรายชั่วโมงตามธรรมชาติของประปาไทย:
    - 00:00 - 04:00: ช่วงใช้น้ำต่ำสุด (Trough)
    - 06:00 - 09:00: ช่วงใช้น้ำสูงสุดช่วงเช้า (Morning Peak)
    - 11:00 - 13:00: ช่วงเที่ยง (Midday Plateau)
    - 17:00 - 20:00: ช่วงใช้น้ำสูงสุดช่วงเย็น (Evening Peak)
    """
    if is_weekend:
        # วันหยุด ตื่นสายกว่า Peak เช้าเลื่อนไป 08:00-10:00
        weights = [
            0.35, 0.30, 0.28, 0.28, 0.32, 0.45,
            0.65, 0.90, 1.25, 1.35, 1.25, 1.15,
            1.10, 1.05, 1.00, 1.00, 1.10, 1.25,
            1.35, 1.30, 1.15, 0.95, 0.70, 0.50
        ]
    else:
        # วันธรรมดา Peak เช้า 06:00-08:30 และ Peak เย็น 17:30-20:00
        weights = [
            0.32, 0.28, 0.25, 0.26, 0.35, 0.60,
            1.15, 1.38, 1.20, 0.95, 0.90, 1.00,
            0.95, 0.90, 0.90, 0.95, 1.10, 1.30,
            1.40, 1.28, 1.05, 0.85, 0.60, 0.42
        ]
    return weights[hour]

def generate_historical_demand_data(days=180, output_path=None):
    """
    สร้างข้อมูลย้อนหลัง 180 วัน (ประมาณ 6 เดือน) ที่ความถี่รายชั่วโมง
    """
    np.random.seed(42)
    end_time = datetime(2026, 11, 25, 12, 0, 0)
    start_time = end_time - timedelta(days=days)
    
    # โซนหลักตามฐานข้อมูลและ TOR
    zones_meta = {
        "DMA-01": {
            "name": "โซนเทศบาลนครหาดใหญ่ 1",
            "province": "สงขลา",
            "base_demand": 210.0,
            "base_pressure": 2.4,
            "temp_mean": 29.5
        },
        "DMA-02": {
            "name": "โซนคลองแห",
            "province": "สงขลา",
            "base_demand": 160.0,
            "base_pressure": 2.2,
            "temp_mean": 29.5
        },
        "DMA-03": {
            "name": "โซนนิคมอุตสาหกรรมภาคใต้",
            "province": "สงขลา",
            "base_demand": 280.0,
            "base_pressure": 2.6,
            "temp_mean": 29.8
        },
        "DMA-SK-01": {
            "name": "โซนเทศบาลนครสงขลา",
            "province": "สงขลา",
            "base_demand": 145.0,  # ลบ.ม./ชม.
            "base_pressure": 2.2,
            "temp_mean": 29.5
        },
        "DMA-PB-01": {
            "name": "โซนเทศบาลเมืองเพชรบุรี",
            "province": "เพชรบุรี",
            "base_demand": 115.0,
            "base_pressure": 2.5,
            "temp_mean": 30.2
        },
        "DMA-KN-01": {
            "name": "โซนเทศบาลเมืองกาญจนบุรี",
            "province": "กาญจนบุรี",
            "base_demand": 175.0,
            "base_pressure": 2.8,
            "temp_mean": 31.0
        },
        "DMA-BKK-01": {
            "name": "โซนดินแดง-ห้วยขวาง",
            "province": "กรุงเทพมหานคร",
            "base_demand": 350.0,
            "base_pressure": 2.0,
            "temp_mean": 32.0
        }
    }
    
    # สร้าง time index
    timestamps = pd.date_range(start=start_time, end=end_time, freq="h")
    records = []
    
    for dma_code, meta in zones_meta.items():
        base_demand = meta["base_demand"]
        base_pressure = meta["base_pressure"]
        temp_mean = meta["temp_mean"]
        
        for ts in timestamps:
            hour = ts.hour
            dow = ts.weekday()
            is_wknd = 1 if dow >= 5 else 0
            
            # 1. Hourly Pattern Weight
            hourly_factor = get_hourly_pattern(hour, is_wknd)
            
            # 2. Day-of-week factor (เสาร์-อาทิตย์ กาญจน์/ชะอำใช้น้ำเยอะขึ้นเพราะท่องเที่ยว)
            if dma_code in ["DMA-KN-01", "DMA-PB-01"] and is_wknd:
                dow_factor = 1.15
            elif dow == 0:  # วันจันทร์
                dow_factor = 1.05
            else:
                dow_factor = 1.0
                
            # 3. Weather / Temp variation
            temp = temp_mean + 4.0 * math.sin(2 * math.pi * (hour - 9) / 24) + np.random.normal(0, 1.0)
            temp_factor = 1.0 + 0.008 * (temp - temp_mean)
            
            # 4. Long-term trend & Seasonal Noise
            day_of_year = ts.timetuple().tm_yday
            seasonal_factor = 1.0 + 0.05 * math.sin(2 * math.pi * day_of_year / 365)
            
            # Random gaussian noise (~3%)
            noise = np.random.normal(1.0, 0.035)
            
            demand = base_demand * hourly_factor * dow_factor * temp_factor * seasonal_factor * noise
            demand = max(demand, base_demand * 0.2)  # ค่าต่ำสุดไม่ต่ำกว่า 20%
            
            # Pressure inversely related to demand (แรงดันตกช่วงใช้น้ำเยอะ)
            pressure = base_pressure - (demand / base_demand - 1.0) * 0.3 + np.random.normal(0, 0.05)
            pressure = max(0.5, round(pressure, 2))
            
            records.append({
                "timestamp": ts.strftime("%Y-%m-%d %H:%M:%S"),
                "dma_code": dma_code,
                "province": meta["province"],
                "actual_demand": round(demand, 2),
                "pressure": pressure,
                "temperature": round(temp, 1),
                "hour": hour,
                "day_of_week": dow,
                "is_weekend": is_wknd
            })
            
    df = pd.DataFrame(records)
    
    if output_path:
        os.makedirs(os.path.dirname(output_path), exist_ok=True)
        df.to_csv(output_path, index=False)
        print(f"[Data Generator] Successfully created dataset with {len(df)} records at: {output_path}")
        
    return df

if __name__ == "__main__":
    out_dir = os.path.join(os.path.dirname(__file__), "data")
    csv_file = os.path.join(out_dir, "water_demand_historical.csv")
    generate_historical_demand_data(days=180, output_path=csv_file)
