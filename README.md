# Page Hoppers - Reading Tracker for Children

A web application that helps parents track their children's reading progress both in and out of school. Built with React/Next.js frontend and Go backend.

## Features

- **Parent Authentication**: Register and login for parents
- **Child Management**: Add and manage multiple children per parent
- **Child Login**: Simple PIN-based login for children
- **Reading Tracking**: Track books read by children (coming soon)
- **Modern UI**: Playful, child-friendly design with custom color palette

## Tech Stack

- **Frontend**: React 19, Next.js 15, TypeScript, Tailwind CSS v4
- **Backend**: Go, GORM, PostgreSQL, JWT authentication
- **Database**: PostgreSQL (via Docker)
- **Testing**: Go testing package, in-memory SQLite for tests

## Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop/) (for PostgreSQL and backend)
- [Go](https://golang.org/dl/) (v1.18+)
- [Node.js](https://nodejs.org/) (v18+)
- [npm](https://www.npmjs.com/) or [yarn](https://yarnpkg.com/)

## Local Development Setup

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd page-hoppers
```

### 2. Set Up the Database and Backend with Docker Compose

Start both PostgreSQL and the backend server using Docker Compose:

```bash
docker-compose up
```

- This will start the database and backend in containers.
- The backend will be available at `http://localhost:8080`.
- The database will be available at `localhost:5432` for other tools.

**To run in the background:**
```bash
docker-compose up -d
```

**Database Details:**
- Host: `localhost`
- Port: `5432`
- User: `pagehoppers_user`
- Password: `password314`
- Database: `pagehoppers_db`

### 3. Backend Setup (Local Development Option)

If you prefer to run the backend locally (outside Docker):

Navigate to the backend directory:

```bash
cd page-hoppers-backend
```

Install Go dependencies:

```bash
go mod tidy
```

Create environment file:

```bash
# Create .env file
cat > .env << EOF
DATABASE_URL=postgres://pagehoppers_user:password314@localhost:5432/pagehoppers_db?sslmode=disable
JWT_SECRET=your-super-secret-key-replace-in-production
PORT=8080
EOF
```

Run database migrations:

```bash
go run scripts/migrate.go
```

Create a test parent user:

```bash
go run scripts/create_test_parent.go
```

Start the backend server:

```bash
go run main.go
```

The backend will be available at `http://localhost:8080`

### 4. Alternative: Use the Start Script

You can also use the provided script to start the database and backend for development:

```bash
./start.sh
```

### 5. Frontend Setup

Open a new terminal and navigate to the frontend directory:

```bash
cd page-hoppers-frontend
```

Install dependencies:

```bash
npm install
```

Start the development server:

```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## Running the Application

### Start Everything

1. **Database & Backend (Docker Compose):**
   ```bash
   docker-compose up
   ```

2. **Backend (Local Option):**
   ```bash
   cd page-hoppers-backend
   go run main.go
   ```

3. **Frontend:**
   ```bash
   cd page-hoppers-frontend
   npm run dev
   ```

### Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080

### Test Credentials

After running the setup scripts, you can log in with:
- **Email**: `parent@example.com`
- **Password**: `testpassword`

## Testing

### Backend Tests

Run all backend tests:

```bash
cd page-hoppers-backend
go test ./tests/...
```

Run specific test files:

```bash
# Run auth handler tests
go test ./tests/handlers/

# Run with verbose output
go test -v ./tests/...

# Run with coverage
go test -cover ./tests/...
```

### Test Structure

```
page-hoppers-backend/tests/
├── helpers.go           # Common test utilities
├── handlers/
│   └── auth_test.go     # Authentication tests
├── models/              # Model tests (future)
└── integration/         # Integration tests (future)
```

### Frontend Tests

Frontend testing setup is planned for future implementation.

## API Endpoints

### Public Endpoints
- `POST /api/auth/parent/register` - Parent registration
- `POST /api/auth/parent/login` - Parent login
- `POST /api/auth/child/login` - Child login

### Protected Endpoints (require JWT token)
- `GET /api/children` - Get parent's children
- `POST /api/children` - Create a new child

## Development Workflow

1. **Database**: Use Docker Compose for consistent PostgreSQL setup
2. **Backend**: Go with hot reload (restart on changes)
3. **Frontend**: Next.js with hot reload
4. **Testing**: Run tests before committing changes

## Troubleshooting

### Database Connection Issues
- Ensure Docker Compose is running: `docker ps`
- Check database URL in `.env` file
- Restart containers: `docker-compose restart`

### Backend Issues
- Check if port 8080 is available
- Verify database connection
- Check logs for error messages

### Frontend Issues
- Clear browser cache
- Check if port 3000 is available
- Verify backend is running

### Test Issues
- Ensure all dependencies are installed
- Check that test database can be created
- Verify import paths are correct

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run all tests
6. Submit a pull request

## License

[Add your license here] 