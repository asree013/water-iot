package controller

import (
	"api/modules/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DMAZoneController struct {
	dmaZoneUseCase entities.TDMAZoneUseCase
}

func NewDMAZoneController(r fiber.Router, dmaZoneUseCase entities.TDMAZoneUseCase) *DMAZoneController {
	controller := &DMAZoneController{
		dmaZoneUseCase: dmaZoneUseCase,
	}
	r.Get("/", controller.GetDMAZone)
	r.Get("/:id", controller.GetDMAZoneById)
	r.Post("/", controller.CreateDMAZone)
	r.Put("/:id", controller.UpdateDMAZoneById)
	r.Delete("/:id", controller.DeleteDMAZoneById)
	return controller
}

// GetDMAZone godoc
// @Summary      Get DMA Zones
// @Description  ดึงรายการพื้นที่จ่ายน้ำย่อย (DMA Zones) พร้อมรองรับ Pagination
// @Tags         DMA Zone
// @Produce      json
// @Param        page   query     int  false  "เลขหน้า (Page number)" default(1)
// @Param        limit  query     int  false  "จำนวนรายการต่อหน้า (Items per page)" default(10)
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/dma-zone [get]
func (d *DMAZoneController) GetDMAZone(c *fiber.Ctx) error {
	var pagePtr, limitPtr *int
	if page := c.QueryInt("page", 0); page > 0 {
		pagePtr = &page
	}
	if limit := c.QueryInt("limit", 0); limit > 0 {
		limitPtr = &limit
	}

	data, err := d.dmaZoneUseCase.FindDatas(pagePtr, limitPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch DMA Zone Successfully",
		"data":    data,
	})
}

// GetDMAZoneById godoc
// @Summary      Get DMA Zone By ID
// @Description  ดึงรายการพื้นที่จ่ายน้ำย่อย (DMA Zones) โดยเลือกข้อมูลจาก ID
// @Tags         DMA Zone
// @Produce      json
// @Param        id   path      int  true  "DMA Zone ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/dma-zone/{id} [get]
func (d *DMAZoneController) GetDMAZoneById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	data, err := d.dmaZoneUseCase.FindDataById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch DMA Zone Successfully",
		"data":    data,
	})
}

// CreateDMAZone godoc
// @Summary      Create DMA Zone
// @Description  เพิ่มข้อมูล DMA Zone
// @Tags         DMA Zone
// @Accept       json
// @Produce      json
// @Param        data  body      entities.TDMAZone  true  "DMA Zone Data Payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/dma-zone [post]
func (d *DMAZoneController) CreateDMAZone(c *fiber.Ctx) error {
	req := new(entities.TDMAZone)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": fiber.ErrBadRequest.Message,
			"error":   err.Error(),
		})
	}
	if _, err := d.dmaZoneUseCase.AddData(*req); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "DMA Zone Created Successfully",
		"data":    req,
	})
}

// UpdateDMAZoneById godoc
// @Summary      Update DMA Zone By ID
// @Description  แก้ไขข้อมูลพื้นที่จ่ายน้ำย่อย (DMA Zone) ตาม ID
// @Tags         DMA Zone
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "DMA Zone ID"
// @Param        data  body      entities.TDMAZone  true  "DMA Zone Data Payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/dma-zone/{id} [put]
func (d *DMAZoneController) UpdateDMAZoneById(c *fiber.Ctx) error {
	req := new(entities.TDMAZone)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": fiber.ErrBadRequest.Message,
			"error":   err.Error(),
		})
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": "id is required",
			"error":   err.Error(),
		})
	}
	if _, err := d.dmaZoneUseCase.FindDataById(uint(id)); err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	data, err := d.dmaZoneUseCase.UpdateData(uint(id), *req)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "DMA Zone Updated Successfully",
		"data":    data,
	})

}

// DeleteDMAZoneById godoc
// @Summary      Delete DMA Zone By ID
// @Description  ลบข้อมูลพื้นที่จ่ายน้ำย่อย (DMA Zone) ตาม ID
// @Tags         DMA Zone
// @Produce      json
// @Param        id   path      int  true  "DMA Zone ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/dma-zone/{id} [delete]
func (d *DMAZoneController) DeleteDMAZoneById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": fiber.ErrBadRequest.Message,
			"error":   err.Error(),
		})
	}
	data, err := d.dmaZoneUseCase.FindDataById(uint(id))
	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	if err := d.dmaZoneUseCase.DeleteData(uint(id)); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "DMA Zone Deleted Successfully",
		"data":    data,
	})
}
