import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import viteTsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  base: "./",
  plugins: [react(), tailwindcss(), viteTsconfigPaths()],
  server: {
    proxy: {
      "/api/v1": "http://localhost:8080",
    },
  },
});
