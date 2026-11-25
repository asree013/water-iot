package controller

import (
	"api/modules/entities"

	"github.com/gofiber/fiber/v2"
)

type PumpController struct {
	useCase entities.TPumpUseCase
}

func NewPumpController(r fiber.Router, useCase entities.TPumpUseCase) *PumpController {
	controller := &PumpController{
		useCase: useCase,
	}
	r.Get("/health", controller.GetFleetHealth)
	r.Get("/health/:pump_id", controller.GetPumpHealthByID)
	r.Get("/summary", controller.GetFleetSummary)
	return controller
}

// GetFleetHealth godoc
// @Summary      Get Fleet Pump Health Scores & PTTF
// @Description  ดึงคะแนนสุขภาพ (Health Score 0-100%) เวลาคาดการณ์ก่อนชำรุด (PTTF) และระดับความเสี่ยงของเครื่องสูบน้ำทุกตัวในระบบตามเกณฑ์ ISO 10816-3 และ Multi-sensor Degradation
// @Tags         PumpHealth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/pump/health [get]
func (p *PumpController) GetFleetHealth(c *fiber.Ctx) error {
	scores, summary, err := p.useCase.GetFleetHealth()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to fetch fleet pump health",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch Fleet Pump Health Successfully",
		"summary": summary,
		"data":    scores,
	})
}

// GetPumpHealthByID godoc
// @Summary      Get Single Pump Health & Telemetry Detail
// @Description  ดึงข้อมูลสุขภาพ ค่าเซนเซอร์ SCADA (Vibration, Bearing/Motor Temp, Efficiency, Operating Hours) และ PTTF เจาะจงตัวเครื่องสูบน้ำ
// @Tags         PumpHealth
// @Produce      json
// @Param        pump_id path      string  true  "รหัสเครื่องสูบน้ำ เช่น PUMP-SK-01, PUMP-BKK-02, PUMP-KN-01"
// @Success      200     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Router       /api/v1/pump/health/{pump_id} [get]
func (p *PumpController) GetPumpHealthByID(c *fiber.Ctx) error {
	pumpID := c.Params("pump_id")
	score, err := p.useCase.GetPumpHealthByID(pumpID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":   fiber.StatusNotFound,
			"message":  "Pump health record not found for: " + pumpID,
			"error":    err.Error(),
			"pump_id":  pumpID,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch Pump Health Successfully",
		"data":    score,
	})
}

// GetFleetSummary godoc
// @Summary      Get Pump Fleet Overview Summary
// @Description  ดึงสถิติภาพรวมสถานะเครื่องสูบน้ำทั้งระบบ (จำนวนปั๊มปกติ/เตือน/วิกฤต, คะแนนเฉลี่ย) สำหรับ Executive Dashboard
// @Tags         PumpHealth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/pump/summary [get]
func (p *PumpController) GetFleetSummary(c *fiber.Ctx) error {
	summary, err := p.useCase.GetFleetSummary()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to fetch fleet summary",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch Pump Fleet Summary Successfully",
		"summary": summary,
	})
}
