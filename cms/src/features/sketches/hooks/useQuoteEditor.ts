import { useCallback, useMemo, useReducer } from "react";
import { useNotifications } from "@/components/ui/notifications";
import { TranscriptLine, Quote } from "@/types/api";
import {
  QuoteFieldsData,
  mapQuoteToQuoteFields,
  quoteFieldsSchema,
  zodErrorToFieldErrors,
} from "../forms/quoteFields.schema";
import { QuoteFieldsErrors } from "../forms/quoteFields.schema";
import { parseHMS } from "@/lib/utils";
import { useUpdateQuotes } from "../api/updateQuotes";

type QuoteUI = QuoteFieldsData & {
  startMs: number;
  clientId: string;
};

type EditorState = {
  sketchId: number;

  transcript: TranscriptLine[];

  quotesByKey: Record<string, QuoteUI>;
  baselineByKey: Record<string, QuoteUI>; // snapshot of last loaded/saved for diff
  deletedIds: Set<number>;

  // UI
  errorsByKey: Record<string, QuoteFieldsErrors>;
  selectedTranscriptIds: Set<number>;
  selectedQuoteKeys: Set<string>;

  saving: boolean;
  error?: string;
};

type Action =
  | { type: "LOAD"; transcript: TranscriptLine[]; quotes: Quote[] }
  | { type: "ADD_QUOTE"; quote: QuoteUI }
  | { type: "ADD_FROM_TRANSCRIPT"; transcriptIds: number[] }
  | { type: "UPDATE_QUOTE"; new: QuoteUI }
  | { type: "DELETE_QUOTE"; key: string }
  | { type: "VALIDATE_QUOTE"; key: string }
  | {
      type: "SET_ERRORS";
      errorsByKey: Record<string, QuoteFieldsErrors | undefined>;
    }
  | { type: "RESTORE_QUOTE"; id: number }
  | { type: "SAVE_START" }
  | { type: "SAVE_SUCCESS"; quotes: Quote[] }
  | { type: "SAVE_ERROR"; error: string };

function isEqual(a?: QuoteUI, b?: QuoteUI) {
  if (!a || !b) return false;
  return (
    a.id === b.id &&
    a.startTimeMs === b.startTimeMs &&
    (a.endTimeMs ?? null) === (b.endTimeMs ?? null) &&
    a.text === b.text &&
    a.cast.join(",") === b.cast.join(",") &&
    a.tags.join(",") === b.tags.join(",")
  );
}

function reducer(state: EditorState, action: Action): EditorState {
  switch (action.type) {
    case "LOAD": {
      const quotesByKey: EditorState["quotesByKey"] = {};
      const baselineByKey: EditorState["baselineByKey"] = {};
      for (const q of action.quotes) {
        const q2: QuoteUI = {
          ...mapQuoteToQuoteFields(q),
          clientId: crypto.randomUUID(),
          startMs: q.startTimeMs / 1000,
        };
        quotesByKey[q2.clientId] = q2;
        baselineByKey[q2.clientId] = q2;
      }
      return {
        ...state,
        transcript: action.transcript,
        quotesByKey,
        baselineByKey,
        deletedIds: new Set(),
        error: undefined,
      };
    }

    case "ADD_QUOTE": {
      const quotesByKey = { ...state.quotesByKey };
      action.quote.clientId = crypto.randomUUID();
      quotesByKey[action.quote.clientId] = action.quote;
      return {
        ...state,
        quotesByKey,
      };
    }

    //case 'ADD_FROM_TRANSCRIPT': {
    //  const quotesByKey = { ...state.quotesByKey };
    //  for (const tid of action.transcriptIds) {
    //    const tl = state.transcript.find(t => t.id === tid);
    //    if (!tl) continue;
    //    const clientId = crypto.randomUUID();
    //    const q: QuoteLine = {
    //      clientId,
    //      startMs: tl.startMs,
    //      endMs: tl.endMs ?? null,
    //      text: tl.text,
    //      castMemberIds: [],
    //      tagIds: [],
    //      source: 'transcript',
    //    };
    //    quotesByKey[`tmp:${clientId}`] = q;
    //  }
    //  return { ...state, quotesByKey };
    //}

    case "UPDATE_QUOTE": {
      const cur = state.quotesByKey[action.new.clientId];
      if (!cur || !cur.clientId) return state;
      return {
        ...state,
        quotesByKey: {
          ...state.quotesByKey,
          [action.new.clientId]: { ...cur, ...action.new },
        },
      };
    }

    case "DELETE_QUOTE": {
      const q = state.quotesByKey[action.key];
      if (!q) return state;
      const quotesByKey = { ...state.quotesByKey };
      delete quotesByKey[action.key];
      const deletedIds = new Set(state.deletedIds);
      if (q.id != null) deletedIds.add(q.id);
      return { ...state, quotesByKey, deletedIds };
    }

    case "VALIDATE_QUOTE": {
      const quote = state.quotesByKey[action.key];
      if (!quote) return;

      const result = quoteFieldsSchema.safeParse(quote);
      if (!result.success) {
        const fieldErrors = zodErrorToFieldErrors<typeof quoteFieldsSchema>(
          result.error,
        );
        return {
          ...state,
          errorsByKey: {
            ...state.errorsByKey,
            [action.key]: fieldErrors,
          },
        };
      }
      const quotesByKey = { ...state.quotesByKey };
      // ensure startMs is updated to reflect in the
      // ordering of the quuotes
      quote.startMs = parseHMS(quote.startTimeMs);
      quotesByKey[action.key] = quote;

      const errorsByKey = { ...state.errorsByKey };
      delete errorsByKey[action.key];
      return {
        ...state,
        quotesByKey,
        errorsByKey,
      };
    }

    // used by save callback when errors are detected
    // before sending request
    case "SET_ERRORS": {
      return {
        ...state,
        saving: false,
        errorsByKey: action.errorsByKey,
      };
    }

    case "SAVE_START":
      return { ...state, saving: true, error: undefined };

    case "SAVE_SUCCESS": {
      const quotesByKey: EditorState["quotesByKey"] = {};
      const baselineByKey: EditorState["baselineByKey"] = {};
      for (const q of action.quotes) {
        const q2: QuoteUI = {
          ...mapQuoteToQuoteFields(q),
          clientId: crypto.randomUUID(),
          startMs: q.startTimeMs / 1000,
        };
        quotesByKey[q2.clientId] = q2;
        baselineByKey[q2.clientId] = q2;
      }
      return {
        ...state,
        saving: false,
        quotesByKey,
        baselineByKey,
        deletedIds: new Set(),
        error: undefined,
      };
    }

    case "SAVE_ERROR":
      return { ...state, saving: false, error: action.error };

    default:
      return state;
  }
}

