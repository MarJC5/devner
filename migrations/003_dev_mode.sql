-- Dev-time classification: 'server' (HMR dev server with reverse-proxy
-- routing), 'watch' (asset compile in watch mode alongside a PHP
-- backend), or '' (no dev-time process). dev_command stores the pnpm
-- script name ("dev" / "watch" / "start") that the project's
-- package.json exposes, so `devner dev start` knows which script to run
-- without re-sniffing the project every time.
ALTER TABLE projects ADD COLUMN dev_mode TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN dev_command TEXT NOT NULL DEFAULT '';
