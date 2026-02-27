import { useEffect } from "react";
import { useQuotes } from "@/features/sketches/api/getQuotes";

import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { QuoteUI } from "../hooks/useQuoteEditor";
import { TranscriptPanel } from "@/features/sketches/components/transcriptPanel";
import { QuotesPanel } from "@/features/sketches/components/quotesPanel";
import { useQuoteEditor } from "../hooks/useQuoteEditor";

export function EditQuotesPage({ sketchId }: { sketchId: number }) {
  const {
    data,
    isLoading: quotesLoading,
    isError: quotesError,
  } = useQuotes({
    id: Number(sketchId),
  });

  const { state, dispatch, quoteKeysSorted, dirtyKeys, hasChanges } =
    useQuoteEditor({
      sketchId,
      transcript: [],
      selectedTranscriptIds: new Set(),
      selectedQuoteKeys: new Set(),
      saving: false,
    });

  useEffect(() => {
    if (!data) return;
    dispatch({
      type: "LOAD",
      transcript: data.transcript,
      quotes: data.quotes,
    });
  }, [data, dispatch]);

  const actions = {
    addQuote: (q: QuoteUI) => dispatch({ type: "ADD_QUOTE", quote: q }),
    updateQuote: (q: QuoteUI) => dispatch({ type: "UPDATE_QUOTE", new: q }),
    deleteQuote: (q: QuoteUI) => {
      dispatch({ type: "DELETE_QUOTE", key: q.clientId });
    },
    validateQuote: (q: QuoteUI) => {
      dispatch({ type: "VALIDATE_QUOTE", key: q.clientId });
    },
    selectQuote: () => {},
  };

  return (
    <div className="w-full flex">
      {quotesLoading ? (
        <div className="h-screen mx-auto">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : quotesError ? (
        <div className="mx-auto">Error getting quote data</div>
      ) : (
        <>
          <div className="w-1/2 flex flex-col">
            <QuotesPanel
              quoteKeys={quoteKeysSorted}
              quotesByKey={state.quotesByKey}
              errorsByKey={state.errorsByKey}
              onAddQuote={actions.addQuote}
              onUpdateQuote={actions.updateQuote}
              onDeleteQuote={actions.deleteQuote}
              onBlur={actions.validateQuote}
              selectedKeys={state.selectedQuoteKeys}
              onSelectKeys={actions.selectQuote}
            />

            {hasChanges && (
              <Button className="sticky bottom-6 self-center text-lg text-white rounded-lg">
                Save
              </Button>
            )}
          </div>
          <div className="w-1/2">
            <h1 className="text-center mb-5">Transcript</h1>
            <TranscriptPanel transcript={data.transcript} />
          </div>
        </>
      )}
    </div>
  );
}
