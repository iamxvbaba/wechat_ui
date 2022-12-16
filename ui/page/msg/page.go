package msg

import (
	"wechat_ui/app"
)

const PageID = "msg"

type Page struct {
	*app.GenericPageModal
}

func (p *Page) OnNavigatedTo() {
}

func (p *Page) OnNavigatedFrom() {
}

func NewPage() *Page {
	page := &Page{
		GenericPageModal: app.NewGenericPageModal(PageID),
	}

	return page
}

func (p *Page) HandleUserInteractions() {

}
