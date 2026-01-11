import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { Creator, Meta } from "@/types/api";

export const getCreators = ({
  page = 1,
  pageSize = 10,
  search,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<{ creators: Creator[]; meta: Meta }> => {
  return api.get("/creators", {
    params: {
      page,
      pageSize,
      q: search,
    },
  });
};

type UseCreatorsOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getCreators>;
};

export const useCreators = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseCreatorsOptions) => {
  return useQuery<{ creators: Creator[]; meta: Meta }>({
    queryKey: ["creators", page, pageSize, search],
    queryFn: () =>
      getCreators({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
