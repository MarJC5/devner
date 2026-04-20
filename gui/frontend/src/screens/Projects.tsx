import { useCallback, useEffect, useState } from "react";
import {
  ListProjects,
  OpenInEditor,
  OpenProjectURL,
  DeleteProject,
  StartDevServer,
  StopDevServer,
  DevServerStatus,
  ListEditors,
} from "../../wailsjs/go/main/App";
import { store } from "../../wailsjs/go/models";
import {
  Folder,
  DotsThreeVertical,
  ArrowSquareOut,
  Code,
  Play,
  Stop,
  Trash,
  Plus,
  MagnifyingGlass,
} from "@phosphor-icons/react";
import { ActionMenu } from "../components/ActionMenu";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { CreateProjectModal } from "../components/CreateProjectModal";
import { cn } from "../lib/cn";
import { useT } from "../lib/i18n";

const NODE_LIKE = ["node", "nextjs", "nuxt", "astro", "sveltekit", "vite"];

export function Projects() {
  const t = useT();
  const [projects, setProjects] = useState<store.Project[]>([]);
  const [devStates, setDevStates] = useState<Record<string, boolean>>({});
  const [editors, setEditors] = useState<string[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [confirm, setConfirm] = useState<null | { name: string }>(null);
  const [busy, setBusy] = useState<Record<string, string>>({});
  const [createOpen, setCreateOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [dbFilter, setDbFilter] = useState<"all" | "mysql" | "postgres" | "none">("all");

  const refresh = useCallback(async () => {
    try {
      const [list, eds] = await Promise.all([ListProjects(), ListEditors()]);
      setProjects(list ?? []);
      setEditors(eds ?? []);
      // Probe dev server state for Node-like projects only.
      const states: Record<string, boolean> = {};
      await Promise.all(
        (list ?? [])
          .filter((p) => NODE_LIKE.includes(p.Type))
          .map(async (p) => {
            try {
              const st = await DevServerStatus(p.Name);
              states[p.Name] = st.running;
            } catch {
              // Container down or devserver doesn't exist yet — treat as stopped.
              states[p.Name] = false;
            }
          })
      );
      setDevStates(states);
    } catch (e) {
      setErr(String(e));
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  // Actions go through a tiny helper that tracks per-project "busy"
  // state so we can grey the card out while a long op (dev start)
  // runs, and surface errors without leaking them to the whole page.
  const run = async (name: string, label: string, fn: () => Promise<any>) => {
    setBusy((b) => ({ ...b, [name]: label }));
    setErr(null);
    try {
      await fn();
      await refresh();
    } catch (e) {
      setErr(`${name}: ${e}`);
    } finally {
      setBusy((b) => {
        const { [name]: _, ...rest } = b;
        return rest;
      });
    }
  };

  const preferredEditor = editors[0] ?? "code";

  // Apply search + DB filter locally. Both are cheap on the small sets
  // devner typically holds (10-50 projects), so no debouncing needed.
  const q = search.trim().toLowerCase();
  const filteredProjects = projects.filter((p) => {
    if (dbFilter === "none" && p.DBEngine !== "") return false;
    if (dbFilter === "mysql" && p.DBEngine !== "mysql") return false;
    if (dbFilter === "postgres" && p.DBEngine !== "postgres") return false;
    if (!q) return true;
    return (
      p.Name.toLowerCase().includes(q) ||
      p.Type.toLowerCase().includes(q) ||
      p.Domain.toLowerCase().includes(q) ||
      (p.DBName || "").toLowerCase().includes(q)
    );
  });

  return (
    <div className="p-6 animate-fade-in text-left">
      <div className="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-left">{t("projects.title")}</h1>
          <p className="mt-1 text-left text-sm text-muted">{t("projects.subtitle")}</p>
        </div>
        <button
          type="button"
          onClick={() => setCreateOpen(true)}
          className="no-drag flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-sm text-primary-foreground transition-opacity hover:opacity-90"
        >
          <Plus className="h-3.5 w-3.5" weight="bold" />
          {t("projects.new")}
        </button>
      </div>

      {/* Search + DB filter bar. Only visible once the store has at
          least one project — no point letting the user filter nothing. */}
      {projects.length > 0 && (
        <div className="mb-4 flex flex-wrap items-center gap-2">
          <div className="relative flex-1 min-w-[200px]">
            <MagnifyingGlass className="pointer-events-none absolute z-10 left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder={t("projects.search")}
              className="panel-input no-drag h-9 w-full rounded-md pl-8 pr-2.5 text-sm text-fg outline-none focus:border-primary/60"
            />
          </div>
          <FilterPill active={dbFilter === "all"} onClick={() => setDbFilter("all")}>
            {t("projects.filter.all")}
          </FilterPill>
          <FilterPill active={dbFilter === "mysql"} onClick={() => setDbFilter("mysql")}>
            {t("projects.filter.mysql")}
          </FilterPill>
          <FilterPill active={dbFilter === "postgres"} onClick={() => setDbFilter("postgres")}>
            {t("projects.filter.postgres")}
          </FilterPill>
          <FilterPill active={dbFilter === "none"} onClick={() => setDbFilter("none")}>
            {t("projects.filter.none")}
          </FilterPill>
        </div>
      )}

      {err && (
        <div className="mb-4 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
          {err}
        </div>
      )}

      {projects.length === 0 ? (
        <div className="panel-card rounded-lg px-6 py-10 text-center text-muted">
          {t("projects.empty")}
        </div>
      ) : filteredProjects.length === 0 ? (
        <div className="panel-card rounded-lg px-6 py-10 text-center text-muted">
          {t("projects.noMatch")} <code className="font-mono">{q || dbFilter}</code>.
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
          {filteredProjects.map((p) => {
            const isNode = NODE_LIKE.includes(p.Type);
            const devRunning = devStates[p.Name] ?? false;
            const busyLabel = busy[p.Name];

            return (
              <article
                key={p.Name}
                onClick={() =>
                  !busyLabel &&
                  run(p.Name, t("action.opening"), () =>
                    OpenInEditor(p.Name, preferredEditor)
                  )
                }
                className="panel-card group cursor-pointer rounded-lg p-4 transition-colors relative"
                title={t("projects.menu.openInEditor", { editor: preferredEditor })}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2 min-w-0 flex-1">
                    <Folder className="h-4 w-4 text-muted shrink-0" />
                    <h3 className="font-semibold truncate">{p.Name}</h3>
                    {isNode && (
                      <span
                        className={`h-1.5 w-1.5 rounded-full shrink-0 ${
                          devRunning ? "bg-success" : "bg-muted/40"
                        }`}
                        title={devRunning ? "dev running" : "dev stopped"}
                      />
                    )}
                  </div>
                  <ActionMenu
                    items={[
                      {
                        label: t("projects.menu.openInEditor", { editor: preferredEditor }),
                        icon: Code,
                        onSelect: () =>
                          run(p.Name, t("action.opening"), () =>
                            OpenInEditor(p.Name, preferredEditor)
                          ),
                      },
                      {
                        label: t("projects.menu.openUrl"),
                        icon: ArrowSquareOut,
                        onSelect: () =>
                          run(p.Name, t("action.opening"), () =>
                            OpenProjectURL(p.Name)
                          ),
                      },
                      ...(isNode
                        ? [
                            devRunning
                              ? {
                                  label: t("projects.menu.stopDev"),
                                  icon: Stop,
                                  onSelect: () =>
                                    run(p.Name, t("action.stopping"), () =>
                                      StopDevServer(p.Name)
                                    ),
                                }
                              : {
                                  label: t("projects.menu.startDev"),
                                  icon: Play,
                                  onSelect: () =>
                                    run(p.Name, t("action.starting"), () =>
                                      StartDevServer(p.Name, "")
                                    ),
                                },
                          ]
                        : []),
                      {
                        label: t("projects.menu.delete"),
                        icon: Trash,
                        danger: true,
                        onSelect: () => setConfirm({ name: p.Name }),
                      },
                    ]}
                  >
                    {() => <DotsThreeVertical className="h-5 w-5" weight="bold" />}
                  </ActionMenu>
                </div>

                <div className="mt-2 flex flex-wrap items-center gap-1.5 text-xs">
                  <span className="rounded bg-ink/5 px-1.5 py-0.5 text-muted">
                    {p.Type}
                  </span>
                  {p.DBEngine && (
                    <span className="rounded bg-ink/5 px-1.5 py-0.5 text-muted">
                      {p.DBEngine}/{p.DBName}
                    </span>
                  )}
                  {p.DevPort > 0 && (
                    <span className="rounded bg-ink/5 px-1.5 py-0.5 text-muted">
                      :{p.DevPort}
                    </span>
                  )}
                </div>

                <div className="mt-3 w-full truncate text-left text-sm text-primary">
                  https://{p.Domain}
                </div>

                {busyLabel && (
                  <div className="absolute inset-0 flex items-center justify-center rounded-lg bg-black/50 backdrop-blur-sm text-xs text-muted">
                    {busyLabel}
                  </div>
                )}
              </article>
            );
          })}
        </div>
      )}

      <ConfirmDialog
        open={confirm !== null}
        title={t("projects.confirm.deleteTitle", { name: confirm?.name ?? "" })}
        description={t("projects.confirm.deleteDescription")}
        confirmLabel={t("action.delete")}
        cancelLabel={t("action.cancel")}
        danger
        onCancel={() => setConfirm(null)}
        onConfirm={() => {
          if (!confirm) return;
          const name = confirm.name;
          setConfirm(null);
          run(name, t("action.deleting"), () => DeleteProject(name));
        }}
      />

      <CreateProjectModal
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={() => {
          refresh();
        }}
      />
    </div>
  );
}

function FilterPill({
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
