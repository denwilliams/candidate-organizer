#!/bin/bash
set -e

# Get schema name from environment variable, default to 'public'
SCHEMA_NAME="${POSTGRES_SCHEMA:-public}"

echo "Running database migrations in schema: $SCHEMA_NAME"

# Run the SQL migration with the search_path set to the target schema
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Set search_path to the target schema
    SET search_path TO $SCHEMA_NAME;

    -- Show current search path for verification
    SHOW search_path;

    -- Enable UUID extension (in public schema so it's available everywhere)
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA public;

    -- Users table
    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        email VARCHAR(255) NOT NULL UNIQUE,
        name VARCHAR(255) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'user', -- 'admin' or 'user'
        workspace_domain VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX idx_users_email ON users(email);
    CREATE INDEX idx_users_role ON users(role);

    -- Job postings table
    CREATE TABLE job_postings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        requirements TEXT,
        location VARCHAR(255),
        salary_range VARCHAR(100),
        status VARCHAR(50) NOT NULL DEFAULT 'draft', -- 'draft', 'open', 'closed'
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX idx_job_postings_status ON job_postings(status);
    CREATE INDEX idx_job_postings_created_by ON job_postings(created_by);
    CREATE INDEX idx_job_postings_created_at ON job_postings(created_at DESC);

    -- Candidates table
    CREATE TABLE candidates (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255),
        phone VARCHAR(50),
        resume_url VARCHAR(500),
        parsed_data JSONB,
        status VARCHAR(50) NOT NULL DEFAULT 'applied', -- 'applied', 'screened', 'interviewing', 'offered', 'rejected'
        salary_expectation VARCHAR(100),
        job_posting_id UUID REFERENCES job_postings(id) ON DELETE SET NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX idx_candidates_status ON candidates(status);
    CREATE INDEX idx_candidates_job_posting_id ON candidates(job_posting_id);
    CREATE INDEX idx_candidates_created_by ON candidates(created_by);
    CREATE INDEX idx_candidates_email ON candidates(email);
    CREATE INDEX idx_candidates_created_at ON candidates(created_at DESC);

    -- Comments table
    CREATE TABLE comments (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        content TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX idx_comments_candidate_id ON comments(candidate_id);
    CREATE INDEX idx_comments_user_id ON comments(user_id);
    CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

    -- Candidate attributes table
    CREATE TABLE candidate_attributes (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
        attribute_key VARCHAR(255) NOT NULL,
        attribute_value TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(candidate_id, attribute_key)
    );

    CREATE INDEX idx_candidate_attributes_candidate_id ON candidate_attributes(candidate_id);
    CREATE INDEX idx_candidate_attributes_key ON candidate_attributes(attribute_key);

    -- AI summaries cache table (optional, for caching AI-generated summaries)
    CREATE TABLE ai_summaries (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
        job_posting_id UUID REFERENCES job_postings(id) ON DELETE CASCADE,
        summary TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(candidate_id, job_posting_id)
    );

    CREATE INDEX idx_ai_summaries_candidate_id ON ai_summaries(candidate_id);

    -- Function to update updated_at timestamp
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS \$\$
    BEGIN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
    \$\$ language 'plpgsql';

    -- Triggers to automatically update updated_at
    CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

    CREATE TRIGGER update_job_postings_updated_at BEFORE UPDATE ON job_postings
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

    CREATE TRIGGER update_candidates_updated_at BEFORE UPDATE ON candidates
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

    CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

    CREATE TRIGGER update_candidate_attributes_updated_at BEFORE UPDATE ON candidate_attributes
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

    CREATE TRIGGER update_ai_summaries_updated_at BEFORE UPDATE ON ai_summaries
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
EOSQL

echo "Database migration completed successfully in schema: $SCHEMA_NAME"
