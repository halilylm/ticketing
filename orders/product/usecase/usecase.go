package usecase

import (
	"context"

	"github.com/halilylm/gommon/logger"
	"github.com/halilylm/secondhand/orders/domain"
)

type product struct {
	productRepo domain.ProductRepository
	logger      logger.Logger
}

func NewProduct(ticketRepo domain.ProductRepository, logger logger.Logger) Product {
	return &product{productRepo: ticketRepo, logger: logger}
}

type Product interface {
	CreateProduct(ctx context.Context, ticket *domain.Product) (*domain.Product, error)
	UpdateProduct(ctx context.Context, ticket *domain.Product) (*domain.Product, error)
}

func (p *product) CreateProduct(ctx context.Context, ticket *domain.Product) (*domain.Product, error) {
	createdTicket, err := p.productRepo.Insert(ctx, ticket)
	if err != nil {
		p.logger.Error(err)
	}
	return createdTicket, nil
}

func (p *product) UpdateProduct(ctx context.Context, ticket *domain.Product) (*domain.Product, error) {
	updatedTicket, err := p.productRepo.Update(ctx, ticket)
	if err != nil {
		p.logger.Error(err)
	}
	return updatedTicket, nil
}
