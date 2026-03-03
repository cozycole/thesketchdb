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
  anchorTranscriptId: number;
  selectedQuoteKeys: Set<string>;

  saving: boolean;
  error?: string;
};

type Action =
  | { type: "LOAD"; transcript: TranscriptLine[]; quotes: Quote[] }
  | { type: "ADD_QUOTES"; quotes: QuoteUI[] }
  | { type: "ADD_FROM_TRANSCRIPT"; transcriptIds: number[] }
  | { type: "UPDATE_QUOTE"; new: QuoteUI }
  | { type: "DELETE_QUOTE"; key: string }
  | { type: "VALIDATE_QUOTE"; key: string }
  | {
      type: "SET_ERRORS";
      errorsByKey: Record<string, QuoteFieldsErrors | undefined>;
    }
  | { type: "SELECT_TRANSCRIPT"; id: number; e: MouseEvent }
  | { type: "SELECT_QUOTE"; key: string; e: PointerEvent }
  | { type: "CLEAR_SELECTED" }
  | { type: "MERGE_QUOTES" }
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
        transcript: action.transcript.sort(
          (a, b) => a.lineNumber - b.lineNumber,
        ),
        quotesByKey,
        baselineByKey,
        deletedIds: new Set(),
        error: undefined,
      };
    }

    case "ADD_QUOTES": {
      // this assumes the quote has been validated
      const quotesByKey = { ...state.quotesByKey };

      for (const q of action.quotes) {
        q.clientId = crypto.randomUUID();
        q.startMs = parseHMS(q.startTimeMs);
        quotesByKey[q.clientId] = q;
      }
      return {
        ...state,
        quotesByKey,
      };
    }

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
          startMs: q.startTimeMs ? q.startTimeMs / 1000 : 0,
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

    case "SELECT_QUOTE": {
      const next = new Set(state.selectedQuoteKeys);

      if (action.e.ctrlKey) {
        next.add(action.key);
      }

      return {
        ...state,
        selectedQuoteKeys: next,
      };
    }

    case "SELECT_TRANSCRIPT": {
      const orderedIds = state.transcript.map((t) => t.id);

      const next = new Set<number>(state.selectedTranscriptIds);

      const clickedId = action.id;
      const prevAnchor = state.anchorTranscriptId;

      const isToggle = action.e ? action.e.ctrlKey || action.e.metaKey : false; // cmd on mac
      const isRange = action.e ? action.e.shiftKey : false;

      // shift: select inclusive range between anchor and clicked (based on transcript order)
      if (isRange && prevAnchor != null) {
        const a = orderedIds.indexOf(prevAnchor);
        const b = orderedIds.indexOf(clickedId);

        // if either id isn't found, fall back to single-select
        if (a !== -1 && b !== -1) {
          const [lo, hi] = a < b ? [a, b] : [b, a];
          next.clear();
          for (let i = lo; i <= hi; i++) next.add(orderedIds[i]!);

          return {
            ...state,
            selectedTranscriptIds: next,
            // keep anchor stable during shift so repeated shift-click extends from same anchor
            anchorTranscriptId: prevAnchor,
          };
        }
      }

      // ctrl/cmd: toggle clicked id without clearing
      if (isToggle) {
        if (next.has(clickedId)) next.delete(clickedId);
        else next.add(clickedId);

        return {
          ...state,
          selectedTranscriptIds: next,
          // update anchor to last interacted item
          anchorTranscriptId: clickedId,
        };
      }

      // plain click: single selection
      next.clear();
      next.add(clickedId);

      return {
        ...state,
        selectedTranscriptIds: next,
        anchorTranscriptId: clickedId,
      };
    }

    case "CLEAR_SELECTED": {
      return {
        ...state,
        selectedQuoteKeys: new Set(),
        selectedTranscriptIds: new Set(),
      };
    }

    case "MERGE_QUOTES": {
      const keys = Array.from(state.selectedQuoteKeys);
      if (keys.length < 2) return state;

      // Get selected quote objects
      const selected = keys.map((k) => state.quotesByKey[k]).filter(Boolean);

      if (selected.length < 2) return state;

      // Sort by start time (ascending)
      selected.sort((a, b) => Number(a.startMs) - Number(b.startMs));

      // Compute merged fields
      const startTimeMs = selected[0].startTimeMs;
      const startTime = selected[0].startMs;

      const endTimes = selected
        .map((q) => q.endTimeMs)
        .filter((v): v is string => Boolean(v));

      const endTimeMs =
        endTimes.length > 0
          ? endTimes.sort((a, b) => Number(a) - Number(b))[endTimes.length - 1]
          : undefined;

      const text = selected.map((q) => q.text).join(" ");

      // Merge cast (dedupe by id)
      const castMap = new Map<
        number,
        { id: number; label: string; image: string }
      >();
      for (const q of selected) {
        for (const c of q.cast ?? []) {
          castMap.set(c.id, c);
        }
      }

      const cast = Array.from(castMap.values());

      // Merge tags (dedupe by id)
      const tagMap = new Map<number, { id: number; label: string }>();
      for (const q of selected) {
        for (const t of q.tags ?? []) {
          tagMap.set(t.id, t);
        }
      }

      const tags = Array.from(tagMap.values());

      // Create merged quote (new clientId)
      const mergedKey = crypto.randomUUID();

      const mergedQuote = {
        id: undefined,
        text,
        startMs: startTime,
        startTimeMs,
        endTimeMs,
        cast,
        tags,
        clientId: mergedKey,
      };
      console.log("MERGED QUOTE");
      console.log(mergedQuote);

      // Build new quotesByKey
      const nextQuotesByKey = { ...state.quotesByKey };

      // Remove old
      for (const k of keys) {
        delete nextQuotesByKey[k];
      }

      // Add to deleted IDs
      const nextDeletedIds = new Set(state.deletedIds);
      for (const q of selected) {
        if (q.id) nextDeletedIds.add(q.id);
      }

      // Insert merged
      nextQuotesByKey[mergedKey] = mergedQuote;

      return {
        ...state,
        quotesByKey: nextQuotesByKey,
        selectedQuoteKeys: new Set(),
        deletedIds: nextDeletedIds,
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
    | "quotesByKey"
    | "baselineByKey"
    | "deletedIds"
    | "errorsByKey"
    | "quoteKeys"
    | "anchorTranscriptId"
  >,
) {
  const { addNotification } = useNotifications();
  const [state, dispatch] = useReducer(reducer, {
    ...initial,
    errorsByKey: {},
    quotesByKey: {},
    baselineByKey: {},
    deletedIds: new Set<number>(),
    anchorTranscriptId: undefined,
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
