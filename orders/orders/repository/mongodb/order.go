package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/halilylm/gommon/events/common/types"
	"github.com/halilylm/ticketing/orders/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(collection *mongo.Collection) domain.OrderRepository {
	return &orderRepository{collection: collection}
}

func (o *orderRepository) Insert(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	order.ID = uuid.NewString()
	_, err := o.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *orderRepository) IsReserved(ctx context.Context, ticketID string) bool {
	count, _ := o.collection.CountDocuments(ctx, bson.M{"ticket_id": ticketID, "status": bson.M{"$in": []types.OrderStatus{types.Created, types.AwaitingPayment, types.Complete}}})
	return count > 0
}

// FindByID finds a order by its id
func (o *orderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	var foundOrder domain.Order
	res := o.collection.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundOrder); err != nil {
		return nil, err
	}
	return &foundOrder, nil
}

func (o *orderRepository) Delete(ctx context.Context, id string) error {
	res, err := o.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (o *orderRepository) ListUserOrders(ctx context.Context, userID string) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0)
	cur, err := o.collection.Find(ctx, bson.M{"user_id": userID})
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var order domain.Order
		if err := cur.Decode(&order); err != nil {
			continue
		}
		orders = append(orders, &order)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *orderRepository) UpdateStatus(ctx context.Context, id string, status types.OrderStatus) (*domain.Order, error) {
	var updatedOrder domain.Order

	res := o.collection.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{"$set": map[string]any{
		"status": status,
	}}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&updatedOrder); err != nil {
		return nil, err
	}
	return &updatedOrder, nil
}
