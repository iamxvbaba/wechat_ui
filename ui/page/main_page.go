package page

import (
	"context"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"wechat_ui/app"
	"wechat_ui/ui/components"
	"wechat_ui/ui/page/contact"
	"wechat_ui/ui/page/msg"
	"wechat_ui/ui/page/start"
	"wechat_ui/ui/v"
)

const (
	MainPageID = "Main"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	NavDrawerWidth          = unit.Dp(160)
	NavDrawerMinimizedWidth = unit.Dp(72)
)

type NavHandler struct {
	Clickable     *widget.Clickable
	Image         *v.Image
	ImageInactive *v.Image
	Title         string
	PageID        string
}

type MainPage struct {
	*app.MasterPage
	ctx       context.Context
	ctxCancel context.CancelFunc
	drawerNav components.NavDrawer
}

func NewMainPage() *MainPage {
	mp := &MainPage{
		MasterPage: app.NewMasterPage(MainPageID),
	}

	mp.initNavItems()

	return mp
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (mp *MainPage) ID() string {
	return MainPageID
}

func (mp *MainPage) initNavItems() {
	navItems := []components.NavHandler{
		{
			Clickable:     v.NewClickable(false),
			Image:         v.MsgIcon,
			ImageInactive: v.MsgIconInactive,
			Title:         "消息",
			PageID:        msg.PageID,
		},
		{
			Clickable:     v.NewClickable(false),
			Image:         v.ContactIcon,
			ImageInactive: v.ContactIconInactive,
			Title:         "通讯录",
			PageID:        contact.PageID,
		},
	}
	mp.drawerNav = components.NewNavDrawer(mp.CurrentPageID(), navItems)
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (mp *MainPage) OnNavigatedTo() {
	mp.ctx, mp.ctxCancel = context.WithCancel(context.TODO())

	// 给一个启动页面
	if mp.CurrentPage() == nil {
		mp.Display(start.NewPage())
	}

	mp.CurrentPage().OnNavigatedTo()

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (mp *MainPage) HandleUserInteractions() {
	if mp.CurrentPage() != nil {
		mp.CurrentPage().HandleUserInteractions()
	}

	mp.drawerNav.CurrentPage = mp.CurrentPageID()

	// 加载左侧导航栏
	for _, item := range mp.drawerNav.DrawerNavItems {
		for item.Clickable.Clicked() {
			var pg app.Page
			switch item.PageID {
			case contact.PageID:
				pg = contact.NewPage()
			case msg.PageID:
				pg = msg.NewPage()
			}

			if pg == nil || mp.ID() == mp.CurrentPageID() {
				continue
			}
			mp.Display(pg)
		}
	}
}

// KeysToHandle 监听的键盘事件
func (mp *MainPage) KeysToHandle() key.Set {
	if currentPage := mp.CurrentPage(); currentPage != nil {

	}
	return ""
}

// HandleKeyPress 处理键盘事件.
func (mp *MainPage) HandleKeyPress(evt *key.Event) {
	if currentPage := mp.CurrentPage(); currentPage != nil {

	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (mp *MainPage) OnNavigatedFrom() {
	// Also disappear all child pages.
	if mp.CurrentPage() != nil {
		mp.CurrentPage().OnNavigatedFrom()
	}

	mp.ctxCancel()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (mp *MainPage) Layout(gtx C) D {
	return mp.layoutDesktop(gtx)
}

func (mp *MainPage) layoutDesktop(gtx C) D {
	return v.LinearLayout{
		Width:       v.MatchParent,
		Height:      v.MatchParent,
		Orientation: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(mp.drawerNav.Layout),
		layout.Rigid(func(gtx C) D {
			if mp.CurrentPage() == nil {
				return D{}
			}
			return mp.CurrentPage().Layout(gtx)
		}),
	)
}
