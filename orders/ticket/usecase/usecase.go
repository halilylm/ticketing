package usecase

import (
	"context"

	"github.com/halilylm/gommon/logger"
	"github.com/halilylm/ticketing/orders/domain"
)

type ticket struct {
	ticketRepo domain.TicketRepository
	logger     logger.Logger
}

func NewTicket(ticketRepo domain.TicketRepository, logger logger.Logger) Ticket {
	return &ticket{ticketRepo: ticketRepo, logger: logger}
}

type Ticket interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
	UpdateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
}

func (t *ticket) CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	createdTicket, err := t.ticketRepo.Insert(ctx, ticket)
	if err != nil {
		t.logger.Error(err)
	}
	return createdTicket, nil
}

func (t *ticket) UpdateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	updatedTicket, err := t.ticketRepo.Update(ctx, ticket)
	if err != nil {
		t.logger.Error(err)
	}
	return updatedTicket, nil
}
