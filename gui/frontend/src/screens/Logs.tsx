import { useCallback, useEffect, useRef, useState } from "react";
import { TailLogs } from "../../wailsjs/go/main/App";
import { ArrowClockwise } from "@phosphor-icons/react";
import { cn } from "../lib/cn";
import { useT } from "../lib/i18n";

const SERVICES = ["frankenphp", "mysql", "postgres", "redis", "mailpit", "adminer"];

// Logs — one tab per stack service. Fetches 200 lines on tab switch +
// has a manual refresh button. Live tailing via event stream would be
// a V2; for now a snapshot per click is plenty given the CLI does the
// same.
export function Logs() {
  const t = useT();
  const [service, setService] = useState(SERVICES[0]);
  const [text, setText] = useState("");
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const viewportRef = useRef<HTMLDivElement>(null);

  const load = useCallback(async (svc: string) => {
    setLoading(true);
    setErr(null);
    try {
      const out = await TailLogs(svc, 500);
      setText(out);
    } catch (e) {
      setErr(String(e));
      setText("");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load(service);
  }, [service, load]);

  // Scroll to the bottom when content changes so the freshest lines
  // are visible.
  useEffect(() => {
    const el = viewportRef.current;
    if (!el) return;
    el.scrollTop = el.scrollHeight;
  }, [text]);

  return (
    <div className="flex h-full flex-col p-6 animate-fade-in text-left">
      <div className="mb-4 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-left">{t("logs.title")}</h1>
          <p className="mt-1 text-left text-sm text-muted">{t("logs.subtitle")}</p>
        </div>
        <button
          type="button"
          onClick={() => load(service)}
          disabled={loading}
          className="panel-input no-drag flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm disabled:opacity-50"
        >
          <ArrowClockwise className={cn("h-3.5 w-3.5", loading && "animate-spin")} />
          {t("action.refresh")}
        </button>
      </div>

      <div className="mb-3 flex gap-1 no-drag">
        {SERVICES.map((s) => (
          <button
            key={s}
            type="button"
            onClick={() => setService(s)}
            className={cn(
              "rounded-md px-3 py-1 text-xs transition-colors",
              service === s
                ? "bg-primary/90 text-primary-foreground"
                : "bg-ink/5 text-muted hover:bg-ink/10 hover:text-fg"
            )}
          >
            {s}
          </button>
        ))}
      </div>

      {err && (
        <div className="mb-3 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-xs text-danger">
          {err}
        </div>
      )}

      <div
        ref={viewportRef}
        className="panel-card flex-1 min-h-0 overflow-y-auto rounded-lg p-3 font-mono text-[11px] leading-relaxed text-fg/75 no-drag"
      >
        {loading && text === "" ? (
          <div className="text-center text-muted py-4">loading…</div>
        ) : text === "" ? (
          <div className="text-center text-muted py-4">(no output — container might be stopped)</div>
        ) : (
          <pre className="whitespace-pre-wrap break-words">{text}</pre>
        )}
      </div>
    </div>
  );
}
