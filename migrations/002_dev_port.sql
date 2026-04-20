-- Node-like projects (Next.js, Astro, Vite, plain Node) need a stable
-- internal port so Caddy can reverse_proxy *.localhost to the running
-- dev server. dev_port = 0 means "no dev server" (PHP projects,
-- static-only Node). Allocated from the 3100-3999 range by
-- store.AllocateDevPort.
ALTER TABLE projects ADD COLUMN dev_port INTEGER NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_projects_dev_port ON projects(dev_port);
