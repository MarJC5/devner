package main

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// buildAppMenu assembles the macOS top-left menu bar. Wails supports
// this natively, so unlike the menubar-right systray (which needs a
// third-party lib that doesn't cooperate with Wails v2's main run
// loop) this works without any CGO shenanigans.
//
// Structure on macOS:
//   [Apple]  Devner  Stack  View  Help
// Windows / Linux get the same menu as a traditional window menu bar.
func (a *App) buildAppMenu() *menu.Menu {
	m := menu.NewMenu()

	// macOS "App" menu (automatic on mac — Wails adds About/Quit etc.
	// when we supply an app-name menu). Keeps the familiar shortcuts.
	m.Append(menu.AppMenu())

	// Edit menu — gives us Cmd+C / Cmd+V on text fields in the webview
	// for free; omitting it breaks clipboard on macOS.
	m.Append(menu.EditMenu())

	// Stack menu — shortcuts for the most common stack actions so the
	// user can start/stop without opening the main window.
	stackMenu := m.AddSubmenu("Stack")
	stackMenu.AddText("Start", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
		if a.deps != nil {
			go func() { _ = a.deps.Runtime.Up(a.ctx) }()
		}
	})
	stackMenu.AddText("Stop", keys.Combo("s", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
		if a.deps != nil {
			go func() { _ = a.deps.Runtime.Stop(a.ctx) }()
		}
	})
	stackMenu.AddText("Rebuild", nil, func(_ *menu.CallbackData) {
		if a.deps != nil {
			go func() { _ = a.deps.Runtime.Rebuild(a.ctx) }()
		}
	})

	// View menu — window-level navigation helpers. Shortcuts mirror
	// the sidebar Cmd+1..7 bindings so muscle memory carries over.
	view := m.AddSubmenu("View")
	view.AddText("Projects", keys.CmdOrCtrl("1"), func(_ *menu.CallbackData) {
		wailsruntime.EventsEmit(a.ctx, "nav:goto", "/projects")
	})
	view.AddText("Stack", keys.CmdOrCtrl("2"), func(_ *menu.CallbackData) {
		wailsruntime.EventsEmit(a.ctx, "nav:goto", "/stack")
	})
	view.AddText("Chat", keys.CmdOrCtrl("3"), func(_ *menu.CallbackData) {
		wailsruntime.EventsEmit(a.ctx, "nav:goto", "/chat")
	})
	view.AddText("Logs", keys.CmdOrCtrl("4"), func(_ *menu.CallbackData) {
		wailsruntime.EventsEmit(a.ctx, "nav:goto", "/logs")
	})
	view.AddText("Databases", keys.CmdOrCtrl("5"), func(_ *menu.CallbackData) {
		wailsruntime.EventsEmit(a.ctx, "nav:goto", "/databases")
	})
	view.AddText("Hosts", keys.CmdOrCtrl("6"), func(_ *menu.CallbackData) {
		wailsruntime.EventsEmit(a.ctx, "nav:goto", "/hosts")
	})
	view.AddText("Settings", keys.CmdOrCtrl("7"), func(_ *menu.CallbackData) {
		wailsruntime.EventsEmit(a.ctx, "nav:goto", "/settings")
	})

	return m
}
