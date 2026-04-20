import { useCallback, useEffect, useState } from "react";
import { ListDatabases, CreateDatabase, DropDatabase } from "../../wailsjs/go/main/App";
import { main } from "../../wailsjs/go/models";
import { Database as DbIcon, Plus, Trash, CaretDown, MagnifyingGlass } from "@phosphor-icons/react";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { cn } from "../lib/cn";
import { useT } from "../lib/i18n";

const NAME_RE = /^[a-z][a-z0-9_]{0,62}$/;

export function Databases() {
  const t = useT();
  const [rows, setRows] = useState<main.DatabaseInfo[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [confirm, setConfirm] = useState<null | main.DatabaseInfo>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const [created, setCreated] = useState<main.CreateDatabaseResult | null>(null);
  const [search, setSearch] = useState("");
  const [engineFilter, setEngineFilter] = useState<"all" | "mysql" | "postgres">("all");

  const refresh = useCallback(async () => {
    try {
      const list = await ListDatabases();
      setRows(list ?? []);
    } catch (e) {
      setErr(String(e));
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const q = search.trim().toLowerCase();
  const filtered = rows.filter((db) => {
    if (engineFilter !== "all" && db.engine !== engineFilter) return false;
    if (!q) return true;
    return db.name.toLowerCase().includes(q);
  });

  const drop = async (db: main.DatabaseInfo) => {
    setErr(null);
    try {
      await DropDatabase(db.engine, db.name);
      await refresh();
    } catch (e) {
      setErr(String(e));
    }
  };

  return (
    <div className="p-6 animate-fade-in text-left">
      <div className="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-left">{t("db.title")}</h1>
          <p className="mt-1 text-left text-sm text-muted">
            {t("db.subtitle")}
          </p>
        </div>
        <button
          type="button"
          onClick={() => setCreateOpen(true)}
          className="no-drag flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-sm text-primary-foreground transition-opacity hover:opacity-90"
        >
          <Plus className="h-3.5 w-3.5" weight="bold" />
          {t("db.new")}
        </button>
      </div>

      {err && (
        <div className="mb-4 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
          {err}
        </div>
      )}

      {created && (
        <div className="mb-4 rounded-lg border border-success/40 bg-success/10 p-3 text-sm">
          <div className="font-semibold text-success">
            Created {created.engine}/{created.database}
          </div>
          <div className="mt-1 text-xs text-fg/80 font-mono">
            user: {created.user} · password: {created.password} · host: {created.host}:{created.port}
          </div>
          <button
            type="button"
            onClick={() => setCreated(null)}
            className="mt-2 text-[11px] text-muted hover:text-fg"
          >
            Dismiss
          </button>
        </div>
      )}

      {/* Search + engine filter. Same widget shape as Projects. */}
      {rows.length > 0 && (
        <div className="mb-4 flex flex-wrap items-center gap-2">
          <div className="relative flex-1 min-w-[200px]">
            <MagnifyingGlass className="pointer-events-none absolute z-10 left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder={t("db.search")}
              className="panel-input no-drag h-9 w-full rounded-md pl-8 pr-2.5 text-sm text-fg outline-none focus:border-primary/60"
            />
          </div>
          <DbFilterPill active={engineFilter === "all"} onClick={() => setEngineFilter("all")}>
            {t("projects.filter.all")}
          </DbFilterPill>
          <DbFilterPill active={engineFilter === "mysql"} onClick={() => setEngineFilter("mysql")}>
            {t("projects.filter.mysql")}
          </DbFilterPill>
          <DbFilterPill active={engineFilter === "postgres"} onClick={() => setEngineFilter("postgres")}>
            {t("projects.filter.postgres")}
          </DbFilterPill>
        </div>
      )}

      {rows.length === 0 ? (
        <div className="panel-card rounded-lg px-6 py-10 text-center text-muted">
          {t("db.empty")}
        </div>
      ) : filtered.length === 0 ? (
        <div className="panel-card rounded-lg px-6 py-10 text-center text-muted">
          {t("db.noMatch")} <code className="font-mono">{q || engineFilter}</code>.
        </div>
      ) : (
        <div className="panel-card rounded-lg overflow-hidden">
          {filtered.map((db, i) => (
            <div
              key={`${db.engine}-${db.name}`}
              className={cn(
                "flex items-center gap-3 px-4 py-2.5 text-sm",
                i > 0 && "border-t border-ink/5"
              )}
            >
              <DbIcon className="h-4 w-4 text-muted shrink-0" />
              <span className="rounded bg-ink/5 px-1.5 py-0.5 text-[11px] text-muted shrink-0">
                {db.engine}
              </span>
              <span className="flex-1 font-mono truncate">{db.name}</span>
              <button
                type="button"
                onClick={() => setConfirm(db)}
                className="text-muted hover:text-danger transition-colors"
                title="Drop"
              >
                <Trash className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>
      )}

      <ConfirmDialog
        open={confirm !== null}
        title={confirm ? t("db.confirm.drop", { engine: confirm.engine, name: confirm.name }) : ""}
        description={t("db.confirm.dropDescription")}
        confirmLabel={t("action.drop")}
        cancelLabel={t("action.cancel")}
        danger
        onCancel={() => setConfirm(null)}
        onConfirm={() => {
          if (!confirm) return;
          const db = confirm;
          setConfirm(null);
          drop(db);
        }}
      />

      <CreateDatabaseModal
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={(r) => {
          setCreated(r);
          refresh();
        }}
      />
    </div>
  );
}

function DbFilterPill({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "no-drag rounded-md px-3 h-9 text-xs transition-colors",
        active
          ? "bg-primary/90 text-primary-foreground"
          : "panel-input text-muted hover:text-fg"
      )}
    >
      {children}
    </button>
  );
}

function CreateDatabaseModal({
  open,
  onClose,
  onCreated,
}: {
  open: boolean;
  onClose: () => void;
  onCreated: (r: main.CreateDatabaseResult) => void;
}) {
  const t = useT();
  const [engine, setEngine] = useState<"mysql" | "postgres">("mysql");
  const [name, setName] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    setName("");
    setErr(null);
    setSubmitting(false);
  }, [open]);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape" && !submitting) onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, submitting, onClose]);

  if (!open) return null;
  const nameValid = NAME_RE.test(name);

  const submit = async () => {
    if (!nameValid) return;
    setSubmitting(true);
    setErr(null);
    try {
      const r = await CreateDatabase(engine, name);
      if (r) onCreated(r);
      onClose();
    } catch (e) {
      setErr(String(e));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in no-drag"
      onMouseDown={() => !submitting && onClose()}
    >
      <div
        className="panel-card w-full max-w-sm animate-slide-up rounded-lg p-5 shadow-2xl text-left"
        onMouseDown={(e) => e.stopPropagation()}
      >
        <h2 className="mb-4 text-base font-semibold">{t("db.new.title")}</h2>

        <div className="space-y-3">
          <label className="block">
            <div className="mb-1 text-[11px] uppercase tracking-wide text-muted">{t("db.new.engine")}</div>
            <div className="relative">
              <select
                value={engine}
                onChange={(e) => setEngine(e.target.value as any)}
                disabled={submitting}
                className="panel-input no-drag w-full h-9 appearance-none rounded-md px-2.5 pr-8 text-sm text-fg outline-none focus:border-primary/60 cursor-pointer"
              >
                <option value="mysql" className="bg-panel">MySQL</option>
                <option value="postgres" className="bg-panel">PostgreSQL</option>
              </select>
              <CaretDown className="pointer-events-none absolute z-10 right-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted" />
            </div>
          </label>

          <label className="block">
            <div className="mb-1 flex items-center justify-between text-[11px] uppercase tracking-wide text-muted">
              <span>{t("db.new.name")}</span>
              {name && !nameValid && (
                <span className="normal-case tracking-normal text-danger">
                  {t("db.new.nameHint")}
                </span>
              )}
            </div>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              disabled={submitting}
              placeholder="mydb"
              autoFocus
              className={cn(
                "panel-input w-full h-9 rounded-md px-2.5 text-sm text-fg outline-none focus:border-primary/60 disabled:opacity-50",
                name && !nameValid && "border-danger/60 focus:border-danger"
              )}
            />
          </label>
        </div>

        {err && (
          <div className="mt-3 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-xs text-danger">
            {err}
          </div>
        )}

        <div className="mt-5 flex justify-end gap-2">
          <button
            type="button"
            onClick={onClose}
            disabled={submitting}
            className="panel-input rounded-md px-3 py-1.5 text-sm disabled:opacity-50"
          >
            {t("action.cancel")}
          </button>
          <button
            type="button"
            onClick={submit}
            disabled={submitting || !nameValid}
            className="rounded-md bg-primary px-3 py-1.5 text-sm text-primary-foreground hover:opacity-90 disabled:opacity-50"
          >
            {submitting ? t("action.creating") : t("action.create")}
          </button>
        </div>
      </div>
    </div>
  );
}
