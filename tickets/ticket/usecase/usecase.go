package usecase

import (
	"context"
	"github.com/halilylm/gommon/logger"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/ticketing/tickets/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type ticket struct {
	ticketRepo domain.TicketRepository
	logger     logger.Logger
}

func NewTicket(ticketRepo domain.TicketRepository, logger logger.Logger) Ticket {
	return &ticket{ticketRepo: ticketRepo, logger: logger}
}

func (t *ticket) NewTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	createdTicket, err := t.ticketRepo.Insert(ctx, ticket)
	if err != nil {
		t.logger.Error(err)
		return nil, rest.NewInternalServerError()
	}
	return createdTicket, nil
}

func (t *ticket) UpdateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	updatedTicket, err := t.ticketRepo.Update(ctx, ticket)
	if err != nil {
		t.logger.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return updatedTicket, nil
}

func (t *ticket) AvailableTickets(ctx context.Context) ([]*domain.Ticket, error) {
	availableTickets, err := t.ticketRepo.AvailableTickets(ctx)
	if err != nil {
		t.logger.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return availableTickets, nil
}

func (t *ticket) ShowTicket(ctx context.Context, id string) (*domain.Ticket, error) {
	ticket, err := t.ticketRepo.FindByID(ctx, id)
	if err != nil {
		t.logger.Error(err)
		if err == mongo.ErrNoDocuments {
			return nil, rest.NewNotFoundError()
		}
		return nil, rest.NewInternalServerError()
	}
	return ticket, nil
}

// Ticket contract
type Ticket interface {
	NewTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
	UpdateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
	AvailableTickets(ctx context.Context) ([]*domain.Ticket, error)
	ShowTicket(ctx context.Context, id string) (*domain.Ticket, error)
}
