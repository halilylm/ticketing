package mongodb

import (
	"context"

	"github.com/halilylm/ticketing/orders/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ticketRepository struct {
	collection *mongo.Collection
}

func NewTicketRepository(collection *mongo.Collection) domain.TicketRepository {
	return &ticketRepository{collection: collection}
}

func (t *ticketRepository) FindByID(ctx context.Context, id string) (*domain.Ticket, error) {
	var foundTicket domain.Ticket
	res := t.collection.FindOne(ctx, bson.M{"_id": id})
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
	_, err := t.collection.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (t *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	var updatedTicket domain.Ticket

	res := t.collection.FindOneAndUpdate(ctx, bson.M{
		"version": ticket.Version - 1,
		"_id":     ticket.ID,
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
