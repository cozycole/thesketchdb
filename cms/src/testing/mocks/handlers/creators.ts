import { HttpResponse, http } from "msw";

import { env } from "@/config/env";

import { db } from "../db";
import { networkDelay } from "../utils";

export const creatorHandlers = [
  http.get(`${env.API_URL}/creators`, async ({ request }) => {
    await networkDelay();

    try {
      //const { error } = requireAuth(cookies);
      //if (error) {
      //  return HttpResponse.json({ message: error }, { status: 401 });
      //}
      const url = new URL(request.url);
      const page = Number(url.searchParams.get("page") || 1);
      const pageSize = Number(url.searchParams.get("pageSize") || 10);
      const query = String(url.searchParams.get("q") || "");

      const total = db.creators.findMany((q) =>
        q.where({
          name: (name) => name.toLowerCase().includes(query.toLowerCase()),
        }),
      ).length;

      const totalPages = Math.ceil(total / pageSize);

      const creators = db.creators.findMany(
        (q) =>
          q.where({
            name: (name) => name.toLowerCase().includes(query.toLowerCase()),
          }),
        {
          take: pageSize,
          skip: pageSize * (page - 1),
        },
      );

      return HttpResponse.json({
        creators: creators,
        meta: {
          page,
          total,
          totalPages,
        },
      });
    } catch (error) {
      return HttpResponse.json(
        { message: error?.message || "Server Error" },
        { status: 500 },
      );
    }
  }),
];
