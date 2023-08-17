/** @type {import('next').NextConfig} */
const path = require('path')

const nextConfig = {
  output: 'export',
  distDir: process.env.NODE_ENV == 'production' ? '../backend/api/frontend/out': 'out',
  images: {
    unoptimized: true,
  },
  
  reactStrictMode: true,
  swcMinify: true,
  sassOptions: {
    includePaths: [path.join(__dirname, 'src/styles')],
  },
  // typescript: {
  //   // !! WARN !!
  //   // Dangerously allow production builds to successfully complete even if //
  //   // your project has type errors.
  //   // !! WARN !!
  //   ignoreBuildErrors: true,
  // },
}

module.exports = nextConfig
