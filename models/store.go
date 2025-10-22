package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
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

func (s *Store) AddOrder(order *Order) {
	s.lock.Lock()
	defer s.lock.Unlock()

	order.ID = uuid.New().String()
	order.SetTotal()
	order.Status = StatusPending
	order.Version = 1

	now := time.Now().UTC()
	order.CreatedAt = now
	order.UpdatedAt = now

	s.Orders[order.ID] = order
}

func (s *Store) DeleteOrder(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	order, ok := s.Orders[id]
	if !ok || order.IsSoftDeleted() {
		return fmt.Errorf("order with ID %s not found", id)
	}

	order.SoftDelete()
	return nil
}

func (s *Store) GetOrder(id string) (*Order, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	order, orderExists := s.Orders[id]

	if orderExists && order.IsSoftDeleted() {
		return nil, false
	}
	return order, orderExists
}

func (s *Store) GetOrders() []*Order {
	s.lock.RLock()
	defer s.lock.RUnlock()

	orders := make([]*Order, 0, len(s.Orders))
	for _, order := range s.Orders {
		if !order.IsSoftDeleted() {
			orders = append(orders, order)
		}
	}
	return orders
}

func (s *Store) UpdateOrderStatus(id string, newStatus OrderStatus) bool {
	order, ok := s.GetOrder(id)
	if !ok {
		log.Printf("Order with ID %s not found\n", id)
		return false
	}

	if !order.CanTransitionTo(newStatus) {
		log.Printf("Invalid status transition from %s to %s for order ID %s\n", order.Status, newStatus, id)
		return false
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	order.Status = newStatus
	order.Version++
	order.UpdatedAt = time.Now().UTC()
	return true
}
