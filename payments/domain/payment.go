package domain

import "context"

// Payment domain
type Payment struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	OrderID  string `json:"order_id" bson:"order_id"`
	StripeID string `json:"stripe_id" bson:"stripe_id"`
}
type PaymentRepository interface {
	Insert(ctx context.Context, payment *Payment) (*Payment, error)
}
