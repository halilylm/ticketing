package usecase

import (
	"context"
	"encoding/json"
	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/gommon/logger"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/secondhand/product/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type product struct {
	productRepo domain.ProductRepository
	logger      logger.Logger
	streaming   events.Streaming
}

func NewProduct(productRepo domain.ProductRepository, logger logger.Logger, streaming events.Streaming) Product {
	return &product{productRepo: productRepo, logger: logger, streaming: streaming}
}

func (p *product) NewProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	createdTicket, err := p.productRepo.Insert(ctx, product)
	if err != nil {
		p.logger.Error(err)
		return nil, rest.NewInternalServerError()
	}
	msg := messages.ProductCreatedEvent{
		ID:      createdTicket.ID,
		Version: createdTicket.Version,
		Title:   createdTicket.Title,
		Price:   createdTicket.Price,
		UserID:  createdTicket.UserID,
	}
	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error(err)
	}
	if err := p.streaming.Publish(messages.ProductCreated, encodedMsg); err != nil {
		p.logger.Error(err)
	}
	return createdTicket, nil
}

func (p *product) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	updatedTicket, err := p.productRepo.Update(ctx, product)
	if err != nil {
		p.logger.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	msg := messages.ProductUpdatedEvent{
		ID:      updatedTicket.ID,
		Version: updatedTicket.Version,
		Title:   updatedTicket.Title,
		Price:   updatedTicket.Price,
		UserID:  updatedTicket.UserID,
	}
	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error(err)
	}
	if err := p.streaming.Publish(messages.ProductUpdated, encodedMsg); err != nil {
		p.logger.Error(err)
	}
	return updatedTicket, nil
}

func (p *product) AvailableProducts(ctx context.Context) ([]*domain.Product, error) {
	availableTickets, err := p.productRepo.AvailableProducts(ctx)
	if err != nil {
		p.logger.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return availableTickets, nil
}

func (p *product) ShowProduct(ctx context.Context, id string) (*domain.Product, error) {
	ticket, err := p.productRepo.FindByID(ctx, id)
	if err != nil {
		p.logger.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return ticket, nil
}

// Product contract
type Product interface {
	NewProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	AvailableProducts(ctx context.Context) ([]*domain.Product, error)
	ShowProduct(ctx context.Context, id string) (*domain.Product, error)
}
