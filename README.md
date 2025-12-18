# candidate-organizer

Web application to help organize and triage job candidates.

- Ad job postings with relevant details
- Add candidates by uploading resumes (PDF) or manual entry
- Parse resumes and extract key information (name, contact info, skills, experience)
- Allow custom setting of candidate attributes that cant be parsed from resumes
- Set current status of candidates (e.g., Applied, Screened, Interviewing, Offered, Rejected)
- Allow users to add multiple comments per candidate for notes during the hiring process
- AI summary of each candidate that highlights key overlaps with job requirements as well as potential gaps, concerns, or red flags
- AI chat to discuss candidates and job postings, eg "What are the top 3 candidates for this job posting?"
- Google Authentication locked to a single Google Workspace domain for secure access
- Salary Expectation Tracking for each candidate, with this data hidden from most users and only visible to admins
- First user to sign up becomes the admin, with ability to promote other users to admin status
- It is expected that each "organization" will have their own installation of the app, so no multi-tenancy is required

Future Features:
- Filter and sort candidates based on various criteria
- Export candidate data to CSV for reporting or sharing
- Embeddings generation for candidates and job postings for advanced search and matching capabilities

## Tech Stack

- React with Next js for Frontend
- Shadcn UI for component library
- Backend: Golang
- Postgres for Database
- NO ORM, create manual repository interfaces and have concrete implementations for Postgres with handwritten SQL queries
- Bundle into single Docker container for easy deployment where the React project is served via the Golang backend
- When developing locally allow running React frontend and Golang backend separately for faster development cycle
- Use Docker Compose to orchestrate local development environment with Postgres and backend/frontend services

## LLM Integration

- Use OpenAI GPT-5.2 models for AI features
- Use streaming API for chat feature to allow real-time responses
