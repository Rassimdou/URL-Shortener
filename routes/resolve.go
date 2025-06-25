package routes

import (
	"context"
	"time"

	"github.com/Rassimdou/URL-Shortener/api/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	client, err := database.CreateClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to the database",
		})
	}
	defer client.Disconnect(context.Background())

	urlsCollection := database.GetURLCollection(client)
	statsCollection := database.GetStatsCollection(client)

	// Find the URL document
	var result struct {
		OriginalURL string `bson:"originalUrl"`
	}
	filter := bson.M{"shortCode": url}
	err = urlsCollection.FindOne(database.Ctx, filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found in the database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}

	// Update click count
	update := bson.M{
		"$inc": bson.M{"clicks": 1},
		"$set": bson.M{"lastAccessed": time.Now()},
	}
	urlsCollection.UpdateOne(database.Ctx, filter, update)

	// Update global counter
	statsCollection.UpdateOne(
		database.Ctx,
		bson.M{"_id": "totalClicks"},
		bson.M{"$inc": bson.M{"value": 1}},
		options.Update().SetUpsert(true),
	)

	return c.Redirect(result.OriginalURL, 301)
}
