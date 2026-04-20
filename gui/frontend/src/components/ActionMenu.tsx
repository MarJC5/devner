import { useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { cn } from "../lib/cn";

// Position the popover with CSS `right: window.width - triggerRight`
// rather than with `transform: translateX(-100%)`. The latter conflicts
// with the slide-up animation's `transform: translateY(…)` —
// transforms can't be composed across the animation and the inline
// style, so the popover would briefly render at the wrong x-position
// before snapping. Using `right` means the x is fixed at mount; only y
// animates.

export type MenuItem = {
  label: string;
  icon?: React.ComponentType<{ className?: string }>;
  onSelect: () => void;
  danger?: boolean;
  disabled?: boolean;
  shortcut?: string;
};

type Props = {
  items: MenuItem[];
  children: (open: boolean) => React.ReactNode;
};

type Position = { top: number; right: number };

export function ActionMenu({ items, children }: Props) {
  const [open, setOpen] = useState(false);
  const [pos, setPos] = useState<Position | null>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const popoverRef = useRef<HTMLDivElement>(null);

  // Compute position synchronously on click so the popover never
  // renders with a stale/missing position. useEffect-based measurement
  // would cause a one-frame jump.
  const openMenu = () => {
    if (!triggerRef.current) return;
    const b = triggerRef.current.getBoundingClientRect();
    setPos({
      top: b.bottom + 4,
      right: window.innerWidth - b.right,
    });
    setOpen(true);
  };

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      const t = e.target as Node;
      if (
        triggerRef.current?.contains(t) ||
        popoverRef.current?.contains(t)
      ) {
        return;
      }
      setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    window.addEventListener("mousedown", onClick);
    window.addEventListener("keydown", onKey);
    return () => {
      window.removeEventListener("mousedown", onClick);
      window.removeEventListener("keydown", onKey);
    };
  }, [open]);

  return (
    <>
      <button
        ref={triggerRef}
        type="button"
        onClick={(e) => {
          e.stopPropagation();
          if (open) setOpen(false);
          else openMenu();
        }}
        className="no-drag text-muted hover:text-fg transition-colors"
      >
        {children(open)}
      </button>

      {open && pos &&
        createPortal(
          <div
            ref={popoverRef}
            style={{
              position: "fixed",
              top: pos.top,
              right: pos.right,
              zIndex: 9999,
            }}
            className="panel-card no-drag min-w-[200px] animate-slide-up rounded-md py-1 shadow-2xl"
            onMouseDown={(e) => e.stopPropagation()}
          >
            {items.map((item, i) => (
              <button
                key={i}
                type="button"
                disabled={item.disabled}
                onClick={(e) => {
                  e.stopPropagation();
                  setOpen(false);
                  item.onSelect();
                }}
                className={cn(
                  "flex w-full items-center gap-2 px-3 py-1.5 text-left text-sm transition-colors",
                  item.disabled && "opacity-40 cursor-not-allowed",
                  !item.disabled && item.danger && "text-danger hover:bg-danger/10",
                  !item.disabled && !item.danger && "text-fg hover:bg-ink/5"
                )}
              >
                {item.icon && <item.icon className="h-4 w-4" />}
                <span className="flex-1">{item.label}</span>
                {item.shortcut && (
                  <span className="text-[10px] text-muted">{item.shortcut}</span>
                )}
              </button>
            ))}
          </div>,
          document.body
        )}
    </>
  );
}
