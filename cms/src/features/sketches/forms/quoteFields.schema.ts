import * as z from "zod";

import { Quote } from "@/types/api";
import { formatHMS, parseHMS } from "@/lib/utils";
import { toCastOption } from "../api/castOptionAdapters";

const hmsField = z.string().superRefine((val, ctx) => {
  try {
    parseHMS(val);
  } catch (e) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: e instanceof Error ? e.message : "Invalid timestamp",
    });
  }
});

export const quoteFieldsSchema = z.object({
  id: z.number().optional(),
  text: z.string().nonempty("Quote can't be empty"),
  startTimeMs: hmsField.nonempty("Required"),
  endTimeMs: hmsField.optional(),
  cast: z
    .array(
      z.object({
        id: z.number(),
        label: z.string(),
        image: z.string().optional(),
      }),
    )
    .nonempty("Must attribute cast member"),
  tags: z
    .array(
      z.object({
        id: z.number(),
        label: z.string(),
      }),
    )
    .optional(),
});

export type QuoteFieldsData = z.infer<typeof quoteFieldsSchema>;
export type QuoteFieldsErrors = Partial<
  Record<keyof z.infer<typeof quoteFieldsSchema> | "global", string>
>;

export function zodErrorToFieldErrors<T extends z.ZodTypeAny>(
  err: z.ZodError<z.infer<T>>,
): Partial<Record<keyof z.infer<T>, string>> {
  const out: Record<string, string> = {};

  for (const issue of err.issues) {
    const key = issue.path[0];
    if (typeof key === "string" && out[key] == null) {
      out[key] = issue.message; // first error per field
    }
  }

  return out as Partial<Record<keyof z.infer<T>, string>>;
}

export function mapQuoteToQuoteFields(q: Quote): QuoteFieldsData {
  return {
    id: q.id ?? undefined,
    text: q.text ?? "",
    startTimeMs: q.startTimeMs !== null ? formatHMS(q.startTimeMs / 1000) : "",
    endTimeMs: q.endTimeMs ? formatHMS(q.endTimeMs / 1000) : "",
    cast: q.castMembers
      ? q.castMembers
          .sort((a, b) => a.position - b.position)
          .map((c) => toCastOption(c))
      : [],
    tags: q.tags
      ? q.tags
          .sort((a, b) => a.id - b.id)
          .map((t) => ({
            id: t.id,
            label: t.category ? `${t.category.name} / ${t.name}` : t.name,
          }))
      : [],
  };
}
