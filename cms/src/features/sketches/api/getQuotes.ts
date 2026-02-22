import { useQuery } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { QueryConfig } from "@/lib/react-query";
import { TranscriptLine } from "@/types/api";

export const getQuotes = ({
  id,
}: {
  id: number;
}): Promise<{ transcript: TranscriptLine[] }> => {
  return api.get(`admin/sketch/${id}/quotes`);
};

type UseQuoteOptions = {
  id: number;
  queryConfig?: QueryConfig<typeof getQuotes>;
};

export const quotesQueryOptions = ({ id }: { id: number }) => ({
  queryKey: ["quotes", id],
  queryFn: () => getQuotes({ id }),
});

export const useQuotes = ({ id, queryConfig }: UseQuoteOptions) => {
  return useQuery<{ transcript: TranscriptLine[] }>({
    ...quotesQueryOptions({ id }),
    placeholderData: (prev) => prev,
    ...queryConfig,
  });
};
