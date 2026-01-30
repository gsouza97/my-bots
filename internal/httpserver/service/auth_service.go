package service

import (
	"context"
	"fmt"

	"github.com/gsouza97/my-bots/internal/dto"
	"github.com/gsouza97/my-bots/internal/repository"
)

type AuthService struct {
	userRepository repository.UserRepository
	token          string
}

func NewAuthService(userRepository repository.UserRepository, token string) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		token:          token,
	}
}

func (s *AuthService) Authenticate(input dto.LoginInput) (string, error) {
	ctx := context.Background()

	user, err := s.userRepository.FindAdminUser(ctx)
	if err != nil {
		return "", err
	}

	if user.Username != input.Username || user.Password != input.Password {
		return "", fmt.Errorf("Unauthorized")
	}

	return s.token, nil
}
