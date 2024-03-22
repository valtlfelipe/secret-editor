package main

import (
	"embed"
	"log"
	r "runtime"

	"github.com/valtlfelipe/secret-editor/backend/services"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	// "github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	pref := services.NewPreferences()
	err := pref.LoadPreferences()
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	// Create an instance of the app structure
	app := NewApp(*pref)

	appMenu := menu.NewMenu()

	if r.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
	}

	FileMenu := appMenu.AddSubmenu("File")
	FileMenu.AddText("New", keys.CmdOrCtrl("n"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(app.ctx, "command:new")
	})
	FileMenu.AddSeparator()
	FileMenu.AddText("Open", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(app.ctx, "command:open")
	})
	FileMenu.AddSeparator()
	FileMenu.AddText("Save", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(app.ctx, "command:save")
	})

	if r.GOOS == "darwin" {
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Secret Editor",
		Width:  1024,
		Height: 768,
		Menu:   appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 249, G: 250, B: 251, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		// Mac: &mac.Options{
		// 	About: &mac.AboutInfo{
		// 		Title:   "My Application",
		// 		Message: "© 2021 Me",
		// 	},
		// },
	})

	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
}
