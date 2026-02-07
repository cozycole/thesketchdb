import { Pen } from "lucide-react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller } from "react-hook-form";
import { useNotifications } from "@/components/ui/notifications";

import { useNavigate } from "react-router";

import { makeCreatorLoadOptions } from "@/features/creators/api/creatorOptionAdapters";
import { makeEpisodeLoadOptions } from "@/features/shows/api/showOptionAdapters";
import { makeRecurringLoadOptions } from "@/features/recurring/api/recurringOptionAdapters";
import { makeSeriesLoadOptions } from "@/features/series/api/seriesOptionAdapters";

import { Sketch } from "@/types/api";

import { Button } from "@/components/ui/button";
import { DatePicker } from "@/components/ui/datePicker";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

import { AsyncSearchSelect } from "@/components/ui/asyncSearchSelect";
import { ImageUploadField } from "@/components/ui/imageUpload";

import { useUpdateSketch } from "../api/updateSketch";
import { useCreateSketch } from "../api/createSketch";

import { sketchFormSchema, sketchToFormDefaults } from "./sketchForm.schema";
import { buildImageUrl } from "@/lib/utils";

interface SketchFormProps {
  mode: "create" | "update";
  existingData?: Sketch;
}

export function SketchForm({ mode, existingData }: SketchFormProps) {
  const { addNotification } = useNotifications();
  const navigate = useNavigate();
  const defaultValues = sketchToFormDefaults(mode, existingData);

  const form = useForm({
    resolver: zodResolver(sketchFormSchema),
    defaultValues: defaultValues,
  });

  const createSketchMutation = useCreateSketch({
    mutationConfig: {
      onSuccess: (sketch) => {
        addNotification({
          type: "success",
          title: "Sketch Created",
        });
        navigate(`/admin/sketch/${sketch.id}`);
      },
      onError: (err) => {
        const data = (err as any)?.response?.data;
        const fields = data?.error;
        if (!fields) return;

        Object.entries(fields).forEach(([name, message]) => {
          form.setError(name as any, {
            type: "server",
            message: String(message),
          });
        });
      },
    },
  });
  const updateSketchMutation = useUpdateSketch({
    mutationConfig: {
      onSuccess: (sketch) => {
        addNotification({
          type: "success",
          title: "Sketch Updated",
        });

        console.log("UPDATING: ", sketch);
        form.reset(sketchToFormDefaults("update", sketch), {
          keepDefaultValues: false,
        });
      },
      onError: (err) => {
        const data = (err as any)?.response?.data;
        const fields = data?.error;
        if (!fields) return;

        Object.entries(fields).forEach(([name, message]) => {
          form.setError(name as any, {
            type: "server",
            message: String(message),
          });
        });
      },
    },
  });

  return (
    <form
      id="sketchForm"
      className="space-y-4"
      onSubmit={form.handleSubmit((values) => {
        if (mode === "update") {
          updateSketchMutation.mutate({
            data: values,
            sketchId: existingData.id,
          });
        } else {
          createSketchMutation.mutate({
            data: values,
          });
        }
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
        existingUrl={buildImageUrl(
          "sketch",
          "small",
          existingData?.thumbnailName,
        )}
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
              autoComplete="off"
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
            <FieldLabel htmlFor="uploadDate">Upload Date</FieldLabel>
            <DatePicker
              value={field.value ? new Date(String(field.value)) : undefined}
              onChange={field.onChange}
              name={field.name}
              placeholder="Select upload date"
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />
      <Controller
        name="duration"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="duration">Duration</FieldLabel>
            <div className="flex">
              <Input
                {...field}
                id="duration"
                type="text"
                className="w-16"
                aria-invalid={fieldState.invalid}
                autoComplete="off"
              />
            </div>
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />
      <Controller
        name="popularity"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="popularity">Popularity</FieldLabel>
            <div className="flex">
              <Input
                {...field}
                id="popularity"
                type="text"
                className="w-16"
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
        name="creator"
        popoverSide="top"
        label="Creator"
        loadOptions={makeCreatorLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search creators..."
      />
      <AsyncSearchSelect
        control={form.control}
        name="episode"
        popoverSide="top"
        label="Show Episode"
        loadOptions={makeEpisodeLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search episodes..."
      />
      <Controller
        name="episodeStartTime"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="episodeSketchOrder">
              Episode Start Time
            </FieldLabel>
            <div className="flex">
              <Input
                {...field}
                id="episodeSketchOrder"
                type="text"
                className="w-20"
                aria-invalid={fieldState.invalid}
                autoComplete="off"
              />
            </div>
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
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
        loadOptions={makeRecurringLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search recurring sketches..."
      />

      <AsyncSearchSelect
        control={form.control}
        name="series"
        popoverSide="top"
        label="Multipart Series"
        loadOptions={makeSeriesLoadOptions({ pageSize: 10 })}
        searchPlaceholder="Search multipart series..."
      />
      <Controller
        name="seriesPart"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor="seriesPart">Series Part Number</FieldLabel>
            <div className="flex">
              <Input
                {...field}
                id="seriesPart"
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
