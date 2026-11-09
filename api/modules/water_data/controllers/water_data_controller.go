package controllers

import (
	"api/modules/entities"

	"github.com/gofiber/fiber/v2"
)

type WaterDataController struct {
	waterDataUseCase entities.WaterDataUseCase
}

func NewWaterDataController(r fiber.Router, waterDataUseCase entities.WaterDataUseCase) *WaterDataController {
	controller := &WaterDataController{
		waterDataUseCase: waterDataUseCase,
	}
	r.Post("/telemetry", controller.SaveTelemetry)
	r.Get("/telemetry", controller.GetLatestTelemetry)

	return controller

}

// SaveTelemetry godoc
// @Summary      Save Telemetry Data
// @Description  บันทึกข้อมูล Telemetry จากเซ็นเซอร์ IoT
// @Tags         Water Data
// @Accept       json
// @Produce      json
// @Param        data body entities.WaterDataReq true "Telemetry Data Payload"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/water-data/telemetry [post]
func (c *WaterDataController) SaveTelemetry(ctx *fiber.Ctx) error {
	req := new(entities.WaterDataReq)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":     fiber.ErrBadRequest.Message,
			"error":       err.Error(),
			"status_code": fiber.ErrBadRequest.Code,
		})
	}
	if err := c.waterDataUseCase.SaveTelemetry(req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":     fiber.ErrInternalServerError.Message,
			"error":       err.Error(),
			"status_code": fiber.ErrInternalServerError.Code,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Telemetry data saved Successfully",
		"status_code": fiber.StatusCreated,
	})
}

// GetLatestTelemetry godoc
// @Summary      Get Latest Telemetry Data
// @Description  ดึงข้อมูล Telemetry ทั้งหมดล่าสุด
// @Tags         Water Data
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/water-data/telemetry [get]
func (c *WaterDataController) GetLatestTelemetry(ctx *fiber.Ctx) error {
	data, err := c.waterDataUseCase.GetLatestTelemetry()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":     fiber.ErrInternalServerError.Message,
			"error":       err.Error(),
			"status_code": fiber.ErrInternalServerError.Code,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Fetch Data Success",
		"status_code": fiber.StatusOK,
		"data":        data,
	})
}
