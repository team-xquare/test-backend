# Deployment Platform Backend

A minimal Go backend for the educational deployment platform that triggers GitHub Actions for infrastructure management.

## Features

- **Authentication**: JWT-based auth with refresh tokens
- **Projects**: Create and manage deployment projects
- **Applications**: Deploy applications with various build types
- **Addons**: Deploy database and infrastructure addons
- **GitHub Integration**: GitHub App webhooks and repository dispatch

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout user

### Projects
- `GET /api/v1/projects` - Get user projects
- `POST /api/v1/projects` - Create project
- `GET /api/v1/projects/:id` - Get project details
- `DELETE /api/v1/projects/:id` - Delete project
- `POST /api/v1/projects/:id/applications` - Deploy application
- `POST /api/v1/projects/:id/addons` - Deploy addon

### GitHub
- `POST /api/v1/github/webhook` - GitHub App webhooks
- `GET /api/v1/github/installations` - Get GitHub installations

## Environment Variables

```env
APP_PORT=8080
JWT_SECRET=your-jwt-secret
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=deployment_platform
MYSQL_USERNAME=root
MYSQL_PASSWORD=password
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
GITHUB_APP_ID=your-github-app-id
GITHUB_PRIVATE_KEY=your-github-private-key
GITHUB_WEBHOOK_SECRET=your-webhook-secret
```

## Running

```bash
cd test-backend
go mod tidy
go run cmd/api/main.go
```