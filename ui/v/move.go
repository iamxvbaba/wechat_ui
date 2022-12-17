package v

import (
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// LayoutMove lays out the widget that makes a window movable.
func LayoutMove(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)
	defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()
	system.ActionInputOp(system.ActionMove).Add(gtx.Ops)
	return dims
}
