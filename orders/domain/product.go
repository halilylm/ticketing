package domain

import (
	"context"
)

// Product domain
type Product struct {
	ID      string `json:"id" bson:"_id,omitempty" validate:"required"`
	Title   string `json:"title" bson:"title"`
	Price   int    `json:"price" bson:"price"`
	Version int    `json:"version" bson:"version"`
}

type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*Product, error)
	Insert(ctx context.Context, product *Product) (*Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
}
