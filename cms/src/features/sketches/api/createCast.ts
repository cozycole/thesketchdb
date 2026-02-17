import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";
import { CastMember } from "@/types/api";

import { castQueryOptions } from "./getCast";
import { CastFormData } from "../forms/castForm.schema";

type CreateCastResponse = { castMember: CastMember };

export const createCast = async ({
  sketchId,
  data,
  existingThumbnail,
  existingProfile,
}: {
  sketchId: number;
  data: CastFormData;
  existingThumbnail: string;
  existingProfile: string;
}): Promise<CastMember> => {
  const fd = new FormData();
  fd.append("characterName", data.characterName || "");
  fd.append("castRole", data.castRole || "");
  fd.append("minorRole", String(data.minorRole) || String(false));

  fd.append("personId", data.actor?.id ? String(data.actor.id) : "");
  fd.append("characterId", data.character?.id ? String(data.character.id) : "");

  // if there's been a manual image upload
  if (data.characterThumbnail)
    fd.append("characterThumbnail", data.characterThumbnail);
  if (data.characterProfile)
    fd.append("characterProfile", data.characterProfile);

  // if we're using an auto screenshot
  if (!data.characterThumbnail && existingThumbnail) {
    const response = await fetch(existingThumbnail);
    if (response.ok) {
      const blob = await response.blob();
      fd.append("characterThumbnail", blob);
    }
  }
  if (!data.characterProfile && existingProfile) {
    const response = await fetch(existingProfile);
    if (response.ok) {
      const blob = await response.blob();
      fd.append("characterProfile", blob);
    }
  }

  const res = await api.post<CreateCastResponse>(
    `/admin/sketch/${sketchId}/cast`,
    fd,
  );
  return res.castMember;
};

type UseCreateCastOptions = {
  mutationConfig?: MutationConfig<typeof createCast>;
};

export const useCreateCast = ({
  mutationConfig,
}: UseCreateCastOptions = {}) => {
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
    mutationFn: createCast,
  });
};
