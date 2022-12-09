package usecase

import (
	"context"
	"github.com/halilylm/gommon/events/common/types"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/ticketing/payments/domain"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type payment struct {
	paymentRepo domain.PaymentRepository
	orderRepo   domain.OrderRepository
}

func NewPayment(paymentRepo domain.PaymentRepository, orderRepo domain.OrderRepository) Payment {
	return &payment{paymentRepo: paymentRepo, orderRepo: orderRepo}
}

func (p *payment) NewPayment(ctx context.Context, orderID, userID string, token string) (*domain.Payment, error) {
	// find the order
	order, err := p.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}

	// check permission
	if err := p.havePermission(ctx, order, userID); err != nil {
		return nil, err
	}

	// check if order can be bought
	if order.Status == types.Cancelled {
		return nil, rest.NewBadRequestError(rest.ErrCantBeBoughtCancelledOrder.Error())
	}

	// payment
	stripe.Key = os.Getenv("STRIPE_KEY")
	orderCharge, err := charge.New(&stripe.ChargeParams{
		Amount:      stripe.Int64(int64(order.Price)),
		Currency:    stripe.String(string(stripe.CurrencyEUR)),
		Description: stripe.String(order.ID),
		Source:      &stripe.SourceParams{Token: stripe.String(token)},
	})

	pay := &domain.Payment{
		OrderID:  order.ID,
		StripeID: orderCharge.ID,
	}

	paymentProcess, err := p.paymentRepo.Insert(ctx, pay)
	if err != nil {
		return nil, rest.NewInternalServerError()
	}
	return paymentProcess, nil
}

// havePermission check if user is authorized
// to do this action
func (p *payment) havePermission(ctx context.Context, order *domain.Order, userID string) error {
	if order.UserID != userID {
		return rest.NewUnauthorizedError()
	}
	return nil
}

type Payment interface {
	NewPayment(ctx context.Context, orderID, userID string, token string) (*domain.Payment, error)
}
