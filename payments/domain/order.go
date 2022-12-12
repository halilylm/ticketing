package domain

import (
	"context"
	"github.com/halilylm/gommon/events/common/types"
)

type Order struct {
	ID      string            `json:"id" bson:"_id,omitempty"`
	Version int               `json:"version" bson:"version"`
	UserID  string            `json:"user_id" bson:"user_id"`
	Charge  int               `json:"charge" bson:"charge"`
	Status  types.OrderStatus `json:"status" bson:"status"`
}
type OrderRepository interface {
	Insert(ctx context.Context, order *Order) (*Order, error)
	Update(ctx context.Context, order *Order) (*Order, error)
	FindByIDAndVersion(ctx context.Context, id string, version int) (*Order, error)
}
