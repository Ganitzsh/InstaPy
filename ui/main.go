package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	var err error

	mongoDBClient, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err != nil {
		logrus.Fatalf("Could not initialize database: %v", err)
	}

	mongoDBClient.Connect(context.Background())
}

func main() {
	api := newAPI()

	globBot = newBot()

	if err := api.Start(); err != nil {
		logrus.Fatalf("Could not start the API: %v", err)
	}
}
