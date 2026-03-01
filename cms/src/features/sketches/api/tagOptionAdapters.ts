import { getTags } from "./getTags";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makeTagLoadOptions = (opts?: {
  pageSize?: number;
  tagType?: string;
}) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { tags } = await getTags({
      page: 1,
      pageSize,
      search: q,
      type: opts?.tagType,
    });

    return tags.map((t) => ({
      id: t.id,
      label: t.name,
    }));
  };
};

export const makeTagLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["tags-options", pageSize, q],
      queryFn: () =>
        getTags({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });
    return data.tags.map((t) => ({
      id: t.id,
      label: t.name,
    }));
  };
};
