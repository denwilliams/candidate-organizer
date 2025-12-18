import { api } from '../client';
import { Candidate, PaginatedResponse } from '../types';

export interface CreateCandidateData {
  name: string;
  email: string;
  phone: string;
  resume_url?: string;
  parsed_data?: Record<string, any>;
  status: 'applied' | 'screened' | 'interviewing' | 'offered' | 'rejected';
  salary_expectation?: string;
  job_posting_id: string;
}

export interface UpdateCandidateData extends Partial<CreateCandidateData> {}

export interface CandidateFilters {
  status?: string;
  job_posting_id?: string;
  search?: string;
}

export const candidatesApi = {
  create: (data: CreateCandidateData, token: string) =>
    api.post<Candidate>('/api/candidates', data, token),

  list: (
    page: number = 1,
    pageSize: number = 20,
    filters: CandidateFilters = {},
    token: string
  ) => {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
      ...filters,
    });
    return api.get<PaginatedResponse<Candidate>>(
      `/api/candidates?${params.toString()}`,
      token
    );
  },

  getById: (id: string, token: string) =>
    api.get<Candidate>(`/api/candidates/${id}`, token),

  update: (id: string, data: UpdateCandidateData, token: string) =>
    api.put<Candidate>(`/api/candidates/${id}`, data, token),

  updateStatus: (id: string, status: string, token: string) =>
    api.patch<Candidate>(`/api/candidates/${id}/status`, { status }, token),

  delete: (id: string, token: string) =>
    api.delete(`/api/candidates/${id}`, token),

  uploadResume: (candidateId: string, file: File, token: string) => {
    const formData = new FormData();
    formData.append('resume', file);
    return fetch(
      `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/candidates/${candidateId}/resume`,
      {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      }
    );
  },
};
