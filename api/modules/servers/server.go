package servers

import (
	"api/config"
	_sseUseCase "api/modules/server_send_event/usecase"
	streamSubscribers "api/modules/stream_data/subscribers"
	streamUsecase "api/modules/stream_data/usecase"
	_waterDataRepository "api/modules/water_data/repository"
	_waterDataUseCase "api/modules/water_data/usecases"
	"api/pkg/mqtt"
	"api/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	app        *fiber.App
	db         *gorm.DB
	cfg        *config.Configs
	sseUseCase _sseUseCase.SSEUseCase
}

func NewServers(cfg *config.Configs, db *gorm.DB) *Server {
	return &Server{
		app: fiber.New(fiber.Config{
			AppName: "Water IOT API",
		}),
		cfg:        cfg,
		db:         db,
		sseUseCase: _sseUseCase.NewSSEUseCase(),
	}
}

func (s *Server) Start() {
	s.MapHandler()

	mqttClient, err := mqtt.NewMqttClient(s.cfg.MQTT.BrokerURL)
	if err != nil {
		log.Printf("[Warning] MQTT not connected: %v", err)
	} else {
		defer mqttClient.Disconnect(250)
		// ประกอบ Clean Architecture: Repo -> WaterDataUseCase -> StreamDataUseCase -> Subscriber
		waterDataRepo := _waterDataRepository.NewWaterDataRepository(s.db)
		waterDataUseCase := _waterDataUseCase.NewWaterDataUseCase(waterDataRepo)
		streamDataUseCase := streamUsecase.NewStreamDataUseCase(waterDataUseCase, s.sseUseCase)
		subscriber := streamSubscribers.NewStreamDataSubscriber(mqttClient, streamDataUseCase)
		if err := subscriber.Subscribe("water-iot/sensors/#"); err != nil {
			log.Printf("[Warning] Failed to subscribe to MQTT topic: %v", err)
		}
	}

	fiberUrl, err := utils.ConnectionUrlBuilder("fiber", s.cfg)
	if err != nil {
		log.Fatal("Error Building Fiber Connection URL: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Server is shtting down...")
		if err := s.app.Shutdown(); err != nil {
			log.Fatal("Error during shutdown:", err)
		}
	}()

	log.Printf("Server starting on %s", fiberUrl)
	if err := s.app.Listen(fiberUrl); err != nil {
		log.Fatalf("Server failed to listen: %v", err)
	}

}
