package main

import (
	"embed"

	"github.com/devner/devner/internal/sysenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// macOS-only: when launched as a .app via LaunchServices the PATH
	// is stripped down to /usr/bin:/bin:/usr/sbin:/sbin, which hides
	// `docker` (lives in /usr/local/bin via Docker Desktop). Restore
	// the dev paths before anything shells out.
	sysenv.EnsureDevPath()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Devner",
		Width:     1120,
		Height:    760,
		MinWidth:  900,
		MinHeight: 600,

		// Transparent backgrounds so the OS-level vibrancy/acrylic
		// effect shows through. The actual "color" is painted by the
		// HTML body in App.css with an rgba() for semi-opacity.
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		// macOS — HiddenInset keeps the red/yellow/green traffic lights
		// but removes the title bar chrome so the whole top of the
		// window is part of the React area. WindowIsTranslucent +
		// WebviewIsTransparent = desktop / wallpaper shows through.
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Devner",
				Message: "Local dev environment orchestrator.",
			},
		},

		// Windows — Mica (Win 11) + acrylic (Win 10 fallback) give the
		// same frosted-glass feel.
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},

		// Native macOS top-left menu bar. On Linux / Windows this
		// becomes the window's own menu bar. Built here rather than in
		// OnStartup because Wails wires it up once at construction.
		Menu: app.buildAppMenu(),

		OnStartup:     app.startup,
		OnShutdown:    app.shutdown,
		OnBeforeClose: app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
