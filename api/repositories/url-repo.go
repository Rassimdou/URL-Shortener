package repositories

import (
	"context"
	"time"

	"github.com/Rassimdou/URL-Shortener/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLRepository struct {
	client *mongo.Client
}

func NewURLRepository(client *mongo.Client) *URLRepository {
	return &URLRepository{client: client}
}

func (r *URLRepository) GetURL(shortCode string) (*models.URL, error) {
	collection := r.client.Database("urlshortener").Collection("urls")
	var url models.URL
	err := collection.FindOne(context.Background(), bson.M{"shortCode": shortCode}).Decode(&url)
	return &url, err
}

func (r *URLRepository) CreateURL(url *models.URL) error {
	collection := r.client.Database("urlshortener").Collection("urls")
	_, err := collection.InsertOne(context.Background(), url)
	return err
}

func (r *URLRepository) IncrementClicks(shortCode string) error {
	collection := r.client.Database("urlshortener").Collection("urls")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"shortCode": shortCode},
		bson.M{
			"$inc": bson.M{"clicks": 1},
			"$set": bson.M{"lastAccessed": time.Now()},
		},
	)
	return err
}

func (r *URLRepository) CheckShortCodeExists(shortCode string) (bool, error) {
	collection := r.client.Database("urlshortener").Collection("urls")
	count, err := collection.CountDocuments(context.Background(), bson.M{"shortCode": shortCode})
	return count > 0, err
}
