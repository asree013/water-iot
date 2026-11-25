package controller

import (
	"api/modules/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PipeAssetController struct {
	pipeAssetUseCase entities.TPipeAssetUseCase
}

func NewPipeAssetController(r fiber.Router, pipeAssetUseCase entities.TPipeAssetUseCase) *PipeAssetController {
	controller := &PipeAssetController{
		pipeAssetUseCase: pipeAssetUseCase,
	}
	r.Get("/", controller.GetPipeAssets)
	r.Get("/:id", controller.GetPipeAssetById)
	r.Post("/", controller.CreatePipeAsset)
	r.Put("/:id", controller.UpdatePipeAssetById)
	r.Delete("/:id", controller.DeletePipeAssetById)
	return controller
}

// GetPipeAssets godoc
// @Summary      Get Pipe Assets
// @Description  ดึงรายการสินทรัพย์ท่อประปา (Pipe Assets) พร้อมรองรับ Pagination
// @Tags         Pipe Asset
// @Produce      json
// @Param        page   query     int  false  "เลขหน้า (Page number)" default(1)
// @Param        limit  query     int  false  "จำนวนรายการต่อหน้า (Items per page)" default(10)
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/pipe-asset [get]
func (p *PipeAssetController) GetPipeAssets(c *fiber.Ctx) error {
	var pagePtr, limitPtr *int
	if page := c.QueryInt("page", 0); page > 0 {
		pagePtr = &page
	}
	if limit := c.QueryInt("limit", 0); limit > 0 {
		limitPtr = &limit
	}

	data, err := p.pipeAssetUseCase.FindDatas(pagePtr, limitPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch Pipe Asset Successfully",
		"data":    data,
	})
}

// GetPipeAssetById godoc
// @Summary      Get Pipe Asset By ID
// @Description  ดึงข้อมูลสินทรัพย์ท่อประปา (Pipe Asset) ตาม ID
// @Tags         Pipe Asset
// @Produce      json
// @Param        id   path      int  true  "Pipe Asset ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/pipe-asset/{id} [get]
func (p *PipeAssetController) GetPipeAssetById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
	}
	data, err := p.pipeAssetUseCase.FindDataById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch Pipe Asset Successfully",
		"data":    data,
	})
}

// CreatePipeAsset godoc
// @Summary      Create Pipe Asset
// @Description  เพิ่มข้อมูลสินทรัพย์ท่อประปา (Pipe Asset)
// @Tags         Pipe Asset
// @Accept       json
// @Produce      json
// @Param        data  body      entities.TPipeAssetReq  true  "Pipe Asset Data Payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/pipe-asset [post]
func (p *PipeAssetController) CreatePipeAsset(c *fiber.Ctx) error {
	req := new(entities.TPipeAssetReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": fiber.ErrBadRequest.Message,
			"error":   err.Error(),
		})
	}
	data, err := p.pipeAssetUseCase.AddData(*req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Pipe Asset Created Successfully",
		"data":    data,
	})
}

// UpdatePipeAssetById godoc
// @Summary      Update Pipe Asset By ID
// @Description  แก้ไขข้อมูลสินทรัพย์ท่อประปา (Pipe Asset) ตาม ID
// @Tags         Pipe Asset
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "Pipe Asset ID"
// @Param        data  body      entities.TPipeAssetReq  true  "Pipe Asset Data Payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/pipe-asset/{id} [put]
func (p *PipeAssetController) UpdatePipeAssetById(c *fiber.Ctx) error {
	req := new(entities.TPipeAssetReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": fiber.ErrBadRequest.Message,
			"error":   err.Error(),
		})
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
	}
	if _, err := p.pipeAssetUseCase.FindDataById(uint(id)); err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	data, err := p.pipeAssetUseCase.UpdateData(uint(id), *req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Pipe Asset Updated Successfully",
		"data":    data,
	})
}

// DeletePipeAssetById godoc
// @Summary      Delete Pipe Asset By ID
// @Description  ลบข้อมูลสินทรัพย์ท่อประปา (Pipe Asset) ตาม ID
// @Tags         Pipe Asset
// @Produce      json
// @Param        id   path      int  true  "Pipe Asset ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/pipe-asset/{id} [delete]
func (p *PipeAssetController) DeletePipeAssetById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
	}
	if _, err := p.pipeAssetUseCase.FindDataById(uint(id)); err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	data, err := p.pipeAssetUseCase.DeleteData(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Pipe Asset Deleted Successfully",
		"data":    data,
	})
}
