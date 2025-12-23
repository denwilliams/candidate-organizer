# Candidate Organizer

A comprehensive web application to help organize and triage job candidates with AI-powered insights.

**Note**: This application is designed for single-organization use. Each organization is expected to have their own installation - no multi-tenancy is required.

## Features

### Core Features
- **Job Posting Management**: Create and manage job postings with detailed requirements
- **Candidate Management**: Add candidates via resume upload (PDF) or manual entry
- **Resume Parsing**: Automatically extract key information from resumes (name, contact, skills, experience)
- **Custom Attributes**: Set custom candidate attributes that can't be parsed from resumes
- **Status Tracking**: Track candidate progress through the hiring pipeline (Applied, Screened, Interviewing, Offered, Rejected)
- **Comments System**: Add multiple comments per candidate for collaborative hiring decisions
- **AI-Powered Summaries**: Generate AI summaries highlighting candidate strengths, overlaps with job requirements, and potential concerns
- **AI Chat Assistant**: Discuss candidates and get insights (e.g., "What are the top 3 candidates for this job posting?")
- **Secure Authentication**: Google Workspace authentication with domain restriction
- **Role-Based Access**: Admin and user roles with permission controls
- **Salary Privacy**: Salary expectations visible only to admins

### Future Features
- Advanced filtering and sorting
- CSV export for reporting
- Embeddings generation for candidates and job postings for advanced search and matching capabilities
- Email notifications
- Interview scheduling

## Tech Stack

### Frontend
- **Framework**: Next.js 15 with React 18
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: Shadcn UI (to be added)

### Backend
- **Language**: Go 1.24+
- **Router**: Chi v5
- **Database**: PostgreSQL 16
- **Architecture**: Clean architecture with repository pattern (no ORM, handwritten SQL)

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Database**: PostgreSQL with UUID support

### LLM Integration
- **AI Provider**: OpenAI GPT-5.2 for AI features
- **Chat Implementation**: Streaming API for real-time chat responses
- **Use Cases**: Candidate summaries, intelligent chat assistant for candidate evaluation

## Project Structure

```
candidate-organizer/
├── backend/                 # Go backend application
│   ├── cmd/
│   │   └── server/         # Main application entry point
│   ├── internal/
│   │   ├── api/            # HTTP handlers and routing
│   │   ├── auth/           # Authentication logic
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection
│   │   ├── models/         # Data models
│   │   ├── repository/     # Database repository interfaces and implementations
│   │   └── service/        # Business logic
│   ├── migrations/         # Database migrations
│   └── pkg/                # Public packages
├── frontend/               # Next.js frontend application
│   ├── app/                # Next.js app directory
│   ├── components/         # React components
│   └── lib/                # Utility functions
├── .docker/                # Docker-related files
├── docs/                   # Documentation
├── docker-compose.yml      # Production Docker Compose
├── docker-compose.dev.yml  # Development Docker Compose
└── TODO.md                 # Implementation plan

```

## Getting Started

### Prerequisites

- **Docker** and **Docker Compose** (recommended)
- **Go 1.24+** (for local backend development)
- **Node.js 20+** and **npm** (for local frontend development)
- **PostgreSQL 16** (if running without Docker)

### Quick Start with Docker Compose

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd candidate-organizer
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env and add your Google OAuth credentials
   ```

3. **Start the development environment**
   ```bash
   # Start only the database (for local development)
   docker-compose -f docker-compose.dev.yml up -d

   # OR start everything with Docker
   docker-compose up --build
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health check: http://localhost:8080/health

### Local Development (without Docker)

#### Backend Setup

1. **Install Go dependencies**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL**
   ```bash
   docker-compose -f docker-compose.dev.yml up -d postgres
   ```

4. **Run the backend**
   ```bash
   go run cmd/server/main.go
   ```

#### Frontend Setup

1. **Install dependencies**
   ```bash
   cd frontend
   npm install
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your configuration
   ```

3. **Run the development server**
   ```bash
   npm run dev
   ```

