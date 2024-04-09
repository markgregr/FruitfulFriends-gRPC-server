package gmiddleware

import (
	"context"
	"errors"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/lib/jwt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/services/auth"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidToken   = errors.New("invalid token")
)

type Auth struct {
	apiToken    string
	authService *auth.AuthService
}

func NewAuthInterceptor(token string, authService *auth.AuthService) *Auth {
	return &Auth{
		apiToken:    token,
		authService: authService,
	}
}

func (s *Auth) AuthFunc(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	// disable auth for reflection service
	if strings.HasPrefix(method, "/grpc.reflection.v1alpha.ServerReflection") {
		return ctx, nil
	}

	md := metautils.ExtractIncoming(ctx)

	if strings.HasPrefix(method, "/auth.Auth/Register") ||
		strings.HasPrefix(method, "/auth.Auth/Login") ||
		strings.HasPrefix(method, "/auth.Auth/IsAdmin") {
		return ctx, nil
	}

	appID, err := strconv.Atoi(md.Get("app_id"))
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "app id is invalid")
	}

	accessToken := md.Get("access_token")
	if accessToken == "" {
		return ctx, status.Errorf(codes.Unauthenticated, "access token is empty")
	}

	app, err := s.authService.AppByID(ctx, appID)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "app not find")
	}

	_, err = jwt.ParseAndValidateToken(accessToken, app.Secret)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "token parse failed")
	}

	userID, err := s.authService.Authentication(ctx, accessToken)
	if errors.Is(err, ErrRecordNotFound) || err == ErrInvalidToken {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, "userID", userID)

	return ctx, nil
}
