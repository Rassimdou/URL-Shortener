package routes 

improt (

"time"

)
type request struct {
	URL   			string	  					`json:"url"`
	CustomShort		string						`json:"short"`
	Expiry			time.Duration				`json:"expiry"`	
	
}

type response struct {
	URL					string						`json:"url"`
	CustomShort			string						`json:"short"`
	Expiry				time.Duration				`json:"expiry"`
	XRateRemaining		int							`json:"rate_limit"`
	XRateLimitRest		time.Duration				`json:"rate_limit_reset"`
}


func ShortenURL(c *fiber.Ctx) error{
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse JSON"})

}
//implement rate limiting

//check if URL is valid

if !govalidator.IsURL(body.URL) {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid URL"})
}

//check for domain error 

if !helpers.RemoveDomainError(body.URL) {
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"domain not allowed"})
}

//enforce https ,SSL 

body.URL = helpers.EnforceHTTP(body.URL)

}

