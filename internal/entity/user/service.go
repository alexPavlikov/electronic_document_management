package user

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/electronic_document_management/internal/config"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
	us, err := s.repository.SelectUser(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("%s - failed to get user", config.LOG_ERROR)
	}
	return us, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	users, err := s.repository.SelectUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - failed to get all users", config.LOG_ERROR)
	}
	return users, nil
}
