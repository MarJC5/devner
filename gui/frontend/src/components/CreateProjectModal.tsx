import { useEffect, useMemo, useState } from "react";
import { CreateProject } from "../../wailsjs/go/main/App";
import { app } from "../../wailsjs/go/models";
import { X, CaretDown } from "@phosphor-icons/react";
import { cn } from "../lib/cn";

// Create project — modal with a per-type form. We reuse the shared
// app.CreateProjectRequest DTO that the CLI sends so the backend flow
// (scaffold → DB → post-setup → Caddy) is identical whether the
// project is created here or via `devner new`.

type Props = {
  open: boolean;
  onClose: () => void;
  onCreated: (name: string) => void;
};

const PROJECT_TYPES = [
  "wordpress",
  "laravel",
  "node",
  "nextjs",
  "nuxt",
  "astro",
  "sveltekit",
  "vite",
] as const;
type ProjectType = (typeof PROJECT_TYPES)[number];

const TEMPLATES_BY_TYPE: Partial<Record<ProjectType, string[]>> = {
  vite: ["react-ts", "vue-ts", "svelte-ts", "solid-ts", "preact-ts", "qwik-ts", "lit-ts", "vanilla-ts"],
  astro: ["minimal", "basics", "blog", "portfolio", "starlight"],
  sveltekit: ["skeleton", "minimal", "demo"],
};

const NAME_RE = /^[a-z][a-z0-9-]{0,62}$/;

