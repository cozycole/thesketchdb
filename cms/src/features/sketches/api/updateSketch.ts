import { useMutation, useQueryClient } from "@tanstack/react-query";
import { z } from "zod";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";
import { Sketch } from "@/types/api";

import { sketchQueryOptions } from "./getSketch";
import { SketchFormData } from "../forms/sketchForm.schema";

export const updateSketchInputSchema = z.object({
  title: z.string().min(1, "Required"),
  body: z.string().min(1, "Required"),
});

export type UpdateSketchInput = z.infer<typeof updateSketchInputSchema>;

export const updateSketch = ({
  data,
  sketchId,
}: {
  data: SketchFormData;
  sketchId: number;
}): Promise<Sketch> => {
  const hasFile = data.thumbnail instanceof File;

  if (!hasFile) {
    const { thumbnail, ...json } = data;
    void thumbnail;
    return api.patch(`/admin/sketch/${sketchId}`, json);
  }

  const fd = new FormData();
  if (data.title != null) fd.append("title", data.title);
  if (data.description != null) fd.append("description", data.description);
  if (data.creators) fd.append("creators", JSON.stringify(data.creators));
  if (data.thumbnail) fd.append("thumbnail", data.thumbnail);
  return api.patch(`/sketch/${sketchId}`, data);
};

type UseUpdateSketchOptions = {
  mutationConfig?: MutationConfig<typeof updateSketch>;
};

export const useUpdateSketch = ({
  mutationConfig,
}: UseUpdateSketchOptions = {}) => {
  const queryClient = useQueryClient();

  const { onSuccess, ...restConfig } = mutationConfig || {};

  return useMutation({
    onSuccess: (data, ...args) => {
      queryClient.refetchQueries({
        queryKey: sketchQueryOptions({ id: data.id }).queryKey,
      });
      onSuccess?.(data, ...args);
    },
    ...restConfig,
    mutationFn: updateSketch,
  });
};
