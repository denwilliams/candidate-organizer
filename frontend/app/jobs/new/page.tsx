'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { jobsApi, CreateJobPostingData } from '@/lib/api/services/jobs';
import { Button } from '@/components/ui/button';

export default function NewJobPage() {
  const { user, isLoading, isAuthenticated } = useAuth();
  const router = useRouter();
  const [formData, setFormData] = useState<CreateJobPostingData>({
    title: '',
    description: '',
    requirements: '',
    location: '',
    salary_range: '',
    status: 'draft',
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isLoading, isAuthenticated, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.title || !formData.description) {
      setError('Title and description are required');
      return;
    }

    try {
      setSubmitting(true);
      setError(null);
      await jobsApi.create(formData);
      router.push('/jobs');
    } catch (err: any) {
      setError(err.message || 'Failed to create job posting');
      console.error('Error creating job:', err);
    } finally {
      setSubmitting(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  if (isLoading) {
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
                  className="text-sm font-medium text-gray-600 hover:text-gray-900"
                >
                  Job Postings
                </button>
              </nav>
            </div>
            <div className="flex items-center gap-4">
              <Button onClick={() => router.push('/jobs')} variant="outline" size="sm">
                Cancel
              </Button>
            </div>
          </div>
        </div>
      </nav>

      <main className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white shadow rounded-lg">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-2xl font-bold text-gray-900">Create Job Posting</h2>
          </div>

          <form onSubmit={handleSubmit} className="p-6 space-y-6">
            {error && (
              <div className="rounded-md bg-red-50 p-4">
                <div className="text-sm text-red-700">{error}</div>
              </div>
            )}

            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                Job Title *
              </label>
              <input
                type="text"
                id="title"
                name="title"
                value={formData.title}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="e.g. Senior Software Engineer"
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
                Description *
              </label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                required
                rows={6}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Describe the role, responsibilities, and what you're looking for..."
              />
            </div>

            <div>
              <label htmlFor="requirements" className="block text-sm font-medium text-gray-700 mb-2">
                Requirements
              </label>
              <textarea
                id="requirements"
                name="requirements"
                value={formData.requirements}
                onChange={handleChange}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="List the required skills, experience, and qualifications..."
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="location" className="block text-sm font-medium text-gray-700 mb-2">
                  Location
                </label>
                <input
                  type="text"
                  id="location"
                  name="location"
                  value={formData.location}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="e.g. Remote, San Francisco, CA"
                />
              </div>

              <div>
                <label htmlFor="salary_range" className="block text-sm font-medium text-gray-700 mb-2">
                  Salary Range
                </label>
                <input
                  type="text"
                  id="salary_range"
                  name="salary_range"
                  value={formData.salary_range}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="e.g. $120k - $180k"
                />
              </div>
            </div>

            <div>
              <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-2">
                Status
              </label>
              <select
                id="status"
                name="status"
                value={formData.status}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="draft">Draft</option>
                <option value="open">Open</option>
                <option value="closed">Closed</option>
              </select>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
              <Button
                type="button"
                onClick={() => router.push('/jobs')}
                variant="outline"
                disabled={submitting}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submitting}>
                {submitting ? 'Creating...' : 'Create Job Posting'}
              </Button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
