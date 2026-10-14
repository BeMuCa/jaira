import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// The build is committed (go:embed reads the disk, and 'go install' has no
// Node), and CI rebuilds it and diffs the result. Nothing in it may depend on
// the machine or the clock, which Vite's content hashes already guarantee.
export default defineConfig({
  plugins: [svelte()],
  base: './',
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  build: { outDir: 'dist', emptyOutDir: true },
  server: {
    // 'npm run dev' against a running 'jaira web --port 7420'.
    proxy: { '/api': 'http://127.0.0.1:7420' },
  },
});
