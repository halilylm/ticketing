package mongodb

import (
	"context"
	"github.com/halilylm/ticketing/payments/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderRepository struct {
	collection *mongo.Collection
}

func (o *orderRepository) Insert(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if _, err := o.collection.InsertOne(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (o *orderRepository) Update(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	var updatedOrder domain.Order

	// optimistic concurrency control implemented
	// here version subtracted by one because
	// the main data stored in ticket microservice
	// so first update will be occurred in ticket
	// microservice so here we hold copy of ticket
	// therefore we need to subtract by one
	res := o.collection.FindOneAndUpdate(ctx, bson.M{
		"version": order.Version - 1,
		"_id":     order.ID,
	}, bson.M{"$set": map[string]any{
		"version": order.Version,
		"user_id": order.UserID,
		"price":   order.Price,
		"status":  order.Status,
	}}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&updatedOrder); err != nil {
		return nil, err
	}
	return &updatedOrder, nil
}
