'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { jobsApi } from '@/lib/api/services/jobs';
import { JobPosting } from '@/lib/api/types';
import { Button } from '@/components/ui/button';

export default function JobsPage() {
  const { user, isLoading, isAuthenticated } = useAuth();
  const router = useRouter();
  const [jobs, setJobs] = useState<JobPosting[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isLoading, isAuthenticated, router]);

  useEffect(() => {
    if (isAuthenticated) {
      fetchJobs();
    }
  }, [isAuthenticated]);

  const fetchJobs = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await jobsApi.list();
      setJobs(response.jobs || []);
    } catch (err: any) {
      setError(err.message || 'Failed to load job postings');
      console.error('Error fetching jobs:', err);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadgeColor = (status: string) => {
    switch (status) {
      case 'open':
        return 'bg-green-100 text-green-800';
      case 'closed':
        return 'bg-red-100 text-red-800';
      case 'draft':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (isLoading || loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-lg">Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated || !user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 justify-between items-center">
            <div className="flex items-center gap-8">
              <h1 className="text-xl font-bold cursor-pointer" onClick={() => router.push('/dashboard')}>
                Candidate Organizer
              </h1>
              <nav className="flex gap-4">
                <button
                  onClick={() => router.push('/jobs')}
                  className="text-sm font-medium text-blue-600"
                >
                  Job Postings
                </button>
              </nav>
            </div>
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-700">{user.name}</span>
              <Button onClick={() => router.push('/dashboard')} variant="outline" size="sm">
                Dashboard
              </Button>
            </div>
          </div>
        </div>
      </nav>

      <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
        <div className="sm:flex sm:items-center sm:justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Job Postings</h1>
            <p className="mt-2 text-sm text-gray-700">
              Manage your job postings and track candidates
            </p>
          </div>
          <div className="mt-4 sm:mt-0">
            <Button onClick={() => router.push('/jobs/new')}>
              Create Job Posting
            </Button>
          </div>
        </div>

        {error && (
          <div className="mb-6 rounded-md bg-red-50 p-4">
            <div className="flex">
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">Error loading job postings</h3>
                <div className="mt-2 text-sm text-red-700">
                  <p>{error}</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {!error && jobs.length === 0 && (
          <div className="text-center py-12 bg-white rounded-lg shadow">
            <svg
              className="mx-auto h-12 w-12 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
            <h3 className="mt-2 text-sm font-semibold text-gray-900">No job postings</h3>
            <p className="mt-1 text-sm text-gray-500">Get started by creating a new job posting.</p>
            <div className="mt-6">
              <Button onClick={() => router.push('/jobs/new')}>
                Create Job Posting
              </Button>
            </div>
          </div>
        )}

        {!error && jobs.length > 0 && (
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
              {jobs.map((job) => (
                <li key={job.id}>
                  <div className="px-4 py-4 sm:px-6 hover:bg-gray-50">
                    <div className="flex items-center justify-between">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-3">
                          <h3 className="text-lg font-medium text-gray-900 truncate">
                            {job.title}
                          </h3>
                          <span
                            className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${getStatusBadgeColor(
                              job.status
                            )}`}
                          >
                            {job.status}
                          </span>
                        </div>
                        <div className="mt-2 flex items-center gap-4 text-sm text-gray-500">
                          {job.location && (
                            <span className="flex items-center gap-1">
                              <svg
                                className="h-4 w-4"
                                fill="none"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                              >
                                <path d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                <path d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                              </svg>
                              {job.location}
                            </span>
                          )}
                          {job.salary_range && (
                            <span className="flex items-center gap-1">
                              <svg
                                className="h-4 w-4"
                                fill="none"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                              >
                                <path d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                              </svg>
                              {job.salary_range}
                            </span>
                          )}
                          <span>
                            Created {new Date(job.created_at).toLocaleDateString()}
                          </span>
                        </div>
                        {job.description && (
                          <p className="mt-2 text-sm text-gray-600 line-clamp-2">
                            {job.description}
                          </p>
                        )}
                      </div>
                      <div className="ml-5 flex-shrink-0">
                        <p className="text-sm text-gray-500">View details coming soon</p>
                      </div>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}
      </main>
    </div>
  );
}
