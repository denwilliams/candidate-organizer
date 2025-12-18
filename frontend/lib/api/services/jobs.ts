import { api } from '../client';
import { JobPosting, PaginatedResponse } from '../types';

export interface CreateJobPostingData {
  title: string;
  description: string;
  requirements: string;
  location: string;
  salary_range: string;
  status: 'open' | 'closed' | 'draft';
}

export interface UpdateJobPostingData extends Partial<CreateJobPostingData> {}

export const jobsApi = {
  create: (data: CreateJobPostingData, token: string) =>
    api.post<JobPosting>('/api/jobs', data, token),

  list: (page: number = 1, pageSize: number = 20, token: string) =>
    api.get<PaginatedResponse<JobPosting>>(
      `/api/jobs?page=${page}&page_size=${pageSize}`,
      token
    ),

  getById: (id: string, token: string) =>
    api.get<JobPosting>(`/api/jobs/${id}`, token),

  update: (id: string, data: UpdateJobPostingData, token: string) =>
    api.put<JobPosting>(`/api/jobs/${id}`, data, token),

  delete: (id: string, token: string) =>
    api.delete(`/api/jobs/${id}`, token),
};
