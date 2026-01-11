import * as z from "zod";
import { Sketch } from "@/types/api";

export const sketchFormSchema = z
  .object({
    mode: z.enum(["create", "update"]),
    id: z.number().optional(),
    title: z.string().trim().min(1, "Title is required"),
    url: z.string().optional(),
    description: z.string().optional(),
    uploadDate: z.coerce.date().optional(),

    // optional in schema; required only when mode==create (below)
    thumbnail: z
      .instanceof(File)
      .refine((file) => !file || file.size <= 5_000_000, "Max file size is 5MB")
      .refine(
        (file) => ["image/jpeg", "image/png"].includes(file.type),
        "Only .jpg and .png formats are supported",
      )
      .optional(),

    creators: z
      .array(
        z.object({
          id: z.number(),
          label: z.string(),
          image: z.string(),
        }),
      )
      .optional(),
    showEpisode: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .optional(),
    episodeSketchOrder: z.string().optional(),
    recurring: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .optional(),
    series: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .optional(),
    partNumber: z.string().optional(),
  })
  .superRefine((val, ctx) => {
    if (val.mode === "create" && !val.thumbnail) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["thumbnail"],
        message: "Thumbnail is required",
      });
    }
  });

export type SketchFormData = z.infer<typeof sketchFormSchema>;

// mapper (API -> form defaults)
export function sketchToFormDefaults(
  mode: SketchFormData["mode"],
  sketch?: Sketch,
): SketchFormData {
  return {
    mode,
    id: sketch?.id,
    title: sketch?.title ?? "",
    url: sketch?.url ?? "",
    description: sketch?.description ?? "",
    uploadDate: sketch?.uploadDate ? new Date(sketch.uploadDate) : undefined,
    thumbnail: undefined,
    creators:
      sketch?.creators?.map((c) => ({
        id: c.id,
        label: c.name,
        image: c.profileImage,
      })) ?? undefined,
    showEpisode: sketch?.showEpisode
      ? {
          id: sketch.showEpisode.id,
          label: sketch.showEpisode.name,
          image: sketch.showEpisode.profileImage,
        }
      : undefined,
    episodeSketchOrder:
      sketch?.episodeSketchOrder != null
        ? String(sketch.episodeSketchOrder)
        : "",
    recurring: sketch?.recurring
      ? {
          id: sketch.recurring.id,
          label: sketch.recurring.name,
          image: sketch.recurring.thumbnail,
        }
      : undefined,
    series: sketch?.series
      ? {
          id: sketch.series.id,
          label: sketch.series.name,
          image: sketch.series.thumbnail,
        }
      : undefined,
    partNumber: sketch?.partNumber != null ? String(sketch.partNumber) : "",
  };
}
