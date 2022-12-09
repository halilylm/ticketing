package domain

import "context"

// Ticket domain
type Ticket struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	Title   string  `json:"title" bson:"title" validate:"required"`
	Price   int     `json:"price" bson:"price" validate:"number"`
	UserID  string  `json:"user_id" bson:"user_id"`
	Version int     `json:"version,omitempty" bson:"version" validate:"number"`
	OrderID *string `json:"order_id" bson:"order_id"`
}

// TicketRepository to interact db
type TicketRepository interface {
	Insert(ctx context.Context, ticket *Ticket) (*Ticket, error)
	Update(ctx context.Context, ticket *Ticket) (*Ticket, error)
	FindByID(ctx context.Context, id string) (*Ticket, error)
	AvailableTickets(ctx context.Context) ([]*Ticket, error)
}
