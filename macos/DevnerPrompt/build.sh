#!/usr/bin/env bash
# Build DevnerPrompt.app — a standalone macOS app that binds Cmd+D to a
# floating prompt panel calling `devner agent`.
#
# Requirements:
#   - Xcode Command Line Tools (for swiftc). Install: xcode-select --install
#   - macOS 12+ target (arm64 or x86_64 — autodetected)
#
# Output: ./build/DevnerPrompt.app  (drag it into /Applications/ or run in place)

set -euo pipefail

cd "$(dirname "$0")"

APP_NAME="DevnerPrompt"
BUNDLE_ID="com.devner.prompt"
OUT_DIR="build"
APP_PATH="$OUT_DIR/$APP_NAME.app"
EXEC_PATH="$APP_PATH/Contents/MacOS/$APP_NAME"

rm -rf "$OUT_DIR"
mkdir -p "$APP_PATH/Contents/MacOS"
mkdir -p "$APP_PATH/Contents/Resources"

# Compile — universal binary when possible (arm64 + x86_64). swiftc -target
# with -arch flags gives us a fat binary we don't need to ship twice.
ARCH_FLAGS="-target arm64-apple-macosx12.0"
if [[ "$(uname -m)" == "x86_64" ]]; then
    ARCH_FLAGS="-target x86_64-apple-macosx12.0"
fi

echo "→ compiling $APP_NAME ($ARCH_FLAGS)"
swiftc main.swift \
    $ARCH_FLAGS \
    -O \
    -framework Cocoa \
    -framework Carbon \
    -o "$EXEC_PATH"

# Info.plist — LSUIElement=YES hides the dock icon (prompt is a global
# popup, not a windowed app).
cat > "$APP_PATH/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleDisplayName</key>
    <string>Devner Prompt</string>
    <key>CFBundleExecutable</key>
    <string>$APP_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>$BUNDLE_ID</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>12.0</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
PLIST

chmod +x "$EXEC_PATH"
echo ""
echo "✓ built $APP_PATH"
echo ""
echo "Run it once to register Cmd+D:"
echo "  open $APP_PATH"
echo ""
echo "Install globally:"
echo "  cp -R $APP_PATH /Applications/"
echo "  open /Applications/$APP_NAME.app"
echo ""
echo "Add to Login Items (System Settings > General > Login Items) so the hotkey"
echo "persists across reboots."
