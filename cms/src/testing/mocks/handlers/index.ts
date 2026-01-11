import { HttpResponse, http } from "msw";

import { env } from "@/config/env";

import { networkDelay } from "../utils";

import { sketchHandlers } from "./sketches";
import { creatorHandlers } from "./creators";
import { showHandlers } from "./shows";

export const handlers = [
  ...creatorHandlers,
  ...showHandlers,
  ...sketchHandlers,
  http.get(`${env.API_URL}/healthcheck`, async () => {
    await networkDelay();
    return HttpResponse.json({ ok: true });
  }),
];
