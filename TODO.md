# TODO - Candidate Organizer Implementation Plan

## Phase 1: Project Setup & Infrastructure

### 1.1 Initial Project Structure
- [ ] Create root project directory structure
- [ ] Initialize frontend (Next.js + React + TypeScript)
- [ ] Initialize backend (Golang project with proper module structure)
- [ ] Set up .gitignore for both frontend and backend
- [ ] Create Docker configuration files
- [ ] Create docker-compose.yml for local development

### 1.2 Database Setup
- [ ] Design database schema (ERD)
- [ ] Create migration files for Postgres
  - [ ] Users table (id, email, name, role, workspace_domain, created_at, updated_at)
  - [ ] Job_postings table (id, title, description, requirements, location, salary_range, status, created_at, updated_at, created_by)
  - [ ] Candidates table (id, name, email, phone, resume_url, parsed_data, status, salary_expectation, job_posting_id, created_at, updated_at, created_by)
  - [ ] Comments table (id, candidate_id, user_id, content, created_at, updated_at)
  - [ ] Candidate_attributes table (id, candidate_id, attribute_key, attribute_value, created_at, updated_at)
- [ ] Create SQL migration runner in backend
- [ ] Set up Postgres connection pooling

### 1.3 Backend Foundation
- [ ] Set up Golang project structure (cmd, internal, pkg)
- [ ] Create repository interfaces
  - [ ] UserRepository interface
  - [ ] JobPostingRepository interface
  - [ ] CandidateRepository interface
  - [ ] CommentRepository interface
  - [ ] AttributeRepository interface
- [ ] Implement Postgres repository implementations with handwritten SQL
- [ ] Set up HTTP server (using chi or gin router)
- [ ] Create middleware (logging, CORS, authentication)
- [ ] Implement Google OAuth authentication
  - [ ] Workspace domain validation
  - [ ] JWT token generation and validation
  - [ ] First user becomes admin logic
- [ ] Create error handling utilities
- [ ] Set up environment variable management

### 1.4 Frontend Foundation
- [ ] Set up Next.js project with TypeScript
- [ ] Configure Shadcn UI component library
- [ ] Set up Tailwind CSS
- [ ] Create base layout components
- [ ] Implement authentication flow (Google OAuth)
- [ ] Create protected route wrapper
- [ ] Set up API client/fetch utilities
- [ ] Create global state management (Context API or Zustand)

## Phase 2: Core Features - User Management & Authentication

### 2.1 Authentication
- [ ] Backend: Google OAuth endpoints (/auth/google, /auth/callback)
- [ ] Backend: Token refresh endpoint
- [ ] Backend: User profile endpoint
- [ ] Frontend: Login page with Google sign-in button
- [ ] Frontend: Auth context provider
- [ ] Frontend: Protected route component
- [ ] First user admin assignment logic

### 2.2 User Management (Admin Only)
- [ ] Backend: List all users endpoint
- [ ] Backend: Promote user to admin endpoint
- [ ] Frontend: User management page
- [ ] Frontend: User list with role badges
- [ ] Frontend: Promote to admin button

## Phase 3: Core Features - Job Postings

### 3.1 Job Posting CRUD
- [ ] Backend: Create job posting endpoint
- [ ] Backend: List job postings endpoint (with pagination)
- [ ] Backend: Get single job posting endpoint
- [ ] Backend: Update job posting endpoint
- [ ] Backend: Delete job posting endpoint
- [ ] Frontend: Job postings list page
- [ ] Frontend: Create job posting form
- [ ] Frontend: Edit job posting form
- [ ] Frontend: Job posting detail view
- [ ] Frontend: Delete confirmation modal

## Phase 4: Core Features - Candidates

### 4.1 Candidate Management
- [ ] Backend: Create candidate endpoint (manual entry)
- [ ] Backend: Upload resume endpoint (PDF)
- [ ] Backend: Resume parsing service (extract name, contact, skills, experience)
- [ ] Backend: List candidates endpoint (with pagination and filtering)
- [ ] Backend: Get single candidate endpoint
- [ ] Backend: Update candidate endpoint
- [ ] Backend: Delete candidate endpoint
- [ ] Backend: Update candidate status endpoint
- [ ] Frontend: Candidates list page with filters
- [ ] Frontend: Add candidate form (manual entry)
- [ ] Frontend: Upload resume component
- [ ] Frontend: Candidate detail view
- [ ] Frontend: Edit candidate form
- [ ] Frontend: Status update dropdown
- [ ] Frontend: Delete confirmation modal

