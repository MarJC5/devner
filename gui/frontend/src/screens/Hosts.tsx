import { useCallback, useEffect, useState } from "react";
import { ListHosts, AddHost, RemoveHost } from "../../wailsjs/go/main/App";
import { main } from "../../wailsjs/go/models";
import { Globe, Plus, Trash } from "@phosphor-icons/react";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { cn } from "../lib/cn";
import { useT } from "../lib/i18n";

export function Hosts() {
  const t = useT();
  const [rows, setRows] = useState<main.HostEntryDTO[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [confirm, setConfirm] = useState<null | main.HostEntryDTO>(null);
  const [domain, setDomain] = useState("");
  const [target, setTarget] = useState("127.0.0.1");
  const [adding, setAdding] = useState(false);

  const refresh = useCallback(async () => {
    try {
      const list = await ListHosts();
      setRows(list ?? []);
    } catch (e) {
      setErr(String(e));
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const add = async () => {
    if (!domain) return;
    setAdding(true);
    setErr(null);
    try {
      await AddHost(domain, target || "127.0.0.1");
      setDomain("");
      await refresh();
    } catch (e) {
      setErr(String(e));
    } finally {
      setAdding(false);
    }
  };

  const remove = async (h: main.HostEntryDTO) => {
    try {
      await RemoveHost(h.domain);
      await refresh();
    } catch (e) {
      setErr(String(e));
    }
  };

  return (
    <div className="p-6 animate-fade-in text-left">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold text-left">{t("hosts.title")}</h1>
        <p className="mt-1 text-left text-sm text-muted">{t("hosts.subtitle")}</p>
      </div>

      {err && (
        <div className="mb-4 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
          {err}
        </div>
      )}

      <div className="panel-card mb-4 rounded-lg p-4">
        <div className="mb-2 text-[11px] uppercase tracking-wide text-muted">
          {t("action.add")}
        </div>
        <div className="flex gap-2">
          <input
            type="text"
            value={domain}
            onChange={(e) => setDomain(e.target.value)}
            placeholder={t("hosts.domainPlaceholder")}
            disabled={adding}
            className="panel-input no-drag flex-1 h-9 rounded-md px-2.5 text-sm text-fg outline-none focus:border-primary/60"
          />
          <input
            type="text"
            value={target}
            onChange={(e) => setTarget(e.target.value)}
            placeholder={t("hosts.targetPlaceholder")}
            disabled={adding}
            className="panel-input no-drag w-32 h-9 rounded-md px-2.5 text-sm text-fg outline-none focus:border-primary/60"
          />
          <button
            type="button"
            onClick={add}
            disabled={adding || !domain}
            className="no-drag flex items-center gap-1.5 rounded-md bg-primary px-3 h-9 text-sm text-primary-foreground hover:opacity-90 disabled:opacity-50"
          >
            <Plus className="h-3.5 w-3.5" weight="bold" />
            {adding ? t("action.adding") : t("action.add")}
          </button>
        </div>
      </div>

      {rows.length === 0 ? (
        <div className="panel-card rounded-lg px-6 py-8 text-center text-muted">
          {t("hosts.empty")}
        </div>
      ) : (
        <div className="panel-card rounded-lg overflow-hidden">
          {rows.map((h, i) => (
            <div
              key={h.domain}
              className={cn(
                "flex items-center gap-3 px-4 py-2.5 text-sm",
                i > 0 && "border-t border-ink/5"
              )}
            >
              <Globe className="h-4 w-4 text-muted shrink-0" />
              <span className="flex-1 font-mono truncate">{h.domain}</span>
              <span className="font-mono text-muted">→ {h.target}</span>
              <button
                type="button"
                onClick={() => setConfirm(h)}
                className="text-muted hover:text-danger transition-colors"
                title="Remove"
              >
                <Trash className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>
      )}

      <ConfirmDialog
        open={confirm !== null}
        title={confirm ? `${t("action.delete")} ${confirm.domain} ?` : ""}
        description=""
        confirmLabel={t("action.delete")}
        cancelLabel={t("action.cancel")}
        danger
        onCancel={() => setConfirm(null)}
        onConfirm={() => {
          if (!confirm) return;
          const h = confirm;
          setConfirm(null);
          remove(h);
        }}
      />
    </div>
  );
}
