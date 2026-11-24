import { defineConfig } from 'astro/config';
import vercel from '@astrojs/vercel';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  site: "https://yt-tracks.kumneger.dev",
  output: "static",
  vite: {
    plugins: [tailwindcss()],
  },
  adapter: vercel({
    webAnalytics: {
      enabled: true,
    },
  }),
});
