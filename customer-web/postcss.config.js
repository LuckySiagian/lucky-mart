import { fileURLToPath } from 'node:url';

// Point Tailwind at this folder's config explicitly so styles still build when
// Vite is started from a parent directory (e.g. `vite customer-web`).
export default {
  plugins: {
    tailwindcss: { config: fileURLToPath(new URL('./tailwind.config.js', import.meta.url)) },
    autoprefixer: {},
  },
}
