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

// Accepts "m:ss", "mm:ss", or "h:mm:ss"
// Examples:
// "4:33"     -> 273
// "04:33"    -> 273
// "1:02:03"  -> 3723
export function parseHMS(input: string): number {
  if (!input) return 0;

  const parts = input
    .trim()
    .split(":")
    .map((p) => p.trim());
  if (parts.some((p) => p === "" || isNaN(Number(p)))) {
    throw new Error("Invalid time format");
  }

  if (parts.length === 2) {
    const [m, s] = parts.map(Number);
    if (s >= 60) throw new Error("Seconds must be < 60");
    return m * 60 + s;
  }

  if (parts.length === 3) {
    const [h, m, s] = parts.map(Number);
    if (m >= 60 || s >= 60) throw new Error("Minutes/seconds must be < 60");
    return h * 3600 + m * 60 + s;
  }

  throw new Error("Expected m:ss or h:mm:ss");
}

export function toYYYYMMDD(d: Date): string {
  if (!d) return "";

  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}
