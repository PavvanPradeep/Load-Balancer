package loadbalancer

import (
	"fmt"
	"sync"
)

type Server struct {
	URL         string
	ActiveConns int
	Mutex       sync.Mutex
}

func (s *Server) IncrementConnections() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.ActiveConns++
	fmt.Printf("Incremented: Server %s active connections: %d\n", s.URL, s.ActiveConns)
}

func (s *Server) DecrementConnections() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.ActiveConns--
	fmt.Printf("Decremented: Server %s active connections: %d\n", s.URL, s.ActiveConns)
}
