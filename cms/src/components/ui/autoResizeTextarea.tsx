import { cn } from "@/lib/utils";
import React, {
  forwardRef,
  useCallback,
  useEffect,
  useImperativeHandle,
  useMemo,
  useRef,
} from "react";

/**
 * Auto-resizing textarea with:
 * - debounced resize (layout work is batched)
 * - immediate "buffer" growth so typing never hides the caret/last line
 * - works with controlled or uncontrolled usage
 */
export type AutoResizeTextareaProps =
  React.TextareaHTMLAttributes<HTMLTextAreaElement> & {
    /** Debounce delay for the expensive resize measurement. Default: 60ms */
    resizeDebounceMs?: number;
    /**
     * Extra pixels added to computed scrollHeight so any debounce delay
     * won't hide the last typed line/caret. Default: 24px
     */
    bufferPx?: number;
    /** Optional min/max height (px). If omitted, no clamp. */
    minHeightPx?: number;
    maxHeightPx?: number;
  };

function clamp(n: number, min?: number, max?: number) {
  if (typeof min === "number") n = Math.max(min, n);
  if (typeof max === "number") n = Math.min(max, n);
  return n;
}

/** Simple debounce (no external deps). */
function debounce<T extends (...args: any[]) => void>(fn: T, wait: number) {
  let t: number | undefined;
  const debounced = (...args: Parameters<T>) => {
    if (t !== undefined) window.clearTimeout(t);
    t = window.setTimeout(() => fn(...args), wait);
  };
  debounced.cancel = () => {
    if (t !== undefined) window.clearTimeout(t);
    t = undefined;
  };
  return debounced as T & { cancel: () => void };
}

export const AutoResizeTextarea = forwardRef<
  HTMLTextAreaElement,
  AutoResizeTextareaProps
>(function AutoResizeTextarea(
  {
    resizeDebounceMs = 60,
    bufferPx = 24,
    minHeightPx,
    maxHeightPx,
    style,
    className,
    onInput,
    onChange,
    ...props
  },
  ref,
) {
  const elRef = useRef<HTMLTextAreaElement | null>(null);

  useImperativeHandle(ref, () => elRef.current as HTMLTextAreaElement, []);

  const measureAndApply = useCallback(() => {
    const el = elRef.current;
    if (!el) return;

    // Reset height so scrollHeight reflects full content (including shrink).
    el.style.height = "0px";

    // scrollHeight includes padding but not margins.
    const next = clamp(el.scrollHeight + bufferPx, minHeightPx, maxHeightPx);
    el.style.height = `${next}px`;
  }, [bufferPx, minHeightPx, maxHeightPx]);

  const debouncedMeasure = useMemo(
    () => debounce(measureAndApply, resizeDebounceMs),
    [measureAndApply, resizeDebounceMs],
  );

  // Ensure initial sizing + when external controlled value changes.
  useEffect(() => {
    measureAndApply();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.value, props.defaultValue, measureAndApply]);

  // Keep size correct on container/font changes.
  useEffect(() => {
    const el = elRef.current;
    if (!el) return;

    const ro = new ResizeObserver(() => measureAndApply());
    ro.observe(el);

    return () => ro.disconnect();
  }, [measureAndApply]);

  useEffect(() => () => debouncedMeasure.cancel(), [debouncedMeasure]);

  const handleInput: React.FormEventHandler<HTMLTextAreaElement> = (e) => {
    // Immediate "buffer" growth so user never loses last line while debounce waits.
    const el = e.currentTarget;
    const desiredImmediate = clamp(
      el.scrollHeight + bufferPx,
      minHeightPx,
      maxHeightPx,
    );

    const current = el.getBoundingClientRect().height;
    if (desiredImmediate > current) {
      el.style.height = `${desiredImmediate}px`;
    }

    // Debounced accurate measure (handles shrink + final exact height).
    debouncedMeasure();

    onInput?.(e);
  };

  const handleChange: React.ChangeEventHandler<HTMLTextAreaElement> = (e) => {
    // Some apps only use onChange; keep behavior identical.
    // onInput fires for more input types, but we support both.
    onChange?.(e);
  };

  return (
    <textarea
      {...props}
      ref={elRef}
      rows={1}
      onInput={handleInput}
      onChange={handleChange}
      className={cn(
        // Match Input's base styles
        "placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground",
        "dark:bg-input/30 border-input w-full min-w-0 rounded-md border bg-transparent",
        "px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none",
        "disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
        "focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",
        "aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",
        className, // allow overrides at the call site
      )}
      style={{
        ...style,
        overflow: "hidden", // avoid scrollbar during resizing
        resize: style?.resize ?? "none",
      }}
    />
  );
});
