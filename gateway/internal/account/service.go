package account

import (
	"context"
	"fmt"

	"gateway/internal/mapper"
	"gateway/internal/model"
	accountpb "gateway/pkg/account/go"
	pagination "gateway/pkg/pagination/go"
)

type Service struct {
	client accountpb.AccountClient
}

func New(client accountpb.AccountClient) *Service {
	return &Service{client: client}
}

func (s *Service) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	user, err := s.client.GetUser(ctx, &accountpb.GetUserRequest{UserId: userID})
	if err != nil {
		return model.User{}, fmt.Errorf("failed get user: %w", err)
	}

	return mapper.PbToUser(user.User), nil
}

func (s *Service) GetUsers(ctx context.Context, limit uint32, offset uint32) ([]model.User, error) {
	users, err := s.client.GetUsers(ctx, &accountpb.GetUsersRequest{
		Pagination: &pagination.Pagination{
			Limit:  uint64(limit),
			Offset: uint64(offset),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed get users: %w", err)
	}

	res := make([]model.User, len(users.Users))
	for i, user := range users.Users {
		res[i] = mapper.PbToUser(user)
	}

	return res, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID uint64) error {
	if _, err := s.client.DeleteUser(ctx, &accountpb.DeleteUserRequest{UserId: userID}); err != nil {
		return fmt.Errorf("failed delete user: %w", err)
	}

	return nil
}

func (s *Service) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	res, err := s.client.CreateUser(ctx, &accountpb.CreateUserRequest{
		User: mapper.UserCreateToPb(user),
	})
	if err != nil {
		return model.User{}, fmt.Errorf("failed create user: %w", err)
	}

	return mapper.PbToUser(res.User), nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error {
	if _, err := s.client.UpdateUser(ctx, &accountpb.UpdateUserRequest{
		UserId: userID,
		User:   mapper.UserUpdateToPb(user),
	}); err != nil {
		return fmt.Errorf("failed update user: %w", err)
	}

	return nil
}
