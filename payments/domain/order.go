package domain

import (
	"context"
	"github.com/halilylm/gommon/events/common/types"
)

// Order domain
type Order struct {
	ID      string            `json:"id" bson:"_id,omitempty"`
	Version int               `json:"version" bson:"version"`
	UserID  string            `json:"user_id" bson:"user_id"`
	Price   int               `json:"price" bson:"price"`
	Status  types.OrderStatus `json:"status" bson:"status"`
}

type OrderRepository interface {
	FindByID(ctx context.Context, id string) (*Order, error)
}
