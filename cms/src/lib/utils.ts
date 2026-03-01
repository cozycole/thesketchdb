import { env } from "@/config/env";

import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function buildImageUrl(
  type: string,
  size?: "small" | "medium" | "large" | "",
  name?: string | null,
) {
  if (!name) return null;
  if (size) {
    return `${env.BUCKET_URL}/${type}/${size}/${name}`;
  }
  return `${env.BUCKET_URL}/${type}/${name}`;
}

export function formatHMS(total: number): string {
  if (total === null || total === undefined) {
    return "";
  }
  const h = Math.floor(total / 3600);
  const m = Math.floor((total % 3600) / 60);
  const s = total % 60;

  if (h > 0) {
    return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
  }
  return `${m}:${s.toString().padStart(2, "0")}`;
}

export function parseHMS(input: string): number | undefined {
  if (!input || !input.trim()) {
    return undefined;
  }
  input = input.trim();

  const value = input.trim();

  // mm:ss  (e.g. 0:35, 12:05)
  const mmss = /^(\d+):(\d{2})$/;

  // hh:mm:ss (e.g. 1:00:00, 12:34:56)
  const hhmmss = /^(\d+):(\d{2}):(\d{2})$/;

  let match: RegExpMatchArray | null;

  if ((match = value.match(mmss))) {
    const minutes = Number(match[1]);
    const seconds = Number(match[2]);

    if (seconds >= 60) {
      throw new Error("Seconds must be < 60");
    }

    return minutes * 60 + seconds;
  }

  if ((match = value.match(hhmmss))) {
    const hours = Number(match[1]);
    const minutes = Number(match[2]);
    const seconds = Number(match[3]);

    if (minutes >= 60 || seconds >= 60) {
      throw new Error("Minutes and seconds must be < 60");
    }

    return hours * 3600 + minutes * 60 + seconds;
  }

  throw new Error("Invalid time format. Expected mm:ss or hh:mm:ss");
}

export function toYYYYMMDD(d: Date): string {
  if (!d) return "";

  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}
