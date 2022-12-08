package mongodb

import (
	"context"

	"github.com/halilylm/ticketing/tickets/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ticketRepository struct {
	collection *mongo.Collection
}

// NewTicketRepository returns a new mongo ticket repository
func NewTicketRepository(collection *mongo.Collection) domain.TicketRepository {
	return &ticketRepository{collection: collection}
}

// Insert creates a new ticket in mongodb
func (t *ticketRepository) Insert(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	res, err := t.collection.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}
	tid, _ := res.InsertedID.(primitive.ObjectID)
	ticket.ID = tid.Hex()
	return ticket, nil
}

// Update updates an existing ticket in mongodb
// optimistic concurrency control implemented here
// so there will be version control before update
// and version increment after update
func (t *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	id, _ := primitive.ObjectIDFromHex(ticket.ID)
	var updatedTicket domain.Ticket

	res := t.collection.FindOneAndUpdate(ctx, bson.M{
		"version": ticket.Version,
		"_id":     id,
	}, bson.M{"$set": map[string]any{
		"title":    ticket.Title,
		"version":  ticket.Version + 1,
		"price":    ticket.Price,
		"order_id": ticket.OrderID,
	}})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&updatedTicket); err != nil {
		return nil, err
	}
	return &updatedTicket, nil
}

// FindByID finds a ticket by its id
func (t *ticketRepository) FindByID(ctx context.Context, id string) (*domain.Ticket, error) {
	tid, _ := primitive.ObjectIDFromHex(id)
	var foundTicket domain.Ticket
	res := t.collection.FindOne(ctx, bson.M{"_id": tid})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundTicket); err != nil {
		return nil, err
	}
	return &foundTicket, nil
}

// AvailableTickets lists all available(unreserved) tickets
func (t *ticketRepository) AvailableTickets(ctx context.Context) ([]*domain.Ticket, error) {
	tickets := make([]*domain.Ticket, 0)
	cur, err := t.collection.Find(ctx, bson.M{"order_id": nil})
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var ticket domain.Ticket
		if err := cur.Decode(&ticket); err != nil {
			continue
		}
		tickets = append(tickets, &ticket)
	}
	return tickets, nil
}
