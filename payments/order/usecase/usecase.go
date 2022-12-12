package usecase

import (
	"context"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/ticketing/payments/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type order struct {
	orderRepo domain.OrderRepository
}

type Order interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	FindOrder(ctx context.Context, id string, version int) (*domain.Order, error)
}

func (o *order) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	createdOrder, err := o.orderRepo.Insert(ctx, order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return createdOrder, nil
}

func (o *order) UpdateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	updatedOrder, err := o.orderRepo.Update(ctx, order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return updatedOrder, nil
}

func (o *order) FindOrder(ctx context.Context, id string, version int) (*domain.Order, error) {
	foundOrder, err := o.orderRepo.FindByIDAndVersion(ctx, id, version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return foundOrder, nil
}
