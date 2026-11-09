package servers

import (
	_waterDataHttp "api/modules/water_data/controllers"
	_waterDataRepository "api/modules/water_data/repository"
	_waterDataUseCase "api/modules/water_data/usecases"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (s *Server) MapHandler() error {
	s.app.Use((recover.New()))
	s.app.Use(logger.New())
	s.app.Use(cors.New())

	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "ok",
			"message":     "Water IOT API is Running",
			"status_code": fiber.StatusOK,
		})
	})

	s.app.Get("/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Water IoT API Reference",
			},
			DarkMode: true,
		})
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/html")
		return c.SendString(htmlContent)
	})

	v1 := s.app.Group("/api/v1")

	waterDataGroup := v1.Group("/water-data")
	waterDataRepository := _waterDataRepository.NewWaterDataRepository(s.db)
	waterDataUseCalse := _waterDataUseCase.NewWaterDataUseCase(waterDataRepository)
	_waterDataHttp.NewWaterDataController(waterDataGroup, waterDataUseCalse)

	s.app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "error, end point not found",
			"result":      nil,
		})
	})

	return nil
}
