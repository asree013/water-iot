package usecase

import (
	"log"
	"sync"
)

type SSEUseCase interface {
	Register(client chan []byte)
	Unregister(client chan []byte)
	Broadcast(message []byte)
}

type sseUseCase struct {
	clients map[chan []byte]bool
	mu      sync.RWMutex
}

// Broadcast implements [SSEUseCase].
func (s *sseUseCase) Broadcast(message []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		select {
		case client <- message:
		default:
		}
	}
}

// Register implements [SSEUseCase].
func (s *sseUseCase) Register(client chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client] = true
	log.Printf("[SSE Hub] Client connected. Total clients: %d", len(s.clients))
}

// Unregister implements [SSEUseCase].
func (s *sseUseCase) Unregister(client chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
		close(client)
		log.Printf("[SSE Hub] Client disconnected. Total clients: %d", len(s.clients))
	}
}

func NewSSEUseCase() SSEUseCase {
	return &sseUseCase{
		clients: make(map[chan []byte]bool),
	}
}
