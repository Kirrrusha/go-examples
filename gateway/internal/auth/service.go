package auth

import (
	"context"
	"fmt"

	"gateway/internal/mapper"
	"gateway/internal/model"
	authpb "gateway/pkg/auth/go"
)

type Service struct {
	client authpb.AuthClient
}

func NewService(client authpb.AuthClient) *Service {
	return &Service{client: client}
}

func (s *Service) Login(ctx context.Context, loginOrEmail, password string) (model.TokenPair, error) {
	res, err := s.client.Login(ctx, &authpb.LoginRequest{
		LoginOrEmail: loginOrEmail,
		Password:     password,
	})
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("failed login: %w", err)
	}

	return mapper.PbToTokenPair(res), nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	res, err := s.client.Refresh(ctx, &authpb.RefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("failed refresh token: %w", err)
	}

	return mapper.PbToTokenPair(res), nil
}

func (s *Service) Verify(ctx context.Context, accessToken string) (uint64, error) {
	res, err := s.client.Validate(ctx, &authpb.ValidateRequest{AccessToken: accessToken})
	if err != nil {
		return 0, fmt.Errorf("failed verify access token: %w", err)
	}

	return res.UserId, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if _, err := s.client.Logout(ctx, &authpb.RefreshRequest{RefreshToken: refreshToken}); err != nil {
		return fmt.Errorf("failed logout: %w", err)
	}

	return nil
}

func (s *Service) Register(ctx context.Context, userID uint64, login string, email string, password string) error {
	_, err := s.client.Register(ctx, &authpb.RegisterRequest{
		Id:       userID,
		Login:    login,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID uint64) error {
	if _, err := s.client.DeleteUser(ctx, &authpb.DeleteRequest{Id: userID}); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
