package usecase

import (
	"api/modules/entities"
	"encoding/json"
	"fmt"
	"log"
)

type StreamDataUseCase interface {
	ConsumeTelemetry(payload []byte) error
}

type streamDataUseCase struct {
	waterDataUseCase entities.WaterDataUseCase
}

// Contructer
func NewStreamDataUseCase(waterDataUseCase entities.WaterDataUseCase) StreamDataUseCase {
	return &streamDataUseCase{
		waterDataUseCase: waterDataUseCase,
	}
}

// ConsumeTelemetry implements [StreamDataUseCase].
func (s *streamDataUseCase) ConsumeTelemetry(payload []byte) error {
	var req entities.WaterDataReq
	if err := json.Unmarshal(payload, &req); err != nil {
		return fmt.Errorf("Invalid json payload: %w", err)
	}

	if err := s.waterDataUseCase.SaveTelemetry(&req); err != nil {
		return fmt.Errorf("failed to save telemetry to db: %w", err)
	}
	log.Printf("[StreamData] Saved: Zone=%s, Pipe=%s, Pressure=%.2f, Flow=%.2f", req.Zone, req.Pressure, req.FlowRate)
	return nil
}
