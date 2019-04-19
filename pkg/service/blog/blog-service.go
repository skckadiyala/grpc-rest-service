package blog

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/skckadiyala/blog-svc/api/proto/blogpb"
	"github.com/skckadiyala/blog-svc/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server struct to implement the Blog Server service interface
type Server struct {
}

// Collection mongoDB collection variable
var Collection *mongo.Collection

type blogItem struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	AuthorID string               `bson:"author_id"`
	Content  string               `bson:"content"`
	Title    string               `bson:"title"`
	Email    string               `bson:"email"`
	Created  *timestamp.Timestamp `bson:"Created"`
	Modified *timestamp.Timestamp `bson:"Updated"`
}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
		Email:    data.Email,
		Created:  data.Created,
		Modified: data.Created,
	}
}

// CreateBlog Implementation of the service
func (*Server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	logger.Log.Info("Create blog request")
	blog := req.GetBlog()

	t := time.Now().In(time.UTC)
	curTime, _ := ptypes.TimestampProto(t)

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
		Email:    blog.GetEmail(),
		Created:  curTime,
		Modified: curTime,
	}

	res, err := Collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
			Email:    blog.GetEmail(),
			Created:  blog.GetCreated(),
			Modified: blog.GetModified(),
		},
	}, nil

}

// ReadBlog Implementation of the service
func (*Server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	logger.Log.Info("Read Blog Request")

	blogID := req.GetId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := Collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

// UpdateBlog Implementation of the service
func (*Server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	logger.Log.Info("Update blog request")

	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := Collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}
	t := time.Now().In(time.UTC)
	curTime, err := ptypes.TimestampProto(t)
	if err != nil {
		log.Fatalf("error while creating date %v", err)
	}

	// we update our internal struct
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()
	data.Modified = curTime

	_, updateErr := Collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr),
		)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

// DeleteBlog Implementation of the service
func (*Server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	logger.Log.Info("Delete blog request")

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	filter := bson.M{"_id": oid}

	res, err := Collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog in MongoDB: %v", err),
		)
	}

	return &blogpb.DeleteBlogResponse{Id: req.GetId()}, nil
}
