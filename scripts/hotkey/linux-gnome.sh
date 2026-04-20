#!/usr/bin/env bash
# Bind Super+D on GNOME to launch `devner prompt` in a new gnome-terminal.
# Uses dconf, so no plugin / daemon needed beyond GNOME itself.
#
# Usage:
#   bash scripts/hotkey/linux-gnome.sh
#
# KDE users: System Settings > Shortcuts > Custom Shortcuts, action
# "Command/URL", command: "konsole -e devner prompt", trigger Super+D.

set -euo pipefail

if ! command -v gsettings >/dev/null; then
  echo "This helper expects GNOME (gsettings not found)." >&2
  exit 1
fi

DEVNER=$(command -v devner || echo "/usr/local/bin/devner")
CMD="gnome-terminal -- $DEVNER prompt"
NAME="Devner Prompt"
BIND="<Super>d"
SLOT="custom-devner"

SCHEMA="org.gnome.settings-daemon.plugins.media-keys"
LIST_KEY="custom-keybindings"
CHILD_SCHEMA="org.gnome.settings-daemon.plugins.media-keys.custom-keybinding"
CHILD_PATH="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/$SLOT/"

# Add the slot to the list if absent.
current=$(gsettings get $SCHEMA $LIST_KEY)
if [[ "$current" != *"$CHILD_PATH"* ]]; then
  new=$(printf "%s" "$current" | sed -e "s|\\]$|, '$CHILD_PATH']|" -e "s|^@as \\[\\]$|['$CHILD_PATH']|")
  gsettings set $SCHEMA $LIST_KEY "$new"
fi

gsettings set $CHILD_SCHEMA:$CHILD_PATH name "$NAME"
gsettings set $CHILD_SCHEMA:$CHILD_PATH command "$CMD"
gsettings set $CHILD_SCHEMA:$CHILD_PATH binding "$BIND"

echo "✓ bound $BIND → $CMD"
echo "  test it: press Super+D"
