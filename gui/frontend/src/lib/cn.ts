import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

// Helper used by every component to merge Tailwind classes safely. clsx
// handles conditional class strings; tailwind-merge resolves conflicts
// like `p-2` + `p-4` → keeps the later one.
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
