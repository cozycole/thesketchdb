import { env } from "@/config/env";

import { seedSketches } from "./seed/sketches";
import { seedCreators } from "./seed/creators";
import { seedShows } from "./seed/shows";

export const enableMocking = async () => {
  if (env.ENABLE_API_MOCKING) {
    const { worker } = await import("./browser");
    const { initializeDb } = await import("./db");
    await initializeDb();
    await seedCreators();
    await seedShows();
    await seedSketches();
    return worker.start();
  }
};
