package main

import (
	"gioui.org/app"
	"log"
	"os"
	"wechat_ui/ui"
)

func main() {
	win, err := ui.CreateWindow()
	if err != nil {
		log.Printf("Could not initialize window: %s\ns", err)
		return
	}

	go func() {
		win.HandleEvents() // blocks until the app window is closed
		os.Exit(0)
	}()

	// Start the GUI frontend.
	app.Main()
}
