import { listSeries } from "./listSeries";
import { buildImageUrl } from "@/lib/utils";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makeSeriesLoadOptions = (opts?: { pageSize?: number }) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { series } = await listSeries({
      page: 1,
      pageSize,
      search: q,
    });

    return series.map((s) => ({
      id: s.id,
      label: s.title,
      image: buildImageUrl("series", "small", s.thumbnailName),
    }));
  };
};

export const makeSeriesLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["seriesSketches-options", pageSize, q],
      queryFn: () =>
        listSeries({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });

    return data.series.map((s) => ({
      id: s.id,
      label: s.title,
      image: buildImageUrl("series", "small", s.thumbnailName),
    }));
  };
};
