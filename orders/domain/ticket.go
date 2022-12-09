package domain

import (
	"context"
)

// Ticket domain
type Ticket struct {
	ID      string `json:"id" bson:"_id,omitempty" validate:"required"`
	Title   string `json:"title" bson:"title"`
	Price   int    `json:"price" bson:"price"`
	Version int    `json:"version" bson:"version"`
}

type TicketRepository interface {
	FindByID(ctx context.Context, id string) (*Ticket, error)
	Insert(ctx context.Context, ticket *Ticket) (*Ticket, error)
	Update(ctx context.Context, ticket *Ticket) (*Ticket, error)
}
