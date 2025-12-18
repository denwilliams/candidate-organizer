export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'user';
  workspace_domain: string;
  created_at: string;
  updated_at: string;
}

export interface JobPosting {
  id: string;
  title: string;
  description: string;
  requirements: string;
  location: string;
  salary_range: string;
  status: 'open' | 'closed' | 'draft';
  created_at: string;
  updated_at: string;
  created_by: string;
}

export interface Candidate {
  id: string;
  name: string;
  email: string;
  phone: string;
  resume_url: string;
  parsed_data: Record<string, any>;
  status: 'applied' | 'screened' | 'interviewing' | 'offered' | 'rejected';
  salary_expectation?: string;
  job_posting_id: string;
  created_at: string;
  updated_at: string;
  created_by: string;
}

export interface Comment {
  id: string;
  candidate_id: string;
  user_id: string;
  user_name: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface CandidateAttribute {
  id: string;
  candidate_id: string;
  attribute_key: string;
  attribute_value: string;
  created_at: string;
  updated_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
}
