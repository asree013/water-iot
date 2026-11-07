package servers

import (
	"api/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	App    *fiber.App
	Ctx    *fiber.Ctx
	Db     *gorm.DB
	Config *config.Configs
}
