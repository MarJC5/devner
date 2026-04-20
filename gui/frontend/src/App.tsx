import { useEffect } from "react";
import { HashRouter, Routes, Route, Navigate, useNavigate } from "react-router-dom";
import "./App.css";
import { I18nProvider } from "./lib/i18n";
import { ThemeProvider } from "./lib/theme";
import { Layout } from "./Layout";
import { Projects } from "./screens/Projects";
import { Stack } from "./screens/Stack";
import { Chat } from "./screens/Chat";
import { Logs } from "./screens/Logs";
import { Databases } from "./screens/Databases";
import { Hosts } from "./screens/Hosts";
import { Settings } from "./screens/Settings";

function Shortcuts() {
  const nav = useNavigate();
  useEffect(() => {
    const ROUTES = ["projects", "stack", "chat", "logs", "databases", "hosts", "settings"];
    const onKey = (e: KeyboardEvent) => {
      if (!(e.metaKey || e.ctrlKey)) return;
      const n = parseInt(e.key, 10);
      if (!isNaN(n) && n >= 1 && n <= ROUTES.length) {
        e.preventDefault();
        nav("/" + ROUTES[n - 1]);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [nav]);

  // Native app-menu entries (View > Projects / Stack / …) send a
  // "nav:goto" event via Wails with the target route. We forward it to
  // react-router so menu clicks feel like sidebar clicks.
  useEffect(() => {
    import("../wailsjs/runtime/runtime").then(({ EventsOn, EventsOff }) => {
      EventsOn("nav:goto", (route: string) => nav(route));
      // EventsOff registered via cleanup — imported each render is OK,
      // it's a cached module.
    });
    return () => {
      import("../wailsjs/runtime/runtime").then(({ EventsOff }) => {
        EventsOff("nav:goto");
      });
    };
  }, [nav]);

  return null;
}

export default function App() {
  return (
    <ThemeProvider>
      <I18nProvider>
        <HashRouter>
        <Shortcuts />
        <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Navigate to="/projects" replace />} />
          <Route path="projects" element={<Projects />} />
          <Route path="stack" element={<Stack />} />
          <Route path="chat" element={<Chat />} />
          <Route path="logs" element={<Logs />} />
          <Route path="databases" element={<Databases />} />
          <Route path="hosts" element={<Hosts />} />
          <Route path="settings" element={<Settings />} />
          <Route path="*" element={<Navigate to="/projects" replace />} />
        </Route>
        </Routes>
      </HashRouter>
      </I18nProvider>
    </ThemeProvider>
  );
}
