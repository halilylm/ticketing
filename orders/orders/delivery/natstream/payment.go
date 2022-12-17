package natstream

import (
	"context"
	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/gommon/events/common/types"
	"github.com/halilylm/secondhand/orders/orders/usecase"
	"log"
	"sync"
	"time"
)

type PaymentConsumerGroup struct {
	stream  events.Streaming
	orderUC usecase.Order
	groupID string
}

func NewPaymentConsumerGroup(
	stream events.Streaming,
	orderUC usecase.Order,
	groupID string,
) *PaymentConsumerGroup {
	return &PaymentConsumerGroup{
		stream:  stream,
		orderUC: orderUC,
		groupID: groupID,
	}
}

func (pcg *PaymentConsumerGroup) consumeCreatedPayments(
	workersNum int,
	topic string,
) {
	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("%d started working\n", workerID)
			deliveredEvents, err := pcg.stream.Consume(topic, pcg.groupID, true, time.Minute)
			if err != nil {
				log.Fatal(err)
			}
			for event := range deliveredEvents {
				var deliveredEvent messages.PaymentCreatedEvent
				if err := event.Unmarshal(&deliveredEvent); err != nil {
					log.Println(err)
					continue
				}
				log.Println(deliveredEvent.OrderID)
				updatedOrder, err := pcg.orderUC.UpdateStatus(context.TODO(), deliveredEvent.OrderID, types.Complete)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(updatedOrder.Status)
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}

func (pcg *PaymentConsumerGroup) RunConsumers() {
	go pcg.consumeCreatedPayments(10, messages.PaymentCreated)
}
