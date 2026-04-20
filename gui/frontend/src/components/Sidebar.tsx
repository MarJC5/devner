import { NavLink } from "react-router-dom";
import {
  FolderOpen,
  Stack as StackIcon,
  ChatCircle,
  FileText,
  Database,
  Globe,
  Gear,
} from "@phosphor-icons/react";
import { cn } from "../lib/cn";
import { useT } from "../lib/i18n";

type Entry = { to: string; tKey: string; icon: React.ComponentType<{ className?: string }>; hint: string };

const ENTRIES: Entry[] = [
  { to: "/projects", tKey: "nav.projects", icon: FolderOpen, hint: "⌘1" },
  { to: "/stack", tKey: "nav.stack", icon: StackIcon, hint: "⌘2" },
  { to: "/chat", tKey: "nav.chat", icon: ChatCircle, hint: "⌘3" },
  { to: "/logs", tKey: "nav.logs", icon: FileText, hint: "⌘4" },
  { to: "/databases", tKey: "nav.databases", icon: Database, hint: "⌘5" },
  { to: "/hosts", tKey: "nav.hosts", icon: Globe, hint: "⌘6" },
  { to: "/settings", tKey: "nav.settings", icon: Gear, hint: "⌘7" },
];

export function Sidebar() {
  const t = useT();
  return (
    <aside className="panel-blur w-52 shrink-0 flex flex-col py-3 border-r border-ink/5">
      <div className="px-4 pb-3 text-left text-sm font-semibold tracking-wide text-fg">
        Devner
      </div>

      <nav className="flex flex-col gap-0.5 px-2 no-drag">
        {ENTRIES.map(({ to, tKey, icon: Icon, hint }) => {
          const label = t(tKey);
          return (
          <NavLink
            key={to}
            to={to}
            title={`${label} (${hint})`}
            className={({ isActive }) =>
              cn(
                "group flex items-center gap-2.5 rounded-md px-3 py-1.5 text-sm transition-colors",
                isActive
                  ? "bg-primary/90 text-primary-foreground"
                  : "text-muted hover:bg-ink/5 hover:text-fg"
              )
            }
          >
            <Icon className="h-4 w-4 shrink-0" />
            <span className="flex-1 whitespace-nowrap truncate">{label}</span>
            <span className="text-[10px] opacity-60 group-hover:opacity-100">
              {hint}
            </span>
          </NavLink>
        );
        })}
      </nav>
    </aside>
  );
}
