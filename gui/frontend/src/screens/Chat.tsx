import { useCallback, useEffect, useRef, useState } from "react";
import { marked } from "marked";
import {
  AgentChat,
  AgentConfirm,
  AgentReset,
  ActiveProvider,
  ListProviders,
} from "../../wailsjs/go/main/App";
import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime";
import {
  PaperPlaneRight,
  Trash,
  Warning,
  Lightning,
  CaretDown,
  CaretRight,
  Check,
  X as XIcon,
} from "@phosphor-icons/react";
import { cn } from "../lib/cn";
import { useT } from "../lib/i18n";

// Types that match what the Go side emits. Keeping them local avoids a
// round-trip through wailsjs/go/models for simple shapes.

type Msg =
  | { kind: "user"; text: string }
  | { kind: "assistant"; text: string; streaming?: boolean }
  | { kind: "tool"; name: string; args: string; result?: string; ok?: boolean };

type ConfirmReq = { tool: string; args: string };

export function Chat() {
  const t = useT();
  const [messages, setMessages] = useState<Msg[]>([]);
  const [input, setInput] = useState("");
  const [running, setRunning] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [confirm, setConfirm] = useState<ConfirmReq | null>(null);

  const [provider, setProvider] = useState<string>("");
  const [providers, setProviders] = useState<string[]>([]);

  const sidRef = useRef<string>("");
  const viewportRef = useRef<HTMLDivElement>(null);

  // Track providers for the dropdown + active default on mount.
  useEffect(() => {
    ListProviders()
      .then((list) => setProviders(list ?? []))
      .catch(() => {});
    ActiveProvider()
      .then((p) => setProvider(p ?? ""))
      .catch(() => {});
  }, []);

  // Auto-scroll to the bottom whenever messages change (user + agent).
  useEffect(() => {
    const el = viewportRef.current;
    if (!el) return;
    el.scrollTop = el.scrollHeight;
  }, [messages, confirm]);

  // Subscribe to all events tied to the current session. Re-runs when
  // sidRef changes — but since we set sidRef once per submit, we
  // attach/detach once per chat turn.
  const subscribe = useCallback((sid: string) => {
    EventsOn(`agent:${sid}:delta`, ({ text }: { text: string }) => {
      setMessages((m) => {
        // Append to last assistant msg if it's streaming; otherwise push new.
        const last = m[m.length - 1];
        if (last?.kind === "assistant" && last.streaming) {
          return [...m.slice(0, -1), { ...last, text: last.text + text }];
        }
        return [...m, { kind: "assistant", text, streaming: true }];
      });
    });

    EventsOn(`agent:${sid}:tool_call`, (data: any) => {
      setMessages((m) => {
        // Seal any streaming assistant msg (stream paused by a tool call).
        const sealed = m.map((x, i) =>
          i === m.length - 1 && x.kind === "assistant" ? { ...x, streaming: false } : x
        );
        return [...sealed, { kind: "tool", name: data.name, args: data.args }];
      });
    });

    EventsOn(`agent:${sid}:tool_result`, (data: any) => {
      setMessages((m) => {
        // Attach the result to the most recent matching tool message
        // that doesn't have a result yet. Works for serial tool calls;
        // parallel tools would need matching by id (not implemented).
        const idx = [...m].reverse().findIndex(
          (x) => x.kind === "tool" && x.name === data.name && x.result === undefined
        );
        if (idx < 0) return m;
        const realIdx = m.length - 1 - idx;
        return m.map((x, i) =>
          i === realIdx && x.kind === "tool"
            ? { ...x, result: data.content, ok: data.ok }
            : x
        );
      });
    });

    EventsOn(`agent:${sid}:need_confirm`, (data: any) => {
      setConfirm({ tool: data.tool, args: data.args });
    });

    EventsOn(`agent:${sid}:err`, ({ msg }: { msg: string }) => {
      setErr(msg);
      setRunning(false);
    });

    EventsOn(`agent:${sid}:done`, () => {
      // Seal the final streaming message.
      setMessages((m) =>
        m.map((x, i) =>
          i === m.length - 1 && x.kind === "assistant" ? { ...x, streaming: false } : x
        )
      );
      setRunning(false);
    });
  }, []);

  const unsubscribe = (sid: string) => {
    for (const ev of ["delta", "tool_call", "tool_result", "need_confirm", "err", "done"]) {
      EventsOff(`agent:${sid}:${ev}`);
    }
  };

  // Clean up all subscriptions on unmount so events don't leak between
  // navigations.
  useEffect(() => {
    return () => {
      if (sidRef.current) unsubscribe(sidRef.current);
    };
  }, []);

  const submit = async () => {
    const text = input.trim();
    if (!text || running) return;
    setInput("");
    setErr(null);
    setMessages((m) => [...m, { kind: "user", text }]);
    setRunning(true);

    try {
      const sid = await AgentChat(text, provider, sidRef.current);
      if (sidRef.current !== sid) {
        // First turn or provider changed — wire subscriptions.
        if (sidRef.current) unsubscribe(sidRef.current);
        sidRef.current = sid;
        subscribe(sid);
      }
    } catch (e) {
      setErr(String(e));
      setRunning(false);
    }
  };

  const clearChat = async () => {
    if (sidRef.current) {
      await AgentReset(sidRef.current);
      unsubscribe(sidRef.current);
    }
    sidRef.current = "";
    setMessages([]);
    setErr(null);
    setConfirm(null);
  };

  const replyConfirm = async (approve: boolean) => {
    if (!sidRef.current || !confirm) return;
    setConfirm(null);
    try {
      await AgentConfirm(sidRef.current, approve);
    } catch (e) {
      setErr(String(e));
    }
  };

  return (
    <div className="flex h-full flex-col p-6 animate-fade-in text-left">
      <div className="mb-4 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-left">{t("chat.title")}</h1>
          <p className="mt-1 text-left text-sm text-muted">{t("chat.subtitle")}</p>
        </div>
        <div className="flex items-center gap-2">
          <ProviderSelect
            value={provider}
            onChange={setProvider}
            options={providers}
            disabled={running}
          />
          <button
            type="button"
            onClick={clearChat}
            disabled={running || messages.length === 0}
            className="panel-input no-drag flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm text-fg disabled:opacity-40"
          >
            <Trash className="h-3.5 w-3.5" />
            {t("action.clear")}
          </button>
        </div>
      </div>

      {err && (
        <div className="mb-3 flex items-center gap-2 rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
          <Warning className="h-4 w-4 shrink-0" />
          <span className="truncate">{err}</span>
        </div>
      )}

      <div
        ref={viewportRef}
        className="panel-card mb-3 flex-1 min-h-0 overflow-y-auto rounded-lg p-4 no-drag"
      >
        {messages.length === 0 && (
          <div className="flex h-full items-center justify-center text-center text-muted">
            <div className="max-w-xs space-y-2">
              <Lightning className="mx-auto h-8 w-8 opacity-50" />
              <p className="text-sm">{t("chat.empty.hint1")}</p>
              <p className="text-xs opacity-60">{t("chat.empty.hint2")}</p>
            </div>
          </div>
        )}

        {messages.map((msg, i) => (
          <Message key={i} msg={msg} />
        ))}

        {confirm && (
          <div className="mt-4 rounded-md border border-danger/40 bg-danger/10 p-3">
            <div className="flex items-center gap-2 text-sm font-semibold text-danger">
              <Warning className="h-4 w-4" />
              {t("chat.confirm.title")}
            </div>
            <div className="mt-1.5 text-xs text-danger/90">
              {t("chat.confirm.body")}{" "}
              <code className="font-mono">
                {confirm.tool}({truncate(confirm.args, 120)})
              </code>
            </div>
            <div className="mt-3 flex justify-end gap-2">
              <button
                type="button"
                onClick={() => replyConfirm(false)}
                className="panel-input rounded-md px-3 py-1.5 text-xs"
              >
                {t("action.cancel")}
              </button>
              <button
                type="button"
                onClick={() => replyConfirm(true)}
                className="rounded-md bg-danger px-3 py-1.5 text-xs text-white transition-opacity hover:opacity-90"
              >
                {t("action.approve")}
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Input composer: textarea + Send button share the same height via
          items-stretch + fixed min-h-[44px]. Textarea is single-line by
          default; Shift+Enter triggers a visual newline but the box
          won't grow (keeps the composer compact). */}
      <div className="flex items-stretch gap-2">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              submit();
            }
          }}
          placeholder={t("chat.placeholder")}
          disabled={running}
          rows={1}
          className="panel-input no-drag flex-1 resize-none rounded-md px-3 py-2.5 text-sm leading-5 text-fg outline-none focus:border-primary/60 disabled:opacity-50"
        />
        <button
          type="button"
          onClick={submit}
          disabled={running || !input.trim()}
          className="no-drag flex items-center gap-1.5 rounded-md bg-primary px-4 text-sm text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          <PaperPlaneRight className="h-3.5 w-3.5" weight="fill" />
          {running ? "Thinking…" : "Send"}
        </button>
      </div>
    </div>
  );
}

