package appwidget

import (
	"gioui.org/widget"
	chatwidget "wechat_ui/ui/pkg/widget"
)

// Room selector state.
type Room struct {
	widget.Clickable
	Image  chatwidget.CachedImage
	Active bool
}
