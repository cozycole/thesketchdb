import { useState } from "react";

import { QuoteIcon, Plus } from "lucide-react";
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

import { formatHMS } from "@/lib/utils";

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

  onAddQuote: (q: QuoteUI) => void;
  onUpdateQuote: (q: QuoteUI) => void;
  onDeleteQuote: (q: QuoteUI) => void;
  onBlur: (q: QuoteUI) => void;

  // optional UI
  selectedKeys: Set<string>;
  onSelectKeys: (keys: Set<string>) => void;
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
  //selectedKeys,
  //onSelectKeys,
}: QuotesPanelProps) {
  const [quoteDialogOpen, setDialogOpen] = useState(false);
  const [newQuoteDraft, setNewQuoteDraft] =
    useState<QuoteUI>(createEmptyQuote());
  const [quoteDraftErrors, setQuoteDraftErrors] = useState<QuoteFieldsErrors>(
    {},
  );

  return (
    <>
      <div>
        <Button
          className="ml-2 text-white font-bold"
          onClick={() => setDialogOpen(true)}
        >
          <Plus className="mr-2 h-4 w-4" />
          Add Quote
        </Button>
      </div>
      <div className="grid grid-cols-1 gap-4 my-2 mx-2">
        {quoteKeys.length === 0 ? (
          <p className="col-span-full mt-10 text-center text-muted-foreground">
            <QuoteIcon className="mx-auto mb-6" />
            No quotes yet. Create one or drag a transcript line here
          </p>
        ) : (
          quoteKeys.map(
            (k) =>
              quotesByKey[k] && (
                <div key={k}>
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
                const result = quoteFieldsSchema.safeParse(newQuoteDraft);

                if (!result.success) {
                  const fieldErrors = zodErrorToFieldErrors<
                    typeof quoteFieldsSchema
                  >(result.error);
                  setQuoteDraftErrors(fieldErrors);
                  return;
                }
                onAddQuote(newQuoteDraft);
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