function Message({ msg }: { msg: Msg }) {
  if (msg.kind === "user") {
    return (
      <div className="mb-4 flex justify-end">
        <div className="max-w-[80%] rounded-lg bg-primary/20 border border-primary/30 px-3 py-2 text-sm text-fg">
          {msg.text}
        </div>
      </div>
    );
  }
  if (msg.kind === "assistant") {
    return (
      <div className="mb-4">
        <div
          className="prose prose-sm prose-invert max-w-none prose-pre:bg-black/40 prose-pre:text-fg prose-code:text-accent prose-headings:text-fg prose-a:text-primary"
          // Users cannot post to the LLM, responses are trusted enough
          // for markdown — but still, marked has XSS protection.
          dangerouslySetInnerHTML={{ __html: renderMarkdown(msg.text) }}
        />
        {msg.streaming && (
          <span className="inline-block h-3 w-[2px] bg-primary animate-pulse" />
        )}
      </div>
    );
  }
  // tool call — collapsible. Collapsed by default; auto-expanded on
  // error so the user sees the failure without an extra click. While
  // still pending (no result yet), show a pulsing dot and no caret.
  return <ToolMessage msg={msg} />;
}

function ToolMessage({ msg }: { msg: Extract<Msg, { kind: "tool" }> }) {
  const pending = msg.result === undefined;
  const errored = msg.ok === false;
  // Default expanded only on error; user can collapse/expand manually.
  const [open, setOpen] = useState(errored);
  // If the message arrives as an error later (streamed result), open it.
  useEffect(() => {
    if (errored) setOpen(true);
  }, [errored]);

  const toggle = () => {
    if (pending) return;
    setOpen((v) => !v);
  };

  const statusDot = pending ? (
    <span className="h-1.5 w-1.5 rounded-full bg-accent animate-pulse" />
  ) : errored ? (
    <XIcon className="h-3 w-3 text-danger" weight="bold" />
  ) : (
    <Check className="h-3 w-3 text-success" weight="bold" />
  );

  return (
    <div className="panel-input mb-2 rounded-md text-xs">
      <button
        type="button"
        onClick={toggle}
        disabled={pending}
        className={cn(
          "flex w-full items-center gap-2 px-2.5 py-1.5 font-mono text-left",
          !pending && "cursor-pointer hover:bg-ink/5"
        )}
      >
        {!pending ? (
          open ? (
            <CaretDown className="h-3 w-3 text-muted shrink-0" />
          ) : (
            <CaretRight className="h-3 w-3 text-muted shrink-0" />
          )
        ) : (
          <Lightning className="h-3 w-3 text-accent shrink-0" weight="fill" />
        )}
        <span className="font-semibold text-accent">{msg.name}</span>
        <span className="truncate text-muted font-normal flex-1">
          {truncate(stripBraces(msg.args), 120)}
        </span>
        {statusDot}
      </button>

      {open && msg.result !== undefined && (
        <div
          className={cn(
            "border-t border-ink/10 bg-black/30 p-2 font-mono text-[11px] whitespace-pre-wrap break-words",
            errored ? "text-danger" : "text-fg/70"
          )}
        >
          {truncate(msg.result, 2000)}
        </div>
      )}
    </div>
  );
}

