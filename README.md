# URL Shortener

A simple and secure URL shortening service built with Go and Gin framework.

## Features

- **Shorten URLs**: Generate random short codes for long URLs
- **Custom codes**: Create shortened URLs with custom codes
- **Analytics**: Track click counts and access statistics
- **Duplicate prevention**: Reuses existing codes for the same URL
- **URL validation**: Ensures proper URL format
- **Rate limiting protection**: Secure code generation with collision handling

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/shorten` | Create a short URL with random code |
| `POST` | `/shorten/custom` | Create a short URL with custom code |
| `GET` | `/:code` | Redirect to original URL |
| `GET` | `/stats/:code` | Get URL statistics |
| `DELETE` | `/delete/:code` | Delete a short URL |
| `GET` | `/health` | Health check |

## Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd url-shortener

# Install dependencies
go mod tidy

# Set environment variables
export BASE_URL=http://localhost:8080
export DATABASE_URL=your_postgres_connection_string

# Run the application
go run main.go
```

## Usage Examples

### Shorten a URL
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

### Custom short code
```bash
curl -X POST http://localhost:8080/shorten/custom \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "custom_code": "my-link"}'
```

### Get statistics
```bash
curl http://localhost:8080/stats/abc123
```

## Requirements

- Go 1.19+
- PostgreSQL database
- Environment variables:
  - `BASE_URL`: Base URL for shortened links (default: http://localhost:8080)
  - `DATABASE_URL`: PostgreSQL connection string

## Database Schema

The service requires a PostgreSQL table with the following structure:

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    last_accessed TIMESTAMP,
    click_count INTEGER DEFAULT 0
);
```
