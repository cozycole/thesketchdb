import { getCastMembers } from "./getCastMembers";
import { buildImageUrl } from "@/lib/utils";
import type { SelectEntity } from "@/components/ui/asyncSearchSelect";
import type { QueryClient } from "@tanstack/react-query";
import { CastMember } from "@/types/api";

export function toCastOption(c: CastMember): SelectEntity {
  let label = c.characterName;
  if (c.actor) {
    if (!c.characterName) {
      label = `${c.actor.first} ${c.actor.last}`;
    } else {
      label += ` (${c.actor.first} ${c.actor.last})`;
    }
  }
  let profile = buildImageUrl("cast/profile", "small", c.profileImage);
  if (!c.profileImage && c.actor) {
    profile = buildImageUrl("person", "small", c.actor.profileImage);
  }

  return {
    id: c.id,
    label: label,
    image: profile,
  };
}

export const makeCastMemberLoadOptions = (opts?: {
  pageSize?: number;
  sketchId?: number;
}) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const { castMembers } = await getCastMembers({
      page: 1,
      pageSize,
      search: q,
      sketch: opts?.sketchId,
    });

    return castMembers.map((c) => toCastOption(c));
  };
};

export const makeCastLoadOptionsRQ = (
  queryClient: QueryClient,
  opts?: { pageSize?: number },
) => {
  const pageSize = opts?.pageSize ?? 10;

  return async (q: string): Promise<SelectEntity[]> => {
    const data = await queryClient.fetchQuery({
      queryKey: ["cast-options", pageSize, q],
      queryFn: () =>
        getCastMembers({
          page: 1,
          pageSize,
          search: q,
        }),
      staleTime: 30_000,
    });
    return data.castMembers.map((c) => toCastOption(c));
  };
};
