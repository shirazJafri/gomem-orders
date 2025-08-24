package models

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusDelivered OrderStatus = "DELIVERED"
	StatusCanceled  OrderStatus = "CANCELED"
)

var StatusTransitions = map[OrderStatus]map[OrderStatus]bool{
	StatusPending:   {StatusConfirmed: true, StatusCanceled: true},
	StatusConfirmed: {StatusShipped: true, StatusCanceled: true},
	StatusShipped:   {StatusDelivered: true},
	StatusDelivered: {},
	StatusCanceled:  {},
}

type OrderLine struct {
	ProductID      string `json:"product_id"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
	LineTotalCents int64  `json:"line_total_cents"`
}

func (orderLine *OrderLine) SetLineTotal() {
	orderLine.LineTotalCents = int64(orderLine.Quantity) * orderLine.UnitPriceCents
}

type Order struct {
	ID         string            `json:"id"`
	CustomerID string            `json:"customer_id"`
	Currency   string            `json:"currency"`
	Lines      []OrderLine       `json:"lines"`
	Attributes map[string]string `json:"attributes,omitempty"`
	TotalCents int64             `json:"total_cents"`
	Status     OrderStatus       `json:"status"`
	Version    int64             `json:"version"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  *time.Time        `json:"deleted_at,omitempty"`
}

type CreateOrderResponse struct {
	ID         string      `json:"id"`
	Status     OrderStatus `json:"status"`
	TotalCents int64       `json:"total_cents"`
}

func (order *Order) SetTotal() {
	order.TotalCents = 0
	for i := range order.Lines {
		order.Lines[i].SetLineTotal()
		order.TotalCents += order.Lines[i].LineTotalCents
	}
}

func (order *Order) CanTransitionTo(newStatus OrderStatus) bool {
	allowed, ok := StatusTransitions[order.Status][newStatus]
	return allowed && ok
}

func (order *Order) SoftDelete() bool {
	if order.DeletedAt != nil {
		return false
	}

	now := time.Now().UTC()
	order.DeletedAt = &now
	order.UpdatedAt = now
	order.Version++
	return true
}
