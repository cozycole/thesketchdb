import { HttpResponse, http } from "msw";

import { env } from "@/config/env";

import { db } from "../db";
import { networkDelay } from "../utils";

export const sketchHandlers = [
  http.get(`${env.API_URL}/sketches`, async ({ request }) => {
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

      const total = db.sketches.findMany((q) =>
        q.where({
          title: (title) => title.toLowerCase().includes(query.toLowerCase()),
        }),
      ).length;

      const totalPages = Math.ceil(total / pageSize);

      const sketches = db.sketches.findMany(
        (q) =>
          q.where({
            title: (title) => title.toLowerCase().includes(query.toLowerCase()),
          }),
        {
          take: pageSize,
          skip: pageSize * (page - 1),
        },
      );

      return HttpResponse.json({
        sketches: sketches,
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
  http.get<{ sketchId: string }>(
    `${env.API_URL}/admin/sketch/:sketchId`,
    async ({ params }) => {
      await networkDelay();

      try {
        //const { error } = requireAuth(cookies);
        //if (error) {
        //  return HttpResponse.json({ message: error }, { status: 401 });
        //}
        const { sketchId } = params;

        const sketch = db.sketches.findFirst((q) =>
          q.where({
            id: (id) => id === Number(sketchId),
          }),
        );

        return HttpResponse.json({
          sketch: sketch,
        });
      } catch (error) {
        return HttpResponse.json(
          { message: error?.message || "Server Error" },
          { status: 500 },
        );
      }
    },
  ),
  http.patch<{ sketchId: string }>(
    `${env.API_URL}/admin/sketch/:sketchId`,
    async ({ request, params }) => {
      await networkDelay();

      try {
        //const { error } = requireAuth(cookies);
        //if (error) {
        //  return HttpResponse.json({ message: error }, { status: 401 });
        //}

        const data = await request.json();
        console.log(data);
        return HttpResponse.json({
          status: 200,
          message: "Sketch successfully added",
        });
      } catch (error) {
        return HttpResponse.json(
          { message: error?.message || "Server Error" },
          { status: 500 },
        );
      }
    },
  ),
];
