package controller

import (
	"api/modules/entities"

	"github.com/gofiber/fiber/v2"
)

type ForecastController struct {
	useCase entities.TForecastUseCase
}

func NewForecastController(r fiber.Router, useCase entities.TForecastUseCase) *ForecastController {
	controller := &ForecastController{
		useCase: useCase,
	}
	r.Get("/demand/:dma_code", controller.GetDemandForecast)
	r.Get("/metrics", controller.GetModelMetrics)
	return controller
}

// GetDemandForecast godoc
// @Summary      Get Water Demand Forecast
// @Description  ดึงผลการพยากรณ์ความต้องการใช้น้ำล่วงหน้า 24 ชม. หรือ 7 วัน พร้อมช่วงความเชื่อมั่น 80% (80% Confidence Interval) และคำแนะนำการสูบน้ำ Peak/Off-Peak ตามเกณฑ์ TOR ภาคผนวก ๔
// @Tags         Forecast
// @Produce      json
// @Param        dma_code path      string  true   "รหัสโซน DMA เช่น DMA-SK-01, DMA-PB-01, DMA-KN-01, DMA-BKK-01"
// @Param        range    query     string  false  "ช่วงเวลาพยากรณ์: '24h' (24 ชม. ล่วงหน้า) หรือ '7d' (7 วันล่วงหน้า)" default(24h)
// @Success      200      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /api/v1/forecast/demand/{dma_code} [get]
func (f *ForecastController) GetDemandForecast(c *fiber.Ctx) error {
	dmaCode := c.Params("dma_code")
	timeRange := c.Query("range", "24h")

	data, summary, err := f.useCase.GetDemandForecast(dmaCode, timeRange)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":   fiber.StatusNotFound,
			"message":  "Forecast not found for DMA zone: " + dmaCode,
			"error":    err.Error(),
			"dma_code": dmaCode,
			"range":    timeRange,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   fiber.StatusOK,
		"message":  "Fetch Demand Forecast Successfully",
		"dma_code": dmaCode,
		"range":    timeRange,
		"summary":  summary,
		"data":     data,
	})
}

// GetModelMetrics godoc
// @Summary      Get Model Accuracy & Evaluation Metrics
// @Description  ดึงผลการประเมินความแม่นยำของโมเดล Machine Learning Demand Forecasting (R2, Accuracy %, RMSE, MAE) เทียบกับเกณฑ์ TOR (ความแม่นยำ >= 85%) สามารถระบุ dma_code เพื่อดูเฉพาะโซน หรือเว้นว่างเพื่อดูภาพรวมทุกโซน
// @Tags         Forecast
// @Produce      json
// @Param        dma_code query     string  false  "รหัสโซน DMA เช่น DMA-03, DMA-SK-01 (หากไม่ระบุจะแสดงภาพรวมและแจกแจงทุกโซน)"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/forecast/metrics [get]
func (f *ForecastController) GetModelMetrics(c *fiber.Ctx) error {
	dmaCode := c.Query("dma_code")
	metrics, err := f.useCase.GetModelMetrics(dmaCode)
	if err != nil {
		if dmaCode != "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":   fiber.StatusNotFound,
				"message":  "Metrics not found for DMA zone: " + dmaCode,
				"error":    err.Error(),
				"dma_code": dmaCode,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to load model accuracy summary",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":            fiber.StatusOK,
		"message":           "Fetch Model Accuracy Summary Successfully",
		"tor_standard":      "Accuracy >= 85%, Confidence Interval 80%, 24h / 7d Horizons",
		"model_performance": metrics,
	})
}
