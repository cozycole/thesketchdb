import * as React from "react";
import {
  Controller,
  useWatch,
  type Control,
  type FieldValues,
  type Path,
} from "react-hook-form";

import { Button } from "@/components/ui/button";
import {
  Field,
  FieldLabel,
  FieldDescription,
  FieldError,
} from "@/components/ui/field";

type ImageUploadFieldProps<T extends FieldValues> = {
  control: Control<T>;
  name: Path<T>;

  label?: string;
  description?: string;

  /** Existing image URL from API (update mode), e.g. sketch.thumbnail.small */
  existingUrl?: string | null;

  /** Width cap for preview container */
  maxPreviewWidthPx?: number;

  accept?: string; // default "image/*"
  disabled?: boolean;
};

export function ImageUploadField<T extends FieldValues>({
  control,
  name,
  label = "Image",
  description,
  existingUrl = null,
  maxPreviewWidthPx = 420,
  accept = "image/jpeg",
  disabled = false,
}: ImageUploadFieldProps<T>) {
  const inputRef = React.useRef<HTMLInputElement | null>(null);
  const [filePreviewUrl, setPreviewUrl] = React.useState<string | null>(null);
  const [imageError, setImageError] = React.useState<string | null>(null);
  const [loadingState, setLoadingState] = React.useState<
    "idle" | "loading" | "loaded" | "error"
  >("idle");

  const file: File | null = useWatch({ control, name }) ?? null;

  React.useEffect(() => {
    if (!file) {
      setPreviewUrl(null);
      setImageError(null);
      setLoadingState("idle");
      return;
    }

    if (!file.type.startsWith("image/")) {
      setImageError(`Invalid file type: ${file.type || "unknown"}`);
      setLoadingState("error");
      return;
    }

    const MAX_SIZE = 10 * 1024 * 1024;
    if (file.size > MAX_SIZE) {
      setImageError(
        `File too large: ${(file.size / 1024 / 1024).toFixed(2)}MB (max 10MB)`,
      );
      setLoadingState("error");
      return;
    }

    setLoadingState("loading");
    setImageError(null);

    try {
      const url = URL.createObjectURL(file);
      setPreviewUrl(url);

      return () => {
        URL.revokeObjectURL(url);
      };
    } catch (err) {
      console.error("Error creating blob URL:", err);
      setImageError("Failed to create preview");
      setLoadingState("error");
    }
  }, [file]);

  const displayUrl = filePreviewUrl ?? existingUrl ?? null;

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => {
        const file = (field.value ?? null) as File | null;

        return (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor={`${String(name)}-file`}>{label}</FieldLabel>

            {description ? (
              <FieldDescription>{description}</FieldDescription>
            ) : null}

            <input
              id={`${String(name)}-file`}
              ref={inputRef}
              type="file"
              accept={accept}
              className="hidden"
              disabled={disabled}
              onChange={(e) => {
                const f = e.target.files?.[0] ?? null;
                setImageError(null);
                setLoadingState("idle");
                field.onChange(f); // store File | null in RHF
              }}
              onBlur={field.onBlur}
              name={field.name}
            />
            <div>
              <div className="mt-2 flex items-center gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  disabled={disabled}
                  onClick={() => inputRef.current?.click()}
                >
                  {displayUrl ? "Replace image" : "Upload image"}
                </Button>
                {existingUrl && file ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    disabled={disabled}
                    onClick={() => {
                      field.onChange(null);
                      setImageError(null);
                      setLoadingState("idle");
                      if (inputRef.current) inputRef.current.value = "";
                    }}
                  >
                    Keep existing
                  </Button>
                ) : null}
              </div>

              {imageError && (
                <div className="mt-2 p-2 bg-red-50 border border-red-200 rounded-md text-sm text-red-600">
                  {imageError}
                </div>
              )}

              {displayUrl ? (
                <div className="mt-2 w-fit rounded-md p-2">
                  <img
                    src={displayUrl}
                    alt={file ? "New image preview" : "Current image"}
                    style={{
                      maxWidth: `${maxPreviewWidthPx}px`,
                      width: "100%",
                      height: "auto",
                    }}
                    className="block border p-2 rounded-md"
                    onError={(e) => {
                      console.error("Image load error:", {
                        src: e.currentTarget.src,
                        error: e,
                      });
                      setImageError(
                        "Failed to load image preview - file may be corrupted or have an unsupported format",
                      );
                      setLoadingState("error");
                    }}
                  />
                  <div className="mt-2 text-xs text-muted-foreground">
                    {file ? (
                      <>
                        New image selected (not saved yet).
                        <br />
                        File: {file.name} ({(file.size / 1024).toFixed(1)}KB)
                      </>
                    ) : (
                      "Current image."
                    )}
                  </div>
                </div>
              ) : (
                <div className="mt-2 text-sm text-muted-foreground">
                  No image set.
                </div>
              )}
            </div>

            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        );
      }}
    />
  );
}
