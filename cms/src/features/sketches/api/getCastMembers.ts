import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { CastMember, Meta } from "@/types/api";

export const getCastMembers = ({
  page = 1,
  pageSize = 10,
  search,
  sketch,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
  sketch?: number;
}): Promise<{ castMembers: CastMember[]; meta: Meta }> => {
  return api.get("/cast", {
    params: {
      page,
      pageSize,
      q: search,
      sketch: sketch,
    },
  });
};

type UseCastMembersOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getCastMembers>;
};

export const useCastMembers = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseCastMembersOptions) => {
  return useQuery<{ castMembers: CastMember[]; meta: Meta }>({
    queryKey: ["cast-members", page, pageSize, search],
    queryFn: () =>
      getCastMembers({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
