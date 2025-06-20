# Go Authentication Service

A robust, production-ready authentication microservice built with Go, following clean architecture principles and best security practices.  
Features include JWT-based authentication, email/phone verification, password reset, role and permission-based authorization, and more.

## Project Structure

```
.
├── cmd/                    # Main application entry points
│   └── server/             # HTTP server
├── internal/               # Private application code
│   ├── auth/               # JWT, OAuth, crypto, and cookie logic
│   ├── api/                # API handlers, middleware, and routes
│   ├── domain/             # Domain models and enums
│   ├── service/            # Business logic (auth, email, sms, etc.)
│   └── infrastructure/     # Database, migrations, repository
├── pkg/                    # Public utilities (token, password, etc.)
├── configs/                # Configuration files
├── templates/              # Email templates (welcome, verification, reset)
├── docker/                 # Dockerfile for containerization
├── scripts/                # Helper scripts
├── test/                   # Test files and examples
├── .github/                # GitHub workflows and configs
├── .golangci.yml           # Linting configuration
├── auth.http               # Example HTTP requests for API testing
├── Makefile                # Common development commands
└── README.md               # Project documentation
```

## Setup

1. **Clone the repository**
2. **Copy and edit the config:**
   ```bash
   cp config.example.json config.json
   # Edit config.json with your settings
   ```
3. **Install dependencies:**
   ```bash
   go mod download
   ```
4. **Run the service:**
   ```bash
   make run
   ```
   Or for development mode:
   ```bash
   make run-dev
   ```

## Makefile Commands

- `make build` — Build the application binary into `bin/`
- `make run` — Run the application
- `make run-dev` — Run in development mode
- `make test` — Run all tests
- `make lint` — Run code linter (`golangci-lint`)
- `make clean` — Remove build artifacts
- `make docker-build` — Build the Docker image using the Dockerfile
- `make docker-run` — Run the Docker container (exposes port 8080)
- `make migrate-up` — Run database migrations
- `make migrate-down` — Roll back database migrations
- `make generate-mocks` — Generate mocks for testing
- `make health` — Health check endpoint

## Running with Docker

You can build and run the service in a container using the provided Dockerfile:

```bash
make docker-build
make docker-run
```

Or manually:

```bash
docker build -t auth:latest -f docker/Dockerfile .
docker run -p 8080:8080 auth:latest
```

## API Testing

- Use the `auth.http` file with [VS Code REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) or JetBrains HTTP Client to test endpoints.
- Example requests for signup, login, verification, etc. are included.

## Code Quality

- Lint your code with:
  ```bash
  make lint
  ```
- Configuration is in `.golangci.yml`.

## Email Templates

- `templates/welcome.html`
- `templates/email_verification.html`
- `templates/password_reset.html`

## Contributing

- Use `make test` to run tests.
- Use `make lint` to check code quality.
- Follow Go best practices and keep PRs focused.

## License

MIT
