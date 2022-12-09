package mongodb

import (
	"context"
	"github.com/halilylm/ticketing/payments/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(collection *mongo.Collection) domain.OrderRepository {
	return &orderRepository{collection: collection}
}

func (o *orderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	var foundOrder domain.Order
	oid, _ := primitive.ObjectIDFromHex(id)
	res := o.collection.FindOne(ctx, bson.M{"_id": oid})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundOrder); err != nil {
		return nil, err
	}
	return &foundOrder, nil
}
