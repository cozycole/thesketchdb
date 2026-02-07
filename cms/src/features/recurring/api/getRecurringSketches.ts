import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { RecurringRef, Meta } from "@/types/api";

export const getRecurringSketches = ({
  page = 1,
  pageSize = 10,
  search,
}: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<{ recurringSketches: RecurringRef[]; meta: Meta }> => {
  return api.get("/recurring-sketches", {
    params: {
      page,
      pageSize,
      q: search,
    },
  });
};

type UseRecurringSketchesOptions = {
  page: number;
  pageSize: number;
  search?: string;
  queryConfig?: QueryConfig<typeof getRecurringSketches>;
};

export const useRecurringSketches = ({
  page,
  pageSize,
  search,
  queryConfig,
}: UseRecurringSketchesOptions) => {
  return useQuery<{ recurringSketches: RecurringRef[]; meta: Meta }>({
    queryKey: ["recurringSketches", page, pageSize, search],
    queryFn: () =>
      getRecurringSketches({
        page,
        pageSize,
        search,
      }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
