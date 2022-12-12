package natstream

import (
	"context"
	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/ticketing/payments/domain"
	"github.com/halilylm/ticketing/payments/order/usecase"
	"log"
	"sync"
	"time"
)

type OrderConsumerGroup struct {
	stream  events.Streaming
	orderUC usecase.Order
	groupID string
}

func NewOrderConsumerGroup(
	stream events.Streaming,
	orderUC usecase.Order,
	groupID string,
) *OrderConsumerGroup {
	return &OrderConsumerGroup{
		stream:  stream,
		orderUC: orderUC,
		groupID: groupID,
	}
}

func (ocg *OrderConsumerGroup) consumeCreatedOrders(workersNum int, topic string) {
	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("%d started working\n", workerID)
			deliveredEvents, err := ocg.stream.Consume(topic, ocg.groupID, true, time.Minute)
			if err != nil {
				log.Fatal(err)
			}
			for event := range deliveredEvents {
				var deliveredEvent messages.OrderCreatedEvent
				if err := event.Unmarshal(&deliveredEvent); err != nil {
					log.Println(err)
					continue
				}
				order := domain.Order{
					ID:      deliveredEvent.ID,
					Version: deliveredEvent.Version,
					UserID:  deliveredEvent.UserID,
					Price:   deliveredEvent.Version,
					Status:  deliveredEvent.Status,
				}
				createdTicket, err := ocg.orderUC.CreateOrder(context.TODO(), foundTicket)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
				if err := ocg.stream.Publish(messages.TicketUpdated, updatedTicket.Marshal()); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}

func (ocg *OrderConsumerGroup) consumeCancelledOrders(workersNum int, topic string) {
	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("%d started working\n", workerID)
			deliveredEvents, err := ocg.stream.Consume(topic, ocg.groupID, true, time.Minute)
			if err != nil {
				log.Fatal(err)
			}
			for event := range deliveredEvents {
				var deliveredEvent messages.OrderCancelledEvent
				if err := event.Unmarshal(&deliveredEvent); err != nil {
					log.Println(err)
					continue
				}
				foundTicket, err := ocg.ticketUC.ShowTicket(context.TODO(), deliveredEvent.TicketID)
				if err != nil {
					log.Println(err)
					continue
				}
				foundTicket.OrderID = nil
				updatedTicket, err := ocg.ticketUC.UpdateTicket(context.TODO(), foundTicket)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
				if err := ocg.stream.Publish(messages.TicketUpdated, updatedTicket.Marshal()); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}

func (ocg *OrderConsumerGroup) RunConsumers() {
	go ocg.consumeCreatedOrders(10, messages.OrderCreated)
	go ocg.consumeCancelledOrders(10, messages.OrderCancelled)
}
