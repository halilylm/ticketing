package mongodb

import (
	"context"
	"github.com/halilylm/ticketing/orders/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ticketRepository struct {
	collection *mongo.Collection
}

func NewTicketRepository(collection *mongo.Collection) domain.TicketRepository {
	return &ticketRepository{collection: collection}
}

func (t *ticketRepository) FindByID(ctx context.Context, id string) (*domain.Ticket, error) {
	uid, _ := primitive.ObjectIDFromHex(id)
	var foundTicket domain.Ticket
	res := t.collection.FindOne(ctx, bson.M{"_id": uid})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundTicket); err != nil {
		return nil, err
	}
	return &foundTicket, nil
}

// Insert creates a new ticket in mongodb
func (t *ticketRepository) Insert(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	res, err := t.collection.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}
	uid, _ := res.InsertedID.(primitive.ObjectID)
	ticket.ID = uid.Hex()
	return ticket, nil
}

func (t *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	id, _ := primitive.ObjectIDFromHex(ticket.ID)
	var updatedTicket domain.Ticket

	res := t.collection.FindOneAndUpdate(ctx, bson.M{
		"version": ticket.Version - 1,
		"_id":     id,
	}, bson.M{"$set": map[string]any{
		"title":   ticket.Title,
		"version": ticket.Version,
		"price":   ticket.Price,
	}})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&updatedTicket); err != nil {
		return nil, err
	}
	return &updatedTicket, nil
}
