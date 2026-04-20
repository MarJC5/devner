import { useCallback, useEffect, useState } from "react";
import {
  GetConfigInfo,
  OpenConfigFile,
  CertsStatus,
  InstallCertsCA,
} from "../../wailsjs/go/main/App";
import { main, network } from "../../wailsjs/go/models";
import {
  Gear,
  FileText,
  Certificate,
  CheckCircle,
  XCircle,
  ArrowClockwise,
  CaretDown,
  Translate,
  Palette,
  Sun,
  Moon,
  Monitor,
} from "@phosphor-icons/react";
import { LANGS, useI18n, type Lang } from "../lib/i18n";
import { useTheme } from "../lib/theme";
import { cn } from "../lib/cn";

export function Settings() {
  const { t, lang, setLang } = useI18n();
  const { mode: themeMode, setMode: setThemeMode } = useTheme();
  const [cfg, setCfg] = useState<main.ConfigInfo | null>(null);
  const [certs, setCerts] = useState<network.CertsStatus | null>(null);
  const [installing, setInstalling] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  // Info banner shown when install was delegated to Terminal.app — not
  // an error per se, but we want visual distinction from success too.
  const [info, setInfo] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    try {
      const [c, s] = await Promise.all([GetConfigInfo(), CertsStatus()]);
      setCfg(c);
      setCerts(s);
    } catch (e) {
      setErr(String(e));
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const installCerts = async () => {
    setInstalling(true);
    setErr(null);
    setInfo(null);
    try {
      await InstallCertsCA();
      await refresh();
    } catch (e) {
      const msg = String(e);
      // On macOS we deliberately hand off to Terminal.app — the Go
      // side returns a "needs terminal" error which isn't really a
      // failure. Detect the sentinel string and surface as an info
      // banner instead of red.
      if (msg.includes("Terminal.app")) {
        setInfo(msg);
      } else {
        setErr(msg);
      }
    } finally {
      setInstalling(false);
    }
  };

  return (
    <div className="p-6 animate-fade-in text-left">
      <div className="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-left">{t("settings.title")}</h1>
          <p className="mt-1 text-left text-sm text-muted">
            {t("settings.subtitle")}
          </p>
        </div>
        <button
          type="button"
          onClick={refresh}
          className="panel-input no-drag flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm"
        >
          <ArrowClockwise className="h-3.5 w-3.5" />
          {t("action.refresh")}
        </button>
      </div>

      {err && (
        <div className="mb-4 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
          {err}
        </div>
      )}

      {info && (
        <div className="mb-4 rounded-md border border-primary/40 bg-primary/10 px-3 py-2 text-sm text-primary-muted">
          {info}
        </div>
      )}

      {/* Appearance — theme + language. Pills instead of a select so
          all three options stay visible, and Sun/Moon/Monitor make
          the mapping obvious without reading labels. */}
      <Section icon={Palette} title={t("settings.section.appearance")}>
        <div className="mb-1 text-[11px] uppercase tracking-wide text-muted">
          {t("settings.theme.label")}
        </div>
        <div className="flex gap-2">
          <ThemePill
            active={themeMode === "light"}
            onClick={() => setThemeMode("light")}
            icon={Sun}
            label={t("settings.theme.light")}
          />
          <ThemePill
            active={themeMode === "dark"}
            onClick={() => setThemeMode("dark")}
            icon={Moon}
            label={t("settings.theme.dark")}
          />
          <ThemePill
            active={themeMode === "system"}
            onClick={() => setThemeMode("system")}
            icon={Monitor}
            label={t("settings.theme.system")}
          />
        </div>
        <p className="mt-2 text-[11px] text-muted">{t("settings.theme.hint")}</p>
      </Section>

      {/* Language */}
      <Section icon={Translate} title={t("settings.section.language")}>
        <label className="block max-w-xs">
          <div className="mb-1 text-[11px] uppercase tracking-wide text-muted">
            {t("settings.language.label")}
          </div>
          <div className="relative">
            <select
              value={lang}
              onChange={(e) => setLang(e.target.value as Lang)}
              className="panel-input no-drag w-full h-9 appearance-none rounded-md px-2.5 pr-8 text-sm text-fg outline-none focus:border-primary/60 cursor-pointer"
            >
              {LANGS.map((l) => (
                <option key={l.code} value={l.code} className="bg-panel">
                  {l.label}
                </option>
              ))}
            </select>
            <CaretDown className="pointer-events-none absolute z-10 right-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted" />
          </div>
        </label>
        <p className="mt-2 text-[11px] text-muted">{t("settings.language.hint")}</p>
      </Section>

      {/* Config paths */}
      <Section icon={Gear} title={t("settings.section.config")}>
        {cfg ? (
          <>
            <Row label={t("settings.config.file")} value={cfg.config_path} />
            <Row label={t("settings.config.dataDir")} value={cfg.data_dir} />
            <Row label={t("settings.config.projectsDir")} value={cfg.projects_dir} />
            <div className="mt-3">
              <button
                type="button"
                onClick={() => OpenConfigFile()}
                className="panel-input no-drag flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs"
              >
                <FileText className="h-3.5 w-3.5" />
                {t("settings.config.openInEditor")}
              </button>
            </div>
          </>
        ) : (
          <div className="text-muted text-sm">{t("settings.loading")}</div>
        )}
      </Section>

      {/* LLM providers */}
      <Section icon={Gear} title={t("settings.section.providers")}>
        {cfg ? (
          <div className="space-y-2">
            {Object.entries(cfg.providers).map(([name, p]) => {
              const isActive = name === cfg.active_provider;
              return (
                <div
                  key={name}
                  className="panel-input flex items-center gap-3 rounded-md p-3"
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold">{name}</span>
                      {isActive && (
                        <span className="rounded bg-primary/30 px-1.5 py-0.5 text-[10px] uppercase tracking-wide text-primary">
                          {t("settings.providers.active")}
                        </span>
                      )}
                    </div>
                    <div className="mt-0.5 font-mono text-[11px] text-muted truncate">
                      {p.kind} · {p.model}
                    </div>
                  </div>
                  <div className="flex items-center gap-1 text-[11px]">
                    {p.has_key ? (
                      <>
                        <CheckCircle className="h-3.5 w-3.5 text-success" />
                        <span className="text-success">{t("settings.providers.apiKey")}</span>
                      </>
                    ) : (
                      <>
                        <XCircle className="h-3.5 w-3.5 text-muted" />
                        <span className="text-muted">{t("settings.providers.noKey")}</span>
                      </>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        ) : (
          <div className="text-muted text-sm">{t("settings.loading")}</div>
        )}
      </Section>

      {/* Certs */}
      <Section icon={Certificate} title={t("settings.section.certs")}>
        {certs ? (
          <>
            <div className="flex items-center gap-2 text-sm">
              {certs.MkcertInstalled ? (
                <>
                  <CheckCircle className="h-4 w-4 text-success" />
                  <span>{t("settings.certs.installed")}</span>
                  {certs.CAROOTPath && (
                    <span className="font-mono text-[11px] text-muted">
                      · {certs.CAROOTPath.trim()}
                    </span>
                  )}
                </>
              ) : (
                <>
                  <XCircle className="h-4 w-4 text-muted" />
                  <span>{t("settings.certs.notInstalled")}</span>
                </>
              )}
            </div>
            <div className="mt-3">
              <button
                type="button"
                onClick={installCerts}
                disabled={installing}
                className="panel-input no-drag flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs disabled:opacity-50"
              >
                {installing ? t("settings.certs.installing") : t("settings.certs.install")}
              </button>
            </div>
            <p className="mt-2 text-[11px] text-muted">{t("settings.certs.hint")}</p>
          </>
        ) : (
          <div className="text-muted text-sm">{t("settings.loading")}</div>
        )}
      </Section>
    </div>
  );
}

function ThemePill({
  active,
  onClick,
  icon: Icon,
  label,
}: {
  active: boolean;
  onClick: () => void;
  icon: React.ComponentType<any>;
  label: string;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "no-drag flex items-center gap-1.5 rounded-md px-3 h-9 text-xs transition-colors",
        active
          ? "bg-primary/90 text-primary-foreground"
          : "panel-input text-muted hover:text-fg"
      )}
    >
      <Icon className="h-3.5 w-3.5" />
      {label}
    </button>
  );
}

function Section({
  icon: Icon,
  title,
  children,
}: {
  icon: React.ComponentType<{ className?: string }>;
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section className="panel-card mb-4 rounded-lg p-4">
      <h2 className="mb-3 flex items-center gap-2 text-sm font-semibold text-fg">
        <Icon className="h-4 w-4 text-muted" />
        {title}
      </h2>
      {children}
    </section>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-baseline gap-3 py-0.5 text-sm">
      <span className="w-28 shrink-0 text-[11px] uppercase tracking-wide text-muted">
        {label}
      </span>
      <span className="font-mono text-[12px] text-fg/80 truncate">{value}</span>
    </div>
  );
}
