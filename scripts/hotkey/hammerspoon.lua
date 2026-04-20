-- Hammerspoon binding for a global Cmd+D hotkey that opens `devner prompt`
-- in a new iTerm2 window (or Terminal.app if iTerm is missing).
--
-- Install:
--   1. Install Hammerspoon: brew install --cask hammerspoon
--   2. Open Hammerspoon, grant Accessibility permission.
--   3. Append this file's contents to ~/.hammerspoon/init.lua
--   4. Reload config from the Hammerspoon menu bar.

local devner = "/usr/local/bin/devner"  -- adjust if `make install` used a different dir

hs.hotkey.bind({"cmd"}, "d", function()
    -- Prefer iTerm2 if installed.
    local iterm = hs.application.find("iTerm2")
    if iterm then
        hs.osascript.applescript([[
            tell application "iTerm"
                create window with default profile
                tell current session of current window
                    write text "]] .. devner .. [[ prompt"
                end tell
            end tell
        ]])
    else
        hs.osascript.applescript([[
            tell application "Terminal"
                activate
                do script "]] .. devner .. [[ prompt"
            end tell
        ]])
    end
end)
