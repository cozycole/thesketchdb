import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";
import { Sketch } from "@/types/api";

import { sketchQueryOptions } from "./getSketch";
import { SketchFormData } from "../forms/sketchForm.schema";
import { parseHMS, toYYYYMMDD } from "@/lib/utils";

type SketchResponse = { sketch: Sketch };

export const updateSketch = async ({
  data,
  sketchId,
}: {
  data: SketchFormData;
  sketchId: number;
}): Promise<Sketch> => {
  const fd = new FormData();
  fd.append("id", String(data.id));
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
  const res = await api.put<SketchResponse>(`/admin/sketch/${sketchId}`, fd);
  return res.sketch;
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
