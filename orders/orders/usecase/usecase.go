package usecase

import (
	"context"
	"encoding/json"
	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/gommon/events/common/types"
	"github.com/halilylm/gommon/logger"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/secondhand/orders/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type order struct {
	ticketRepo domain.TicketRepository
	orderRepo  domain.OrderRepository
	logger     logger.Logger
	stream     events.Streaming
}

func NewOrder(ticketRepo domain.TicketRepository, orderRepo domain.OrderRepository, logger logger.Logger, stream events.Streaming) Order {
	return &order{ticketRepo: ticketRepo, orderRepo: orderRepo, logger: logger, stream: stream}
}

func (o *order) NewOrder(ctx context.Context, ticketID, userID string) (*domain.Order, error) {
	// find the ticket
	ticket, err := o.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}

	// check if this ticket already is reserved
	if isReserved := o.orderRepo.IsReserved(ctx, ticket.ID); isReserved {
		return nil, rest.NewBadRequestError(rest.ErrTicketAlreadyReserved.Error())
	}

	// generate the order
	order := &domain.Order{
		UserID:   userID,
		Status:   types.Created,
		TicketID: ticketID,
		Version:  0,
	}

	createdOrder, err := o.orderRepo.Insert(ctx, order)
	if err != nil {
		return nil, rest.NewInternalServerError()
	}
	msg := messages.OrderCreatedEvent{
		ID:       createdOrder.ID,
		Version:  createdOrder.Version,
		Status:   createdOrder.Status,
		UserID:   createdOrder.UserID,
		TicketID: ticketID,
		Charge:   ticket.Price,
	}
	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		o.logger.Error(err)
	}
	if err := o.stream.Publish(messages.OrderCreated, encodedMsg); err != nil {
		return nil, rest.NewInternalServerError()
	}

	return createdOrder, nil
}

func (o *order) ShowOrder(ctx context.Context, id, userID string) (*domain.Order, error) {
	// check permission
	if err := o.havePermission(ctx, id, userID); err != nil {
		return nil, err
	}

	order, err := o.orderRepo.FindByID(ctx, id)
	if err != nil {
		o.logger.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return order, nil
}

func (o *order) DeleteOrder(ctx context.Context, id, userID string) error {
	// find the ticket
	foundOrder, err := o.orderRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest.NewNotFoundError()
		}
		return rest.NewInternalServerError()
	}
	// check permission
	if err := o.havePermission(ctx, id, userID); err != nil {
		return err
	}
	if err := o.orderRepo.Delete(ctx, id); err != nil {
		if err == mongo.ErrNoDocuments {
			return rest.NewNotFoundError()
		}
		return rest.NewInternalServerError()
	}
	msg := messages.OrderCancelledEvent{
		ID:       foundOrder.ID,
		Version:  foundOrder.Version,
		TicketID: foundOrder.TicketID,
	}
	b, _ := json.Marshal(msg)
	if err := o.stream.Publish(messages.OrderCancelled, b); err != nil {
		log.Println(err)
	}
	return nil
}

func (o *order) ListUserOrders(ctx context.Context, userID string) ([]*domain.Order, error) {
	orders, err := o.orderRepo.ListUserOrders(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return orders, nil
}

func (o *order) UpdateStatus(ctx context.Context, id string, status types.OrderStatus) (*domain.Order, error) {
	updatedOrder, err := o.orderRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return updatedOrder, nil
}

// havePermission check if user is authorized
// to do this action
func (o *order) havePermission(ctx context.Context, orderID, userID string) error {
	order, err := o.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest.NewNotFoundError()
		}
		return rest.NewInternalServerError()
	}
	if order.UserID != userID {
		return rest.NewUnauthorizedError()
	}
	return nil
}

type Order interface {
	NewOrder(ctx context.Context, ticketID, userID string) (*domain.Order, error)
	ShowOrder(ctx context.Context, id, userID string) (*domain.Order, error)
	DeleteOrder(ctx context.Context, id, userID string) error
	ListUserOrders(ctx context.Context, userID string) ([]*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status types.OrderStatus) (*domain.Order, error)
}
