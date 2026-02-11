import { getSketches } from "./getSketches";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makeSketchLoadOptions = (opts?: { pageSize?: number }) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { sketches } = await getSketches({
      page: 1,
      pageSize,
      search: q,
    });

    return sketches.map((s) => ({
      id: s.id,
      label: s.title,
      image: s.thumbnailName,
    }));
  };
};

export const makeSketchLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["sketch-options", pageSize, q],
      queryFn: () =>
        getSketches({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });

    return data.sketches.map((s) => ({
      id: s.id,
      label: s.title,
      image: s.thumbnailName,
    }));
  };
};
