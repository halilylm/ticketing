package domain

import (
	"context"
	"encoding/json"
	"github.com/halilylm/gommon/events/common/types"
)

// Order domain
type Order struct {
	ID        string            `json:"id" bson:"_id,omitempty"`
	UserID    string            `json:"user_id" bson:"user_id"`
	Status    types.OrderStatus `json:"status" bson:"status"`
	ProductID string            `json:"product_id" bson:"product_id"`
	Version   int               `json:"version,omitempty" bson:"version"`
}

func (o *Order) Marshal() []byte {
	b, _ := json.Marshal(o)
	return b
}

type OrderRepository interface {
	Insert(ctx context.Context, order *Order) (*Order, error)
	IsReserved(ctx context.Context, ticketID string) bool
	FindByID(ctx context.Context, id string) (*Order, error)
	Delete(ctx context.Context, id string) error
	ListUserOrders(ctx context.Context, userID string) ([]*Order, error)
	UpdateStatus(ctx context.Context, id string, status types.OrderStatus) (*Order, error)
}
