; AutoHotkey v2 script: binds Win+D to open `devner prompt` in Windows Terminal.
;
; Install:
;   1. Install AutoHotkey v2: https://www.autohotkey.com/
;   2. Install Windows Terminal (wt.exe) from the Microsoft Store.
;   3. Save this file as %USERPROFILE%\devner.ahk
;   4. Put a shortcut to it in %APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup
;      so it loads at login.
;
; Note: Win+D is normally "show desktop" in Windows. If that conflicts,
; change `#d` below to for example `^+d` (Ctrl+Shift+D).

#Requires AutoHotkey v2.0

#d::
{
    ; Adjust the path if devner.exe is elsewhere.
    Run 'wt.exe new-tab --title "Devner" -- devner.exe prompt'
}
