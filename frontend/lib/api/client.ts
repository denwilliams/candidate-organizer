const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface ApiError {
  error: string;
  message?: string;
  code: number;
}

export class ApiRequestError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    public data?: ApiError
  ) {
    super(data?.message || statusText);
    this.name = 'ApiRequestError';
  }
}

interface RequestOptions extends RequestInit {
  token?: string;
}

async function handleResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get('content-type');
  const isJson = contentType?.includes('application/json');

  if (!response.ok) {
    const errorData = isJson ? await response.json() : null;
    throw new ApiRequestError(response.status, response.statusText, errorData);
  }

  if (response.status === 204) {
    return {} as T;
  }

  if (!isJson) {
    throw new Error('Expected JSON response from API');
  }

  const data = await response.json();
  return data.data || data;
}

export async function apiRequest<T = any>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> {
  const { token, ...fetchOptions } = options;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(fetchOptions.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const url = `${API_BASE_URL}${endpoint}`;

  const response = await fetch(url, {
    ...fetchOptions,
    headers,
    credentials: 'include', // Important: include cookies for authentication
  });

  return handleResponse<T>(response);
}

export const api = {
  get: <T = any>(endpoint: string, token?: string) =>
    apiRequest<T>(endpoint, { method: 'GET', token }),

  post: <T = any>(endpoint: string, data?: any, token?: string) =>
    apiRequest<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
      token,
    }),

  put: <T = any>(endpoint: string, data?: any, token?: string) =>
    apiRequest<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
      token,
    }),

  patch: <T = any>(endpoint: string, data?: any, token?: string) =>
    apiRequest<T>(endpoint, {
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
      token,
    }),

  delete: <T = any>(endpoint: string, token?: string) =>
    apiRequest<T>(endpoint, { method: 'DELETE', token }),
};
