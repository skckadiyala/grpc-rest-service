package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://my-user:my-password@localhost:17017/systest"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("systest").Collection("qux")
	res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println(id)

}
