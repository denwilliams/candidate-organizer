'use client';

import { Suspense, useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { authService } from '@/lib/api/services/auth';

function AuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleCallback = async () => {
      // Check for error in URL params
      const errorParam = searchParams.get('error');
      if (errorParam) {
        setError(decodeURIComponent(errorParam));
        return;
      }

      // Check for success
      const success = searchParams.get('success');
      if (success === 'true') {
        try {
          // Get user profile from the backend (auth cookie should be set)
          const authData = await authService.getAuthFromCookies();

          if (authData) {
            // Login with the user data
            // Note: Since we're using HTTP-only cookies, the token is managed by the browser
            // We'll use a placeholder token in localStorage, but actual auth is via cookies
            login('cookie-auth', authData.user);

            // Redirect to dashboard
            router.push('/dashboard');
          } else {
            setError('Failed to authenticate. Please try again.');
          }
        } catch (err) {
          console.error('Authentication error:', err);
          setError('An error occurred during authentication. Please try again.');
        }
      }
    };

    handleCallback();
  }, [searchParams, login, router]);

  if (error) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center">
        <div className="w-full max-w-md space-y-8 rounded-xl bg-white p-10 shadow-lg">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-red-600">Authentication Error</h1>
            <p className="mt-4 text-gray-600">{error}</p>
            <button
              onClick={() => router.push('/login')}
              className="mt-6 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              Back to Login
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent motion-reduce:animate-[spin_1.5s_linear_infinite]" />
        <p className="mt-4 text-lg">Completing authentication...</p>
      </div>
    </div>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center">
          <div className="text-center">
            <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent motion-reduce:animate-[spin_1.5s_linear_infinite]" />
            <p className="mt-4 text-lg">Loading...</p>
          </div>
        </div>
      }
    >
      <AuthCallbackContent />
    </Suspense>
  );
}
