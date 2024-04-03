package gserver

import (
	"context"
	"fmt"
	"runtime/debug"

	grpcmdlw "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StdUnaryMiddleware(log *log.Entry, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	arr := []grpc.UnaryServerInterceptor{
		grpcctxtags.UnaryServerInterceptor(),
		grpclogrus.UnaryServerInterceptor(log),
		grpclogrus.PayloadUnaryServerInterceptor(log, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}),
		grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(func(i interface{}) error {
			log.WithField("panic_stack", string(debug.Stack())).
				Error("grpc panic")

			return fmt.Errorf("%#v", i)
		})),
	}
	arr = append(arr, interceptors...)

	return grpc.UnaryInterceptor(
		grpcmdlw.ChainUnaryServer(arr...),
	)
}

func StdStreamMiddleware(log *log.Entry, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	arr := []grpc.StreamServerInterceptor{
		grpcctxtags.StreamServerInterceptor(),
		grpclogrus.StreamServerInterceptor(log),
		grpclogrus.PayloadStreamServerInterceptor(log, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}),
		grpcrecovery.StreamServerInterceptor(),
	}
	arr = append(arr, interceptors...)

	return grpc.StreamInterceptor(
		grpcmdlw.ChainStreamServer(arr...),
	)
}
