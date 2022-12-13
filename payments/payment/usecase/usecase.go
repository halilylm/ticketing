package usecase

import (
	"context"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/ticketing/payments/domain"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"go.mongodb.org/mongo-driver/mongo"
)

type payment struct {
	paymentRepo domain.PaymentRepository
	orderRepo   domain.OrderRepository
}

func NewPayment(paymentRepo domain.PaymentRepository, orderRepo domain.OrderRepository) Payment {
	return &payment{paymentRepo: paymentRepo, orderRepo: orderRepo}
}

type Payment interface {
	Pay(ctx context.Context, orderID string) (*domain.Payment, error)
}

func (p *payment) Pay(ctx context.Context, orderID string) (*domain.Payment, error) {
	order, err := p.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	stripe.Key = "sk_test_51FHv9PKXyE07Xs5iFGatWlPOKCxIx6e6mXP5IezkIuFDUdaIyzW95hKDWGabXbzCgF9OHS9rfn93NFEVPkdFOOze00jIw79wFW"
	charged, err := charge.New(&stripe.ChargeParams{
		Amount:      stripe.Int64(int64(order.Charge * 100)),
		Currency:    stripe.String(string(stripe.CurrencyEUR)),
		Description: stripe.String(order.ID),
		Customer:    stripe.String(order.UserID),
		Source:      &stripe.SourceParams{Token: stripe.String("tok_visa")}})
	paid := &domain.Payment{
		OrderID:  order.ID,
		StripeID: charged.ID,
	}
	createdPayment, err := p.paymentRepo.Insert(ctx, paid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return createdPayment, nil
}
