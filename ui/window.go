package ui

import (
	"wechat_ui/app"
	"wechat_ui/ui/assets"
	"wechat_ui/ui/page"
	"wechat_ui/ui/v"
	"wechat_ui/ui/values"

	giouiApp "gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

type Window struct {
	*giouiApp.Window
	navigator app.WindowNavigator
}

type (
	C = layout.Context
	D = layout.Dimensions
)
type WriteClipboard struct {
	Text string
}

func CreateWindow() (*Window, error) {
	giouiWindow := giouiApp.NewWindow(giouiApp.MinSize(values.AppWidth, values.AppHeight),
		giouiApp.Title("wechat"),
		giouiApp.Decorated(false)) // giouiApp.Decorated(false) 去掉程序顶部默认装饰

	// 全局
	assets.Window = giouiWindow

	win := &Window{
		Window:    giouiWindow,
		navigator: app.NewSimpleWindowNavigator(giouiWindow.Invalidate),
	}

	return win, nil
}

func (win *Window) HandleEvents() {

	for {
		e := <-win.Events()
		switch evt := e.(type) {

		case system.DestroyEvent:
			win.navigator.CloseAllPages()
			return // exits the loop, caller will exit the program.

		case system.FrameEvent:
			ops := win.handleFrameEvent(evt)
			evt.Frame(ops)
		default:
			//log.Printf("Unhandled window event %v\n", e)
		}
	}
}

// handleFrameEvent 处理事件
func (win *Window) handleFrameEvent(evt system.FrameEvent) *op.Ops {
	switch {
	case win.navigator.CurrentPage() == nil:
		// 直接进入主页面.
		win.navigator.Display(page.NewMainPage())

	default:
		// 应用程序窗口可能已经接收到一些触发此 FrameEvent 的用户交互，例如按键、按钮单击等。
		// 在重新显示 UI 组件之前处理此类交互。这可确保根据用户刚刚执行的操作向用户显示正确的界面。.
		win.handleRelevantKeyPresses(evt)
		win.navigator.CurrentPage().HandleUserInteractions()
		if modal := win.navigator.TopModal(); modal != nil {
			modal.Handle()
		}
	}

	// 将窗口的 UI 组件绘制到屏幕上
	ops := win.prepareToDisplayUI(evt)

	// 监听键盘事件
	win.addKeyEventRequestsToOps(ops)

	return ops
}

// handleRelevantKeyPresses 检查任何打开的模态框或页面是否是 load.KeyEventHandler
// 以及提供的 system.FrameEvent 是否包含模式或页面的按键事件
func (win *Window) handleRelevantKeyPresses(evt system.FrameEvent) {
	handleKeyPressFor := func(tag string, maybeHandler interface{}) {
		for _, event := range evt.Queue.Events(tag) {
			if keyEvent, isKeyEvent := event.(key.Event); isKeyEvent && keyEvent.State == key.Press {

			}
		}
	}

	if modal := win.navigator.TopModal(); modal != nil {
		handleKeyPressFor(modal.ID(), modal)
	} else {
		handleKeyPressFor(win.navigator.CurrentPageID(), win.navigator.CurrentPage())
	}
}

func (win *Window) prepareToDisplayUI(evt system.FrameEvent) *op.Ops {
	backgroundWidget := layout.Expanded(func(gtx C) D {
		return v.Fill(gtx, values.Gray4)
	})

	currentPageWidget := layout.Stacked(func(gtx C) D {
		if modal := win.navigator.TopModal(); modal != nil {
			gtx = gtx.Disabled()
		}
		return win.navigator.CurrentPage().Layout(gtx)
	})

	topModalLayout := layout.Stacked(func(gtx C) D {
		modal := win.navigator.TopModal()
		if modal == nil {
			return layout.Dimensions{}
		}
		return modal.Layout(gtx)
	})

	ops := &op.Ops{}
	gtx := layout.NewContext(ops, evt)
	layout.Stack{Alignment: layout.N}.Layout(
		gtx,
		backgroundWidget,
		currentPageWidget,
		topModalLayout,
	)

	return ops
}

func (win *Window) addKeyEventRequestsToOps(ops *op.Ops) {
}
