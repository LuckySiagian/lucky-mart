import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

const apiTarget = process.env.API_PROXY_TARGET || 'http://localhost:8080';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3002,
    proxy: {
      '/api': { target: apiTarget, changeOrigin: true },
      // Uploaded photos are served by the backend.
      '/uploads': { target: apiTarget, changeOrigin: true },
    },
  },
});
