import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { CastMember, CastScreenshot } from "@/types/api";

export const getCast = ({
  id,
}: {
  id: number;
}): Promise<{ cast: CastMember[]; screenshots: CastScreenshot[] }> => {
  return api.get(`admin/sketch/${id}/cast`);
};

type UseCastOptions = {
  id: number;
  queryConfig?: QueryConfig<typeof getCast>;
};

export const castQueryOptions = ({ id }: { id: number }) => ({
  queryKey: ["cast", id],
  queryFn: () => getCast({ id }),
});

export const useCast = ({ id, queryConfig }: UseCastOptions) => {
  return useQuery<{ cast: CastMember[]; screenshots: CastScreenshot[] }>({
    ...castQueryOptions({ id }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
