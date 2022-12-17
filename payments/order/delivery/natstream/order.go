package natstream

import (
	"context"
	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/gommon/events/common/types"
	"github.com/halilylm/secondhand/payments/domain"
	"github.com/halilylm/secondhand/payments/order/usecase"
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
					Charge:  deliveredEvent.Charge,
					Status:  deliveredEvent.Status,
				}
				createdOrder, err := ocg.orderUC.CreateOrder(context.TODO(), &order)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
				log.Println(createdOrder.ID)
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
				foundOrder, err := ocg.orderUC.FindOrder(context.TODO(), deliveredEvent.ID, deliveredEvent.Version)
				if err != nil {
					log.Println(err)
					continue
				}
				foundOrder.Status = types.Cancelled
				updatedTicket, err := ocg.orderUC.UpdateOrder(context.TODO(), foundOrder)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
				log.Println(updatedTicket.ID + " is updated")
			}
		}(i)
	}
}

func (ocg *OrderConsumerGroup) RunConsumers() {
	go ocg.consumeCreatedOrders(10, messages.OrderCreated)
	go ocg.consumeCancelledOrders(10, messages.OrderCancelled)
}
