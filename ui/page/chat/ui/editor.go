package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"strings"
	"wechat_ui/ui/v"
)

var (
	leftIcons = []*v.Image{
		v.Emoticon,
		v.File,
		v.Screenshot,
	}

	rightIcons = []*v.Image{
		v.Circle,
		v.Call,
	}
)

// layoutEditor lays out the message editor.
func (ui *UI) layoutEditor2(gtx C) D {
	if ui.AddBtn.Clicked() {
		active := ui.Rooms.Active()
		text := strings.TrimSpace(active.Editor.Text())
		if text != "" {
			active.SendLocal(text)
			active.Editor.SetText("")
		}
	}
	if ui.DeleteBtn.Clicked() {
		serial := ui.ContextMenuTarget.Serial()
		ui.Rooms.Active().DeleteRow(serial)
	}
	active := ui.Rooms.Active()
	editor := &active.Editor
	for _, e := range editor.Events() {
		switch e.(type) {
		case widget.SubmitEvent:
			text := strings.TrimSpace(editor.Text())
			if text != "" {
				active.SendLocal(text)
				editor.SetText("")
			}
		}
	}
	editor.Submit = true

	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 分割线
		layout.Rigid(v.NewSeparator(component.WithAlpha(th.Fg, 50)).Layout),
		// 工具栏
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							list := layout.List{Axis: layout.Horizontal, Alignment: layout.Start}
							return list.Layout(gtx, len(leftIcons), func(gtx C, index int) D {
								return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
									return leftIcons[index].Layout20dp(gtx)
								})
							})
						}),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							list := layout.List{Axis: layout.Horizontal, Alignment: layout.Start}
							return list.Layout(gtx, len(rightIcons), func(gtx C, index int) D {
								return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
									return rightIcons[index].Layout20dp(gtx)
								})
							})
						}),
				)
			})
		}),
		// 输入框
		layout.Rigid(func(gtx C) D {
			// 限定最低宽度
			height := gtx.Dp(unit.Dp(50))
			gtx.Constraints.Max.Y = height
			gtx.Constraints.Min.Y = height
			// 限定输入框长度
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, material.Editor(th.Theme, editor, "Send a message").Layout)
		}),
		// 发送按钮
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			in := layout.Inset{Bottom: unit.Dp(8), Right: unit.Dp(8)}
			return in.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// .E 放到最右边
				return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th.Theme, &ui.AddBtn, "Send(S)")
					return btn.Layout(gtx)
				})
			})
		}),
	)
}
