# StrikePad

A modern web application with Go backend and React frontend.

## Project Structure

```
strikepad/
├── be/                 # Backend (Go)
│   ├── internal/       # Internal packages
│   ├── migrations/     # Database migrations
│   ├── main.go        # Application entry point
│   └── Makefile       # Build commands
├── fe/                 # Frontend (React + TypeScript)
│   ├── src/           # Source code
│   ├── public/        # Static assets
│   └── package.json   # Dependencies
└── .github/           # GitHub Actions workflows
    └── workflows/
```

## Backend (Go)

### Features

- **Authentication**: Email/password authentication with bcrypt hashing
- **Validation**: Custom password complexity validation with go-playground/validator
- **Error Handling**: Unified error response system with E000-format error codes
- **Database**: PostgreSQL with GORM ORM and Atlas migrations
- **Logging**: Structured logging with slog and hourly rotation
- **Testing**: Comprehensive test coverage with sqlmock

### Quick Start

```bash
cd be
make deps          # Install dependencies
make migrate-apply # Run database migrations
make run          # Start the server
```

### Development

```bash
make dev          # Hot reload with Air
make test         # Run tests
make test-coverage # Run tests with coverage
make lint         # Run linter
```

### API Endpoints

- `POST /api/auth/signup` - User registration
- `POST /api/auth/login` - User authentication
- `GET /health` - Health check

### Error Codes

The API uses a unified error code system:

- **E001-E099**: General errors (validation, server errors)
- **E100-E199**: Authentication errors (login, registration)
- **E200-E299**: Validation errors (field validation)
- **E300-E399**: Business logic errors

See [Error Codes Documentation](be/README_ERROR_CODES.md) for details.

## Frontend (React)

### Features

- **React 18** with TypeScript
- **Vite** for fast development
- **Tailwind CSS** for styling
- **Chart.js** for data visualization
- **Jest** for testing

### Quick Start

```bash
cd fe
npm install       # Install dependencies
npm run dev       # Start development server
```

### Development

```bash
npm run build     # Build for production
npm run test      # Run tests
npm run lint      # Run linter
```

## Database

The project uses PostgreSQL with Atlas for migrations.

### Setup

1. Install PostgreSQL
2. Create database: `createdb strikepad`
3. Run migrations: `cd be && make migrate-apply`

### Migrations

```bash
cd be
make migrate-status    # Check migration status
make migrate-apply     # Apply migrations
make migrate-diff      # Create new migration
```

## Docker

Use Docker Compose for local development:

```bash
docker-compose up -d   # Start PostgreSQL
```

## CI/CD

### GitHub Actions

- **Pull Request Tests**: Runs tests and reports coverage on PRs to `develop`
- **Coverage Reports**: Uploads coverage to Codecov and comments on PRs

### Coverage Requirements

- Minimum coverage threshold: 60%
- Target coverage: 80%+

## Development Workflow

1. Create feature branch from `develop`
2. Make changes and add tests
3. Run tests locally: `make test-coverage`
4. Create pull request to `develop`
5. GitHub Actions will run tests and report coverage
6. Merge after review and passing tests

## Environment Variables

### Backend (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=strikepad
DB_SSLMODE=disable
ENV=development
```

### Frontend (.env.local)

```env
VITE_API_URL=http://localhost:8080
```

## License

MIT License