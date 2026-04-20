import { useEffect, useState } from "react";
import { Outlet } from "react-router-dom";
import { EventsOn, EventsOff } from "../wailsjs/runtime/runtime";
import { main } from "../wailsjs/go/models";
import { Sidebar } from "./components/Sidebar";
import { StatusBar } from "./components/StatusBar";

// A dedicated draggable strip at the very top. This is where macOS's
// hidden-inset traffic lights live (red/yellow/green). Needs to be
// tall enough that the buttons clear the content below — 44px is the
// standard macOS titlebar height. Everything else defaults to no-drag
// so clicks in the sidebar / content reach the correct target.
function DragBar() {
  return (
    <div
      style={{ ["--wails-draggable" as any]: "drag" }}
      className="h-11 shrink-0 border-b border-ink/5"
    />
  );
}

export function Layout() {
  const [rows, setRows] = useState<main.ContainerStatus[]>([]);

  useEffect(() => {
    EventsOn("stack:tick", (data: main.ContainerStatus[]) =>
      setRows(data ?? [])
    );
    return () => EventsOff("stack:tick");
  }, []);

  const running = rows.filter((r) => r.state === "running").length;
  const total = 6;

  return (
    // The outer container is no-drag so nested buttons / inputs don't
    // inherit drag from the body. The DragBar re-enables drag for a
    // thin strip across the top.
    <div
      style={{ ["--wails-draggable" as any]: "no-drag" }}
      className="flex h-screen flex-col text-fg"
    >
      <DragBar />
      <div className="flex flex-1 min-h-0">
        <Sidebar />
        <main className="flex-1 overflow-y-auto">
          <Outlet />
        </main>
      </div>
      <StatusBar
        stackHealthy={running === total}
        runningCount={running}
        totalCount={total}
        version="v0.1.0-dev"
      />
    </div>
  );
}