## Configuration

### Environment Variables

#### Backend (.env in backend/)
```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/candidate_organizer?sslmode=disable
POSTGRES_SCHEMA=public  # PostgreSQL schema name (optional, defaults to 'public')
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/callback
WORKSPACE_DOMAIN=yourcompany.com
JWT_SECRET=your-super-secret-jwt-key
FRONTEND_URL=http://localhost:3000
OPENAI_API_KEY=your-openai-api-key  # Optional, for AI features
```

#### Frontend (.env.local in frontend/)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a new project or select an existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URIs:
   - `http://localhost:8080/api/v1/auth/callback` (development)
   - Your production callback URL
6. Copy the Client ID and Client Secret to your `.env` file

## Database

### Migrations

The database schema is initialized automatically when you start the PostgreSQL container. The migration files are located in `backend/migrations/`.

### Custom PostgreSQL Schemas

The application supports running on custom PostgreSQL schemas, allowing you to:
- Run multiple isolated instances on the same database
- Organize database objects by environment (e.g., `dev`, `staging`, `prod`)
- Implement multi-tenancy at the schema level

To use a custom schema:

1. **Set the environment variable** in your `.env` file or docker-compose:
   ```env
   POSTGRES_SCHEMA=my_custom_schema
   ```

2. **Start the application** - The schema will be created automatically during migration

3. **Example use cases**:
   ```env
   # Different environments
   POSTGRES_SCHEMA=app_production
   POSTGRES_SCHEMA=app_staging

   # Multi-tenant setups
   POSTGRES_SCHEMA=tenant_acme
   POSTGRES_SCHEMA=tenant_widgets_inc
   ```

**Note**: If not specified, the application defaults to the `public` schema.

### Schema Overview

- **users**: User accounts with role-based access
- **job_postings**: Job posting details
- **candidates**: Candidate information and resume data
- **comments**: Comments on candidates
- **candidate_attributes**: Custom attributes for candidates
- **ai_summaries**: Cached AI-generated summaries

## API Documentation

The backend exposes a RESTful API at `http://localhost:8080/api/v1`. Key endpoints:

### Authentication
- `GET /api/v1/auth/google` - Initiate Google OAuth
- `GET /api/v1/auth/callback` - OAuth callback
- `POST /api/v1/auth/refresh` - Refresh access token

### Users
- `GET /api/v1/users` - List all users (admin only)
- `GET /api/v1/users/me` - Get current user
- `POST /api/v1/users/{id}/promote` - Promote user to admin (admin only)

### Job Postings
- `GET /api/v1/jobs` - List job postings
- `POST /api/v1/jobs` - Create job posting
- `GET /api/v1/jobs/{id}` - Get job posting
- `PUT /api/v1/jobs/{id}` - Update job posting
- `DELETE /api/v1/jobs/{id}` - Delete job posting

### Candidates
- `GET /api/v1/candidates` - List candidates
- `POST /api/v1/candidates` - Create candidate
- `POST /api/v1/candidates/upload` - Upload resume
- `GET /api/v1/candidates/{id}` - Get candidate
- `PUT /api/v1/candidates/{id}` - Update candidate
- `DELETE /api/v1/candidates/{id}` - Delete candidate
- `PUT /api/v1/candidates/{id}/status` - Update candidate status

### Comments
- `GET /api/v1/candidates/{id}/comments` - List comments
- `POST /api/v1/candidates/{id}/comments` - Add comment
- `PUT /api/v1/candidates/{id}/comments/{commentId}` - Update comment
- `DELETE /api/v1/candidates/{id}/comments/{commentId}` - Delete comment

### AI Features
- `POST /api/v1/candidates/{id}/summary` - Generate AI summary
- `POST /api/v1/chat` - AI chat assistant

## Development

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### Building for Production

```bash
# Build everything with Docker
docker-compose up --build

# Or build individually
cd backend && go build -o bin/server cmd/server/main.go
cd frontend && npm run build
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please open an issue on GitHub.
