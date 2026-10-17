import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vitest/config'

// Khi phát triển, /api/<module> được chuyển tới dịch vụ chạy trên máy, giống định tuyến của Traefik trên cụm
// (CR-DOC-02 mục 4).
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  server: {
    port: 5173,
    strictPort: true,
    proxy: {
      '/api/runtime': { target: 'http://localhost:8081', rewrite: (p) => p.replace(/^\/api\/runtime/, ''), ws: true },
      '/api/content': { target: 'http://localhost:8082', rewrite: (p) => p.replace(/^\/api\/content/, '') },
      '/api/learning': { target: 'http://localhost:8083', rewrite: (p) => p.replace(/^\/api\/learning/, '') },
    },
  },
  test: {
    environment: 'jsdom',
  },
})
