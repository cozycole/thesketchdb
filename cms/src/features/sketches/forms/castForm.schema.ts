import * as z from "zod";
import { CastMember } from "@/types/api";
import { buildImageUrl } from "@/lib/utils";

export const castFormSchema = z.object({
  id: z.number().optional(),
  characterName: z.string().trim().min(1, "Character name is required"),
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
  // optional in schema; required only when mode==create (below)
  characterThumbnail: z
    .instanceof(File)
    .refine((file) => !file || file.size <= 5_000_000, "Max file size is 5MB")
    .refine(
      (file) => ["image/jpeg"].includes(file.type),
      "Only .jpg are supported",
    )
    .optional(),
  characterProfile: z
    .instanceof(File)
    .refine((file) => !file || file.size <= 5_000_000, "Max file size is 5MB")
    .refine(
      (file) => ["image/jpeg"].includes(file.type),
      "Only .jpg are supported",
    )
    .optional(),
});

export type CastFormData = z.infer<typeof castFormSchema>;

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
    characterThumbnail: undefined,
    characterProfile: undefined,
  };
}
