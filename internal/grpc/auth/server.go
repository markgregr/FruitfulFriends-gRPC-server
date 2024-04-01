package auth

import (
	"context"
	"github.com/go-playground/validator/v10"
	ssov1 "github.com/markgregr/FruitfulFriends-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		//TODO:....
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		//TODO...
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		//TODO...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
