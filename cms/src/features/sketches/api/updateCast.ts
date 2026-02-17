import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";
import { CastMember } from "@/types/api";

import { castQueryOptions } from "./getCast";
import { CastFormData } from "../forms/castForm.schema";

type UpdateCastResponse = { castMember: CastMember };

export const updateCast = async ({
  sketchId,
  data,
}: {
  sketchId: number;
  data: CastFormData;
}): Promise<CastMember> => {
  const fd = new FormData();
  fd.append("id", String(data.id));
  fd.append("characterName", data.characterName || "");
  fd.append("castRole", data.castRole || "");
  fd.append("minorRole", String(data.minorRole) || String(false));

  fd.append("personId", data.actor?.id ? String(data.actor.id) : "");
  fd.append("characterId", data.character?.id ? String(data.character.id) : "");

  if (data.characterThumbnail)
    fd.append("characterThumbnail", data.characterThumbnail);
  if (data.characterProfile)
    fd.append("characterProfile", data.characterProfile);

  const res = await api.put<UpdateCastResponse>(
    `/admin/sketch/${sketchId}/cast/${data.id}`,
    fd,
  );
  return res.castMember;
};

type UseUpdateCastOptions = {
  mutationConfig?: MutationConfig<typeof updateCast>;
};

export const useUpdateCast = ({
  mutationConfig,
}: UseUpdateCastOptions = {}) => {
  const queryClient = useQueryClient();

  const { onSuccess, ...restConfig } = mutationConfig || {};

  return useMutation({
    onSuccess: (data, variables, ...args) => {
      queryClient.invalidateQueries({
        queryKey: castQueryOptions({ id: variables.sketchId }).queryKey,
      });
      onSuccess?.(data, variables, ...args);
    },
    ...restConfig,
    mutationFn: updateCast,
  });
};
