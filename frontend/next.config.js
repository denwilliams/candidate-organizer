/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'export',
  // Disable image optimization for static export
  images: {
    unoptimized: true,
  },
  // Remove rewrites since API will be on same origin
  basePath: '',
  trailingSlash: true,
};

module.exports = nextConfig;
