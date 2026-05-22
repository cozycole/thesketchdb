import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";

import { sketchQueryOptions } from "./getSketch";

export const deleteSketch = async ({ sketchId }: { sketchId: number }) => {
  await api.delete(`/admin/sketch/${sketchId}`);
  return null;
};

type UseDeleteSketchOptions = {
  mutationConfig?: MutationConfig<typeof deleteSketch>;
};

export const useDeleteSketch = ({
  mutationConfig,
}: UseDeleteSketchOptions = {}) => {
  const queryClient = useQueryClient();

  const { onSuccess, ...restConfig } = mutationConfig || {};

  return useMutation({
    onSuccess: (data, variables, ...args) => {
      queryClient.removeQueries({
        queryKey: sketchQueryOptions({ id: variables.sketchId }).queryKey,
      });
      queryClient.invalidateQueries({
        queryKey: ["sketches"],
      });
      onSuccess?.(data, variables, ...args);
    },
    ...restConfig,
    mutationFn: deleteSketch,
  });
};
