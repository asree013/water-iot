package servers

import (
	_waterDataHttp "api/modules/water_data/controllers"
	_waterDataRepository "api/modules/water_data/repository"
	_waterDataUseCase "api/modules/water_data/usecases"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandler() error {
	v1 := s.App.Group("/api/v1")

	waterDataGroup := v1.Group("/water-data")
	waterDataRepository := _waterDataRepository.NewWaterDataRepository(s.Db)
	waterDataUseCalse := _waterDataUseCase.NewWaterDataUseCase(waterDataRepository)
	_waterDataHttp.NewWaterDataController(waterDataGroup, waterDataUseCalse)

	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "error, end point not found",
			"result":      nil,
		})
	})

	return nil
}
