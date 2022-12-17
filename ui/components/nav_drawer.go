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
	DrawerNavItems  []NavHandler
	DrawerUtilItems []NavHandler

	CurrentPage string

	axis      layout.Axis
	textSize  unit.Sp
	leftInset unit.Dp
	width     unit.Dp
	alignment layout.Alignment
	direction layout.Direction
}

func NewNavDrawer(curPage string, navItems, utilItems []NavHandler) NavDrawer {
	nd := NavDrawer{}
	nd.axis = layout.Vertical
	nd.textSize = values.TextSize12
	nd.leftInset = values.MarginPadding0
	nd.width = navDrawerMinimizedWidth
	nd.alignment = layout.Middle
	nd.direction = layout.Center
	nd.CurrentPage = curPage
	nd.DrawerNavItems = navItems
	nd.DrawerUtilItems = utilItems
	return nd
}
func (nd *NavDrawer) Layout(gtx C) D {
	gtx.Constraints.Max.X = gtx.Dp(nd.width)
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

	// 填充背景色
	v.Fill(gtx, values.DarkGray)

	return layout.Flex{
		Axis:    nd.axis,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		// 页面导航类
		layout.Rigid(func(gtx C) D {
			list := layout.List{Axis: nd.axis, Alignment: nd.alignment}
			return list.Layout(gtx, len(nd.DrawerNavItems), func(gtx C, i int) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				item := nd.DrawerNavItems[i]
				img := item.ImageInactive
				if item.PageID == nd.CurrentPage {
					img = item.Image
				}
				return item.Clickable.Button.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
						return nd.direction.Layout(gtx, img.Layout20dp)
					})
				})
			})
		}),

		// 占位并且可移动窗口
		layout.Flexed(1, func(gtx C) D {
			return v.LayoutMove(gtx, func(gtx C) D {
				return D{Size: gtx.Constraints.Max}
			})
		}),

		// 工具类
		layout.Rigid(func(gtx C) D {
			list := layout.List{Axis: nd.axis, Alignment: nd.alignment}
			return list.Layout(gtx, len(nd.DrawerUtilItems), func(gtx C, i int) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				item := nd.DrawerUtilItems[i]
				img := item.ImageInactive
				return item.Clickable.Button.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
						return nd.direction.Layout(gtx, img.Layout20dp)
					})
				})
			})
		}),
	)
}
