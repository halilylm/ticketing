package domain

import (
	"context"
)

// Product domain
type Product struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	Title   string  `json:"title" bson:"title" validate:"required"`
	Price   int     `json:"price" bson:"price" validate:"number"`
	UserID  string  `json:"user_id" bson:"user_id"`
	Version int     `json:"version,omitempty" bson:"version" validate:"number"`
	OrderID *string `json:"order_id" bson:"order_id"`
}

// ProductRepository to interact db
type ProductRepository interface {
	Insert(ctx context.Context, product *Product) (*Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	FindByID(ctx context.Context, id string) (*Product, error)
	AvailableProducts(ctx context.Context) ([]*Product, error)
}
