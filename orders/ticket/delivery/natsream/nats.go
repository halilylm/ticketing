package natsream

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/ticketing/orders/domain"
	"github.com/halilylm/ticketing/orders/ticket/usecase"
)

type TicketConsumerGroup struct {
	stream   events.Streaming
	ticketUC usecase.Ticket
	groupID  string
}

func NewTicketConsumerGroup(
	stream events.Streaming,
	ticketUC usecase.Ticket,
	groupID string,
) *TicketConsumerGroup {
	return &TicketConsumerGroup{
		stream:   stream,
		ticketUC: ticketUC,
		groupID:  groupID,
	}
}

func (tcg *TicketConsumerGroup) consumeCreatedTickets(
	workersNum int,
	topic string,
) {
	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("%d started working\n", workerID)
			events, err := tcg.stream.Consume(topic, tcg.groupID, true, time.Minute)
			if err != nil {
				log.Fatal(err)
			}
			for event := range events {
				var deliveredEvent messages.TicketCreatedEvent
				if err := event.Unmarshal(&deliveredEvent); err != nil {
					log.Println(err)
					continue
				}
				ticket := domain.Ticket{
					ID:      deliveredEvent.ID,
					Title:   deliveredEvent.Title,
					Price:   deliveredEvent.Price,
					Version: deliveredEvent.Version,
				}
				createdTicket, err := tcg.ticketUC.CreateTicket(context.TODO(), &ticket)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(createdTicket.ID)
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}

func (tcg *TicketConsumerGroup) consumeUpdatedTickets(
	workersNum int,
	topic string,
) {
	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("%d started working\n", workerID)
			events, err := tcg.stream.Consume(topic, tcg.groupID, true, time.Minute)
			if err != nil {
				log.Fatal(err)
			}
			for event := range events {
				var deliveredEvent messages.TicketUpdatedEvent
				if err := event.Unmarshal(&deliveredEvent); err != nil {
					log.Println(err)
					continue
				}
				ticket := domain.Ticket{
					ID:      deliveredEvent.ID,
					Title:   deliveredEvent.Title,
					Price:   deliveredEvent.Price,
					Version: deliveredEvent.Version,
				}
				createdTicket, err := tcg.ticketUC.CreateTicket(context.TODO(), &ticket)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(createdTicket.ID)
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}

func (tcg *TicketConsumerGroup) RunConsumers() {
	go tcg.consumeCreatedTickets(10, messages.TicketCreated)
	go tcg.consumeUpdatedTickets(10, messages.TicketUpdated)
}
