/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  // Optimize images (can enable optimization with standalone)
  images: {
    unoptimized: true, // Keep unoptimized for now, can enable later
  },
  // Proxy API requests to Go backend
  async rewrites() {
    return [
      {
        source: '/api/v1/:path*',
        destination: 'http://localhost:8080/api/v1/:path*',
      },
    ];
  },
};

module.exports = nextConfig;
