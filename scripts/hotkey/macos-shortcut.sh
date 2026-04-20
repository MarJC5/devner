#!/usr/bin/env bash
# Creates a macOS "Shortcut" via shortcuts(1) that runs `devner prompt` in
# a new Terminal window. You then bind Cmd+D to it in System Settings >
# Keyboard > Keyboard Shortcuts > Services (or Shortcuts.app itself).
#
# Usage:
#   bash scripts/hotkey/macos-shortcut.sh
#
# After running, open Shortcuts.app, find "Devner Prompt", right-click
# → "Set Keyboard Shortcut" → Cmd+D.

set -euo pipefail

if ! command -v osascript >/dev/null; then
  echo "This helper only works on macOS." >&2
  exit 1
fi

DEVNER=$(command -v devner || echo "/usr/local/bin/devner")
if [ ! -x "$DEVNER" ]; then
  echo "devner not on PATH. Install it first (`make install` in the repo)." >&2
  exit 1
fi

cat <<APPLESCRIPT | osascript
tell application "Terminal"
    activate
    do script "$DEVNER prompt"
end tell
APPLESCRIPT

echo ""
echo "✓ A Terminal window just ran \`devner prompt\`."
echo ""
echo "To bind Cmd+D permanently, the cleanest way is:"
echo "  1. Open Raycast (or Alfred) — both have built-in global hotkey support."
echo "     Create a Script Command: $DEVNER prompt"
echo "     Assign Cmd+D."
echo ""
echo "  2. OR use Apple Shortcuts.app:"
echo "     File > New Shortcut"
echo "     Add action 'Run Shell Script' with content:"
echo "       $DEVNER prompt"
echo "     Name it 'Devner Prompt'."
echo "     In System Settings > Keyboard > Keyboard Shortcuts > Services,"
echo "     find your shortcut and bind Cmd+D."
echo ""
echo "  3. OR Hammerspoon — see scripts/hotkey/hammerspoon.lua"
