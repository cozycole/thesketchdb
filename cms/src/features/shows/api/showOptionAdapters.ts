import { getEpisodes } from "./getEpisodes";
import { buildImageUrl } from "@/lib/utils";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makeEpisodeLoadOptions = (opts?: { pageSize?: number }) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { episodes } = await getEpisodes({
      page: 1,
      pageSize,
      search: q,
    });

    return episodes.map((e) => ({
      id: e.id,
      label: `${e.season.show.name} S${e.season.number} E${e.number}`,
      image: buildImageUrl("show", "small", e.season.show.profileImage),
    }));
  };
};

export const makeEpisodeLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["episode-options", pageSize, q],
      queryFn: () =>
        getEpisodes({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });

    return data.episodes.map((e) => ({
      id: e.id,
      label: `${e.season.show.name} S${e.season.number} E${e.number}`,
      image: e.season.show.profileImage,
    }));
  };
};