export function CreateProjectModal({ open, onClose, onCreated }: Props) {
  const [type, setType] = useState<ProjectType>("laravel");
  const [name, setName] = useState("");
  const [db, setDb] = useState<"" | "mysql" | "postgres">("");
  const [template, setTemplate] = useState("");
  const [wpInstall, setWpInstall] = useState(false);
  const [wpTitle, setWpTitle] = useState("");
  const [wpAdmin, setWpAdmin] = useState("admin");
  const [wpPass, setWpPass] = useState("admin");
  const [wpEmail, setWpEmail] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const supportsTemplate = !!TEMPLATES_BY_TYPE[type];
  const isWordPress = type === "wordpress";
  const nameValid = NAME_RE.test(name);

  // Reset on open so stale state from a previous session doesn't leak.
  useEffect(() => {
    if (!open) return;
    setName("");
    setDb("");
    setTemplate("");
    setWpInstall(false);
    setErr(null);
    setSubmitting(false);
  }, [open]);

  // Reset template when type changes (previous template may not apply).
  useEffect(() => {
    setTemplate("");
  }, [type]);

  // Esc / Enter keyboard shortcuts at the modal scope.
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape" && !submitting) onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, submitting, onClose]);

  const submit = async () => {
    if (!nameValid) return;
    setSubmitting(true);
    setErr(null);
    try {
      const req: app.CreateProjectRequest = {
        Name: name,
        Type: type,
        DBEngine: db,
        Template: supportsTemplate ? template : "",
        WPInstall: isWordPress && wpInstall,
        WPTitle: wpTitle,
        WPAdmin: wpAdmin,
        WPPassword: wpPass,
        WPEmail: wpEmail,
      };
      await CreateProject(req);
      onCreated(name);
      onClose();
    } catch (e) {
      setErr(String(e));
    } finally {
      setSubmitting(false);
    }
  };

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in no-drag"
      onMouseDown={() => !submitting && onClose()}
    >
      <div
        className="panel-card w-full max-w-md animate-slide-up rounded-lg p-5 shadow-2xl text-left"
        onMouseDown={(e) => e.stopPropagation()}
      >
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-base font-semibold">New project</h2>
          <button
            onClick={onClose}
            disabled={submitting}
            className="text-muted hover:text-fg disabled:opacity-40"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="space-y-3">
          <Field label="Type">
            <Select
              value={type}
              onChange={(v) => setType(v as ProjectType)}
              disabled={submitting}
            >
              {PROJECT_TYPES.map((t) => (
                <option key={t} value={t} className="bg-panel">
                  {t}
                </option>
              ))}
            </Select>
          </Field>

          <Field label="Name" hint={nameHint(name)}>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              disabled={submitting}
              placeholder="myapp"
              className={cn(
                inputCls,
                name && !nameValid && "border-danger/60 focus:border-danger"
              )}
              autoFocus
            />
          </Field>

          <Field label="Database">
            <Select
              value={db}
              onChange={(v) => setDb(v as typeof db)}
              disabled={submitting}
            >
              <option value="" className="bg-panel">
                None
              </option>
              <option value="mysql" className="bg-panel">MySQL</option>
              <option value="postgres" className="bg-panel">PostgreSQL</option>
            </Select>
          </Field>

          {supportsTemplate && (
            <Field label="Template">
              <Select
                value={template}
                onChange={setTemplate}
                disabled={submitting}
              >
                <option value="" className="bg-panel">
                  default
                </option>
                {TEMPLATES_BY_TYPE[type]!.map((t) => (
                  <option key={t} value={t} className="bg-panel">
                    {t}
                  </option>
                ))}
              </Select>
            </Field>
          )}

          {isWordPress && (
            <>
              <div className="flex items-center gap-2 pt-1">
                <input
                  id="wp-install"
                  type="checkbox"
                  checked={wpInstall}
                  onChange={(e) => setWpInstall(e.target.checked)}
                  disabled={submitting}
                  className="accent-primary"
                />
                <label htmlFor="wp-install" className="text-sm text-fg">
                  Run <code className="font-mono">wp core install</code>
                </label>
              </div>

              {wpInstall && (
                <div className="grid grid-cols-2 gap-2 animate-fade-in">
                  <Field label="Site title">
                    <input
                      type="text"
                      value={wpTitle}
                      onChange={(e) => setWpTitle(e.target.value)}
                      placeholder={name || "My site"}
                      disabled={submitting}
                      className={inputCls}
                    />
                  </Field>
                  <Field label="Admin user">
                    <input
                      type="text"
                      value={wpAdmin}
                      onChange={(e) => setWpAdmin(e.target.value)}
                      disabled={submitting}
                      className={inputCls}
                    />
                  </Field>
                  <Field label="Admin password">
                    <input
                      type="text"
                      value={wpPass}
                      onChange={(e) => setWpPass(e.target.value)}
                      disabled={submitting}
                      className={inputCls}
                    />
                  </Field>
                  <Field label="Admin email">
                    <input
                      type="email"
                      value={wpEmail}
                      onChange={(e) => setWpEmail(e.target.value)}
                      placeholder="admin@localhost.test"
                      disabled={submitting}
                      className={inputCls}
                    />
                  </Field>
                </div>
              )}
            </>
          )}
        </div>

        {err && (
          <div className="mt-3 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-xs text-danger">
            {err}
          </div>
        )}

        <div className="mt-5 flex items-center justify-between">
          <span className="text-[11px] text-muted">
            {submitting
              ? "this can take ~30s for the first build…"
              : nameValid
              ? `will be reachable at https://${name}.localhost`
              : ""}
          </span>
          <div className="flex gap-2">
            <button
              type="button"
              onClick={onClose}
              disabled={submitting}
              className="panel-input rounded-md px-3 py-1.5 text-sm disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={submit}
              disabled={submitting || !nameValid}
              className="rounded-md bg-primary px-3 py-1.5 text-sm text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-50"
            >
              {submitting ? "Creating…" : "Create"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

// Shared base class so inputs and selects have identical box metrics.
// h-9 (36px) + px-2.5 py-1.5 lines them up exactly regardless of the
// native control's intrinsic height — macOS selects default to ~24px
// otherwise, which is ugly next to a text input.
const fieldBase =
  "panel-input w-full h-9 rounded-md px-2.5 text-sm text-fg outline-none focus:border-primary/60 disabled:opacity-50";

const inputCls = `${fieldBase} py-1.5`;

// Select wrapper: native <select> stripped of its appearance + custom
// caret on the right. Caret is pure decoration (pointer-events: none).
function Select({
  value,
  onChange,
  disabled,
  children,
}: {
  value: string;
  onChange: (v: string) => void;
  disabled?: boolean;
  children: React.ReactNode;
}) {
  return (
    <div className="relative">
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        className={cn(fieldBase, "appearance-none pr-8 cursor-pointer")}
      >
        {children}
      </select>
      <CaretDown className="pointer-events-none absolute z-10 right-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted" />
    </div>
  );
}

function Field({
  label,
  hint,
  children,
}: {
  label: string;
  hint?: string;
  children: React.ReactNode;
}) {
  return (
    <label className="block">
      <div className="mb-1 flex items-center justify-between text-[11px] uppercase tracking-wide text-muted">
        <span>{label}</span>
        {hint && <span className="normal-case tracking-normal">{hint}</span>}
      </div>
      {children}
    </label>
  );
}

function nameHint(name: string): string | undefined {
  if (!name) return undefined;
  if (NAME_RE.test(name)) return undefined;
  return "lowercase letters, digits, hyphens; must start with a letter";
}
