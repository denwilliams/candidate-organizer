'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter, useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { jobsApi } from '@/lib/api/services/jobs';
import { JobPosting } from '@/lib/api/types';
import { Button } from '@/components/ui/button';

// Required for static export - returns empty array since routes are dynamic
export function generateStaticParams() {
  return [];
}

export default function JobDetailPage() {
  const { user, isLoading, isAuthenticated } = useAuth();
  const router = useRouter();
  const params = useParams();
  const jobId = params.id as string;
  const [job, setJob] = useState<JobPosting | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isLoading, isAuthenticated, router]);

  useEffect(() => {
    if (isAuthenticated && jobId) {
      fetchJob();
    }
  }, [isAuthenticated, jobId]);

  const fetchJob = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await jobsApi.getById(jobId);
      setJob(response.job);
    } catch (err: any) {
      setError(err.message || 'Failed to load job posting');
      console.error('Error fetching job:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this job posting?')) {
      return;
    }

    try {
      await jobsApi.delete(jobId);
      router.push('/jobs');
    } catch (err: any) {
      alert(err.message || 'Failed to delete job posting');
      console.error('Error deleting job:', err);
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

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Error</h2>
          <p className="text-sm text-red-600 mb-6">{error}</p>
          <Button onClick={() => router.push('/jobs')} className="w-full">
            Back to Job Postings
          </Button>
        </div>
      </div>
    );
  }

  if (!job) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Job Not Found</h2>
          <p className="text-sm text-gray-600 mb-6">
            The job posting you're looking for doesn't exist or has been deleted.
          </p>
          <Button onClick={() => router.push('/jobs')} className="w-full">
            Back to Job Postings
          </Button>
        </div>
      </div>
    );
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
                  className="text-sm font-medium text-gray-600 hover:text-gray-900"
                >
                  Job Postings
                </button>
              </nav>
            </div>
            <div className="flex items-center gap-4">
              <Button onClick={() => router.push('/jobs')} variant="outline" size="sm">
                Back
              </Button>
            </div>
          </div>
        </div>
      </nav>

      <main className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white shadow rounded-lg overflow-hidden">
          <div className="p-6 border-b border-gray-200">
            <div className="flex justify-between items-start mb-4">
              <div className="flex-1">
                <div className="flex items-center gap-3 mb-2">
                  <h1 className="text-3xl font-bold text-gray-900">{job.title}</h1>
                  <span
                    className={`inline-flex items-center rounded-full px-3 py-1 text-sm font-medium ${getStatusBadgeColor(
                      job.status
                    )}`}
                  >
                    {job.status}
                  </span>
                </div>
                <div className="flex flex-wrap gap-4 text-sm text-gray-600">
                  {job.location && (
                    <span className="flex items-center gap-1">
                      <svg
                        className="h-5 w-5"
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
                        className="h-5 w-5"
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
                </div>
              </div>
              <div className="flex gap-2">
                <Button onClick={() => router.push(`/jobs/${job.id}/edit`)} variant="outline">
                  Edit
                </Button>
                <Button onClick={handleDelete} variant="outline">
                  Delete
                </Button>
              </div>
            </div>
          </div>

          <div className="p-6 space-y-6">
            {job.description && (
              <div>
                <h2 className="text-lg font-semibold text-gray-900 mb-3">Description</h2>
                <div className="text-gray-700 whitespace-pre-wrap">{job.description}</div>
              </div>
            )}

            {job.requirements && (
              <div>
                <h2 className="text-lg font-semibold text-gray-900 mb-3">Requirements</h2>
                <div className="text-gray-700 whitespace-pre-wrap">{job.requirements}</div>
              </div>
            )}

            <div className="pt-6 border-t border-gray-200">
              <div className="text-sm text-gray-500 space-y-1">
                <p>Created: {new Date(job.created_at).toLocaleString()}</p>
                <p>Last updated: {new Date(job.updated_at).toLocaleString()}</p>
                <p>Created by: {job.created_by}</p>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-6 bg-white shadow rounded-lg p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Candidates</h2>
          <p className="text-sm text-gray-600">
            Candidate management coming soon...
          </p>
        </div>
      </main>
    </div>
  );
}
