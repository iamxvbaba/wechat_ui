package start

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"wechat_ui/ui/assets"
)

func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	return layout.Center.Layout(gtx, material.Body2(assets.Theme, "start page").Layout)
}
