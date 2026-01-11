import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { Episode, Meta } from "@/types/api";

export const getEpisodes = ({
  page = 1,
  pageSize = 10,
  search,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<{ episodes: Episode[]; meta: Meta }> => {
  return api.get("/episodes", {
    params: {
      page,
      pageSize,
      q: search,
    },
  });
};

type UseEpisodesOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getEpisodes>;
};

export const useEpisodes = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseEpisodesOptions) => {
  return useQuery<{ episodes: Episode[]; meta: Meta }>({
    queryKey: ["episodes", page, pageSize, search],
    queryFn: () =>
      getEpisodes({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
