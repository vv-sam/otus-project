package grpc_services

import (
	"context"

	api "github.com/vv-sam/otus-project/proto/grpc/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService interface {
	Login(username, password string) bool
	GenerateToken(username string) (string, error)
}

type AuthService struct {
	api.UnimplementedAuthServiceServer
	userService userService
}

func NewAuthService(userService userService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if !s.userService.Login(req.Username, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	token, err := s.userService.GenerateToken(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &api.LoginResponse{Token: token}, nil
}
