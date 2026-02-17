import { getPeople } from "./getPeople";
import { buildImageUrl } from "@/lib/utils";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makePersonLoadOptions = (opts?: { pageSize?: number }) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { people } = await getPeople({
      page: 1,
      pageSize,
      search: q,
    });

    return people.map((s) => ({
      id: s.id,
      label: s.first + " " + s.last,
      image: buildImageUrl("person", "small", s.profileImage),
    }));
  };
};

export const makePersonLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["person-options", pageSize, q],
      queryFn: () =>
        getPeople({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });

    return data.people.map((s) => ({
      id: s.id,
      label: s.first + " " + s.last,
      image: s.profileImage,
    }));
  };
};
