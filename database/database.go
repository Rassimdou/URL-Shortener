package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctxn = context.Background()

func CreateClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetMaxPoolSize(100).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Add these functions to your database.go file
func GetURLCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("urlshortener").Collection("urls")
}

func GetStatsCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("urlshortener").Collection("stats")
}
