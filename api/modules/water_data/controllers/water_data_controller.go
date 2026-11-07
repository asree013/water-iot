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
