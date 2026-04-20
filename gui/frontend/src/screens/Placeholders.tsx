// All screens have real implementations now. This file stays so any
// future "not built yet" screen can reuse the styling without
// re-inventing it.

export function NotBuilt({ title, hint }: { title: string; hint: string }) {
  return (
    <div className="p-6 animate-fade-in">
      <h1 className="text-2xl font-semibold">{title}</h1>
      <p className="mt-2 text-sm text-muted">{hint}</p>
      <div className="mt-6 rounded-lg border border-dashed border-ink/10 bg-ink/5 px-6 py-12 text-center text-muted">
        Coming soon.
      </div>
    </div>
  );
}
