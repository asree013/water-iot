"""
main.py
Water Demand Forecasting Pipeline (ML Tier 3 - TOR กปภ.)
รันกระบวนการ Train โมเดลและสร้างผลการพยากรณ์ความต้องการใช้น้ำ (24 ชั่วโมง และ 7 วัน)
พร้อมช่วงความเชื่อมั่น 80% (Confidence Interval 80%)
"""

import os
import json
import pandas as pd
from data_generator import generate_historical_demand_data
from demand_forecaster import WaterDemandForecaster

ZONES = ["DMA-01", "DMA-02", "DMA-03", "DMA-SK-01", "DMA-PB-01", "DMA-KN-01", "DMA-BKK-01"]

def run_demand_forecasting_pipeline():
    base_dir = os.path.dirname(__file__)
    data_dir = os.path.join(base_dir, "data")
    model_dir = os.path.join(base_dir, "models")
    os.makedirs(data_dir, exist_ok=True)
    os.makedirs(model_dir, exist_ok=True)

    csv_path = os.path.join(data_dir, "water_demand_historical.csv")
    if not os.path.exists(csv_path):
        print("[Pipeline] Generating historical time-series dataset...")
        df_raw = generate_historical_demand_data(days=180, output_path=csv_path)
    else:
        print(f"[Pipeline] Loading historical dataset from: {csv_path}")
        df_raw = pd.read_csv(csv_path)

    all_24h_forecasts = []
    all_7d_forecasts = []
    model_summaries = {}

    print("=" * 65)
    print("      WATER DEMAND FORECASTING (24h / 7d) - TOR COMPLIANT")
    print("=" * 65)

    for dma_code in ZONES:
        print(f"\n>>> Processing Zone: {dma_code}")
        forecaster = WaterDemandForecaster(dma_code=dma_code)
        
        # 1. Train Model
        metrics = forecaster.train(df_raw)
        model_summaries[dma_code] = metrics
        forecaster.save_model(model_dir)

        # 2. Generate 7-day (168 hours) Multi-step Forecast with 80% CI
        forecast_168h = forecaster.forecast_multistep(df_raw, horizon_hours=168)

        # 3. แยก 24 ชั่วโมงแรกสำหรับควบคุมการสูบน้ำ Peak/Off-Peak
        forecast_24h = forecast_168h[:24]

        all_24h_forecasts.extend(forecast_24h)
        all_7d_forecasts.extend(forecast_168h)

    # บันทึกผลลัพธ์ลง JSON Files
    out_24h_path = os.path.join(data_dir, "demand_forecast_24h.json")
    out_7d_path = os.path.join(data_dir, "demand_forecast_7d.json")
    out_summary_path = os.path.join(data_dir, "model_accuracy_summary.json")

    with open(out_24h_path, "w", encoding="utf-8") as f:
        json.dump(all_24h_forecasts, f, ensure_ascii=False, indent=2)

    with open(out_7d_path, "w", encoding="utf-8") as f:
        json.dump(all_7d_forecasts, f, ensure_ascii=False, indent=2)

    with open(out_summary_path, "w", encoding="utf-8") as f:
        json.dump(model_summaries, f, ensure_ascii=False, indent=2)

    print("\n" + "=" * 65)
    print("                    EVALUATION SUMMARY")
    print("=" * 65)
    for zone, met in model_summaries.items():
        status = "PASSED" if met["tor_standard_met"] else "FAILED"
        print(f"[{zone}] Accuracy: {met['accuracy_percent']}% | R2: {met['r2_score']} | Status: {status}")

    print("\n[Pipeline] Forecast artifacts created successfully:")
    print(f"  - 24h Forecast: {out_24h_path} ({len(all_24h_forecasts)} hourly records)")
    print(f"  - 7d Forecast:  {out_7d_path} ({len(all_7d_forecasts)} hourly records)")
    print(f"  - Accuracy:    {out_summary_path}")

    return all_24h_forecasts, all_7d_forecasts, model_summaries

if __name__ == "__main__":
    run_demand_forecasting_pipeline()
