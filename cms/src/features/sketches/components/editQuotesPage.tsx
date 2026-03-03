import { useEffect, useRef, useState } from "react";
import { useQuotes } from "@/features/sketches/api/getQuotes";

import { SaveIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { Button } from "@/components/ui/button";
import { QuoteUI } from "../hooks/useQuoteEditor";
import { TranscriptPanel } from "@/features/sketches/components/transcriptPanel";
import { QuotesPanel } from "@/features/sketches/components/quotesPanel";
import { useQuoteEditor } from "../hooks/useQuoteEditor";

import {
  DragStartEvent,
  DndContext,
  PointerSensor,
  DragOverlay,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { formatHMS } from "@/lib/utils";

type DragBundle = { kind: "transcript"; ids: number[] } | null;

export function EditQuotesPage({ sketchId }: { sketchId: number }) {
  const {
    data,
    isLoading: quotesLoading,
    isError: quotesError,
  } = useQuotes({
    id: Number(sketchId),
  });

  const { state, dispatch, quoteKeysSorted, hasChanges, save } = useQuoteEditor(
    {
      sketchId,
      transcript: [],
      selectedTranscriptIds: new Set(),
      selectedQuoteKeys: new Set(),
      saving: false,
    },
  );

  useEffect(() => {
    if (!data) return;
    dispatch({
      type: "LOAD",
      transcript: data.transcript,
      quotes: data.quotes,
    });
  }, [data, dispatch]);

  const actions = {
    addQuote: (qs: QuoteUI[]) => dispatch({ type: "ADD_QUOTES", quotes: qs }),
    updateQuote: (q: QuoteUI) => dispatch({ type: "UPDATE_QUOTE", new: q }),
    deleteQuote: (q: QuoteUI) => {
      dispatch({ type: "DELETE_QUOTE", key: q.clientId });
    },
    validateQuote: (q: QuoteUI) => {
      dispatch({ type: "VALIDATE_QUOTE", key: q.clientId });
    },
    selectQuote: (key: string, e: PointerEvent) => {
      dispatch({ type: "SELECT_QUOTE", key, e });
    },
    selectTranscriptLine: (id: number, e: MouseEvent) => {
      dispatch({ type: "SELECT_TRANSCRIPT", id, e });
    },

    mergeQuotes: () => {
      dispatch({ type: "MERGE_QUOTES" });
    },
    clearSelected: () => {
      dispatch({ type: "CLEAR_SELECTED" });
    },
  };

  // remove selection on off click
  const editorRef = useRef<HTMLDivElement | null>(null);
  useEffect(() => {
    function handlePointerDown(e: PointerEvent) {
      const el = editorRef.current;
      if (!el) return;

      // If click is outside transcript panel → clear selection
      if (!el.contains(e.target as Node)) {
        dispatch({ type: "CLEAR_SELECTED" });
      }
    }

    document.addEventListener("pointerdown", handlePointerDown);
    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
    };
  }, [dispatch]);

  // transcript to quotes drag and drop set up
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
  );

  const [dragBundle, setDragBundle] = useState<DragBundle>(null);
  const onDragStart = (e: DragStartEvent) => {
    const data = e.active.data.current as
      | { type: "transcriptLine"; lineId: number }
      | undefined;

    if (!data || data.type !== "transcriptLine") return;

    const activeId = data.lineId;

    // If the thing you started dragging is NOT selected, make it the only selection
    // (optional, but matches typical UX).
    if (!state.selectedTranscriptIds.has(activeId)) {
      actions.selectTranscriptLine(activeId, null);
      setDragBundle({ kind: "transcript", ids: [activeId] });
      return;
    }

    // Otherwise drag everything selected
    setDragBundle({
      kind: "transcript",
      ids: Array.from(state.selectedTranscriptIds),
    });
  };

  const onDragEnd = () => {
    setDragBundle(null);
    actions.clearSelected();
    // convert the transcript lines to quotes
    const quotes: QuoteUI[] = [];
    for (const id of state.selectedTranscriptIds) {
      const t = state.transcript.find((e) => e.id === id);
      quotes.push({
        id: undefined,
        startTimeMs: formatHMS(Math.round(t.startMs / 1000)),
        endTimeMs: formatHMS(Math.round(t.endMs / 1000)),
        text: t.text,
        cast: [],
        tags: [],
      } as QuoteUI);
    }
    dispatch({ type: "ADD_QUOTES", quotes: quotes });
  };

  return (
    <DndContext
      sensors={sensors}
      onDragStart={onDragStart}
      onDragEnd={onDragEnd}
    >
      <div className="w-full flex" ref={editorRef}>
        {quotesLoading ? (
          <div className="h-screen mx-auto">
            <Spinner className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : quotesError ? (
          <div className="mx-auto">Error getting quote data</div>
        ) : (
          <>
            <div className="w-1/2 flex flex-col select-none">
              <QuotesPanel
                sketchId={sketchId}
                quoteKeys={quoteKeysSorted}
                quotesByKey={state.quotesByKey}
                errorsByKey={state.errorsByKey}
                onAddQuote={actions.addQuote}
                onUpdateQuote={actions.updateQuote}
                onDeleteQuote={actions.deleteQuote}
                onBlur={actions.validateQuote}
                selectedKeys={state.selectedQuoteKeys}
                onSelectKey={actions.selectQuote}
                onClearSelected={actions.clearSelected}
                onMerge={actions.mergeQuotes}
              />

              {hasChanges && (
                <Button
                  className="sticky mt-3 bottom-6 self-center text-lg text-white rounded-lg "
                  onClick={() => save()}
                >
                  {state.saving ? <Spinner /> : <SaveIcon />}
                  Save
                </Button>
              )}
            </div>
            <div className="w-1/2 select-none">
              <h1 className="text-center mb-5">Transcript</h1>
              <TranscriptPanel
                transcript={data.transcript}
                onSelect={actions.selectTranscriptLine}
                selectedIds={state.selectedTranscriptIds}
              />
            </div>
            {/* Optional: render a grouped drag preview */}
            <DragOverlay dropAnimation={null}>
              {dragBundle?.kind === "transcript" ? (
                <div className="rounded-md border bg-white p-2 text-sm shadow">
                  {dragBundle.ids.length} line
                  {dragBundle.ids.length === 1 ? "" : "s"}
                </div>
              ) : null}
            </DragOverlay>
          </>
        )}
      </div>
    </DndContext>
  );
}
