import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { Tag, Meta } from "@/types/api";

export const getTags = ({
  page = 1,
  pageSize = 10,
  search,
  type,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
  type?: string;
}): Promise<{ tags: Tag[]; meta: Meta }> => {
  return api.get("/tags", {
    params: {
      page,
      pageSize,
      q: search,
      type: type,
    },
  });
};

type UseTagsOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getTags>;
};

export const useTags = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseTagsOptions) => {
  return useQuery<{ tags: Tag[]; meta: Meta }>({
    queryKey: ["tags", page, pageSize, search],
    queryFn: () =>
      getTags({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
