package models

import "time"

// User represents a user in the system
type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Role            string    `json:"role"` // "admin" or "user"
	WorkspaceDomain string    `json:"workspace_domain"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// JobPosting represents a job posting
type JobPosting struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Requirements  string    `json:"requirements"`
	Location      string    `json:"location"`
	SalaryRange   string    `json:"salary_range"`
	Status        string    `json:"status"` // "open", "closed", "draft"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     string    `json:"created_by"`
}

// Candidate represents a job candidate
type Candidate struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Email             string            `json:"email"`
	Phone             string            `json:"phone"`
	ResumeURL         string            `json:"resume_url"`
	ParsedData        map[string]interface{} `json:"parsed_data"`
	Status            string            `json:"status"` // "applied", "screened", "interviewing", "offered", "rejected"
	SalaryExpectation string            `json:"salary_expectation,omitempty"` // Only visible to admins
	JobPostingID      string            `json:"job_posting_id"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	CreatedBy         string            `json:"created_by"`
}

// Comment represents a comment on a candidate
type Comment struct {
	ID          string    `json:"id"`
	CandidateID string    `json:"candidate_id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"` // Denormalized for convenience
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CandidateAttribute represents a custom attribute for a candidate
type CandidateAttribute struct {
	ID             string    `json:"id"`
	CandidateID    string    `json:"candidate_id"`
	AttributeKey   string    `json:"attribute_key"`
	AttributeValue string    `json:"attribute_value"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