### 4.2 Custom Candidate Attributes
- [ ] Backend: Add custom attribute endpoint
- [ ] Backend: Update custom attribute endpoint
- [ ] Backend: Delete custom attribute endpoint
- [ ] Frontend: Custom attributes section in candidate detail
- [ ] Frontend: Add/edit custom attribute form
- [ ] Frontend: Delete attribute confirmation

### 4.3 Salary Expectation (Admin Only)
- [ ] Backend: Role-based access control for salary data
- [ ] Frontend: Conditional rendering of salary field based on user role
- [ ] Frontend: Admin badge/indicator for salary visibility

## Phase 5: Core Features - Comments

### 5.1 Comment System
- [ ] Backend: Add comment endpoint
- [ ] Backend: List comments for candidate endpoint
- [ ] Backend: Update comment endpoint
- [ ] Backend: Delete comment endpoint
- [ ] Frontend: Comments section in candidate detail
- [ ] Frontend: Add comment form
- [ ] Frontend: Comment list with timestamps and authors
- [ ] Frontend: Edit/delete comment actions (for comment owner)

## Phase 6: AI Features

### 6.1 AI Candidate Summary
- [ ] Backend: Integrate AI service (OpenAI API or similar)
- [ ] Backend: Generate candidate summary endpoint
- [ ] Backend: Cache AI summaries to reduce API costs
- [ ] Frontend: AI summary section in candidate detail
- [ ] Frontend: Regenerate summary button
- [ ] Frontend: Loading state for AI generation

### 6.2 AI Chat Interface
- [ ] Backend: AI chat endpoint (streaming response)
- [ ] Backend: Context building from job postings and candidates
- [ ] Frontend: Chat interface component
- [ ] Frontend: Chat history display
- [ ] Frontend: Message input with send button
- [ ] Frontend: Streaming response handling
- [ ] Frontend: Example questions/prompts

## Phase 7: Future Features

### 7.1 Filtering & Sorting
- [ ] Backend: Advanced filtering query builder
- [ ] Backend: Sorting parameters
- [ ] Frontend: Filter panel with multiple criteria
- [ ] Frontend: Sort controls on candidate list
- [ ] Frontend: Filter tags/chips display
- [ ] Frontend: Clear all filters button

### 7.2 CSV Export
- [ ] Backend: Export candidates to CSV endpoint
- [ ] Backend: CSV generation utility
- [ ] Frontend: Export button on candidates list
- [ ] Frontend: Export progress indicator
- [ ] Frontend: Download completed file

## Phase 8: DevOps & Deployment

### 8.1 Docker Configuration
- [ ] Create Dockerfile for frontend (multi-stage build)
- [ ] Create Dockerfile for backend
- [ ] Create production Dockerfile that bundles both (backend serves frontend static files)
- [ ] Optimize Docker images for size
- [ ] Create .dockerignore files

### 8.2 Docker Compose
- [ ] docker-compose.yml for local development
  - [ ] Postgres service
  - [ ] Backend service with hot reload
  - [ ] Frontend service with hot reload
  - [ ] Volume mounts for development
- [ ] docker-compose.prod.yml for production-like testing
- [ ] Environment variable templates (.env.example)

### 8.3 Documentation & Testing
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Setup instructions in README
- [ ] Environment variables documentation
- [ ] Backend unit tests for repositories
- [ ] Backend integration tests for API endpoints
- [ ] Frontend component tests
- [ ] E2E tests for critical flows

## Phase 9: Polish & Optimization

### 9.1 UI/UX Improvements
- [ ] Responsive design for mobile devices
- [ ] Loading states for all async operations
- [ ] Error handling and user-friendly error messages
- [ ] Toast notifications for actions
- [ ] Form validation with helpful feedback
- [ ] Keyboard shortcuts for power users
- [ ] Dark mode support (optional)

### 9.2 Performance Optimization
- [ ] Database indexing for common queries
- [ ] API response caching where appropriate
- [ ] Frontend code splitting
- [ ] Image optimization for resumes
- [ ] Lazy loading for large lists
- [ ] Database connection pooling optimization

### 9.3 Security Hardening
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention
- [ ] CSRF protection
- [ ] Rate limiting on API endpoints
- [ ] Input validation and sanitization
- [ ] Secure file upload handling
- [ ] Environment variable security audit
- [ ] HTTPS enforcement in production

## Current Status

**Phase 1: Project Setup & Infrastructure** - NOT STARTED

---

## Notes

- Prioritize MVP features first (Phases 1-5) before Future Features (Phase 7)
- AI features (Phase 6) require API key configuration and cost consideration
- All database operations use handwritten SQL (NO ORM as per requirements)
- Google Workspace domain restriction should be configurable via environment variable
- Resume parsing may require external service or library (consider open-source options)
