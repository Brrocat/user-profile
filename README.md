# User Profile Service

Microservice for managing user profiles in the car-sharing platform.

## Features

- User profile management (CRUD operations)
- Profile caching with Redis for performance
- Data validation
- gRPC API

## Technology Stack

- **Language**: Go 1.25+
- **Framework**: gRPC
- **Database**: PostgreSQL
- **Cache**: Redis
- **Validation**: go-playground/validator

## API

### gRPC Methods

- `GetUserProfile` - Retrieve user profile by user ID
- `CreateUserProfile` - Create new user profile
- `UpdateUserProfile` - Update existing user profile
- `DeleteUserProfile` - Delete user profile

### Protobuf

See `car-sharing-protos/proto/userprofile/user_profile.proto` for detailed API specification.

## Configuration

### Environment Variables

- `ENV` - Environment (development/production)
- `PORT` - gRPC server port (default: 50052)
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `CACHE_TTL` - Cache time-to-live duration (default: 1h)

## Running Locally

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis

### Steps

1. Clone the repository
2. Copy `.env.example` to `.env` and configure variables
3. Run database migrations:
   ```bash
   psql -d user_profile_db -f migrations/001_create_user_profiles_table.sql# user-profile
# user-profile
