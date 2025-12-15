import { env } from "@/config/env";
import { seedSketches } from "./seed/sketches";

export const enableMocking = async () => {
  if (env.ENABLE_API_MOCKING) {
    const { worker } = await import("./browser");
    const { initializeDb } = await import("./db");
    await initializeDb();
    seedSketches();
    return worker.start();
  }
};
