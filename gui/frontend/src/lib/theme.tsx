import { createContext, useContext, useEffect, useState, type ReactNode } from "react";

// Theme system. Three user-facing modes:
//   - "light"  → force light
//   - "dark"   → force dark
//   - "system" → follow the OS (via prefers-color-scheme media query)
// The *applied* theme is always light or dark — "system" is resolved
// once at mount and re-resolved when the OS preference changes.
//
// The user's mode choice persists across reboots via localStorage
// (webview-local, per-app). Default is "system" so the app blends
// with the rest of the OS out of the box.

export type ThemeMode = "light" | "dark" | "system";
export type AppliedTheme = "light" | "dark";

const STORAGE_KEY = "devner.theme";

type Ctx = {
  mode: ThemeMode;
  setMode: (m: ThemeMode) => void;
  applied: AppliedTheme;
};

const ThemeCtx = createContext<Ctx | null>(null);

function readInitial(): ThemeMode {
  try {
    const v = localStorage.getItem(STORAGE_KEY);
    if (v === "light" || v === "dark" || v === "system") return v;
  } catch {
    // localStorage may be blocked in some webview configs — fall
    // through to the default.
  }
  return "system";
}

function systemPrefers(): AppliedTheme {
  if (typeof window === "undefined") return "dark";
  return window.matchMedia?.("(prefers-color-scheme: light)").matches ? "light" : "dark";
}

function resolve(mode: ThemeMode): AppliedTheme {
  return mode === "system" ? systemPrefers() : mode;
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setModeState] = useState<ThemeMode>(readInitial);
  const [applied, setApplied] = useState<AppliedTheme>(() => resolve(readInitial()));

  // Apply the theme as a class on <html>. We toggle `.light` (not
  // `.dark`) because our CSS defaults to dark on :root and only
  // overrides tokens when .light is present.
  useEffect(() => {
    const next = resolve(mode);
    setApplied(next);
    document.documentElement.classList.toggle("light", next === "light");
    try {
      localStorage.setItem(STORAGE_KEY, mode);
    } catch {
      // ignore
    }
  }, [mode]);

  // Follow the OS when "system" is selected. Re-resolve on media
  // query change; no-op when the user forced light/dark.
  useEffect(() => {
    if (mode !== "system") return;
    const mql = window.matchMedia?.("(prefers-color-scheme: light)");
    if (!mql) return;
    const onChange = () => {
      const next: AppliedTheme = mql.matches ? "light" : "dark";
      setApplied(next);
      document.documentElement.classList.toggle("light", next === "light");
    };
    mql.addEventListener("change", onChange);
    return () => mql.removeEventListener("change", onChange);
  }, [mode]);

  return (
    <ThemeCtx.Provider value={{ mode, setMode: setModeState, applied }}>
      {children}
    </ThemeCtx.Provider>
  );
}

export function useTheme(): Ctx {
  const v = useContext(ThemeCtx);
  if (!v) throw new Error("useTheme must be used inside <ThemeProvider>");
  return v;
}
