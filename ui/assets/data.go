package assets

import (
	"bytes"
	"embed"
	"gioui.org/app"
	"gioui.org/widget/material"
	"image"
	"strings"
)

//go:embed *
var Resources embed.FS

var (
	Window   *app.Window
	IconList map[string]image.Image
	Theme    *material.Theme
)

func init() {
	icons, err := Icons()
	if err != nil {
		panic("Error loading icons")
	}
	IconList = icons

	fc := FontCollection()
	//text.NewShaper(fc)
	Theme = material.NewTheme(fc)
}

func Icons() (map[string]image.Image, error) {
	entries, err := Resources.ReadDir("icons")
	if err != nil {
		return nil, err
	}

	icons := make(map[string]image.Image)
	for _, entry := range entries {

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}

		imgBytes, err := Resources.ReadFile("icons/" + entry.Name())
		if err != nil {
			return nil, err
		}

		img, _, err := image.Decode(bytes.NewBuffer(imgBytes))
		if err != nil {
			return nil, err
		}

		split := strings.Split(entry.Name(), ".")
		icons[split[0]] = img
	}

	return icons, nil
}
