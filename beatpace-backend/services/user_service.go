package services

import (
	"context"

	"github.com/yimango/beatpace-backend/model"
	"github.com/yimango/beatpace-backend/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *userService) GetUserBySpotifyID(ctx context.Context, spotifyID string) (*model.User, error) {
	return s.userRepo.GetUserBySpotifyID(ctx, spotifyID)
} 