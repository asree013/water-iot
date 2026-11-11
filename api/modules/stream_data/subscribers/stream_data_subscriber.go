package subscribers

import (
	"api/modules/stream_data/usecase"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StreamDataSubscriber struct {
	client  mqtt.Client
	useCase usecase.StreamDataUseCase
}

func NewStreamDataSubscriber(client mqtt.Client, useCase usecase.StreamDataUseCase) *StreamDataSubscriber {
	return &StreamDataSubscriber{
		client:  client,
		useCase: useCase,
	}
}

func (s *StreamDataSubscriber) Subscribe(topic string) error {
	if token := s.client.Subscribe(topic, 0, s.handlerMessage); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("[StreamData] Subscribed to topic: %s", topic)
	return nil
}

func (s *StreamDataSubscriber) handlerMessage(client mqtt.Client, message mqtt.Message) {
	if err := s.useCase.ConsumeTelemetry(message.Payload()); err != nil {
		log.Printf("[StreamData] Failed to consume telemetry: %v", err)
	}
}
