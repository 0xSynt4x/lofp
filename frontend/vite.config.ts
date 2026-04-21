import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 4992,
    host: '0.0.0.0',
    allowedHosts: ['mtlgp1k93w-4992.cnb.run', 'mtlgp1k93w-4993.cnb.run', 'localhost'],
    proxy: {
      '/api': 'http://localhost:4993',
      '/healthz': 'http://localhost:4993',
      '/ws': {
        target: 'http://localhost:4993',
        ws: true,
      },
    },
  },
})
