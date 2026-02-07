import * as z from "zod";
import { Sketch } from "@/types/api";
import { buildImageUrl, formatHMS } from "@/lib/utils";

export const sketchFormSchema = z
  .object({
    mode: z.enum(["create", "update"]),
    id: z.number().optional(),
    title: z.string().trim().min(1, "Title is required"),
    url: z.string().optional(),
    description: z.string().optional(),
    uploadDate: z.coerce.date().optional(),
    duration: z.string().optional(),
    popularity: z.string().optional(),
    // optional in schema; required only when mode==create (below)
    thumbnail: z
      .instanceof(File)
      .refine((file) => !file || file.size <= 5_000_000, "Max file size is 5MB")
      .refine(
        (file) => ["image/jpeg"].includes(file.type),
        "Only .jpg are supported",
      )
      .optional(),
    // will need to accomadate multiple creators as an array at some point
    creator: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .optional(),
    episode: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .optional(),
    episodeSketchOrder: z.string().optional(),
    episodeStartTime: z.string().optional(),
    recurring: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .optional()
      .nullable(),
    series: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .optional()
      .nullable(),
    seriesPart: z.string().optional(),
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
  const e = sketch?.episode;
  const c = sketch?.creator;
  return {
    mode,
    id: sketch?.id ?? 0,
    title: sketch?.title ?? "",
    url: sketch?.url ?? "",
    description: sketch?.description ?? "",
    duration: formatHMS(sketch?.duration),
    popularity: sketch?.popularity ? String(sketch?.popularity) : "",
    uploadDate: sketch?.uploadDate ? new Date(sketch.uploadDate) : undefined,
    thumbnail: undefined,
    creator: c
      ? {
          id: c.id,
          label: c.name,
          image: buildImageUrl("creator", "small", c.profileImage),
        }
      : undefined,
    episode: e
      ? {
          id: e.id,
          label: `${e.season.show.name} S${e.season.number} E${e.number}`,
          image: buildImageUrl("show", "small", e.season?.show?.profileImage),
        }
      : undefined,
    episodeStartTime: formatHMS(sketch?.episodeStart),
    episodeSketchOrder:
      sketch?.episodeSketchOrder != null
        ? String(sketch.episodeSketchOrder)
        : "",
    recurring: sketch?.recurring
      ? {
          id: sketch.recurring.id,
          label: sketch.recurring.title,
          image: buildImageUrl(
            "recurring",
            "small",
            sketch.recurring.thumbnailName,
          ),
        }
      : undefined,
    series: sketch?.series
      ? {
          id: sketch.series.id,
          label: sketch.series.title,
          image: buildImageUrl("series", "small", sketch.series.thumbnailName),
        }
      : undefined,
    seriesPart: sketch?.seriesPart ? String(sketch.seriesPart) : "",
  };
}
