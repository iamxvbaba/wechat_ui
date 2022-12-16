package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"wechat_ui/ui/v"
	"wechat_ui/ui/values"
)

var (
	navDrawerMinimizedWidth = values.MarginPadding48
)

type NavHandler struct {
	Clickable     *v.Clickable
	Image         *v.Image
	ImageInactive *v.Image
	Title         string
	PageID        string
}

type NavDrawer struct {
	DrawerNavItems []NavHandler
	CurrentPage    string

	axis      layout.Axis
	textSize  unit.Sp
	leftInset unit.Dp
	width     unit.Dp
	alignment layout.Alignment
	direction layout.Direction
}

func NewNavDrawer(curPage string, navItems []NavHandler) NavDrawer {
	nd := NavDrawer{}
	nd.axis = layout.Vertical
	nd.textSize = values.TextSize12
	nd.leftInset = values.MarginPadding0
	nd.width = navDrawerMinimizedWidth
	nd.alignment = layout.Middle
	nd.direction = layout.Center
	nd.CurrentPage = curPage
	nd.DrawerNavItems = navItems
	return nd
}
func (nd *NavDrawer) Layout(gtx layout.Context) layout.Dimensions {
	return v.LinearLayout{
		Width:       gtx.Dp(nd.width),
		Height:      v.MatchParent,
		Orientation: layout.Vertical,
		Background:  values.DarkGray, //values.Surface,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(nd.DrawerNavItems), func(gtx C, i int) D {
				mGtx := gtx
				background := values.DarkGray //values.Surface

				//if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
				//	background = values.Gray5
				//}
				return v.LinearLayout{
					Orientation: nd.axis,
					Width:       v.MatchParent,
					Height:      v.WrapContent,
					Padding:     layout.UniformInset(values.MarginPadding10),
					Alignment:   nd.alignment,
					Direction:   nd.direction,
					Background:  background,
					Clickable:   nd.DrawerNavItems[i].Clickable,
				}.Layout(mGtx,
					layout.Rigid(func(gtx C) D {
						img := nd.DrawerNavItems[i].ImageInactive

						if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
							img = nd.DrawerNavItems[i].Image
						}
						return img.Layout18dp(gtx)
					}),
					//layout.Rigid(func(gtx C) D {
					//	return layout.Inset{
					//		Left: nd.leftInset,
					//	}.Layout(gtx, func(gtx C) D {
					//		textColor := values.GrayText1
					//		if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
					//			textColor = values.DeepBlue
					//		}
					//		txt := v.NewLabel(nd.textSize, nd.DrawerNavItems[i].Title)
					//		txt.Color = textColor
					//		return txt.Layout(gtx)
					//	})
					//}),
				)
			})
		}),
	)
}
