import { api } from '../client';
import { CandidateAttribute } from '../types';

export interface CreateAttributeData {
  candidate_id: string;
  attribute_key: string;
  attribute_value: string;
}

export interface UpdateAttributeData {
  attribute_key: string;
  attribute_value: string;
}

export const attributesApi = {
  create: (data: CreateAttributeData, token: string) =>
    api.post<CandidateAttribute>('/api/attributes', data, token),

  listByCandidate: (candidateId: string, token: string) =>
    api.get<CandidateAttribute[]>(
      `/api/candidates/${candidateId}/attributes`,
      token
    ),

  update: (id: string, data: UpdateAttributeData, token: string) =>
    api.put<CandidateAttribute>(`/api/attributes/${id}`, data, token),

  delete: (id: string, token: string) =>
    api.delete(`/api/attributes/${id}`, token),
};
