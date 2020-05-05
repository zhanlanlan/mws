package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// O 全局mongo链接
var O *mongo.Client

func init() {

	if O == nil {
		url := "mongodb://mws_mongo:mws_mongo@127.0.0.1:27017/mws"
		// Set client options
		clientOptions := options.Client().ApplyURI(url)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			log.Fatal(err)
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		O = client
	}

}
