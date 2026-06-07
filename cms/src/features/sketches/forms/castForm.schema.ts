import * as z from "zod";
import { CastMember } from "@/types/api";
import { buildImageUrl } from "@/lib/utils";

export const castFormSchema = z
  .object({
    id: z.number().optional(),
    characterName: z.string().optional(),
    castRole: z.string().optional(),
    minorRole: z.boolean().optional(),
    actor: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .nullable()
      .optional(),
    character: z
      .object({
        id: z.number(),
        label: z.string(),
        image: z.string(),
      })
      .nullable()
      .optional(),
    cropBorder: z.boolean().optional(),
    // optional in schema; required only when mode==create (below)
    characterThumbnail: z
      .instanceof(File)
      .refine((file) => !file || file.size <= 5_000_000, "Max file size is 5MB")
      .refine(
        (file) => ["image/jpeg"].includes(file.type),
        "Only .jpg are supported",
      )
      .optional(),
    existingThumbnailUrl: z.string().optional(),
    characterProfile: z
      .instanceof(File)
      .refine((file) => !file || file.size <= 5_000_000, "Max file size is 5MB")
      .refine(
        (file) => ["image/jpeg"].includes(file.type),
        "Only .jpg are supported",
      )
      .optional(),
    existingProfileUrl: z.string().optional(),
    tags: z
      .array(
        z.object({
          id: z.number(),
          label: z.string(),
        }),
      )
      .optional(),
  })
  .superRefine((cast, ctx) => {
    const hasActor = !!cast.actor?.id;
    const hasImages =
      (cast.characterThumbnail instanceof File || cast.existingThumbnailUrl) &&
      (cast.characterProfile instanceof File || cast.existingProfileUrl);

    const hasNewCharacter = !!cast.characterName?.trim() && hasImages;
    if (!hasActor && !hasNewCharacter) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message:
          "Character name and profile/thumbnail images necessary if actor is not specified",
        path: ["global"],
      });
    }
  });

export type CastFormData = z.infer<typeof castFormSchema> & { global?: any };

// mapper (API -> form defaults)
export function castToFormDefaults(cast?: CastMember): CastFormData {
  const a = cast?.actor;
  const c = cast?.character;
  return {
    id: cast?.id ?? 0,
    characterName: cast?.characterName ?? "",
    castRole: cast?.castRole ?? "",
    minorRole: cast?.minorRole ?? false,
    actor: a
      ? {
          id: a.id,
          label: `${a.first} ${a.last}`,
          image: buildImageUrl("person", "small", a.profileImage),
        }
      : undefined,
    character: c
      ? {
          id: c.id,
          label: c.name,
          image: buildImageUrl("character", "small", c.profileImage),
        }
      : undefined,
    cropBorder: undefined,
    characterThumbnail: undefined,
    characterProfile: undefined,
    tags: cast?.tags.map((t) => ({ id: t.id, label: t.name })),
  };
}
