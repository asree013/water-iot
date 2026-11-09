package servers

import (
	"api/config"
	"api/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	app *fiber.App
	db  *gorm.DB
	cfg *config.Configs
}

func NewServers(cfg *config.Configs, db *gorm.DB) *Server {
	return &Server{
		app: fiber.New(fiber.Config{
			AppName: "Water IOT API",
		}),
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Start() {
	s.MapHandler()

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
