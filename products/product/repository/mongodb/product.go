package mongodb

import (
	"context"
	"github.com/google/uuid"
	"github.com/halilylm/secondhand/product/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) domain.ProductRepository {
	return &productRepository{collection: collection}
}

func (p *productRepository) Insert(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	product.ID = uuid.NewString()
	_, err := p.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productRepository) Update(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	var updatedProduct domain.Product
	res := p.collection.FindOneAndUpdate(ctx, bson.M{
		"version": product.Version,
		"_id":     product.ID,
	}, bson.M{"$set": map[string]any{
		"title":    product.Title,
		"version":  product.Version + 1,
		"price":    product.Price,
		"order_id": product.OrderID,
	}}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&updatedProduct); err != nil {
		return nil, err
	}
	return &updatedProduct, nil
}

func (p *productRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var foundProduct domain.Product
	res := p.collection.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundProduct); err != nil {
		return nil, err
	}
	return &foundProduct, nil
}

func (p *productRepository) AvailableProducts(ctx context.Context) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	cur, err := p.collection.Find(ctx, bson.M{"order_id": nil})
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var product domain.Product
		if err := cur.Decode(&product); err != nil {
			continue
		}
		products = append(products, &product)
	}
	return products, nil
}
