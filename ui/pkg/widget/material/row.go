package material

import (
	"image"
	"image/color"
	"time"
	layout2 "wechat_ui/ui/pkg/layout"
	chatwidget "wechat_ui/ui/pkg/widget"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

// RowStyle configures the presentation of a chat message within
// a vertical list of chat messages.
type RowStyle struct {
	layout2.Row
	// Local indicates that the message was sent by the local user,
	// and should be right-aligned.
	Local bool
	// Time is the timestamp associated with the message.
	Time material.LabelStyle
	// StatusIcon is an optional icon that will be displayed to the right of
	// the message instead of its timestamp.
	StatusIcon *widget.Icon
	// StatusIconColor is the color of the status icon, if any is set.
	StatusIconColor color.NRGBA
	// IconSize defines the size of the StatusIcon (if it is set).
	IconSize unit.Dp
	// StatusMessage defines a warning message to be displayed beneath the
	// chat message.
	StatusMessage material.LabelStyle
	// UserInfoStyle configures how the sender's information is displayed.
	UserInfoStyle
	// MessageStyle configures how the text and its background are presented.
	MessageStyle
	// Interaction holds the interactive state of this message.
	Interaction *chatwidget.Row
	// Menu configures the right-click context menu for this message.
	Menu component.MenuStyle
}

// RowConfig describes the aspects of a chat message relevant for
// displaying it within a widget.
type RowConfig struct {
	Sender  string
	Avatar  image.Image
	Content string
	SentAt  time.Time
	Image   image.Image
	Local   bool
	Status  string
}

// NewRow creates a style type that can lay out the data for a message.
func NewRow(th *material.Theme, interact *chatwidget.Row, menu *component.MenuState, msg RowConfig) RowStyle {
	if interact == nil {
		interact = &chatwidget.Row{}
	}
	if menu == nil {
		menu = &component.MenuState{}
	}
	ms := RowStyle{
		Row: layout2.Row{
			Margin:         layout2.VerticalMargin(),
			InternalMargin: layout2.VerticalMargin(),
			Gutter:         layout2.Gutter(),
			Direction:      layout.W,
		},
		Time:          material.Body2(th, msg.SentAt.Local().Format("15:04")),
		Local:         msg.Local,
		IconSize:      unit.Dp(32),
		UserInfoStyle: UserInfo(th, &interact.UserInfo, msg.Sender, msg.Avatar),
		Interaction:   interact,
		Menu:          component.Menu(th, menu),
		MessageStyle:  Message(th, &interact.Message, msg.Content, msg.Image),
	}
	ms.UserInfoStyle.Local = msg.Local
	if msg.Local {
		ms.Row.Direction = layout.E
	}
	if msg.Status != "" {
		ms.StatusMessage = material.Body2(th, msg.Status)
		ms.StatusMessage.Color = DefaultDangerColor
		ms.StatusIcon = ErrorIcon
		ms.StatusIconColor = DefaultDangerColor
	}
	return ms
}

// Layout the message.
func (c RowStyle) Layout(gtx C) D {
	return c.Row.Layout(gtx,
		layout2.ContentRow(c.UserInfoStyle.Layout),
		layout2.FullRow(nil, c.layoutBubble, c.layoutTimeOrIcon),
		layout2.UnifiedRow(c.layoutStatusMessage),
	)
}

// layoutBubble lays out the chat bubble.
func (c RowStyle) layoutBubble(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return c.MessageStyle.Layout(gtx)
		}),
		layout.Expanded(func(gtx C) D {
			return c.Interaction.ContextArea.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = image.Point{}
				return c.Menu.Layout(gtx)
			})
		}),
	)
}

// layoutTimeOrIcon lays out a status icon if one is set, and
// otherwise lays out the time the messages was sent.
func (c RowStyle) layoutTimeOrIcon(gtx C) D {
	return layout.Center.Layout(gtx, func(gtx C) D {
		if c.StatusIcon == nil {
			return c.Time.Layout(gtx)
		}
		sideLength := gtx.Dp(c.IconSize)
		gtx.Constraints.Max.X = sideLength
		gtx.Constraints.Max.Y = sideLength
		gtx.Constraints.Min = gtx.Constraints.Constrain(gtx.Constraints.Min)
		return c.StatusIcon.Layout(gtx, c.StatusIconColor)
	})
}

// layoutStatusMessage lays out status message text, if any.
func (c RowStyle) layoutStatusMessage(gtx C) D {
	if c.StatusMessage.Text == "" {
		return D{}
	}
	return layout.E.Layout(gtx, c.StatusMessage.Layout)
}
