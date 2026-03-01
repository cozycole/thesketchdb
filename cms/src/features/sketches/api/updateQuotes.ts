import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/lib/api-client";
import { MutationConfig } from "@/lib/react-query";

import { Quote } from "@/types/api";
import { QuoteUI } from "../hooks/useQuoteEditor";

import { quotesQueryOptions } from "./getQuotes";
import { parseHMS } from "@/lib/utils";

type UpdateQuotesResponse = { quotes: Quote[] };

export const updateQuotes = async ({
  sketchId,
  updatedQuotes,
  deletedIds,
}: {
  sketchId: number;
  updatedQuotes: QuoteUI[];
  deletedIds: Set<number>;
}): Promise<Quote[]> => {
  const payload = {
    upsert: updatedQuotes.map((q) => {
      return {
        id: q.id ? q.id : null,
        startTimeMs: q.startMs * 1000,
        endTimeMs: parseHMS(q.endTimeMs) ? parseHMS(q.endTimeMs) * 1000 : null,
        text: q.text ?? "",
        cast: q.cast.map((c) => c.id),
        tags: q.tags.map((t) => t.id),
      };
    }),
    delete: [...deletedIds],
  };

  console.log("SUBMITTING: ", payload);

  const res = await api.put<UpdateQuotesResponse>(
    `/admin/sketch/${sketchId}/quotes`,
    payload,
  );
  return res.quotes;
};

type UseUpdateQuotesOptions = {
  mutationConfig?: MutationConfig<typeof updateQuotes>;
};

export const useUpdateQuotes = ({
  mutationConfig,
}: UseUpdateQuotesOptions = {}) => {
  const queryClient = useQueryClient();

  const { onSuccess, ...restConfig } = mutationConfig || {};

  return useMutation({
    onSuccess: (data, variables, ...args) => {
      queryClient.invalidateQueries({
        queryKey: quotesQueryOptions({ id: variables.sketchId }).queryKey,
      });
      onSuccess?.(data, variables, ...args);
    },
    ...restConfig,
    mutationFn: updateQuotes,
  });
};
