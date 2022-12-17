package mongodb

import (
	"context"

	"github.com/halilylm/secondhand/orders/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepository struct {
	collection *mongo.Collection
}

// NewProductRepository returns a new mongo ticket repository
func NewProductRepository(collection *mongo.Collection) domain.ProductRepository {
	return &productRepository{collection: collection}
}

// FindByID finds a product by id
func (p *productRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var foundProduct domain.Product
	res := p.collection.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundProduct); err != nil {
		return nil, err
	}
	return &foundProduct, nil
}

// Insert creates a new ticket in mongodb
func (p *productRepository) Insert(ctx context.Context, ticket *domain.Product) (*domain.Product, error) {
	_, err := p.collection.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// Update updates an existing ticket in mongodb
func (p *productRepository) Update(ctx context.Context, ticket *domain.Product) (*domain.Product, error) {
	var updatedTicket domain.Product

	// optimistic concurrency control implemented
	// here version subtracted by one because
	// the main data stored in ticket microservice
	// so first update will be occurred in ticket
	// microservice so here we hold copy of ticket
	// therefore we need to subtract by one
	res := p.collection.FindOneAndUpdate(ctx, bson.M{
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
