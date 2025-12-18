import { api } from '../client';
import { JobPosting } from '../types';

export interface CreateJobPostingData {
  title: string;
  description: string;
  requirements: string;
  location: string;
  salary_range: string;
  status: 'open' | 'closed' | 'draft';
}

export interface UpdateJobPostingData extends CreateJobPostingData {}

export interface JobsListResponse {
  jobs: JobPosting[];
  limit: number;
  offset: number;
}

export interface JobResponse {
  job: JobPosting;
  message?: string;
}

export const jobsApi = {
  create: (data: CreateJobPostingData, token?: string) =>
    api.post<JobResponse>('/api/v1/jobs', data, token),

  list: (limit: number = 20, offset: number = 0, token?: string) =>
    api.get<JobsListResponse>(
      `/api/v1/jobs?limit=${limit}&offset=${offset}`,
      token
    ),

  getById: (id: string, token?: string) =>
    api.get<JobResponse>(`/api/v1/jobs/${id}`, token),

  update: (id: string, data: UpdateJobPostingData, token?: string) =>
    api.put<JobResponse>(`/api/v1/jobs/${id}`, data, token),

  delete: (id: string, token?: string) =>
    api.delete<{ message: string }>(`/api/v1/jobs/${id}`, token),
};
