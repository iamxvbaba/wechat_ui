package chat

import (
	"wechat_ui/app"
	"wechat_ui/ui/assets"
	"wechat_ui/ui/page/chat/ui"
)

const PageID = "chat"

type Page struct {
	*app.GenericPageModal
	ui *ui.UI
}

func (p *Page) OnNavigatedTo() {
}

func (p *Page) OnNavigatedFrom() {
}

func NewPage() *Page {
	pm := app.NewGenericPageModal(PageID)
	page := &Page{
		GenericPageModal: pm,
		ui: ui.NewUI(assets.Window.Invalidate, ui.Config{
			Theme:      "light",
			Latency:    1000,
			LoadSize:   30,
			BufferSize: 30,
		}),
	}

	return page
}

func (p *Page) HandleUserInteractions() {
}
