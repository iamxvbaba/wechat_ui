package chat

import (
	"gioui.org/layout"
)

func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	return p.ui.Layout(gtx)
}
