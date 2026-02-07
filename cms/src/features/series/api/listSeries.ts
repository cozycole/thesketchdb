import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { SeriesRef, Meta } from "@/types/api";

export const listSeries = ({
  page = 1,
  pageSize = 10,
  search,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<{ series: SeriesRef[]; meta: Meta }> => {
  return api.get("/sketch-series", {
    params: {
      page,
      pageSize,
      q: search,
    },
  });
};

type UseSketchSeriesOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof listSeries>;
};

export const useSeriesSketches = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseSketchSeriesOptions) => {
  return useQuery<{ series: SeriesRef[]; meta: Meta }>({
    queryKey: ["seriesSketches", page, pageSize, search],
    queryFn: () =>
      listSeries({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
