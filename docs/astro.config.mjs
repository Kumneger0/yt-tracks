import { defineConfig } from 'astro/config';
import vercel from '@astrojs/vercel';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  site: "https://yt-tracks.kumneger.dev",
  output: "static",
  vite: {
    plugins: [tailwindcss()],
    server: {
      allowedHosts: [
        "4ddf-2a09-bac5-3215-2c5a-00-46b-37.ngrok-free.app",
      ],
    },
  },
  adapter: vercel({
    webAnalytics: {
      enabled: true,
    },
  }),
});
