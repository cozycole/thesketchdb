import { useEffect } from "react";
import { Spinner } from "@/components/ui/spinner";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller } from "react-hook-form";
import { useNotifications } from "@/components/ui/notifications";

import { makePersonLoadOptions } from "@/features/people/api/personOptionAdapters";
import { makeCharacterLoadOptions } from "@/features/characters/api/characterOptionAdapters";

import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectTrigger,
  SelectValue,
  SelectItem,
} from "@/components/ui/select";

import { AsyncSearchSelectRHF } from "@/components/ui/asyncSearchSelectRHF";
import { ImageUploadField } from "@/components/ui/imageUpload";

import { useCreateCast } from "../api/createCast";
import { useUpdateCast } from "../api/updateCast";
import { useDeleteCast } from "../api/deleteCast";
import { useCastFormStore } from "../stores/castFormStore";

import { castFormSchema, castToFormDefaults } from "./castForm.schema";
import { buildImageUrl } from "@/lib/utils";

interface CastFormProps {
  sketchId: number;
}

export function CastForm({ sketchId }: CastFormProps) {
  const { addNotification } = useNotifications();
  const {
    isOpen,
    editingCast,
    selectedThumbnail,
    selectedProfileImage,
    closeForm,
  } = useCastFormStore();

  const { control, handleSubmit, reset, setError } = useForm({
    resolver: zodResolver(castFormSchema),
  });

  const { mutate: createMutate, isPending: createPending } = useCreateCast({
    mutationConfig: {
      onSuccess: () => {
        addNotification({
          type: "success",
          title: "Cast member created",
        });
        closeForm();
      },
      onError: (err) => {
        const data = (err as any)?.response?.data;
        const fields = data?.error;
        if (!fields) return;

        Object.entries(fields).forEach(([name, message]) => {
          setError(name as any, {
            type: "server",
            message: String(message),
          });
        });
      },
    },
  });

  const { mutate: updateMutate, isPending: updatePending } = useUpdateCast({
    mutationConfig: {
      onSuccess: () => {
        addNotification({
          type: "success",
          title: "Cast updated",
        });
        closeForm();
      },
      onError: (err) => {
        const data = (err as any)?.response?.data;
        const fields = data?.error;
        if (!fields) return;

        Object.entries(fields).forEach(([name, message]) => {
          setError(name as any, {
            type: "server",
            message: String(message),
          });
        });
      },
    },
  });

  const { mutate: deleteMutate, isPending: deletePending } = useDeleteCast({
    mutationConfig: {
      onSuccess: () => {
        addNotification({
          type: "success",
          title: "Cast member deleted",
        });
        closeForm();
      },
    },
  });

  // This useEffect runs whenever the modal opens or the editingCast/selectedImages change
  // It populates the form with the appropriate data
  useEffect(() => {
    if (isOpen) {
      if (editingCast) {
        // SCENARIO 1: Editing existing cast member
        // Reset form with data from editingCast object
        const defaultValues = castToFormDefaults(editingCast);
        reset({
          id: defaultValues.id,
          characterName: defaultValues.characterName || "",
          actor: defaultValues.actor || undefined,
          character: defaultValues.character || undefined,
          castRole: defaultValues.castRole || "",
          minorRole: defaultValues.minorRole || false,
          characterThumbnail: undefined,
          characterProfile: undefined,
        });
      } else {
        // SCENARIO 2 & 3: Adding new cast member
        // Either from screenshot (with pre-selected images) or manually (all empty)
        // if it's form screenshot, the existingUrls will be filled into the image upload
        // component (see below)
        reset({
          characterName: "",
          castRole: "",
          actor: undefined,
          character: undefined,
          minorRole: false,
          characterThumbnail: undefined,
          characterProfile: undefined,
        });
      }
    }
  }, [isOpen, editingCast, selectedThumbnail, selectedProfileImage, reset]);

  const existingThumbnailUrl = editingCast
    ? buildImageUrl("cast/thumbnail", "small", selectedThumbnail)
    : selectedThumbnail
      ? buildImageUrl("cast_auto_screenshots/thumbnail", "", selectedThumbnail)
      : undefined;

  const existingProfileUrl = editingCast
    ? buildImageUrl("cast/profile", "medium", selectedProfileImage)
    : selectedProfileImage
      ? buildImageUrl("cast_auto_screenshots/profile", "", selectedProfileImage)
      : undefined;
  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && closeForm()}>
      <DialogContent
        className="sm:max-w-3xl overflow-y-scroll max-h-screen"
        onOpenAutoFocus={(e) => e.preventDefault()}
      >
        <DialogHeader>
          <DialogTitle>
            {editingCast ? "Edit Cast Member" : "Add Cast Member"}
          </DialogTitle>
        </DialogHeader>
        <form
          id="sketchForm"
          className="space-y-4"
          onSubmit={handleSubmit((values) => {
            if (values.id) {
              updateMutate({
                sketchId: sketchId,
                data: values,
              });
            } else {
              createMutate({
                sketchId: sketchId,
                data: values,
                existingThumbnail: existingThumbnailUrl,
                existingProfile: existingProfileUrl,
              });
            }
          })}
        >
          <div className="flex gap-4">
            <div className="w-72 space-y-4">
              <Controller
                name="characterName"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="sketchTitle">
                      Character Name
                    </FieldLabel>
                    <Input
                      {...field}
                      id="characterName"
                      aria-invalid={fieldState.invalid}
                      autoComplete="off"
                    />
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
                  </Field>
                )}
              />
              <AsyncSearchSelectRHF
                control={control}
                name="actor"
                label="Actor"
                loadOptions={makePersonLoadOptions({ pageSize: 8 })}
                searchPlaceholder="Search people..."
              />
              <AsyncSearchSelectRHF
                control={control}
                name="character"
                label="Character Link"
                loadOptions={makeCharacterLoadOptions({ pageSize: 8 })}
                searchPlaceholder="Search characters..."
              />
              <Controller
                name="castRole"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="castRole">Cast Role</FieldLabel>
                    <Select
                      value={field.value}
                      onValueChange={field.onChange}
                      disabled={field.disabled}
                      aria-invalid={fieldState.invalid}
                    >
                      <SelectTrigger className="w-full max-w-48">
                        <SelectValue placeholder="Select role" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          <SelectItem value="cast">Cast</SelectItem>
                          <SelectItem value="guest">Guest</SelectItem>
                          <SelectItem value="host">Host</SelectItem>
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
                  </Field>
                )}
              />
              <Controller
                name="minorRole"
                control={control}
                render={({ field, fieldState }) => (
                  <Field
                    data-invalid={fieldState.invalid}
                    orientation="horizontal"
                  >
                    <Checkbox
                      name="minorRole"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                      disabled={field.disabled}
                      aria-invalid={fieldState.invalid}
                    />
                    <FieldLabel htmlFor="minorRole">Minor Role</FieldLabel>
                    {fieldState.invalid && (
                      <FieldError errors={[fieldState.error]} />
                    )}
                  </Field>
                )}
              />
            </div>
            <div>
              <div className="min-h-72">
                <ImageUploadField
                  control={control}
                  name="characterThumbnail"
                  label="Thumbnail"
                  existingUrl={existingThumbnailUrl}
                  maxPreviewWidthPx={350}
                />
              </div>
              <ImageUploadField
                control={control}
                name="characterProfile"
                label="Profile Image"
                existingUrl={existingProfileUrl}
                maxPreviewWidthPx={250}
              />
            </div>
          </div>
          <DialogFooter className="gap-2">
            {editingCast && (
              <Button
                type="button"
                variant="destructive"
                onClick={() =>
                  deleteMutate({ sketchId, castId: editingCast.id })
                }
              >
                {deletePending ? (
                  <>
                    <Spinner /> Deleting
                  </>
                ) : (
                  "Delete"
                )}
              </Button>
            )}
            <Button type="submit" disabled={createPending || updatePending}>
              {createPending || updatePending ? (
                <>
                  <Spinner /> "Saving"
                </>
              ) : (
                "Save"
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
