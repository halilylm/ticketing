package mongodb

import (
	"context"
	"github.com/google/uuid"
	"github.com/halilylm/ticketing/payments/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) domain.PaymentRepository {
	return &paymentRepository{collection: collection}
}

func (p *paymentRepository) Insert(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	payment.ID = uuid.NewString()
	if _, err := p.collection.InsertOne(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}
