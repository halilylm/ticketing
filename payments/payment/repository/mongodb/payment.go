package mongodb

import (
	"context"
	"github.com/halilylm/ticketing/payments/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) domain.PaymentRepository {
	return &paymentRepository{collection: collection}
}

func (p *paymentRepository) Insert(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	res, err := p.collection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}
	tid, _ := res.InsertedID.(primitive.ObjectID)
	payment.ID = tid.Hex()
	return payment, nil
}
