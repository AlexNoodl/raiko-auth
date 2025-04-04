package services

import (
	"context"
	"github.com/sirupsen/logrus"
	"github/alexnoodl/raiko-auth/internal/models"
	pb "github/alexnoodl/raiko-auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcServer struct {
	pb.UnimplementedAuthServiceServer
	AuthService *AuthService
	logger      *logrus.Logger
}

func NewAuthGrpcServer(authService *AuthService, logger *logrus.Logger) *AuthGrpcServer {
	return &AuthGrpcServer{
		AuthService: authService,
		logger:      logger,
	}
}

func (s *AuthGrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.WithField("login", req.Login).Info("gRPC Login request received")

	token, err := s.AuthService.Login(req.Login, req.Password)
	if err != nil {
		s.logger.WithError(err).Error("gRPC Login failed")
		return &pb.LoginResponse{Error: err.Error()}, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (s *AuthGrpcServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"username": req.Username,
		"email":    req.Email,
	}).Info("gRPC Register request received")

	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err := s.AuthService.Register(user)
	if err != nil {
		s.logger.WithError(err).Error("gRPC Register failed")
		return &pb.RegisterResponse{Error: err.Error()}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.RegisterResponse{Message: "Successfully registered user"}, nil
}
