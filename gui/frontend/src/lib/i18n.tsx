import { createContext, useContext, useEffect, useState, type ReactNode } from "react";

// Lightweight i18n. Two locales (en, fr) shipped as static dictionaries.
// State lives in React context, persisted via localStorage so the choice
// survives app relaunch (Wails WKWebView/WebView2 persists localStorage
// per-app). No runtime fallback: a missing key falls back to the key
// itself so it stays visible during development.

export type Lang = "en" | "fr";

export const LANGS: { code: Lang; label: string }[] = [
  { code: "en", label: "English" },
  { code: "fr", label: "Français" },
];

const STORAGE_KEY = "devner.lang";

type Dict = Record<string, string>;

const en: Dict = {
  // Sidebar / nav
  "nav.projects": "Projects",
  "nav.stack": "Stack",
  "nav.chat": "Chat",
  "nav.logs": "Logs",
  "nav.databases": "Databases",
  "nav.hosts": "Hosts",
  "nav.settings": "Settings",

  // Status bar
  "status.stack": "stack",

  // Common actions
  "action.cancel": "Cancel",
  "action.create": "Create",
  "action.creating": "Creating…",
  "action.delete": "Delete",
  "action.drop": "Drop",
  "action.refresh": "Refresh",
  "action.dismiss": "Dismiss",
  "action.add": "Add",
  "action.adding": "Adding…",
  "action.clear": "Clear",
  "action.approve": "Approve",
  "action.opening": "opening…",
  "action.stopping": "stopping…",
  "action.starting": "starting…",
  "action.deleting": "deleting…",

  // Projects screen
  "projects.title": "Projects",
  "projects.subtitle": "Click a card to open in code. Use the ⋮ menu for other actions.",
  "projects.new": "New project",
  "projects.search": "Search by name, type, URL…",
  "projects.filter.all": "All",
  "projects.filter.mysql": "MySQL",
  "projects.filter.postgres": "Postgres",
  "projects.filter.none": "No DB",
  "projects.empty": "No projects yet.",
  "projects.noMatch": "No projects match",
  "projects.menu.openInEditor": "Open in {editor}",
  "projects.menu.openUrl": "Open URL in browser",
  "projects.menu.startDev": "Start dev server",
  "projects.menu.stopDev": "Stop dev server",
  "projects.menu.delete": "Delete project",
  "projects.confirm.deleteTitle": "Delete {name}?",
  "projects.confirm.deleteDescription":
    "Removes project files, drops its database, and unregisters the Caddy entry. This cannot be undone.",

  // Stack screen
  "stack.title": "Stack",
  "stack.subtitle": "FrankenPHP + DBs + Redis + Mailpit + Adminer.",
  "stack.start": "Start",
  "stack.startDoing": "Starting…",
  "stack.stop": "Stop",
  "stack.stopDoing": "Stopping…",
  "stack.rebuild": "Rebuild",
  "stack.rebuildDoing": "Rebuilding…",

  // Chat screen
  "chat.title": "Chat",
  "chat.subtitle": "Ask the agent to create a project, start a dev server, run wp-cli, inspect logs, anything.",
  "chat.empty.hint1": "Ask the agent to create a project, start a dev server, run wp-cli, inspect logs, anything.",
  "chat.empty.hint2": "Shift+Enter for newline · Enter to send",
  "chat.placeholder": "Ask the agent…",
  "chat.confirm.title": "Confirm destructive action",
  "chat.confirm.body": "The agent wants to run",

  // Logs screen
  "logs.title": "Logs",
  "logs.subtitle": "Last 500 lines of each stack service. Click the tab to refresh.",

  // Databases screen
  "db.title": "Databases",
  "db.subtitle": "MySQL + Postgres user databases on the shared stack.",
  "db.new": "New database",
  "db.search": "Search databases by name…",
  "db.empty": "No databases yet.",
  "db.noMatch": "No databases match",
  "db.new.title": "New database",
  "db.new.engine": "Engine",
  "db.new.name": "Name",
  "db.new.nameHint": "lowercase, digits, underscores; must start with a letter",
  "db.confirm.drop": "Drop {engine}/{name}?",
  "db.confirm.dropDescription":
    "All data in this database will be lost. The matching user is also dropped.",

  // Hosts screen
  "hosts.title": "Hosts",
  "hosts.subtitle": "/etc/hosts entries managed by devner.",
  "hosts.domainPlaceholder": "myapp.test",
  "hosts.targetPlaceholder": "127.0.0.1",
  "hosts.empty": "No custom hosts entries. *.localhost resolves natively and does not need one.",

  // Settings screen
  "settings.title": "Settings",
  "settings.subtitle": "Config paths, LLM providers, certificates, language.",
  "settings.section.config": "Config",
  "settings.section.language": "Language",
  "settings.section.providers": "LLM providers",
  "settings.section.certs": "Certificates",
  "settings.config.file": "Config file",
  "settings.config.dataDir": "Data dir",
  "settings.config.projectsDir": "Projects dir",
  "settings.config.openInEditor": "Open config.toml in editor",
  "settings.language.label": "Interface language",
  "settings.language.hint": "Applied immediately, persisted across app restarts.",
  "settings.section.appearance": "Appearance",
  "settings.theme.label": "Theme",
  "settings.theme.hint": "Follow System picks light or dark based on your OS preference.",
  "settings.theme.light": "Light",
  "settings.theme.dark": "Dark",
  "settings.theme.system": "System",
  "settings.providers.active": "active",
  "settings.providers.apiKey": "api key",
  "settings.providers.noKey": "no key",
  "settings.certs.installed": "mkcert installed",
  "settings.certs.notInstalled": "mkcert not installed",
  "settings.certs.install": "Install CA (mkcert -install)",
  "settings.certs.installing": "Installing…",
  "settings.certs.hint":
    "Not needed for browsers (Caddy's internal CA is trusted at runtime). Install only if you need curl / Postman to trust local certs.",
  "settings.loading": "loading…",
};

