import { HttpResponse, http } from "msw";

import { env } from "@/config/env";

import { networkDelay } from "../utils";

import { sketchHandlers } from "./sketches";

export const handlers = [
  ...sketchHandlers,
  http.get(`${env.API_URL}/healthcheck`, async () => {
    await networkDelay();
    return HttpResponse.json({ ok: true });
  }),
];
