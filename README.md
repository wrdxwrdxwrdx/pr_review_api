# PR Review API

A REST API service for managing pull request reviews, teams, and users. Built with Go, Gin framework, and PostgreSQL.

## Prerequisites

- Docker and Docker Compose
- Go 1.23+ (if running locally without Docker)

## Quick Start with Docker

1. **Clone the repository** (if not already done):

   ```bash
   git clone <repository-url>
   cd pr_review_api
   ```

2. **Create environment file**:

   ```bash
   cp .env.example .env
   ```

   Edit `.env` and adjust the values as needed.

3. **Start the services**:

   ```bash
   docker-compose up --build
   ```

   This will:

   - Start a PostgreSQL 16 database container
   - Build and start the Go API server
   - Automatically initialize the database schema
   - Expose the API on `http://localhost:8080` (or your configured `SERVER_PORT`)

4. **Verify the services are running**:

   ```bash
   docker-compose ps
   ```

5. **Access the database** (optional):

   ```bash
   docker exec -it pr_review_db psql -U postgres -d pr_review
   ```

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Authentication Endpoints

#### Register User

Create a new user account and receive a JWT token.

**Endpoint**: `POST /api/v1/auth/register`

**Authentication**: None (public)

**Request Body**:

```json
{
  "user_id": "u1",
  "username": "Alice"
}
```

**Response** (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "u1",
  "expires_in": 86400
}
```

#### Login

Authenticate an existing user and receive a JWT token.

**Endpoint**: `POST /api/v1/auth/login`

**Authentication**: None (public)

**Request Body**:

```json
{
  "user_id": "u1",
  "username": "Alice"
}
```

**Response** (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "u1",
  "expires_in": 86400
}
```

### Team Endpoints

#### Create Team

Create a new team with members. This endpoint creates/updates users and assigns them to the team.

**Endpoint**: `POST /api/v1/team/add`

**Authentication**: JWT Bearer token (required)

**Headers**:

```
Authorization: Bearer <jwt_token>
```

**Request Body**:

```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    }
  ]
}
```

**Response** (201 Created):

```json
{
  "team": {
    "team_name": "backend",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      }
    ]
  }
}
```

#### Get Team

Retrieve team information with all members.

**Endpoint**: `GET /api/v1/team/get?TeamNameQuery=<team_name>`

**Authentication**: Admin token OR JWT Bearer token

**Headers** (choose one):

```
X-Admin-Token: <admin_token>
```

OR

```
Authorization: Bearer <jwt_token>
```

**Query Parameters**:

- `TeamNameQuery` (required): The name of the team

**Example Request**:

```bash
curl -H "X-Admin-Token: admin-secret-token" \
  "http://localhost:8080/api/v1/team/get?TeamNameQuery=backend"
```

**Response** (200 OK):

```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    }
  ]
}
```

### User Endpoints

#### Set User Active Status

Update a user's active status (admin only).

**Endpoint**: `POST /api/v1/users/setIsActive`

**Authentication**: Admin token (required)

**Headers**:

```
X-Admin-Token: <admin_token>
```

**Request Body**:

```json
{
  "user_id": "u2",
  "is_active": false
}
```

**Response** (200 OK):

```json
{
  "user": {
    "user_id": "u2",
    "username": "Bob",
    "team_name": "backend",
    "is_active": false
  }
}
```

#### Get User Reviews

Get all pull requests where a user is assigned as a reviewer.

**Endpoint**: `GET /api/v1/users/getReview?UserIdQuery=<user_id>`

**Authentication**: Admin token OR JWT Bearer token

**Headers** (choose one):

```
X-Admin-Token: <admin_token>
```

OR

```
Authorization: Bearer <jwt_token>
```

**Query Parameters**:

- `UserIdQuery` (required): The user ID

**Note**: Regular users can only view their own reviews. Admins can view any user's reviews.

**Example Request**:

```bash
curl -H "Authorization: Bearer <jwt_token>" \
  "http://localhost:8080/api/v1/users/getReview?UserIdQuery=u2"
```

**Response** (200 OK):

```json
{
  "user_id": "u2",
  "pull_requests": [
    {
      "pull_request_id": "pr-1001",
      "pull_request_name": "Add search",
      "author_id": "u1",
      "status": "OPEN"
    }
  ]
}
```

### Pull Request Endpoints

