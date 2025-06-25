package repositories

import (
	"context"

	"github.com/Rassimdou/URL-Shortener/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatsRepository struct {
	client *mongo.Client
}

func NewStatsRepository(client *mongo.Client) *StatsRepository {
	return &StatsRepository{client: client}
}

func (r *StatsRepository) IncrementTotalClicks() error {
	collection := r.client.Database("urlshortener").Collection("stats")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": "totalClicks"},
		bson.M{"$inc": bson.M{"value": 1}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *StatsRepository) GetRateLimit(ip string) (*models.RateLimit, error) {
	collection := r.client.Database("urlshortener").Collection("stats")
	var limit models.RateLimit
	err := collection.FindOne(context.Background(), bson.M{"_id": ip}).Decode(&limit)
	return &limit, err
}

func (r *StatsRepository) SetRateLimit(limit *models.RateLimit) error {
	collection := r.client.Database("urlshortener").Collection("stats")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": limit.IP},
		bson.M{"$set": limit},
		options.Update().SetUpsert(true),
	)
	return err
}
