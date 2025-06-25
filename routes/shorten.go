package routes

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Rassimdou/URL-Shortener/api/database"
	"github.com/Rassimdou/URL-Shortener/api/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset float64       `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Connect to MongoDB
	client, err := database.CreateClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to the database",
		})
	}
	defer client.Disconnect(context.Background())

	statsCollection := database.GetStatsCollection(client)
	urlsCollection := database.GetURLCollection(client)

	// Rate limiting implementation
	ip := c.IP()
	var rateLimit struct {
		Remaining int       `bson:"remaining"`
		Reset     time.Time `bson:"reset"`
	}

	filter := bson.M{"_id": ip}
	err = statsCollection.FindOne(database.Ctx, filter).Decode(&rateLimit)

	quota, _ := strconv.Atoi(os.Getenv("APIQUOTA"))
	now := time.Now()

	if err == mongo.ErrNoDocuments {
		// New IP - initialize rate limit
		rateLimit = struct {
			Remaining int       `bson:"remaining"`
			Reset     time.Time `bson:"reset"`
		}{
			Remaining: quota - 1,
			Reset:     now.Add(30 * time.Minute),
		}
		_, err = statsCollection.InsertOne(database.Ctx, bson.M{
			"_id":       ip,
			"remaining": rateLimit.Remaining,
			"reset":     rateLimit.Reset,
		})
	} else {
		// Existing IP - check rate limit
		if rateLimit.Reset.Before(now) {
			// Reset period passed
			rateLimit.Remaining = quota - 1
			rateLimit.Reset = now.Add(30 * time.Minute)
		} else if rateLimit.Remaining <= 0 {
			// Rate limit exceeded
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "rate limit exceeded",
				"rate_limit_reset": time.Until(rateLimit.Reset).Minutes(),
			})
		} else {
			// Decrement remaining
			rateLimit.Remaining--
		}

		// Update rate limit
		update := bson.M{
			"$set": bson.M{
				"remaining": rateLimit.Remaining,
				"reset":     rateLimit.Reset,
			},
		}
		statsCollection.UpdateOne(database.Ctx, filter, update)
	}

	// Validate URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "you cant hack the system (:"})
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	// Generate or use custom short code
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// Check if short code exists
	var existingURL struct{}
	err = urlsCollection.FindOne(database.Ctx, bson.M{"shortCode": id}).Decode(&existingURL)
	if err == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "custom short URL already exists",
		})
	}

	// Set default expiry
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// Create URL document
	_, err = urlsCollection.InsertOne(database.Ctx, bson.M{
		"shortCode":   id,
		"originalUrl": body.URL,
		"createdAt":   time.Now(),
		"expiresAt":   time.Now().Add(time.Duration(body.Expiry) * time.Hour),
		"clicks":      0,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to save URL",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  rateLimit.Remaining,
		XRateLimitReset: time.Until(rateLimit.Reset).Minutes(),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
