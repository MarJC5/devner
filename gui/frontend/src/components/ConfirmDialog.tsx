import { useEffect } from "react";
import { cn } from "../lib/cn";

// Modal confirmation. Non-bypassable: user must click Confirm or
// Cancel. Esc = cancel. Backdrop click = cancel. Not closable by
// clicking outside the dialog content.

type Props = {
  open: boolean;
  title: string;
  description?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  danger?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
};

export function ConfirmDialog({
  open,
  title,
  description,
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  danger,
  onConfirm,
  onCancel,
}: Props) {
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onCancel();
      if (e.key === "Enter") onConfirm();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, onCancel, onConfirm]);

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in no-drag"
      onMouseDown={onCancel}
    >
      <div
        className="panel-card w-full max-w-sm animate-slide-up rounded-lg p-5 shadow-2xl"
        onMouseDown={(e) => e.stopPropagation()}
      >
        <h2 className="text-base font-semibold">{title}</h2>
        {description && (
          <p className="mt-2 text-sm text-muted">{description}</p>
        )}
        <div className="mt-5 flex justify-end gap-2">
          <button
            type="button"
            onClick={onCancel}
            className="panel-input rounded-md px-3 py-1.5 text-sm"
          >
            {cancelLabel}
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className={cn(
              "rounded-md px-3 py-1.5 text-sm transition-opacity hover:opacity-90",
              danger
                ? "bg-danger text-white"
                : "bg-primary text-primary-foreground"
            )}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
