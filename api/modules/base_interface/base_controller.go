package base_interface

import "github.com/gofiber/fiber/v2"

type IBseController[T any] interface {
	GetDataController(c *fiber.Ctx) ([]T, error)
	GetDataByIdController(c *fiber.Ctx) (T, error)
	AddDataController(c *fiber.Ctx) (T, error)
	UpdateDataController(c *fiber.Ctx) (T, error)
	DeleteDataController(c *fiber.Ctx) (T, error)
}
