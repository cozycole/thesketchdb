import { getCreators } from "./getCreators";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makeCreatorLoadOptions = (opts?: { pageSize?: number }) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { creators } = await getCreators({
      page: 1,
      pageSize,
      search: q,
    });

    return creators.map((s) => ({
      id: s.id,
      label: s.name,
      image: s.profileImage,
    }));
  };
};

export const makeCreatorLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["creator-options", pageSize, q],
      queryFn: () =>
        getCreators({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });

    return data.creators.map((s) => ({
      id: s.id,
      label: s.name,
      image: s.profileImage,
    }));
  };
};
