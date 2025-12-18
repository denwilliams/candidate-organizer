import { api, apiRequest } from '../client';
import { User } from '../types';

export interface LoginResponse {
  token: string;
  user: User;
}

export const authService = {
  /**
   * Get the OAuth URL to redirect to Google login
   */
  getGoogleAuthUrl: (): string => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    return `${apiUrl}/api/v1/auth/google`;
  },

  /**
   * Get current user profile
   */
  getProfile: (token?: string): Promise<User> => {
    return api.get<User>('/api/v1/auth/me', token);
  },

  /**
   * Refresh authentication token
   */
  refreshToken: (token?: string): Promise<{ token: string }> => {
    return api.post<{ token: string }>('/api/v1/auth/refresh', {}, token);
  },

  /**
   * Logout the current user
   */
  logout: (token?: string): Promise<void> => {
    return api.post<void>('/api/v1/auth/logout', {}, token);
  },

  /**
   * Get auth token and user info from cookies (called after OAuth callback)
   * Note: This requires cookies to be set by the backend
   */
  getAuthFromCookies: async (): Promise<{ token: string; user: User } | null> => {
    try {
      // The backend sets the auth_token in an HTTP-only cookie
      // We need to make a request to get the user profile
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/auth/me`, {
        credentials: 'include', // Important: include cookies
      });

      if (!response.ok) {
        return null;
      }

      const user = await response.json();

      // Extract token from cookie if needed (though we'll rely on cookies for auth)
      // For local storage, we might need to get it differently
      return { token: '', user }; // Token will be in cookie
    } catch (error) {
      console.error('Failed to get auth from cookies:', error);
      return null;
    }
  },
};
