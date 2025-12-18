'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { usersApi } from '@/lib/api/services/users';
import { User } from '@/lib/api/types';
import { Button } from '@/components/ui/button';

export default function UsersPage() {
  const { user, isLoading: authLoading, isAuthenticated } = useAuth();
  const router = useRouter();
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [promotingUserId, setPromotingUserId] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [authLoading, isAuthenticated, router]);

  useEffect(() => {
    if (!authLoading && isAuthenticated && user?.role !== 'admin') {
      router.push('/dashboard');
    }
  }, [authLoading, isAuthenticated, user, router]);

  useEffect(() => {
    if (isAuthenticated && user?.role === 'admin') {
      fetchUsers();
    }
  }, [isAuthenticated, user]);

  const fetchUsers = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await usersApi.listAll();
      setUsers(response.users);
    } catch (err: any) {
      console.error('Error fetching users:', err);
      setError(err.message || 'Failed to fetch users');
    } finally {
      setIsLoading(false);
    }
  };

  const handlePromoteToAdmin = async (userId: string) => {
    try {
      setPromotingUserId(userId);
      setError(null);
      await usersApi.promoteToAdmin(userId);

      // Refresh the users list
      await fetchUsers();
    } catch (err: any) {
      console.error('Error promoting user:', err);
      setError(err.message || 'Failed to promote user');
    } finally {
      setPromotingUserId(null);
    }
  };

  if (authLoading || isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-lg">Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated || !user || user.role !== 'admin') {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 justify-between items-center">
            <div className="flex items-center">
              <h1 className="text-xl font-bold">Candidate Organizer</h1>
            </div>
            <div className="flex items-center gap-4">
              <Button
                onClick={() => router.push('/dashboard')}
                variant="outline"
                size="sm"
              >
                Back to Dashboard
              </Button>
            </div>
          </div>
        </div>
      </nav>

      <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900">User Management</h2>
          <p className="mt-1 text-sm text-gray-500">
            Manage users and their roles
          </p>
        </div>

        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 rounded-md p-4">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg
                  className="h-5 w-5 text-red-400"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"
                  />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">{error}</h3>
              </div>
            </div>
          </div>
        )}

        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {users.map((userItem) => (
              <li key={userItem.id}>
                <div className="px-4 py-4 flex items-center sm:px-6">
                  <div className="min-w-0 flex-1 sm:flex sm:items-center sm:justify-between">
                    <div className="truncate">
                      <div className="flex text-sm">
                        <p className="font-medium text-gray-900 truncate">
                          {userItem.name}
                        </p>
                        <p className="ml-4 flex-shrink-0 text-gray-500">
                          {userItem.email}
                        </p>
                      </div>
                      <div className="mt-2 flex">
                        <div className="flex items-center text-sm text-gray-500">
                          <span>
                            Joined {new Date(userItem.created_at).toLocaleDateString()}
                          </span>
                        </div>
                      </div>
                    </div>
                    <div className="mt-4 flex-shrink-0 sm:mt-0 sm:ml-5">
                      <div className="flex items-center gap-3">
                        {userItem.role === 'admin' ? (
                          <span className="inline-flex items-center rounded-full bg-blue-100 px-3 py-1 text-xs font-medium text-blue-800">
                            Admin
                          </span>
                        ) : (
                          <>
                            <span className="inline-flex items-center rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-800">
                              User
                            </span>
                            <Button
                              onClick={() => handlePromoteToAdmin(userItem.id)}
                              disabled={promotingUserId === userItem.id}
                              size="sm"
                              variant="outline"
                            >
                              {promotingUserId === userItem.id ? (
                                'Promoting...'
                              ) : (
                                'Promote to Admin'
                              )}
                            </Button>
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>

        {users.length === 0 && !isLoading && (
          <div className="text-center py-12">
            <p className="text-gray-500">No users found</p>
          </div>
        )}
      </main>
    </div>
  );
}
