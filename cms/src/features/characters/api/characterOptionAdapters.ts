import { getCharacters } from "./getCharacters";
import { buildImageUrl } from "@/lib/utils";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";

export const makeCharacterLoadOptions = (opts?: { pageSize?: number }) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { characters } = await getCharacters({
      page: 1,
      pageSize,
      search: q,
    });

    return characters.map((s) => ({
      id: s.id,
      label: s.name,
      image: buildImageUrl("character", "small", s.profileImage),
    }));
  };
};

export const makeCharacterLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["character-options", pageSize, q],
      queryFn: () =>
        getCharacters({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });

    return data.characters.map((s) => ({
      id: s.id,
      label: s.name,
      image: s.profileImage,
    }));
  };
};
