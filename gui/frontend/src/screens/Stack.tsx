import { useEffect, useState } from "react";
import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime";
import {
  StackStatus,
  StartStack,
  StopStack,
  RebuildStack,
} from "../../wailsjs/go/main/App";
import { main } from "../../wailsjs/go/models";
import { cn } from "../lib/cn";
import { Play, Stop, ArrowClockwise } from "@phosphor-icons/react";
import { useT } from "../lib/i18n";

// Live stack view. Paints from a one-shot StackStatus() call on mount,
// then the Go side emits "stack:tick" every 3s and we replace the
// state. Subscribe/unsubscribe is tied to mount — so if the user
// navigates away the events still fire but get dropped on the floor.
type Row = main.ContainerStatus;

// Known services the stack is supposed to run. We render a card per
// service even if docker ps didn't return it (state=missing) so the
// user sees at a glance what's down.
const KNOWN = [
  { key: "frankenphp_devner", label: "frankenphp" },
  { key: "mysql_devner", label: "mysql" },
  { key: "postgres_devner", label: "postgres" },
  { key: "redis_devner", label: "redis" },
  { key: "mailpit_devner", label: "mailpit" },
  { key: "adminer_devner", label: "adminer" },
];

export function Stack() {
  const t = useT();
  const [rows, setRows] = useState<Row[]>([]);
  const [loading, setLoading] = useState(true);
  const [acting, setActing] = useState<null | "up" | "stop" | "rebuild">(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    StackStatus()
      .then((r) => setRows(r ?? []))
      .catch((e) => setErr(String(e)))
      .finally(() => setLoading(false));

    EventsOn("stack:tick", (data: Row[]) => setRows(data ?? []));
    return () => EventsOff("stack:tick");
  }, []);

  async function run(kind: "up" | "stop" | "rebuild") {
    setActing(kind);
    setErr(null);
    try {
      if (kind === "up") await StartStack();
      else if (kind === "stop") await StopStack();
      else await RebuildStack();
    } catch (e) {
      setErr(String(e));
    } finally {
      setActing(null);
    }
  }

  // Merge docker-reported rows with the known service list. If a row
  // isn't in docker, we still render a card with state "missing".
  const byKey = new Map(rows.map((r) => [r.name, r]));
  const services = KNOWN.map((s) => ({
    label: s.label,
    row: byKey.get(s.key),
  }));

  const runningCount = rows.filter((r) => r.state === "running").length;

  return (
    <div className="p-6 animate-fade-in text-left">
      <div className="mb-6 flex items-center gap-4">
        <div className="flex-1 text-left">
          <h1 className="text-2xl font-semibold text-left">{t("stack.title")}</h1>
          <p className="text-sm text-muted mt-1 text-left">
            {loading ? t("settings.loading") : `${runningCount}/${KNOWN.length} · ${t("stack.subtitle")}`}
          </p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => run("up")}
            disabled={acting !== null}
            className="no-drag flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-sm text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-50"
          >
            <Play className="h-3.5 w-3.5" />
            {acting === "up" ? t("stack.startDoing") : t("stack.start")}
          </button>
          <button
            onClick={() => run("stop")}
            disabled={acting !== null}
            className="panel-input no-drag flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm text-fg disabled:opacity-50"
          >
            <Stop className="h-3.5 w-3.5" />
            {acting === "stop" ? t("stack.stopDoing") : t("stack.stop")}
          </button>
          <button
            onClick={() => run("rebuild")}
            disabled={acting !== null}
            className="panel-input no-drag flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm text-fg disabled:opacity-50"
          >
            <ArrowClockwise className="h-3.5 w-3.5" />
            {acting === "rebuild" ? t("stack.rebuildDoing") : t("stack.rebuild")}
          </button>
        </div>
      </div>

      {err && (
        <div className="mb-4 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
          {err}
        </div>
      )}

      <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
        {services.map(({ label, row }) => (
          <ServiceCard key={label} label={label} row={row} />
        ))}
      </div>
    </div>
  );
}

function ServiceCard({ label, row }: { label: string; row?: Row }) {
  const state = row?.state ?? "missing";
  const dot = {
    running: "bg-success",
    restarting: "bg-accent",
    exited: "bg-danger",
    created: "bg-muted",
    paused: "bg-muted",
    missing: "bg-muted/40",
  }[state] ?? "bg-muted";

  return (
    <article className="panel-card rounded-lg p-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className={cn("h-2 w-2 rounded-full", dot)} />
          <h3 className="font-semibold">{label}</h3>
        </div>
        <span className="text-[11px] uppercase tracking-wide text-muted">
          {state}
        </span>
      </div>
      <p className="mt-2 text-xs text-muted truncate">
        {row?.image ?? "—"}
      </p>
      <p className="mt-1 text-xs text-muted">
        {row?.status ?? (state === "missing" ? "not running" : "—")}
      </p>
    </article>
  );
}