const fr: Dict = {
  "nav.projects": "Projets",
  "nav.stack": "Stack",
  "nav.chat": "Chat",
  "nav.logs": "Logs",
  "nav.databases": "Databases",
  "nav.hosts": "Hôtes",
  "nav.settings": "Réglages",

  "status.stack": "stack",

  "action.cancel": "Annuler",
  "action.create": "Créer",
  "action.creating": "Création…",
  "action.delete": "Supprimer",
  "action.drop": "Supprimer",
  "action.refresh": "Rafraîchir",
  "action.dismiss": "Ignorer",
  "action.add": "Ajouter",
  "action.adding": "Ajout…",
  "action.clear": "Effacer",
  "action.approve": "Approuver",
  "action.opening": "ouverture…",
  "action.stopping": "arrêt…",
  "action.starting": "démarrage…",
  "action.deleting": "suppression…",

  "projects.title": "Projets",
  "projects.subtitle": "Cliquez sur une carte pour ouvrir dans votre éditeur. Le menu ⋮ liste les autres actions.",
  "projects.new": "Nouveau projet",
  "projects.search": "Rechercher par nom, type, URL…",
  "projects.filter.all": "Tous",
  "projects.filter.mysql": "MySQL",
  "projects.filter.postgres": "Postgres",
  "projects.filter.none": "Sans BDD",
  "projects.empty": "Aucun projet pour le moment.",
  "projects.noMatch": "Aucun projet ne correspond à",
  "projects.menu.openInEditor": "Ouvrir dans {editor}",
  "projects.menu.openUrl": "Ouvrir l'URL dans le navigateur",
  "projects.menu.startDev": "Démarrer le dev server",
  "projects.menu.stopDev": "Arrêter le dev server",
  "projects.menu.delete": "Supprimer le projet",
  "projects.confirm.deleteTitle": "Supprimer {name} ?",
  "projects.confirm.deleteDescription":
    "Supprime les fichiers du projet, sa base de données, et retire l'entrée Caddy. Action irréversible.",

  "stack.title": "Stack",
  "stack.subtitle": "FrankenPHP + BDDs + Redis + Mailpit + Adminer.",
  "stack.start": "Démarrer",
  "stack.startDoing": "Démarrage…",
  "stack.stop": "Arrêter",
  "stack.stopDoing": "Arrêt…",
  "stack.rebuild": "Reconstruire",
  "stack.rebuildDoing": "Reconstruction…",

  "chat.title": "Chat",
  "chat.subtitle": "Demandez à l'agent de créer un projet, lancer un dev server, exécuter wp-cli, consulter les logs, etc.",
  "chat.empty.hint1": "Demandez à l'agent de créer un projet, lancer un dev server, exécuter wp-cli, consulter les logs, etc.",
  "chat.empty.hint2": "Maj+Entrée pour un saut de ligne · Entrée pour envoyer",
  "chat.placeholder": "Demandez à l'agent…",
  "chat.confirm.title": "Confirmer l'action destructive",
  "chat.confirm.body": "L'agent souhaite exécuter",

  "logs.title": "Logs",
  "logs.subtitle": "500 dernières lignes de chaque service. Cliquez l'onglet pour rafraîchir.",

  "db.title": "Bases de données",
  "db.subtitle": "Bases MySQL + Postgres utilisateur sur la stack partagée.",
  "db.new": "Nouvelle base",
  "db.search": "Rechercher une base par nom…",
  "db.empty": "Aucune base pour le moment.",
  "db.noMatch": "Aucune base ne correspond à",
  "db.new.title": "Nouvelle base",
  "db.new.engine": "Moteur",
  "db.new.name": "Nom",
  "db.new.nameHint": "minuscules, chiffres, underscores ; doit débuter par une lettre",
  "db.confirm.drop": "Supprimer {engine}/{name} ?",
  "db.confirm.dropDescription":
    "Toutes les données de cette base seront perdues. L'utilisateur associé est également supprimé.",

  "hosts.title": "Hôtes",
  "hosts.subtitle": "Entrées /etc/hosts gérées par devner.",
  "hosts.domainPlaceholder": "monapp.test",
  "hosts.targetPlaceholder": "127.0.0.1",
  "hosts.empty": "Aucune entrée personnalisée. *.localhost est résolu nativement sans avoir besoin d'entrée.",

  "settings.title": "Réglages",
  "settings.subtitle": "Chemins de config, providers LLM, certificats, langue.",
  "settings.section.config": "Configuration",
  "settings.section.language": "Langue",
  "settings.section.providers": "Providers LLM",
  "settings.section.certs": "Certificats",
  "settings.config.file": "Fichier de config",
  "settings.config.dataDir": "Dossier de données",
  "settings.config.projectsDir": "Dossier des projets",
  "settings.config.openInEditor": "Ouvrir config.toml dans l'éditeur",
  "settings.language.label": "Langue de l'interface",
  "settings.language.hint": "Appliquée immédiatement, conservée entre les sessions.",
  "settings.section.appearance": "Apparence",
  "settings.theme.label": "Thème",
  "settings.theme.hint": "« Système » suit la préférence de votre OS (clair ou sombre).",
  "settings.theme.light": "Clair",
  "settings.theme.dark": "Sombre",
  "settings.theme.system": "Système",
  "settings.providers.active": "actif",
  "settings.providers.apiKey": "clé api",
  "settings.providers.noKey": "pas de clé",
  "settings.certs.installed": "mkcert installé",
  "settings.certs.notInstalled": "mkcert non installé",
  "settings.certs.install": "Installer la CA (mkcert -install)",
  "settings.certs.installing": "Installation…",
  "settings.certs.hint":
    "Inutile pour les navigateurs (la CA interne de Caddy est reconnue au runtime). À installer uniquement si curl / Postman doit faire confiance aux certs locaux.",
  "settings.loading": "chargement…",
};

