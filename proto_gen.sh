protoc --proto_path=api/proto/blogpb --proto_path=third_party --go_out=plugins=grpc:api/proto/blogpb blog.proto
protoc --proto_path=api/proto/blogpb --proto_path=third_party --grpc-gateway_out=logtostderr=true:api/proto/blogpb blog.proto
protoc --proto_path=api/proto/blogpb --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/blog blog.proto
