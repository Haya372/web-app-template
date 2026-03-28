import tailwindcss from '@tailwindcss/vite'
import viteReact from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  plugins: [
    tsconfigPaths({ projects: ['./tsconfig.json'] }),
    tailwindcss(),
    viteReact(),
  ],
  server: {
    // biome-ignore lint/style/noProcessEnv: vite.config.ts runs in Node.js; VITE_PORT is set by mise from .env.local
    port: parseInt(process.env.VITE_PORT ?? '3000', 10),
  },
})
