package servers

import (
	_sseController "api/modules/server_send_event/controller"
	_waterDataHttp "api/modules/water_data/controllers"
	_waterDataRepository "api/modules/water_data/repository"
	_waterDataUseCase "api/modules/water_data/usecases"

	_dmazoneHttp "api/modules/dma_zone/controller"
	_dmazoneRepository "api/modules/dma_zone/repository"
	_dmazoneUseCase "api/modules/dma_zone/usecase"

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
	// route: /api/v1/water-data
	waterDataGroup := v1.Group("/water-data")
	waterDataRepository := _waterDataRepository.NewWaterDataRepository(s.db)
	waterDataUseCase := _waterDataUseCase.NewWaterDataUseCase(waterDataRepository)
	_waterDataHttp.NewWaterDataController(waterDataGroup, waterDataUseCase)

	// route: /api/v1/dma-zone
	dmazoneGroup := v1.Group("/dma-zone")
	dmazoneRepository := _dmazoneRepository.NewDMAZoneRepository(s.db)
	dmazoneUseCase := _dmazoneUseCase.NewDMAZoneUseCase(dmazoneRepository)
	_dmazoneHttp.NewDMAZoneController(dmazoneGroup, dmazoneUseCase)

	// route: /api/v1/sse-water-data
	sseGroup := v1.Group("/sse-water-data")
	_sseController.NewSSEController(sseGroup, s.sseUseCase)

	s.app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "error, endpoint not found",
			"result":      nil,
		})
	})

	return nil
}
