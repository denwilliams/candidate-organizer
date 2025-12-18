import { api } from '../client';
import { User } from '../types';

export const usersApi = {
  getProfile: (token: string) =>
    api.get<User>('/api/users/me', token),

  listAll: (token: string) =>
    api.get<User[]>('/api/users', token),

  promoteToAdmin: (userId: string, token: string) =>
    api.post(`/api/users/${userId}/promote`, undefined, token),
};
