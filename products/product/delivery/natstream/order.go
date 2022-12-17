package natstream

import (
	"context"
	"encoding/json"
	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/gommon/logger"
	"github.com/halilylm/secondhand/product/product/usecase"
	"log"
	"sync"
	"time"
)

type OrderConsumerGroup struct {
	stream    events.Streaming
	productUC usecase.Product
	groupID   string
	logger    logger.Logger
}

func NewOrderConsumerGroup(
	stream events.Streaming,
	productUC usecase.Product,
	groupID string,
) *OrderConsumerGroup {
	return &OrderConsumerGroup{
		stream:    stream,
		productUC: productUC,
		groupID:   groupID,
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
				foundProduct, err := ocg.productUC.ShowProduct(context.TODO(), deliveredEvent.ProductID)
				if err != nil {
					log.Println(err)
					continue
				}
				foundProduct.OrderID = &deliveredEvent.ID
				updatedProduct, err := ocg.productUC.UpdateProduct(context.TODO(), foundProduct)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
				msg := messages.ProductUpdatedEvent{
					ID:      updatedProduct.ID,
					Version: updatedProduct.Version,
					Title:   updatedProduct.Title,
					Price:   updatedProduct.Price,
					UserID:  updatedProduct.UserID,
				}
				encodedMsg, err := json.Marshal(msg)
				if err != nil {
					ocg.logger.Error(err)
				}
				if err := ocg.stream.Publish(messages.ProductUpdated, encodedMsg); err != nil {
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
				foundProduct, err := ocg.productUC.ShowProduct(context.TODO(), deliveredEvent.ProductID)
				if err != nil {
					log.Println(err)
					continue
				}
				foundProduct.OrderID = nil
				updatedProduct, err := ocg.productUC.UpdateProduct(context.TODO(), foundProduct)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
				msg := messages.ProductUpdatedEvent{
					ID:      updatedProduct.ID,
					Version: updatedProduct.Version,
					Title:   updatedProduct.Title,
					Price:   updatedProduct.Price,
					UserID:  updatedProduct.UserID,
				}
				encodedMsg, err := json.Marshal(msg)
				if err != nil {
					ocg.logger.Error(err)
				}
				if err := ocg.stream.Publish(messages.ProductUpdated, encodedMsg); err != nil {
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
