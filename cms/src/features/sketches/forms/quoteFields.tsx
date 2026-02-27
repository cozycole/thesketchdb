import { Trash } from "lucide-react";
import { Button } from "@/components/ui/button";

import { QuoteFieldsData } from "./quoteFields.schema";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { AutoResizeTextarea } from "@/components/ui/autoResizeTextarea";

type QuoteFieldsProps<T extends QuoteFieldsData = QuoteFieldsData> = {
  value: T;
  onChange: (next: T) => void;
  onDelete: (q: T) => void;
  onBlurQuote?: (q: T) => void;
  errors?: Partial<Record<keyof QuoteFieldsData, string>>;
  disabled?: boolean;
};

export function QuoteFields<T extends QuoteFieldsData = QuoteFieldsData>({
  value,
  onChange,
  onDelete,
  onBlurQuote,
  errors,
  disabled,
}: QuoteFieldsProps<T>) {
  const set = <K extends keyof QuoteFieldsData>(
    key: K,
    v: QuoteFieldsData[K],
  ) => onChange({ ...value, [key]: v });
  return (
    <div className="flex flex-col gap-4 p-4 rounded-lg border-gray-200 border-2">
      <div className="flex gap-4">
        <Field>
          <FieldLabel>Start Timestamp</FieldLabel>
          <Input
            type="text"
            className="w-32 bg-white"
            value={value.startTimeMs ?? 0}
            disabled={disabled}
            onChange={(e) => set("startTimeMs", e.target.value)}
            onBlur={() => onBlurQuote?.(value)}
          />
          {errors?.startTimeMs && <FieldError>{errors.startTimeMs}</FieldError>}
        </Field>

        <Field>
          <FieldLabel>End Timestamp</FieldLabel>
          <Input
            type="text"
            className="w-32 bg-white"
            value={value.endTimeMs ?? ""}
            disabled={disabled}
            onChange={(e) =>
              set(
                "endTimeMs",
                e.target.value === "" ? undefined : e.target.value,
              )
            }
            onBlur={() => onBlurQuote?.(value)}
          />
          {errors?.endTimeMs && <FieldError>{errors.endTimeMs}</FieldError>}
        </Field>
        {value["clientId"] && (
          <Button variant="destructive" onClick={() => onDelete(value)}>
            <Trash />
          </Button>
        )}
      </div>

      <Field>
        <FieldLabel>Quote</FieldLabel>
        <AutoResizeTextarea
          className="w-full bg-white"
          placeholder="Insert quote..."
          value={value.text ?? ""}
          disabled={disabled}
          onChange={(e) => set("text", e.target.value)}
          onBlur={() => onBlurQuote?.(value)}
        />
        {errors?.text && <FieldError>{errors.text}</FieldError>}
      </Field>

      {/* cast + tags pickers can be controlled the same way:
          set("cast", nextCast) / set("tags", nextTags) */}
    </div>
  );
}
