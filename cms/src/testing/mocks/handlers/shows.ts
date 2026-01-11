import { HttpResponse, http } from "msw";

import { env } from "@/config/env";

import { db } from "../db";
import { networkDelay } from "../utils";

export const showHandlers = [
  http.get(`${env.API_URL}/episodes`, async ({ request }) => {
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

      const epQuery = extractEpisodeQuery(query);
      console.log(epQuery);

      const allEpisodes = db.episodes.all();
      const totalQueriedEpisodes = [];
      for (const ep of allEpisodes) {
        // show name is defined and it doesn't include the showname query
        if (
          epQuery.showName &&
          !ep.season.show.name.includes(epQuery.showName)
        ) {
          continue;
        }

        if (
          epQuery.seasonNumber &&
          !(ep.season.seasonNumber == epQuery.seasonNumber)
        ) {
          continue;
        }

        if (
          epQuery.episodeNumber &&
          !(ep.episodeNumber == epQuery.episodeNumber)
        ) {
          continue;
        }

        totalQueriedEpisodes.push(ep);
      }
      const total = totalQueriedEpisodes.length;
      const totalPages = Math.ceil(total / pageSize);
      return HttpResponse.json({
        episodes: totalQueriedEpisodes.slice(
          pageSize * (page - 1),
          pageSize * page,
        ),
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

export type EpisodeQuery = {
  showName: string;
  seasonNumber?: number;
  episodeNumber?: number;
};

export function extractEpisodeQuery(input: string): EpisodeQuery {
  let normalized = input.trim().toLowerCase();
  const episodeQuery: EpisodeQuery = { showName: "" };

  // s01e02 / s1e2
  const seRe = /s(\d{1,2})e(\d{1,2})/;
  const seMatch = normalized.match(seRe);

  if (seMatch) {
    episodeQuery.seasonNumber = parseInt(seMatch[1], 10);
    episodeQuery.episodeNumber = parseInt(seMatch[2], 10);

    normalized = normalized.replace(seMatch[0], "").trim();
  } else {
    // s01 / s1
    const sRe = /s(\d{1,2})/;
    const sMatch = normalized.match(sRe);

    if (sMatch) {
      episodeQuery.seasonNumber = parseInt(sMatch[1], 10);
      normalized = normalized.replace(sMatch[0], "").trim();
    }
  }

  episodeQuery.showName = normalized;

  return episodeQuery;
}