// Strip outer `{}` from JSON args when they contain nothing — e.g.
// `{}` renders as empty, `{"name":"demo"}` shows `name: "demo"`. Pure
// cosmetic — the full args are re-shown in the expanded panel if
// someone needs them.
function stripBraces(s: string): string {
  const trimmed = s.trim();
  if (trimmed === "{}" || trimmed === "") return "";
  if (trimmed.startsWith("{") && trimmed.endsWith("}")) {
    return trimmed.slice(1, -1).replace(/^"|"$/g, "").replace(/","/g, ", ");
  }
  return trimmed;
}

function ProviderSelect({
  value,
  onChange,
  options,
  disabled,
}: {
  value: string;
  onChange: (v: string) => void;
  options: string[];
  disabled?: boolean;
}) {
  return (
    <div className="relative">
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        className="panel-input no-drag appearance-none rounded-md px-3 pr-7 py-1.5 text-sm text-fg outline-none focus:border-primary/60 disabled:opacity-50 cursor-pointer"
      >
        {options.map((o) => (
          <option key={o} value={o} className="bg-panel">
            {o}
          </option>
        ))}
      </select>
      <CaretDown className="pointer-events-none absolute z-10 right-2 top-1/2 h-3 w-3 -translate-y-1/2 text-muted" />
    </div>
  );
}

function truncate(s: string, n: number): string {
  if (s.length <= n) return s;
  return s.slice(0, n) + "…";
}

function renderMarkdown(md: string): string {
  // `marked` v18 is sync by default and returns a string.
  try {
    return marked.parse(md, { gfm: true, breaks: true }) as string;
  } catch {
    return escapeHtml(md);
  }
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}
