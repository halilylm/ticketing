package domain

import "context"

// User domain
type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	FirstName string `json:"first_name" bson:"first_name" validate:"required"`
	Surname   string `json:"surname" bson:"surname" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" bson:"password" validate:"required"`
}

func (u *User) HidePassword() {
	u.Password = ""
}

// UserRepository to interact db
type UserRepository interface {
	Insert(ctx context.Context, user *User) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
