import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { CharacterRef, Meta } from "@/types/api";

export const getCharacters = ({
  page = 1,
  pageSize = 10,
  search,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<{ characters: CharacterRef[]; meta: Meta }> => {
  return api.get("/characters", {
    params: {
      page,
      pageSize,
      q: search,
    },
  });
};

type UseCharactersOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getCharacters>;
};

export const useCharacters = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseCharactersOptions) => {
  return useQuery<{ characters: CharacterRef[]; meta: Meta }>({
    queryKey: ["characters", page, pageSize, search],
    queryFn: () =>
      getCharacters({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