const DICTS: Record<Lang, Dict> = { en, fr };

type I18nCtx = {
  lang: Lang;
  setLang: (l: Lang) => void;
  t: (key: string, vars?: Record<string, string>) => string;
};

const Ctx = createContext<I18nCtx | null>(null);

function detectInitial(): Lang {
  try {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved === "en" || saved === "fr") return saved;
  } catch {
    // localStorage can throw in private-browsing webviews — ignore.
  }
  const nav = navigator.language?.toLowerCase() ?? "";
  return nav.startsWith("fr") ? "fr" : "en";
}

export function I18nProvider({ children }: { children: ReactNode }) {
  const [lang, setLangState] = useState<Lang>(detectInitial);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, lang);
    } catch {
      // ignore — see detectInitial
    }
    document.documentElement.lang = lang;
  }, [lang]);

  const t = (key: string, vars?: Record<string, string>): string => {
    let s = DICTS[lang][key] ?? DICTS.en[key] ?? key;
    if (vars) {
      for (const [k, v] of Object.entries(vars)) {
        s = s.replace(new RegExp(`\\{${k}\\}`, "g"), v);
      }
    }
    return s;
  };

  return (
    <Ctx.Provider value={{ lang, setLang: setLangState, t }}>
      {children}
    </Ctx.Provider>
  );
}

export function useI18n(): I18nCtx {
  const v = useContext(Ctx);
  if (!v) throw new Error("useI18n must be used inside <I18nProvider>");
  return v;
}

// Convenience: t only. Most call sites only need translation, not lang.
export function useT() {
  return useI18n().t;
}
