import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";

import { castQueryOptions } from "./getCast";

export const deleteCast = async ({
  sketchId,
  castId,
}: {
  sketchId: number;
  castId: number;
}) => {
  await api.delete(`/admin/sketch/${sketchId}/cast/${castId}`);
  return null;
};

type UseDeleteCastOptions = {
  mutationConfig?: MutationConfig<typeof deleteCast>;
};

export const useDeleteCast = ({
  mutationConfig,
}: UseDeleteCastOptions = {}) => {
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
    mutationFn: deleteCast,
  });
};
