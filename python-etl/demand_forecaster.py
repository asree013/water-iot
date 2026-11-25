"""
demand_forecaster.py
โมเดลพยากรณ์ความต้องการใช้น้ำ (Water Demand Forecasting)
ตามเกณฑ์ TOR โครงการ กปภ. (ภาคผนวก ๑ ข้อ ๑.๔.๑, ภาคผนวก ๓ และ ภาคผนวก ๔ ข้อ ๑.๑)
อ้างอิง Herrera et al. (2010): Time-series Feature Engineering + Machine Learning

ฟีเจอร์หลัก:
1. พยากรณ์ล่วงหน้า 24 ชั่วโมง (24-Hour Horizon, รายชั่วโมง) สำหรับควบคุมการสูบน้ำ Peak/Off-Peak
2. พยากรณ์ล่วงหน้า 7 วัน (7-Day Horizon, 168 ชั่วโมง และสรุปรายวัน)
3. ช่วงความเชื่อมั่น 80% (Confidence Interval 80%: Upper Bound & Lower Bound)
4. ประเมินความแม่นยำ (R2 >= 0.85, MAPE <= 15% ตามเกณฑ์ TOR ความแม่นยำ >= 85%)
"""

import os
import joblib
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
from sklearn.ensemble import GradientBoostingRegressor, RandomForestRegressor
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score

# Z-score สำหรับ 80% Confidence Interval (2-tailed: alpha=0.20 -> z = 1.28155)
Z_80_CI = 1.28155

