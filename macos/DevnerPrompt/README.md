# DevnerPrompt.app

Native macOS floating popup for the devner agent. Press **Cmd+D** anywhere, a panel slides in over the current app, you type, the agent answers. Esc to hide.

No dock icon, no menu bar item, no third-party dependencies. Pure Swift + AppKit + Carbon hotkey API.

## Build

Requires Xcode Command Line Tools (`xcode-select --install`) — no full Xcode needed.

```bash
make mac-app             # from repo root
# or:
bash macos/DevnerPrompt/build.sh
```

Output: `macos/DevnerPrompt/build/DevnerPrompt.app` (~120 KB, arm64 native, universal if you change `build.sh`).

## Install

```bash
make mac-app-install     # copies to /Applications/
open /Applications/DevnerPrompt.app
```

The first `open` registers the Cmd+D hotkey. The app has `LSUIElement=YES` so it runs headlessly — no dock icon.

Make it launch at login:
> System Settings → General → Login Items → "+" → `/Applications/DevnerPrompt.app`

## Uninstall

```bash
rm -rf /Applications/DevnerPrompt.app
# Kill a running copy:
pkill -x DevnerPrompt || true
```

## How it works

- **Global hotkey** via `RegisterEventHotKey` (Carbon). No Accessibility permission required — Carbon hotkeys work app-locally but because we call from an LSUIElement app with no owning window, the registration is effectively system-wide within this process.
- **Panel**: `NSPanel` with `.nonactivatingPanel` + `.floating` level, centered on the screen under the cursor.
- **Agent invocation**: shells out to `/usr/local/bin/devner agent --yes --raw <prompt>` (or `/opt/homebrew/bin/devner`, or `$DEVNER_BIN`). `--raw` disables markdown rendering so we get clean text we can drop into `NSTextView`.
- **Streaming**: `Process` + `Pipe` with `readabilityHandler` that dispatches to the main thread and appends to the text view.

## Config

Environment variable read at launch:

| Var | Default | Purpose |
|---|---|---|
| `DEVNER_BIN` | `/usr/local/bin/devner` (or `/opt/homebrew/bin/devner`) | Path to the devner Go binary. |

To change the hotkey or binding, edit `main.swift`:

```swift
let keyCode: UInt32 = UInt32(kVK_ANSI_D)       // virtual key
let modifiers: UInt32 = UInt32(cmdKey)          // cmdKey | shiftKey | optionKey | controlKey
```

Virtual keys: see `Carbon/Events.h` or <https://developer.apple.com/documentation/appkit/1535851-function_key_unicodes>.
