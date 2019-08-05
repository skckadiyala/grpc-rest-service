package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/skckadiyala/blog-svc/logger"
	"github.com/skckadiyala/blog-svc/pkg/service/blog"
	"github.com/skckadiyala/blog-svc/protocol/grpc"
	"github.com/skckadiyala/blog-svc/protocol/grpc/middleware"
	"github.com/skckadiyala/blog-svc/protocol/rest"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection *mongo.Collection

// func (*server) ListBlog(ctx context.Context, req *blogpb.ListBlogRequest) (*blogpb.ListBlogResponse, error) {
// 	return nil, nil
// }

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// GRPCPort is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	// DB Datastore parameters section
	// MongoDBHost is url of database
	MongoDBHost string

	// MongoDBUser is username to connect to database
	MongoDBUser string
	// MongoDBPassword password to connect to database
	MongoDBPassword string
	// MongoDBSchema is schema of database
	MongoDBSchema string

	MongoDBCollection string

	// Auth Authorization for API Calls
	Auth bool

	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ctx := context.Background()
	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port to bind")
	flag.StringVar(&cfg.MongoDBHost, "db-host", "", "Mongo Database host")
	flag.StringVar(&cfg.MongoDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.MongoDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.MongoDBSchema, "db-schema", "systest", "Database schema")
	flag.StringVar(&cfg.MongoDBCollection, "db-collection", "blog", "Database colletion")

	flag.BoolVar(&cfg.Auth, "api-auth", false, "APIAuthentication")

	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "2006-01-02T15:04:05.999999999Z07:00",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")

	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		fmt.Printf("invalid TCP port for gRPC server: '%s' \n", cfg.GRPCPort)
		os.Exit(0)
	}

	if len(cfg.HTTPPort) == 0 {
		fmt.Printf("invalid TCP port for HTTP gateway: '%s' \n", cfg.HTTPPort)
		os.Exit(0)
	}
	if cfg.Auth {
		middleware.APIAuth = true
	}
	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		fmt.Printf("failed to initialize logger: %v \n", err)
		os.Exit(0)
	}

	logger.Log.Info("Connecting to MongoDB")
	// connect to MongoDB
	var client *mongo.Client
	var err error
	if cfg.MongoDBUser != "" && cfg.MongoDBPassword != "" {

		mgoURI := "mongodb://" + cfg.MongoDBUser + ":" + cfg.MongoDBPassword + "@" + cfg.MongoDBHost + "/" + cfg.MongoDBSchema
		fmt.Println("Mongo URI", mgoURI)
		client, err = mongo.NewClient(options.Client().ApplyURI(mgoURI))
		if err != nil {
			fmt.Printf("failed to created MongoDB client: %v \n", err)
			os.Exit(0)
		}

	} else {
		client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + cfg.MongoDBHost))
		if err != nil {
			fmt.Printf("failed to created MongoDB client: %v \n", err)
			os.Exit(0)
		}
	}

	err = client.Connect(context.TODO())
	if err != nil {
		fmt.Printf("failed to connect MongoDB: %v \n", err)
		os.Exit(0)
	}
	blog.Collection = client.Database(cfg.MongoDBSchema).Collection(cfg.MongoDBCollection)
	logger.Log.Info("Starting Blogging Service ...")
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()
	grpc.RunServer(ctx, blog.Server{}, cfg.GRPCPort)
}
