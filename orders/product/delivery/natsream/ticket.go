package natsream

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/halilylm/gommon/events"
	"github.com/halilylm/gommon/events/common/messages"
	"github.com/halilylm/secondhand/orders/domain"
	"github.com/halilylm/secondhand/orders/product/usecase"
)

type ProductConsumerGroup struct {
	stream    events.Streaming
	productUC usecase.Product
	groupID   string
}

func NewProductConsumerGroup(
	stream events.Streaming,
	productUC usecase.Product,
	groupID string,
) *ProductConsumerGroup {
	return &ProductConsumerGroup{
		stream:    stream,
		productUC: productUC,
		groupID:   groupID,
	}
}

func (tcg *ProductConsumerGroup) consumeCreatedProducts(
	workersNum int,
	topic string,
) {
	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("%d started working\n", workerID)
			deliveredEvents, err := tcg.stream.Consume(topic, tcg.groupID, true, time.Minute)
			if err != nil {
				log.Fatal(err)
			}
			for event := range deliveredEvents {
				var deliveredEvent messages.ProductCreatedEvent
				if err := event.Unmarshal(&deliveredEvent); err != nil {
					log.Println(err)
					continue
				}
				product := domain.Product{
					ID:      deliveredEvent.ID,
					Title:   deliveredEvent.Title,
					Price:   deliveredEvent.Price,
					Version: deliveredEvent.Version,
				}
				createdProduct, err := tcg.productUC.CreateProduct(context.TODO(), &product)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(createdProduct.ID)
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}

func (tcg *ProductConsumerGroup) consumeUpdatedProducts(
	workersNum int,
	topic string,
) {
	wg := &sync.WaitGroup{}
	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go func(workerID int) {
			log.Printf("%d started working\n", workerID)
			deliveredEvents, err := tcg.stream.Consume(topic, tcg.groupID, true, time.Minute)
			if err != nil {
				log.Fatal(err)
			}
			for event := range deliveredEvents {
				var deliveredEvent messages.ProductUpdatedEvent
				if err := event.Unmarshal(&deliveredEvent); err != nil {
					log.Println(err)
					continue
				}
				product := domain.Product{
					ID:      deliveredEvent.ID,
					Title:   deliveredEvent.Title,
					Price:   deliveredEvent.Price,
					Version: deliveredEvent.Version,
				}
				_, err := tcg.productUC.UpdateProduct(context.TODO(), &product)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := event.Ack(); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}

func (tcg *ProductConsumerGroup) RunConsumers() {
	go tcg.consumeCreatedProducts(10, messages.ProductCreated)
	go tcg.consumeUpdatedProducts(10, messages.ProductUpdated)
}
