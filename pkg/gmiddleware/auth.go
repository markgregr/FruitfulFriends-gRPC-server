package gmiddleware

import (
	"context"
	"errors"
	"github.com/markgregr/FruitfulFriends-gRPC-server/internal/services/auth"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiTokenKey = "secret"
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
	//if md.Get(apiTokenKey) != s.apiToken {
	//	return ctx, status.Errorf(codes.Unauthenticated, "wrong api token")
	//}

	if strings.HasPrefix(method, "/auth.Auth/Register") ||
		strings.HasPrefix(method, "/auth.Auth/Login") ||
		strings.HasPrefix(method, "/auth.Auth/IsAdmin") {
		return ctx, nil
	}

	accessToken := md.Get("access_token")
	if accessToken == "" {
		return ctx, status.Errorf(codes.Unauthenticated, "access token is empty")
	}

	user, err := s.authService.AuthByToken(ctx, accessToken)
	if errors.Is(err, ErrRecordNotFound) || err == ErrInvalidToken {
		return ctx, status.Errorf(codes.Unauthenticated, "user not find")
	}
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, "user", user)

	return ctx, nil
}
