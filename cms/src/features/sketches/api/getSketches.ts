import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { Sketch, Meta } from "@/types/api";

export const getSketches = ({
  page = 1,
  pageSize = 10,
  search,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<{ sketches: Sketch[]; meta: Meta }> => {
  return api.get("/sketches", {
    params: {
      page,
      pageSize,
      q: search,
    },
  });
};

type UseSketchesOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getSketches>;
};

export const useSketches = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseSketchesOptions) => {
  return useQuery<{ sketches: Sketch[]; meta: Meta }>({
    queryKey: ["sketches", page, pageSize, search],
    queryFn: () =>
      getSketches({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
