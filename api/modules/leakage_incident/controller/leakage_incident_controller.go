package controller

import (
	"api/modules/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LeakageIncidentController struct {
	incidentUseCase entities.TLeakageIncidentUseCase
}

func NewLeakageIncidentController(r fiber.Router, incidentUseCase entities.TLeakageIncidentUseCase) *LeakageIncidentController {
	controller := &LeakageIncidentController{
		incidentUseCase: incidentUseCase,
	}
	r.Get("/", controller.GetLeakageIncidents)
	r.Get("/:id", controller.GetLeakageIncidentById)
	r.Post("/", controller.CreateLeakageIncident)
	r.Put("/:id", controller.UpdateLeakageIncidentById)
	r.Delete("/:id", controller.DeleteLeakageIncidentById)
	return controller
}

// GetLeakageIncidents godoc
// @Summary      Get Leakage Incidents
// @Description  ดึงรายการประวัติรับแจ้งซ่อมท่อแตกรั่ว (Leakage Incidents) พร้อมรองรับ Pagination
// @Tags         Leakage Incident
// @Produce      json
// @Param        page   query     int  false  "เลขหน้า (Page number)" default(1)
// @Param        limit  query     int  false  "จำนวนรายการต่อหน้า (Items per page)" default(10)
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /api/v1/leakage-incident [get]
func (l *LeakageIncidentController) GetLeakageIncidents(c *fiber.Ctx) error {
	var pagePtr, limitPtr *int
	if page := c.QueryInt("page", 0); page > 0 {
		pagePtr = &page
	}
	if limit := c.QueryInt("limit", 0); limit > 0 {
		limitPtr = &limit
	}

	data, err := l.incidentUseCase.FindDatas(pagePtr, limitPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch Leakage Incident Successfully",
		"data":    data,
	})
}

// GetLeakageIncidentById godoc
// @Summary      Get Leakage Incident By ID
// @Description  ดึงข้อมูลประวัติรับแจ้งซ่อมท่อแตกรั่ว (Leakage Incident) ตาม ID
// @Tags         Leakage Incident
// @Produce      json
// @Param        id   path      int  true  "Leakage Incident ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/leakage-incident/{id} [get]
func (l *LeakageIncidentController) GetLeakageIncidentById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
	}
	data, err := l.incidentUseCase.FindDataById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Fetch Leakage Incident Successfully",
		"data":    data,
	})
}

// CreateLeakageIncident godoc
// @Summary      Create Leakage Incident
// @Description  เพิ่มข้อมูลประวัติรับแจ้งซ่อมท่อแตกรั่ว (Leakage Incident)
// @Tags         Leakage Incident
// @Accept       json
// @Produce      json
// @Param        data  body      entities.TLeakageIncidentReq  true  "Leakage Incident Data Payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/leakage-incident [post]
func (l *LeakageIncidentController) CreateLeakageIncident(c *fiber.Ctx) error {
	req := new(entities.TLeakageIncidentReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": fiber.ErrBadRequest.Message,
			"error":   err.Error(),
		})
	}
	data, err := l.incidentUseCase.AddData(*req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Leakage Incident Created Successfully",
		"data":    data,
	})
}

// UpdateLeakageIncidentById godoc
// @Summary      Update Leakage Incident By ID
// @Description  แก้ไขข้อมูลประวัติรับแจ้งซ่อมท่อแตกรั่ว (Leakage Incident) ตาม ID
// @Tags         Leakage Incident
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Leakage Incident ID"
// @Param        data  body      entities.TLeakageIncidentReq  true  "Leakage Incident Data Payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/leakage-incident/{id} [put]
func (l *LeakageIncidentController) UpdateLeakageIncidentById(c *fiber.Ctx) error {
	req := new(entities.TLeakageIncidentReq)
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
	if _, err := l.incidentUseCase.FindDataById(uint(id)); err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	data, err := l.incidentUseCase.UpdateData(uint(id), *req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Leakage Incident Updated Successfully",
		"data":    data,
	})
}

// DeleteLeakageIncidentById godoc
// @Summary      Delete Leakage Incident By ID
// @Description  ลบข้อมูลประวัติรับแจ้งซ่อมท่อแตกรั่ว (Leakage Incident) ตาม ID
// @Tags         Leakage Incident
// @Produce      json
// @Param        id   path      int  true  "Leakage Incident ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/leakage-incident/{id} [delete]
func (l *LeakageIncidentController) DeleteLeakageIncidentById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
	}
	if _, err := l.incidentUseCase.FindDataById(uint(id)); err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":  fiber.ErrNotFound.Code,
			"message": fiber.ErrNotFound.Message,
			"error":   err.Error(),
		})
	}
	data, err := l.incidentUseCase.DeleteData(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError.Code,
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Leakage Incident Deleted Successfully",
		"data":    data,
	})
}
