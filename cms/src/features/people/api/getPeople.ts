import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { PersonRef, Meta } from "@/types/api";

export const getPeople = ({
  page = 1,
  pageSize = 10,
  search,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<{ people: PersonRef[]; meta: Meta }> => {
  return api.get("/people", {
    params: {
      page,
      pageSize,
      q: search,
    },
  });
};

type UsePeopleOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getPeople>;
};

export const usePeople = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UsePeopleOptions) => {
  return useQuery<{ people: PersonRef[]; meta: Meta }>({
    queryKey: ["people", page, pageSize, search],
    queryFn: () =>
      getPeople({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
