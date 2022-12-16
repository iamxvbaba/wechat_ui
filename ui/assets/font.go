package assets

import (
	"fmt"
	"sync"

	"gioui.org/font/opentype"
	"gioui.org/text"
)

var (
	once       sync.Once
	collection []text.FontFace
)

// FontCollection registers the fonts to used in the app
func FontCollection() []text.FontFace {
	msyh, err := getFontByte("fonts/chinese.msyh.ttf") // 微软雅黑字体
	if err != nil {
		panic(err)
	}
	once.Do(func() {
		register(text.Font{}, msyh)
		n := len(collection)
		collection = collection[:n:n]
	})
	return collection
}

func register(fnt text.Font, fontByte []byte) {
	face, err := opentype.Parse(fontByte)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	fnt.Typeface = "Go"
	collection = append(collection, text.FontFace{Font: fnt, Face: face})
}

func getFontByte(path string) ([]byte, error) {
	return content.ReadFile(path)
}
