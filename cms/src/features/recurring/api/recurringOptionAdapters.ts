import { getRecurringSketches } from "./getRecurringSketches";
import { buildImageUrl } from "@/lib/utils";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makeRecurringLoadOptions = (opts?: { pageSize?: number }) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { recurringSketches } = await getRecurringSketches({
      page: 1,
      pageSize,
      search: q,
    });

    return recurringSketches.map((s) => ({
      id: s.id,
      label: s.title,
      image: buildImageUrl("recurring", "small", s.thumbnailName),
    }));
  };
};

export const makeRecurringLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["recurringSketches-options", pageSize, q],
      queryFn: () =>
        getRecurringSketches({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });

    return data.recurringSketches.map((s) => ({
      id: s.id,
      label: s.title,
      image: buildImageUrl("recurring", "small", s.thumbnailName),
    }));
  };
};
