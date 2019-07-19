package middleware

import (
	"context"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/skckadiyala/blog-svc/auth"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var au auth.Authenticator

var APIAuth bool

// codeToLevel redirects OK to DEBUG level logging instead of INFO
// This is example how you can log several gRPC code results
func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		// It is DEBUG
		return zap.DebugLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

// AddLogging returns grpc.Server config option that turn on logging.
func AddLogging(logger *zap.Logger, opts []grpc.ServerOption) []grpc.ServerOption {
	// Shared options for the logger, with a custom gRPC code to log level function.
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}
	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLogger(logger)

	// Add unary interceptor
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger, o...),
		grpc.UnaryServerInterceptor(unaryInterceptor),
	))

	// Add stream interceptor (added as an example here)
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(logger, o...),
	))

	return opts
}

// unaryInterceptor call Auth with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if APIAuth {
		auth, err := auth.ExtractHeader(ctx, "authorization")
		if err != nil {
			return ctx, err
		}
		if strings.HasPrefix(auth, "Basic ") {
			_, err = au.BasicAuth(ctx)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(auth, "Apikey ") {
			_, err = au.APIKeyAuth(ctx)
			if err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}

	return handler(ctx, req)
}
