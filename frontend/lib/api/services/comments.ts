import { api } from '../client';
import { Comment } from '../types';

export interface CreateCommentData {
  candidate_id: string;
  content: string;
}

export interface UpdateCommentData {
  content: string;
}

export const commentsApi = {
  create: (data: CreateCommentData, token: string) =>
    api.post<Comment>('/api/comments', data, token),

  listByCandidate: (candidateId: string, token: string) =>
    api.get<Comment[]>(`/api/candidates/${candidateId}/comments`, token),

  update: (id: string, data: UpdateCommentData, token: string) =>
    api.put<Comment>(`/api/comments/${id}`, data, token),

  delete: (id: string, token: string) =>
    api.delete(`/api/comments/${id}`, token),
};
