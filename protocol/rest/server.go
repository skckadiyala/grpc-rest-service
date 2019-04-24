package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/skckadiyala/blog-svc/api/proto/blogpb"
	"github.com/skckadiyala/blog-svc/logger"
	"github.com/skckadiyala/blog-svc/protocol/rest/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RunServer to start the REST/HTTP server
func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := blogpb.RegisterBlogServiceHandlerFromEndpoint(ctx, gwmux, "localhost:"+grpcPort, opts); err != nil {
		logger.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}
	var Cert tls.Certificate

	mux := http.NewServeMux()
	mux.HandleFunc("/apidoc/", serveSwagger)
	mux.Handle("/", gwmux)
	curdir, _ := os.Getwd()
	fmt.Println("cur dir", curdir)

	if httpPort == "8443" {
		logger.Log.Info("getting certs REST/HTTP ssl server...")
		certPEM, err := ioutil.ReadFile("/opt/blog/secure/cert.pem")
		if err != nil {
			logger.Log.Fatal("failed to read cert pem file", zap.String("reason", err.Error()))
		}
		keyPEM, err := ioutil.ReadFile("/opt/blog/secure/key.pem")
		if err != nil {
			logger.Log.Fatal("failed to read key pem file", zap.String("reason", err.Error()))
		}
		Cert, err = tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			logger.Log.Fatal("failed to parse key pair:", zap.String("reason", err.Error()))
		}
	}

	srv := &http.Server{
		Addr: ":" + httpPort,
		// add handler with middleware
		Handler: middleware.AddRequestID(
			middleware.AddLogger(logger.Log, mux)), //mux
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{Cert},
		},
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Info("shutting down REST/HTTP server...")
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	if httpPort == "8443" {
		logger.Log.Info("starting HTTPS/REST gateway...")
		return srv.ListenAndServeTLS("", "")
	}
	logger.Log.Info("starting HTTP/REST gateway...")
	return srv.ListenAndServe() // without TLS
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/apidoc/")
	p = path.Join("/opt/blog/swagger-ui/", p)
	// p = path.Join("third_party/swagger-ui/", p)
	fmt.Println("request map ", p)
	http.ServeFile(w, r, p)

}
