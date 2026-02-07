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
  http.post(`${env.API_URL}/admin/sketch`, async ({ request }) => {
    await networkDelay();

    try {
      const latestSketch = db.sketches.findMany(
        (q) => q.where({ id: () => true }),
        { orderBy: { id: "desc" } },
      );

      let id = 1;
      if (latestSketch.length) {
        id = latestSketch[0].id + 1;
      }

      const creator = db.creators.findFirst((q) => q.where({ id: 3 }));
      const title = "Testing Title";
      await db.sketches.create({
        id: id,
        slug: title.toLowerCase().replace(/[^a-z0-9]+/g, "-"),
        title: title,
        url: `https://example.com/sketch-${id}`,
        description: `Newly created sketch ID: ${id}`,
        thumbnailUrl:
          "https://thesketchdb.sfo2.cdn.digitaloceanspaces.com/sketch/medium/02719283-277e-4d99-ab58-40f78f235481.jpg",
        creators: [creator],
      });

      return HttpResponse.json({
        status: 200,
        message: "Sketch successfully added",
        sketch: { id: id },
      });
    } catch (error) {
      return HttpResponse.json(
        { message: error?.message || "Server Error" },
        { status: 500 },
      );
    }
  }),
];