class WaterDemandForecaster:
    def __init__(self, dma_code="DMA-BKK-01"):
        self.dma_code = dma_code
        self.model = None
        self.feature_cols = []
        self.residual_std = 0.0
        self.metrics = {}

    def extract_features(self, df):
        """
        Feature Engineering อ้างอิงตามหลัก Herrera et al. (2010):
        - Lag Features: t-1, t-2, t-24 (วันก่อน), t-168 (สัปดาห์ก่อน)
        - Rolling Window: ค่าเฉลี่ยและส่วนเบี่ยงเบนมาตรฐานรอบ 24 ชั่วโมง
        - Cyclical Calendar: sin/cos ของชั่วโมงและวันในสัปดาห์
        """
        data = df.copy()
        if not pd.api.types.is_datetime64_any_dtype(data["timestamp"]):
            data["timestamp"] = pd.to_datetime(data["timestamp"])
        data = data.sort_values("timestamp").reset_index(drop=True)

        # 1. Time Features
        data["hour"] = data["timestamp"].dt.hour
        data["day_of_week"] = data["timestamp"].dt.dayofweek
        data["is_weekend"] = (data["day_of_week"] >= 5).astype(int)

        # 2. Cyclical Features (ทำให้ 23:00 เชื่อมโยงกับ 00:00)
        data["sin_hour"] = np.sin(2 * np.pi * data["hour"] / 24.0)
        data["cos_hour"] = np.cos(2 * np.pi * data["hour"] / 24.0)
        data["sin_dow"] = np.sin(2 * np.pi * data["day_of_week"] / 7.0)
        data["cos_dow"] = np.cos(2 * np.pi * data["day_of_week"] / 7.0)

        # 3. Lag Features
        data["lag_1"] = data["actual_demand"].shift(1)
        data["lag_2"] = data["actual_demand"].shift(2)
        data["lag_24"] = data["actual_demand"].shift(24)
        data["lag_168"] = data["actual_demand"].shift(168)

        # 4. Rolling Statistics
        data["rolling_mean_24"] = data["actual_demand"].shift(1).rolling(window=24, min_periods=1).mean()
        data["rolling_std_24"] = data["actual_demand"].shift(1).rolling(window=24, min_periods=1).std().fillna(0)
        data["rolling_max_24"] = data["actual_demand"].shift(1).rolling(window=24, min_periods=1).max()
        data["rolling_min_24"] = data["actual_demand"].shift(1).rolling(window=24, min_periods=1).min()

        # Drop rows where lag_168 is NaN
        data = data.dropna().reset_index(drop=True)
        return data

    def train(self, df_raw):
        """
        ฝึกสอนโมเดล Gradient Boosting Regressor สำหรับ DMA Zone
        """
        df_zone = df_raw[df_raw["dma_code"] == self.dma_code].copy()
        if len(df_zone) == 0:
            raise ValueError(f"No data found for dma_code: {self.dma_code}")

        df_feat = self.extract_features(df_zone)

        self.feature_cols = [
            "hour", "day_of_week", "is_weekend",
            "sin_hour", "cos_hour", "sin_dow", "cos_dow",
            "temperature", "pressure",
            "lag_1", "lag_2", "lag_24", "lag_168",
            "rolling_mean_24", "rolling_std_24", "rolling_max_24", "rolling_min_24"
        ]

        X = df_feat[self.feature_cols]
        y = df_feat["actual_demand"]

        # Train/Test Split (80% train, 20% test ตามลำดับเวลา)
        split_idx = int(len(X) * 0.8)
        X_train, X_test = X.iloc[:split_idx], X.iloc[split_idx:]
        y_train, y_test = y.iloc[:split_idx], y.iloc[split_idx:]

        # ฝึกโมเดล GradientBoostingRegressor
        self.model = GradientBoostingRegressor(
            n_estimators=120,
            learning_rate=0.08,
            max_depth=4,
            random_state=42,
        )
        self.model.fit(X_train, y_train)

        # ประเมินผลบน Test set
        y_pred = self.model.predict(X_test)
        residuals = y_test - y_pred
        self.residual_std = float(np.std(residuals))

        # Metrics
        r2 = r2_score(y_test, y_pred)
        rmse = np.sqrt(mean_squared_error(y_test, y_pred))
        mae = mean_absolute_error(y_test, y_pred)
        mape = np.mean(np.abs(residuals / y_test)) * 100.0
        accuracy = max(0.0, 100.0 - mape)

        self.metrics = {
            "r2_score": round(float(r2), 4),
            "rmse": round(float(rmse), 2),
            "mae": round(float(mae), 2),
            "mape": round(float(mape), 2),
            "accuracy_percent": round(float(accuracy), 2),
            "tor_standard_met": bool(accuracy >= 85.0 and r2 >= 0.85)
        }

        print(f"[{self.dma_code}] Training Completed:")
        print(f"  Accuracy: {self.metrics['accuracy_percent']}% (TOR >= 85%)")
        print(f"  R2 Score: {self.metrics['r2_score']} | RMSE: {self.metrics['rmse']} m3/h | MAE: {self.metrics['mae']} m3/h")
        print(f"  TOR Pass: {'PASSED' if self.metrics['tor_standard_met'] else 'FAILED'}")

        return self.metrics

    def forecast_multistep(self, df_recent, horizon_hours=168):
        """
        พยากรณ์ล่วงหน้าแบบ Recursive Multi-step:
        - 24 ชั่วโมง (24h)
        - 7 วัน (168h)
        พร้อมคำนวณช่วงความเชื่อมั่น 80% (Confidence Interval 80%)
        """
        if self.model is None:
            raise ValueError("Model is not trained yet!")

        df_zone = df_recent[df_recent["dma_code"] == self.dma_code].copy()
        df_feat = self.extract_features(df_zone)

        last_known = df_feat.iloc[-1].to_dict()
        current_time = pd.to_datetime(last_known["timestamp"])

        history_demands = list(df_feat["actual_demand"].values)
        forecast_records = []

        # สถิติพื้นฐานของโซน
        mean_temp = df_zone["temperature"].mean()
        mean_pressure = df_zone["pressure"].mean()

        for step in range(1, horizon_hours + 1):
            future_time = current_time + timedelta(hours=step)
            hour = future_time.hour
            dow = future_time.weekday()
            is_wknd = 1 if dow >= 5 else 0

            # Lags จากประวัติที่มีอยู่บวกกับค่าพยากรณ์ก่อนหน้า
            lag_1 = history_demands[-1]
            lag_2 = history_demands[-2] if len(history_demands) >= 2 else lag_1
            lag_24 = history_demands[-24] if len(history_demands) >= 24 else lag_1
            lag_168 = history_demands[-168] if len(history_demands) >= 168 else lag_24

            # Rolling จาก 24 ค่าล่าสุด
            window_24 = history_demands[-24:]
            r_mean = float(np.mean(window_24))
            r_std = float(np.std(window_24))
            r_max = float(np.max(window_24))
            r_min = float(np.min(window_24))

            # อุณหภูมิคาดการณ์ตามช่วงเวลา
            temp_est = mean_temp + 3.5 * np.sin(2 * np.pi * (hour - 9) / 24)

            feature_dict = {
                "hour": hour,
                "day_of_week": dow,
                "is_weekend": is_wknd,
                "sin_hour": np.sin(2 * np.pi * hour / 24.0),
                "cos_hour": np.cos(2 * np.pi * hour / 24.0),
                "sin_dow": np.sin(2 * np.pi * dow / 7.0),
                "cos_dow": np.cos(2 * np.pi * dow / 7.0),
                "temperature": temp_est,
                "pressure": mean_pressure,
                "lag_1": lag_1,
                "lag_2": lag_2,
                "lag_24": lag_24,
                "lag_168": lag_168,
                "rolling_mean_24": r_mean,
                "rolling_std_24": r_std,
                "rolling_max_24": r_max,
                "rolling_min_24": r_min
            }

            feat_df = pd.DataFrame([feature_dict])[self.feature_cols]
            pred_demand = float(self.model.predict(feat_df)[0])
            pred_demand = max(0.0, round(pred_demand, 2))

            # ช่วงความเชื่อมั่น 80% (ความคลาดเคลื่อนสะสมตาม step เล็กน้อย)
            # uncertainty เติบโตตาม sqrt(step/24)
            step_factor = 1.0 + 0.05 * np.sqrt(step / 24.0)
            ci_margin = round(Z_80_CI * self.residual_std * step_factor, 2)
            lower_bound_80 = max(0.0, round(pred_demand - ci_margin, 2))
            upper_bound_80 = round(pred_demand + ci_margin, 2)

            # การประเมินช่วงสูบน้ำตาม TOR ภาคผนวก ๗ (Peak vs Off-Peak)
            is_peak_pump = bool(hour in [6, 7, 8, 17, 18, 19])

            record = {
                "dma_code": self.dma_code,
                "forecast_timestamp": future_time.strftime("%Y-%m-%d %H:%M:%S"),
                "step_ahead": step,
                "predicted_demand": pred_demand,
                "lower_bound_80": lower_bound_80,
                "upper_bound_80": upper_bound_80,
                "ci_margin": ci_margin,
                "hour": hour,
                "is_peak_period": is_peak_pump,
                "pump_operation_recommended": "INCREASE" if is_peak_pump and pred_demand > r_mean else "STANDARD"
            }
            forecast_records.append(record)
            history_demands.append(pred_demand)

        return forecast_records

    def save_model(self, model_dir):
        os.makedirs(model_dir, exist_ok=True)
        path = os.path.join(model_dir, f"demand_model_{self.dma_code}.joblib")
        joblib.dump({
            "model": self.model,
            "feature_cols": self.feature_cols,
            "residual_std": self.residual_std,
            "metrics": self.metrics,
            "dma_code": self.dma_code
        }, path)
        print(f"[{self.dma_code}] Model saved to: {path}")
        return path
