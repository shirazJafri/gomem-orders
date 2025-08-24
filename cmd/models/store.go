package models

import (
	"log"
	"sync"
)

type Store struct {
	Orders map[string]*Order
	lock   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		Orders: make(map[string]*Order),
		lock:   sync.RWMutex{},
	}
}

func (s *Store) DeleteOrder(id string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	order, ok := s.GetOrder(id)
	if !ok {
		log.Printf("Order with ID %s not found\n", id)
		return false
	}

	return order.SoftDelete()
}

func (s *Store) GetOrder(id string) (*Order, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	order, ok := s.Orders[id]

	// Soft deletes should not be returned
	if ok && order.DeletedAt != nil {
		return nil, false
	}
	return order, ok
}
