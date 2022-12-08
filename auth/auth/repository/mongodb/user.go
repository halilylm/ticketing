package mongodb

import (
	"context"

	"github.com/halilylm/ticketing/auth/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository returns a new mongo user repository
func NewUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &userRepository{collection}
}

// Insert creates a new user in mongodb
func (u *userRepository) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	res, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	uid, _ := res.InsertedID.(primitive.ObjectID)
	user.ID = uid.Hex()
	return user, nil
}

// FindByID finds a user by its id
func (u *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	uid, _ := primitive.ObjectIDFromHex(id)
	var foundUser domain.User
	res := u.collection.FindOne(ctx, bson.M{"_id": uid})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundUser); err != nil {
		return nil, err
	}
	return &foundUser, nil
}

// FindByEmail finds a user by its email
func (u *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var foundUser domain.User
	res := u.collection.FindOne(ctx, bson.M{"email": email})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&foundUser); err != nil {
		return nil, err
	}
	return &foundUser, nil
}
