package routes

import (
	"github.com/Rassimdou/URL-Shortener/api/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	//Get URL parameter from request path
	url := c.Params("url")
	// Create Redis client for database 0 (URL storage)
	r := database.CreateClient(0)
	defer r.Close() // Ensure connection closes when function exits

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found in the database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to the database",
		})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	//  Increment counter for statistics
	_ = rInr.Incr(database.Ctx, "counter")

	// Redirect to original URL
	return c.Redirect(value, 301) // 301 = Permanent redirect

}
