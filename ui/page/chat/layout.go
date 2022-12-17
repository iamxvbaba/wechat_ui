package chat

import (
	"gioui.org/layout"
)

func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	//gtx.Constraints.Min.X = gtx.Constraints.Max.X
	//gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	//return layout.Center.Layout(gtx, material.Body2(assets.Theme, "message").Layout)
	return p.ui.Layout(gtx)
}
