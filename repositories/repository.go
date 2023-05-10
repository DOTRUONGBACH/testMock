package repositories

import (
	"context"
	"mock-project/services/customer-service/ent"
	"mock-project/services/customer-service/ent/user"
)

// UserRepository
type UserRepository struct {
	client *ent.Client
}

func (r *UserRepository) FindByEmail(email string) (*ent.User, error) {
	user, err := r.client.User.Query().Where(user.Email(email)).Only(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(name, email, password string) (*ent.User, error) {
	user, err := r.client.User.Create().
		SetEmail(email).
		SetPassword(password).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}
