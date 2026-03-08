import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { SketchVideo } from "@/types/api";

export const getSketchVideos = ({
  id,
}: {
  id: number;
}): Promise<{ videos: SketchVideo[] }> => {
  return api.get(`/admin/sketch/${id}/videos`);
};

type UseSketchVideosOptions = {
  id: number;
  queryConfig?: QueryConfig<typeof getSketchVideos>;
};

export const sketchVideosQueryOptions = ({ id }: { id: number }) => ({
  queryKey: ["sketchVideos", id],
  queryFn: () => getSketchVideos({ id }),
});

export const useSketchVideos = ({
  id,
  queryConfig,
}: UseSketchVideosOptions) => {
  return useQuery<{ videos: SketchVideo[] }>({
    ...sketchVideosQueryOptions({ id }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
