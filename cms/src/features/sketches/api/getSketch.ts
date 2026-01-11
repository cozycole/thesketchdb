import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { Sketch } from "@/types/api";

export const getSketch = ({
  id,
}: {
  id: number;
}): Promise<{ sketch: Sketch }> => {
  return api.get(`/admin/sketch/${id}`);
};

type UseSketchOptions = {
  id: number;
  queryConfig?: QueryConfig<typeof getSketch>;
};

export const sketchQueryOptions = ({ id }: { id: number }) => ({
  queryKey: ["sketch", id],
  queryFn: () => getSketch({ id }),
});

export const useSketch = ({ id, queryConfig }: UseSketchOptions) => {
  return useQuery<{ sketch: Sketch }>({
    ...sketchQueryOptions({ id }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
