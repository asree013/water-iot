package controller

import (
	"api/modules/server_send_event/usecase"
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type SSEController struct {
	sseUseCase usecase.SSEUseCase
}

func NewSSEController(r fiber.Router, sseUseCase usecase.SSEUseCase) *SSEController {
	controller := SSEController{
		sseUseCase: sseUseCase,
	}
	r.Get("/telemetry", controller.StreamTelemetryData)

	return &controller
}

// StreamTelemetry godoc
// @Summary      Stream Real-time Telemetry (SSE)
// @Description  เปิดการเชื่อมต่อ Server-Sent Events เพื่อรับข้อมูลเซ็นเซอร์ IoT แบบสดๆ
// @Tags         SSE
// @Produce      text/event-stream
// @Router       /api/v1/sse-water-data/telemetry [get]
func (s *SSEController) StreamTelemetryData(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	clientChan := make(chan []byte, 10)
	s.sseUseCase.Register(clientChan)

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		defer s.sseUseCase.Unregister(clientChan)

		for msg := range clientChan {
			fmt.Fprintf(w, "data: %s\n\n", msg)

			if err := w.Flush(); err != nil {
				break
			}
		}
	}))
	return nil
}
