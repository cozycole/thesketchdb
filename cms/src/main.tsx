import * as React from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./index.css";

const root = document.getElementById("root");
if (!root) throw new Error("No root element found");

async function prepare() {
  // Only enable mocking in development
  if (import.meta.env.DEV) {
    const { enableMocking } = await import("./testing/mocks");
    return enableMocking();
  }
  return Promise.resolve();
}

prepare().then(() => {
  createRoot(root).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  );
});
