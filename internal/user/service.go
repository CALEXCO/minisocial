package user

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, username, email string) (*User, error) {
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email")
	}

	u := &User{
		Name: username,
		Mail: email,
	}
	return s.store.Create(ctx, u)
}
func (s *Service) GetById(ctx context.Context, id int) (User, error) {
	return s.store.GetById(ctx, id)
}

func (s *Service) UpdateName(ctx context.Context, username string, id int64) (User, error) {

	u := User{
		ID:   id,
		Name: username,
	}
	return s.store.UpdateName(ctx, u)
}

func (s *Service) GetAll(ctx context.Context) ([]User, error) {

	return s.store.GetAll(ctx)
}

/*
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
*/
