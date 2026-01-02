# BetterUptime Go

A Go-based website uptime monitoring system that tracks website availability across multiple regions using Redis streams and PostgreSQL.

## Features

- **User Management**: Register and authenticate users with JWT tokens.
- **Website Monitoring**: Add websites to monitor, track uptime status (UP, DOWN, UNKNOWN).
- **Distributed Processing**: Uses Redis streams for queuing and consumer groups for region-based processing.
- **Real-time Monitoring**: Periodic checks with response time tracking.
- **API-Driven**: RESTful API for user and website management.
- **Scalable**: Supports multiple workers per region for high availability.

## Architecture

- **API (`cmd/api/main.go`)**: HTTP server for user/website endpoints, runs database migrations.
- **Pusher (`cmd/pusher/main.go`)**: Periodically pushes websites to Redis stream.
- **Worker (`cmd/worker/main.go`)**: Consumes messages, fetches website status, stores results.
- **Database**: PostgreSQL for users, websites, regions, and ticks.
- **Queue**: Redis streams for asynchronous processing.

## Usage

### Running Components

1. **API Server**:

   ```bash
   go run cmd/api/main.go -port 8080
   ```

2. **Pusher** (pushes websites every 3 minutes):

   ```bash
   go run cmd/pusher/main.go
   ```

3. **Worker** (per region):
   ```bash
   REGION_ID=us-east-1 WORKER_ID=worker-1 go run cmd/worker/main.go
   ```

### API Endpoints

- `POST /user/register`: Register a new user.
- `POST /user/login`: Login and get JWT token.
- `POST /website`: Add a website (requires auth).
- `GET /status/:website_id`: Particular website info.

Use `Authorization: Bearer <token>` for authenticated requests.

### Monitoring

- Check Redis stream: `redis-cli XLEN Betteruptime:Websites`
- Check DB ticks: `SELECT * FROM website_tick;`
- Logs: Each component outputs to stdout.

## Configuration

- **Environment Variables**:

  - `DATABASE_URL`: PostgreSQL connection string.
  - `REDIS_ADDR`: Redis address.
  - `REGION_ID`: Region for worker (must match DB).
  - `WORKER_ID`: Unique worker ID.

- **Database Schema**: Auto-migrated via Goose.
