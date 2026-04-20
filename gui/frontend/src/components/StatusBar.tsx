import { cn } from "../lib/cn";
import { useT } from "../lib/i18n";

type Props = {
  stackHealthy?: boolean;
  runningCount?: number;
  totalCount?: number;
  provider?: string;
  toolsCount?: number;
  version?: string;
};

export function StatusBar({
  stackHealthy,
  runningCount,
  totalCount,
  provider,
  toolsCount,
  version,
}: Props) {
  const t = useT();
  return (
    <div className="panel-blur no-drag flex h-8 shrink-0 items-center gap-4 border-t border-ink/5 px-3 text-[11px] text-muted">
      {typeof stackHealthy === "boolean" && (
        <span className="flex items-center gap-1.5">
          <span
            className={cn(
              "h-1.5 w-1.5 rounded-full",
              stackHealthy ? "bg-success" : "bg-danger"
            )}
          />
          {t("status.stack")}{" "}
          {typeof runningCount === "number" && typeof totalCount === "number"
            ? `${runningCount}/${totalCount}`
            : stackHealthy ? "up" : "down"}
        </span>
      )}
      {provider && <span>provider: {provider}</span>}
      {typeof toolsCount === "number" && <span>{toolsCount} tools</span>}
      <span className="ml-auto opacity-60">{version ?? "dev"}</span>
    </div>
  );
}
