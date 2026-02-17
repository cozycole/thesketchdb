import { useMutation, useQueryClient } from "@tanstack/react-query";
import { MutationConfig } from "@/lib/react-query";

import { api } from "@/lib/api-client";

import { castQueryOptions } from "./getCast";

export const updateCastOrder = async ({
  sketchId,
  castPositions,
}: {
  sketchId: number;
  castPositions: number[];
}) => {
  await api.put(`/admin/sketch/${sketchId}/cast/order`, {
    castPositions: castPositions,
  });
  return null;
};

type UseUpdateCastOrderOptions = {
  mutationConfig?: MutationConfig<typeof updateCastOrder>;
};

export const useUpdateCastOrder = ({
  mutationConfig,
}: UseUpdateCastOrderOptions = {}) => {
  const queryClient = useQueryClient();

  const { onSuccess, ...restConfig } = mutationConfig || {};

  return useMutation({
    onSuccess: (_data, variables, ...args) => {
      queryClient.refetchQueries({
        queryKey: castQueryOptions({ id: variables.sketchId }).queryKey,
      });
      onSuccess?.(_data, variables, ...args);
    },
    ...restConfig,
    mutationFn: updateCastOrder,
  });
};
