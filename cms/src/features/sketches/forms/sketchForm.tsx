import { Pen } from "lucide-react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller } from "react-hook-form";
import { useNotifications } from "@/components/ui/notifications";

import { makeSketchLoadOptions } from "../api/sketchOptionAdapters";
import { makeCreatorLoadOptions } from "@/features/creators/api/creatorOptionAdapters";
import { makeEpisodeLoadOptions } from "@/features/shows/api/showOptionAdapters";

import { Sketch } from "@/types/api";

import { Button } from "@/components/ui/button";
import { DatePicker } from "@/components/ui/datePicker";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

import { AsyncSearchSelect } from "@/components/ui/asyncSearchSelect";
import { ImageUploadField } from "@/components/ui/imageUpload";

import { useUpdateSketch } from "../api/updateSketch";

import { sketchFormSchema, sketchToFormDefaults } from "./sketchForm.schema";

interface SketchFormProps {
  mode: "create" | "update";
  existingData?: Sketch;
}

export function SketchForm({ mode, existingData }: SketchFormProps) {
  const { addNotification } = useNotifications();
  const updateSketchMutation = useUpdateSketch({
    mutationConfig: {
      onSuccess: () => {
        addNotification({
          type: "success",
          title: mode === "create" ? "Sketch Created" : "Sketch Updated",
        });
      },
    },
  });

  const defaultValues = sketchToFormDefaults(mode, existingData);

  const form = useForm({
    resolver: zodResolver(sketchFormSchema),
    defaultValues: defaultValues,
  });

  return (
    <form
      id="sketchForm"
      className="space-y-4"
      onSubmit={form.handleSubmit((values) => {
        updateSketchMutation.mutate({
          data: values,
          sketchId: existingData.id,
        });
      })}
    >
      <Controller
        name="title"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="sketchTitle">Title</FieldLabel>
            <Input
              {...field}
              id="sketchTitle"
              aria-invalid={fieldState.invalid}
              autoComplete="off"
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />
      <Controller
        name="url"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="sketchUrl">URL</FieldLabel>
            <Input
              {...field}
              id="sketchUrl"
              aria-invalid={fieldState.invalid}
              autoComplete="off"
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />
      <ImageUploadField
        control={form.control}
        name="thumbnail"
        label="Thumbnail"
        existingUrl={existingData?.thumbnailUrl ?? null}
        maxPreviewWidthPx={420}
      />
      <Controller
        name="description"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="sketchDescription">Description</FieldLabel>
            <Textarea
              {...field}
              id="sketchDescription"
              aria-invalid={fieldState.invalid}
              className="min-h-[120px]"
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />
      <Controller
        name="uploadDate"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="sketchUploadDate">Upload Date</FieldLabel>
            <DatePicker
              value={field.value ? new Date(field.value) : undefined}
              onChange={field.onChange}
              name={field.name}
              placeholder="Select upload date"
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />

      <AsyncSearchSelect
        control={form.control}
        name="creators"
        multiple={true}
        popoverSide="top"
        label="Creators"
        loadOptions={makeCreatorLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search creators..."
      />
      <AsyncSearchSelect
        control={form.control}
        name="showEpisode"
        popoverSide="top"
        label="Show Episode"
        loadOptions={makeEpisodeLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search episodes..."
      />
      <Controller
        name="episodeSketchOrder"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="episodeSketchOrder">
              Episode Sketch Order
            </FieldLabel>
            <div className="flex">
              <Input
                {...field}
                id="episodeSketchOrder"
                type="number"
                className="w-20"
                aria-invalid={fieldState.invalid}
                autoComplete="off"
              />
            </div>
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />
      <AsyncSearchSelect
        control={form.control}
        name="recurring"
        popoverSide="top"
        label="Recurring Sketch"
        loadOptions={makeSketchLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search recurring sketches..."
      />

      <AsyncSearchSelect
        control={form.control}
        name="series"
        popoverSide="top"
        label="Multipart Series"
        loadOptions={makeSketchLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search multipart series..."
      />
      <Controller
        name="partNumber"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="partNumber">Series Part Number</FieldLabel>
            <div className="flex">
              <Input
                {...field}
                id="partNumber"
                type="number"
                className="w-20"
                aria-invalid={fieldState.invalid}
                autoComplete="off"
              />
            </div>
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />
      <Button size="sm" className="text-white" type="submit" form="sketchForm">
        <Pen className="size-4" />{" "}
        {mode === "create" ? "Create Sketch" : "Update Sketch"}
      </Button>
    </form>
  );
}
