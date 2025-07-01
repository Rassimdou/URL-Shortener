# URL Shortener

A high-performance URL shortening service built with Go and MongoDB, using the Fiber web framework.

## Features

- 📝 Shorten long URLs to compact, shareable links
- 🔄 Redirect shortened URLs to their original destinations
- ⏳ URL expiration support
- 🔒 Rate limiting to prevent abuse
- 🔄 Automatic click tracking

## Tech Stack

- **Backend**: Go (Fiber framework)
- **Database**: MongoDB
- **Middleware**: Custom rate limiting and logging

## Project Structure

```
.
├── main.go              # Main application entry point
├── models/              # Data models (URL, RateLimit, Stats)
├── routes/              # API route handlers
├── repositories/        # Database operations
├── helpers/            # Utility functions
├── database/           # Database configuration
└── data/              # Static data files
```

## API Endpoints

### Shorten URL

- **POST** `/api/v1`
- **Request Body**:
  ```json
  {
    "url": "https://example.com/long-url",
    "short": "custom-short-code", 
    "expiry": 86400               
  }
  ```
- **Response**:
  ```json
  {
    "url": "https://example.com/long-url",
    "short": "abc123",
    "expiry": 86400,
    "rate_limit": 100,            
    "rate_limit_reset": 1687898484.0  
  }
  ```

### Redirect URL

- **GET** `/:url` (e.g., `/abc123`)
- Automatically redirects to the original URL
- Updates click count and last accessed timestamp
- Returns 404 if URL not found
- Returns 500 for database errors

## Environment Variables

Create a `.env` file with the following variables:

```env
APP_PORT=3000
MONGODB_URI=mongodb://localhost:27017
```

## Installation

1. Install Go (version 1.18 or higher)
2. Install MongoDB
3. Clone the repository
4. Set up environment variables
5. Run `go mod tidy` to install dependencies
6. Start the application with `go run main.go`

## Usage

1. Shorten a URL:
   ```bash
   curl -X POST http://localhost:3000/api/v1 \
   -H "Content-Type: application/json" \
   -d '{"url": "https://example.com/long-url"}'
   ```

2. Access shortened URL:
   ```bash
   http://localhost:3000/abc123
   ```



