package v

import (
	"image"
	"wechat_ui/ui/values"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Clickable struct {
	Button    *widget.Clickable
	style     *values.ClickableStyle
	Hover     bool
	Radius    CornerRadius
	isEnabled bool
}

func NewClickable(hover bool) *Clickable {
	return &Clickable{
		Button:    &widget.Clickable{},
		style:     &values.ClickableStyle{Color: values.SurfaceHighlight, HoverColor: values.Gray5},
		Hover:     hover,
		isEnabled: true,
	}
}

func (cl *Clickable) Style() values.ClickableStyle {
	return *cl.style
}

func (cl *Clickable) ChangeStyle(style *values.ClickableStyle) {
	cl.style = style
}

func (cl *Clickable) Clicked() bool {
	return cl.Button.Clicked()
}

func (cl *Clickable) IsHovered() bool {
	return cl.Button.Hovered()
}

// SetEnabled enables/disables the clickable.
func (cl *Clickable) SetEnabled(enable bool, gtx *layout.Context) layout.Context {
	var mGtx layout.Context
	if gtx != nil && !enable {
		mGtx = gtx.Disabled()
	}

	cl.isEnabled = enable
	return mGtx
}

// Enabled Return clickable enabled/disabled state.
func (cl *Clickable) Enabled() bool {
	return cl.isEnabled
}

func (cl *Clickable) Layout(gtx C, w layout.Widget) D {
	return cl.Button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				tr := gtx.Dp(unit.Dp(cl.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(cl.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(cl.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(cl.Radius.BottomLeft))
				defer clip.RRect{
					Rect: image.Rectangle{Max: image.Point{
						X: gtx.Constraints.Min.X,
						Y: gtx.Constraints.Min.Y,
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()
				clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()

				if cl.Hover && cl.Button.Hovered() {
					paint.Fill(gtx.Ops, cl.style.HoverColor)
				}

				for _, c := range cl.Button.History() {
					drawInk(gtx, c, cl.style.Color)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Stacked(w),
		)
	})
}
