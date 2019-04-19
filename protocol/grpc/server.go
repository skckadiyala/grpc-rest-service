package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/skckadiyala/blog-svc/api/proto/blogpb"
	"github.com/skckadiyala/blog-svc/logger"
	"github.com/skckadiyala/blog-svc/pkg/service/blog"
	"github.com/skckadiyala/blog-svc/protocol/grpc/middleware"
	"google.golang.org/grpc"
)

// RunServer to start the gRPC server
func RunServer(ctx context.Context, blogAPI blog.Server, port string) error {

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	opts = middleware.AddLogging(logger.Log, opts)

	grpcServer := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(grpcServer, &blogAPI)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")
			grpcServer.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start gRPC server
	return grpcServer.Serve(listen)
}
