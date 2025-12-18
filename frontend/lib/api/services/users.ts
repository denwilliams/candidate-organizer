import { api } from '../client';
import { User } from '../types';

export const usersApi = {
  getProfile: (token: string) =>
    api.get<User>('/api/v1/auth/me', token),

  listAll: () =>
    api.get<{ users: User[] }>('/api/v1/users'),

  promoteToAdmin: (userId: string) =>
    api.post<{ message: string; user: User }>(`/api/v1/users/${userId}/promote`, undefined),
};
