import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import viteTsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig(({ command }) => ({
  base: command === "build" ? "/admin/" : "/",
  plugins: [react(), tailwindcss(), viteTsconfigPaths()],
  test: {
    globals: true,
    environment: "jsdom",
  },
  exclude: ["**/testing/**", "**/*.test.*", "**/*.spec.*", "**/mocks/**"],
  server: {
    proxy: {
      "/api/v1": "http://localhost:8080",
    },
  },
}));
