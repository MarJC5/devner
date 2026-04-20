import animate from "tailwindcss-animate";
import typography from "@tailwindcss/typography";

/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["class"],
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        // Theme-aware palette. Values come from CSS vars defined in
        // App.css (:root for dark, :root.light for light). The
        // `<alpha-value>` token lets Tailwind build `bg-fg/50`,
        // `text-muted/80`, etc. automatically.
        bg: "rgb(var(--bg) / <alpha-value>)",
        panel: "rgb(var(--panel) / <alpha-value>)",
        panel2: "rgb(var(--panel2) / <alpha-value>)",
        muted: "rgb(var(--muted) / <alpha-value>)",
        fg: "rgb(var(--fg) / <alpha-value>)",
        // Theme-inverting tint. On dark, "ink" is white → `bg-ink/5`
        // reads as a subtle white veil. On light, it flips to near-
        // black, so the same class gives the equivalent darker veil.
        ink: "rgb(var(--ink) / <alpha-value>)",
        "border-soft": "rgb(var(--border) / <alpha-value>)",

        // Primary + status colors are theme-agnostic (fixed hex).
        primary: {
          DEFAULT: "#4A6CF7",
          foreground: "#FFFFFF",
          muted: "#F0F4FF",
          faint: "#A0ACEE",
          watermark: "#3A58D8",
        },
        accent: "#FFCC00",
        success: "#00C853",
        danger: "#FF5577",
      },
      fontFamily: {
        sans: ["Inter", "ui-sans-serif", "system-ui", "-apple-system", "Segoe UI", "Roboto", "sans-serif"],
        mono: ["JetBrains Mono", "ui-monospace", "SFMono-Regular", "Menlo", "monospace"],
      },
      borderRadius: {
        lg: "12px",
        md: "10px",
        sm: "8px",
      },
      keyframes: {
        "fade-in": { from: { opacity: 0 }, to: { opacity: 1 } },
        "slide-up": { from: { transform: "translateY(8px)", opacity: 0 }, to: { transform: "translateY(0)", opacity: 1 } },
      },
      animation: {
        "fade-in": "fade-in 160ms ease-out",
        "slide-up": "slide-up 220ms ease-out",
      },
    },
  },
  plugins: [animate, typography],
};
