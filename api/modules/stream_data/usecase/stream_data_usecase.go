package usecase

import (
	"api/modules/entities"
	"api/modules/server_send_event/usecase"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type StreamDataUseCase interface {
	ConsumeTelemetry(payload []byte) error
}

type streamDataUseCase struct {
	waterDataUseCase entities.WaterDataUseCase
	sseUseCase       usecase.SSEUseCase
}

// Contructer
func NewStreamDataUseCase(waterDataUseCase entities.WaterDataUseCase, sseUseCase usecase.SSEUseCase) StreamDataUseCase {
	return &streamDataUseCase{
		waterDataUseCase: waterDataUseCase,
		sseUseCase:       sseUseCase,
	}
}

// ConsumeTelemetry implements [StreamDataUseCase].
func (s *streamDataUseCase) ConsumeTelemetry(payload []byte) error {
	s.sseUseCase.Broadcast(payload)
	var req entities.WaterDataReq
	if err := json.Unmarshal(payload, &req); err != nil {
		return fmt.Errorf("Invalid json payload: %w", err)
	}
	if checkSaveData := os.Getenv("AUTO_SAVE_BYTE_DATA"); checkSaveData != "true" {
		log.Printf("[StreamData (Skip Save)] Zone=%s, Pipe=%s, Pressure=%.2f, Flow=%.2f",
			req.Zone, req.Pipe, req.Pressure, req.FlowRate)
		return nil
	}
	if err := s.waterDataUseCase.SaveTelemetry(&req); err != nil {
		return fmt.Errorf("failed to save telemetry to db: %w", err)
	}
	log.Printf("[StreamData] Saved: Zone=%s, Pipe=%s, Pressure=%.2f, Flow=%.2f", req.Zone, req.Pressure, req.FlowRate)
	return nil
}