export function useQuoteEditor(
  initial: Omit<
    EditorState,
    "quotesByKey" | "baselineByKey" | "deletedIds" | "errorsByKey" | "quoteKeys"
  >,
) {
  const { addNotification } = useNotifications();
  const [state, dispatch] = useReducer(reducer, {
    ...initial,
    errorsByKey: {},
    quotesByKey: {},
    baselineByKey: {},
    deletedIds: new Set<number>(),
  });

  // sort the keys on start timestamp
  const quoteKeysSorted = useMemo(() => {
    return Object.keys(state.quotesByKey).sort((a, b) => {
      const qa = state.quotesByKey[a]!;
      const qb = state.quotesByKey[b]!;
      return qa.startMs - qb.startMs;
    });
  }, [state.quotesByKey]);

  const dirtyKeys = useMemo(() => {
    const dirty: string[] = [];
    for (const [k, q] of Object.entries(state.quotesByKey)) {
      const base = state.baselineByKey[k];
      if (!base) {
        dirty.push(k); // new
      } else if (!isEqual(base, q)) {
        dirty.push(k);
      }
    }
    return dirty;
  }, [state.quotesByKey, state.baselineByKey]);

  const hasChanges = dirtyKeys.length > 0 || state.deletedIds.size > 0;
  const { sketchId, saving, quotesByKey, deletedIds } = state;
  const saveMutation = useUpdateQuotes({
    mutationConfig: {
      onSuccess: (data) => {
        addNotification({
          type: "success",
          title: "Quotes saved",
        });
        dispatch({ type: "SAVE_SUCCESS", quotes: data });
      },
      onError: (err) => {
        dispatch({ type: "SAVE_ERROR", error: err.message });
        console.log(err);
      },
    },
  });

  const save = useCallback(async () => {
    if (!hasChanges || saving) return;

    const errorsByKey: Record<string, QuoteFieldsErrors | undefined> = {};
    for (const k of dirtyKeys) {
      const result = quoteFieldsSchema.safeParse(quotesByKey[k]);
      if (!result.success) {
        errorsByKey[k] = zodErrorToFieldErrors<typeof quoteFieldsSchema>(
          result.error,
        );
      }
    }

    dispatch({ type: "SET_ERRORS", errorsByKey });

    const hasErrors = Object.values(errorsByKey).some(Boolean);
    if (hasErrors) return;

    dispatch({ type: "SAVE_START" });

    saveMutation.mutate({
      sketchId: sketchId,
      updatedQuotes: dirtyKeys.map((k) => quotesByKey[k]),
      deletedIds: deletedIds,
    });
  }, [
    hasChanges,
    sketchId,
    saveMutation,
    saving,
    quotesByKey,
    deletedIds,
    dirtyKeys,
  ]);

  return { state, dispatch, quoteKeysSorted, dirtyKeys, hasChanges, save };
}

export type { QuoteUI, EditorState };
