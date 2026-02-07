import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";
import { Sketch } from "@/types/api";

import { sketchQueryOptions } from "./getSketch";
import { SketchFormData } from "../forms/sketchForm.schema";
import { parseHMS, toYYYYMMDD } from "@/lib/utils";

type SketchResponse = { sketch: Sketch };

export const createSketch = async ({
  data,
}: {
  data: SketchFormData;
}): Promise<Sketch> => {
  const fd = new FormData();
  fd.append("title", data.title);
  fd.append("url", data.url);
  fd.append("description", data.description);
  fd.append("duration", String(parseHMS(data.duration)));
  fd.append("uploadDate", toYYYYMMDD(data.uploadDate));
  fd.append("popularity", String(data.popularity));

  fd.append("creatorId", data.creator?.id ? String(data.creator.id) : "");

  fd.append("episodeId", data.episode?.id ? String(data.episode.id) : "");
  fd.append("episodeStart", String(parseHMS(data.episodeStartTime)));
  fd.append("number", String(data.episodeSketchOrder));

  fd.append("seriesId", data.series?.id ? String(data.series.id) : "");
  fd.append("seriesPart", data.seriesPart);

  fd.append("recurringId", data.recurring?.id ? String(data.recurring.id) : "");

  if (data.thumbnail) fd.append("thumbnail", data.thumbnail);
  const res = await api.post<SketchResponse>(`/admin/sketch`, fd);
  return res.sketch;
};

type UseCreateSketchOptions = {
  mutationConfig?: MutationConfig<typeof createSketch>;
};

export const useCreateSketch = ({
  mutationConfig,
}: UseCreateSketchOptions = {}) => {
  const queryClient = useQueryClient();

  const { onSuccess, ...restConfig } = mutationConfig || {};

  return useMutation({
    onSuccess: (data, ...args) => {
      queryClient.refetchQueries({
        queryKey: sketchQueryOptions({ id: data?.id ?? 0 }).queryKey,
      });
      onSuccess?.(data, ...args);
    },
    ...restConfig,
    mutationFn: createSketch,
  });
};
