import { useState } from "react";
import {
  DndContext,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
  useDroppable,
} from "@dnd-kit/core";

import { QuoteIcon, Plus, MergeIcon, XIcon } from "lucide-react";
import { QuoteUI } from "../hooks/useQuoteEditor";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { QuoteFields } from "../forms/quoteFields";
import {
  QuoteFieldsErrors,
  quoteFieldsSchema,
  zodErrorToFieldErrors,
} from "../forms/quoteFields.schema";

import { cn, formatHMS } from "@/lib/utils";

function createEmptyQuote(): QuoteUI {
  return {
    clientId: "",
    id: undefined,
    text: "",
    startMs: 0,
    startTimeMs: formatHMS(0),
    endTimeMs: undefined,
    cast: [],
    tags: [],
  };
}

type QuotesPanelProps = {
  quoteKeys: string[];
  quotesByKey: Record<string, QuoteUI>;
  errorsByKey: Record<string, QuoteFieldsErrors>;

  sketchId: number;

  onAddQuote: (q: QuoteUI[]) => void;
  onUpdateQuote: (q: QuoteUI) => void;
  onDeleteQuote: (q: QuoteUI) => void;
  onBlur: (q: QuoteUI) => void;

  // optional UI
  selectedKeys: Set<string>;
  onSelectKey: (key: string, e: PointerEvent) => void;
  onClearSelected: () => void;
  onMerge: () => void;
};

export function QuotesPanel({
  quoteKeys,
  quotesByKey,
  errorsByKey,
  sketchId,
  onAddQuote,
  onUpdateQuote,
  onDeleteQuote,
  onBlur,
  selectedKeys,
  onSelectKey,
  onClearSelected,
  onMerge,
}: QuotesPanelProps) {
  const [quoteDialogOpen, setDialogOpen] = useState(false);
  const [newQuoteDraft, setNewQuoteDraft] =
    useState<QuoteUI>(createEmptyQuote());
  const [quoteDraftErrors, setQuoteDraftErrors] = useState<QuoteFieldsErrors>(
    {},
  );

  const { isOver, setNodeRef, active } = useDroppable({
    id: "quotes-dropzone",
  });

  const isTranscriptDragging = active?.data?.current?.type === "transcriptLine";
  const highlight = isOver && isTranscriptDragging;

  return (
    <>
      <div className="flex sticky justify-between mx-2 top-16">
        <Button
          className="text-white font-bold"
          onClick={() => setDialogOpen(true)}
        >
          <Plus className="mr-2 h-4 w-4" />
          Add Quote
        </Button>
        {selectedKeys.size > 1 && (
          <div className="flex gap-2">
            <Button
              className="bg-blue-500 text-white hover:bg-blue-400"
              onClick={onMerge}
            >
              <MergeIcon className="mr-2 h-4 w-4" />
              Merge
            </Button>
            <Button variant="secondary" onClick={onClearSelected}>
              <XIcon className="h-4 w-4" />
            </Button>
          </div>
        )}
      </div>
      <div
        ref={setNodeRef}
        className={cn(
          "flex flex-col gap-3 mt-2 rounded-lg transition-colors",
          highlight && "ring-2 ring-orange-500 bg-orange-50/40",
        )}
      >
        {quoteKeys.length === 0 ? (
          <p className="col-span-full mt-10 text-center text-muted-foreground">
            <QuoteIcon className="mx-auto mb-6" />
            No quotes yet. Create one or drag a transcript line here
          </p>
        ) : (
          quoteKeys.map(
            (k) =>
              quotesByKey[k] && (
                <div
                  key={k}
                  className={
                    selectedKeys.has(k)
                      ? "rounded-lg ring-2 ring-orange-400"
                      : ""
                  }
                  onPointerDown={(e) => onSelectKey(k, e)}
                >
                  <QuoteFields
                    value={quotesByKey[k]}
                    sketchId={sketchId}
                    errors={errorsByKey[k]}
                    onChange={onUpdateQuote}
                    onDelete={onDeleteQuote}
                    onBlurQuote={onBlur}
                  />
                </div>
              ),
          )
        )}
      </div>
      <Dialog open={quoteDialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent onOpenAutoFocus={(e) => e.preventDefault()}>
          <DialogHeader>
            <DialogTitle>Add Quote</DialogTitle>
          </DialogHeader>
          <QuoteFields<QuoteUI>
            value={newQuoteDraft}
            sketchId={sketchId}
            onChange={setNewQuoteDraft}
            onDelete={onDeleteQuote}
            errors={quoteDraftErrors}
          />
          <DialogFooter>
            <Button
              className="text-white ml-auto"
              onClick={() => {
                // this logic is necessary because the quote panel handles
                // the state of the quoteDraft before adding it
                const result = quoteFieldsSchema.safeParse(newQuoteDraft);

                if (!result.success) {
                  const fieldErrors = zodErrorToFieldErrors<
                    typeof quoteFieldsSchema
                  >(result.error);
                  setQuoteDraftErrors(fieldErrors);
                  return;
                }

                onAddQuote([newQuoteDraft]);
                setDialogOpen(false);
                setNewQuoteDraft(createEmptyQuote());
                setQuoteDraftErrors({});
              }}
            >
              Add
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