#### Create Pull Request

Create a new pull request and automatically assign up to 2 reviewers from the author's team.

**Endpoint**: `POST /api/v1/pullRequest/create`

**Authentication**: Admin token (required)

**Headers**:

```
X-Admin-Token: <admin_token>
```

**Request Body**:

```json
{
  "pull_request_id": "pr-1001",
  "pull_request_name": "Add search",
  "author_id": "u1"
}
```

**Response** (201 Created):

```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u2", "u3"],
    "createdAt": "2025-01-15T10:30:00Z",
    "mergedAt": null
  }
}
```

#### Merge Pull Request

Mark a pull request as merged (idempotent operation).

**Endpoint**: `POST /api/v1/pullRequest/merge`

**Authentication**: Admin token (required)

**Headers**:

```
X-Admin-Token: <admin_token>
```

**Request Body**:

```json
{
  "pull_request_id": "pr-1001"
}
```

**Response** (200 OK):

```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1",
    "status": "MERGED",
    "assigned_reviewers": ["u2", "u3"],
    "createdAt": "2025-01-15T10:30:00Z",
    "mergedAt": "2025-01-15T14:20:00Z"
  }
}
```

#### Reassign Reviewer

Replace a specific reviewer with another active user from the same team.

**Endpoint**: `POST /api/v1/pullRequest/reassign`

**Authentication**: Admin token (required)

**Headers**:

```
X-Admin-Token: <admin_token>
```

**Request Body**:

```json
{
  "pull_request_id": "pr-1001",
  "old_reviewer_id": "u2"
}
```

**Response** (201 Created):

```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u3", "u5"],
    "createdAt": "2025-01-15T10:30:00Z",
    "mergedAt": null
  },
  "replaced_by": "u5"
}
```

## Example Workflow

Here's a complete example workflow using curl:

1. **Register a user**:

   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "u1",
       "username": "Alice"
     }'
   ```

   Save the `token` from the response.

2. **Create a team** (using the JWT token from step 1):

   ```bash
   curl -X POST http://localhost:8080/api/v1/team/add \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token_from_step_1>" \
     -d '{
       "team_name": "backend",
       "members": [
         {"user_id": "u1", "username": "Alice", "is_active": true},
         {"user_id": "u2", "username": "Bob", "is_active": true},
         {"user_id": "u3", "username": "Charlie", "is_active": true}
       ]
     }'
   ```

3. **Create a pull request** (using admin token):

   ```bash
   curl -X POST http://localhost:8080/api/v1/pullRequest/create \
     -H "Content-Type: application/json" \
     -H "X-Admin-Token: admin-secret-token" \
     -d '{
       "pull_request_id": "pr-1001",
       "pull_request_name": "Add search feature",
       "author_id": "u1"
     }'
   ```

4. **Get user reviews** (using JWT token):

   ```bash
   curl -X GET "http://localhost:8080/api/v1/users/getReview?UserIdQuery=u2" \
     -H "Authorization: Bearer <token_from_step_1>"
   ```

5. **Merge the pull request** (using admin token):

   ```bash
   curl -X POST http://localhost:8080/api/v1/pullRequest/merge \
     -H "Content-Type: application/json" \
     -H "X-Admin-Token: admin-secret-token" \
     -d '{
       "pull_request_id": "pr-1001"
     }'
   ```

## Project Structure

```
pr_review_api/
├── api/
│   └── openapi.yaml          # OpenAPI specification
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── domain/               # Domain entities and errors
│   ├── handlers/             # HTTP request handlers
│   ├── middleware/           # HTTP middleware (auth, logging, etc.)
│   ├── repository/           # Data access layer
│   │   ├── postgres/         # PostgreSQL implementation
│   │   └── interfaces/       # Repository interfaces
│   └── services/             # Business logic layer
├── pkg/
│   ├── auth/                 # JWT authentication
│   └── validator/            # Validation utilities
├── Dockerfile                # Docker image definition
├── docker-compose.yml        # Docker Compose configuration
├── go.mod                   # Go module dependencies
└── README.md                 # This file
```

## Database Schema

The application automatically initializes the following database schema:

- **teams**: Stores team information
- **users**: Stores user information with team relationships
- **pull_requests**: Stores pull request information with reviewer assignments

All tables include automatic timestamp management for `created_at` and `updated_at` fields.