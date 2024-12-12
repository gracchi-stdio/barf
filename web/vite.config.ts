import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    host: '0.0.0.0', // docker; all interfaces
    strictPort: true,
    port: 3000,
    watch: {
      usePolling: true,
    }
  },
  build: {
    target: 'esnext',
  },
});
