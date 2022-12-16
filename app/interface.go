package app

import (
	"gioui.org/layout"
)

// Page 定义控制窗口上显示的UI组件的外观和功能的方法。
type Page interface {
	// ID 是标识页面的唯一字符串，可用于将此页面与其他页面区分开来。.
	ID() string
	// OnAttachedToNavigator 导航发生时调用即当页面或者模态框推入窗口的显示时
	// navigator 参数是用于显示内容的 PageNavigator 或 WindowNavigator 对象。这是在调用 OnNavigatedTo() 之前调用的。
	OnAttachedToNavigator(navigator PageNavigator)
	// OnNavigatedTo 在页面即将显示时调用，可用于初始化仅在页面显示时相关的页面功能。
	//在调用 HandleUserInteractions() 和 Layout() 之前调用它。
	OnNavigatedTo()
	// HandleUserInteractions
	// 在 Layout() 之前调用以确定最近是否在页面上发生了任何用户交互，并且可用于在页面的 UI 组件显示前不久更新它们。
	HandleUserInteractions()
	// Layout 组件绘制到提供的布局上下文中，最终在屏幕上绘制
	Layout(layout.Context) layout.Dimensions
	// OnNavigatedFrom 当页面即将从显示的窗口中移除时调用.
	// 理想情况下，此方法应用于禁用不显示页面时不相关的功能.
	// NOTE: 页面可能会重新显示在应用程序的窗口中，在这种情况下，将再次调用 OnNavigatedTo()。
	// 此方法不应破坏 UI 组件，除非它们将在 OnNavigatedTo() 方法中重新创建。
	OnNavigatedFrom()
}

type Modal interface {
	// ID 是标识模态的唯一字符串，可用于将此模态与其他模态区分开来。
	ID() string
	// OnAttachedToNavigator 导航发生时调用即当页面或者模态框推入窗口的显示时
	// navigator 参数是用于显示内容的 PageNavigator 或 WindowNavigator 对象。这在 OnResume() 被调用之前被调用
	OnAttachedToNavigator(navigator PageNavigator)
	// OnResume 调用OnResume来初始化数据并准备好要显示的 UI 元素.
	OnResume()
	//Handle 在 Layout() 之前调用，以确定最近是否在模态框上发生了任何用户交互，并且可用于在页面的 UI 组件显示前不久更新它们。
	Handle()
	// Layout 将模式的 UI 组件绘制到提供的布局上下文中，最终在屏幕上绘制。
	Layout(gtx layout.Context) layout.Dimensions
	// OnDismiss 在模态关闭后调用.
	// NOTE: 模式可能会重新显示在应用程序的窗口上，在这种情况下将再次调用 OnResume()。
	// 此方法不应破坏 UI 组件，除非它们将在 OnResume() 方法中重新创建。
	OnDismiss()
}

// Closable 应该由想要知道何时关闭以执行一些清理操作的 Page 和 Modal 来实现。
type Closable interface {
	// OnClosed 调用 OnClosed 以指示页面或模式的特定实例已被解除并且将不再显示。
	OnClosed()
}

// PageNavigator 定义了在窗口或 MasterPage 中的页面之间导航的方法。
type PageNavigator interface {
	// CurrentPage 返回位于堆栈顶部的页面。如果堆栈为空，则返回 nil。
	CurrentPage() Page
	// CurrentPageID 返回当前页面的 ID，如果没有显示页面，则返回空字符串.
	CurrentPageID() string
	// Display 使指定页面显示在父窗口上或页面上. 同一页面的所有其他实例将被关闭并从后台堆栈中删除
	Display(page Page)
	// CloseCurrentPage 关闭堆栈顶部的页面并准备好显示下一页。
	CloseCurrentPage()
	// ClosePagesAfter 关闭堆栈顶部的所有页面，直到（并排除）具有指定 ID 的页面。如果没有找到具有提供的 ID 的页面，
	// 则不会弹出任何页面。在弹出其他页面后，将显示具有指定 ID 的页面。
	ClosePagesAfter(keepPageID string)
	// ClearStackAndDisplay 关闭堆栈中的所有页面并显示指定页面。
	ClearStackAndDisplay(page Page)
	// CloseAllPages 关闭堆栈中的所有页面.
	CloseAllPages()
}

// WindowNavigator 定义了页面导航、显示模态和重新加载整个窗口显示的方法
type WindowNavigator interface {
	PageNavigator
	// ShowModal 在当前页面上显示一个模态。任何以前显示的模态都将被这个新模态隐藏。
	ShowModal(Modal)
	// DismissModal 如果此 WindowNavigator 先前显示了具有指定 ID 的模式，则解除该模式。
	// 如果具有指定 ID 的模式超过 1 个，则仅关闭最顶层的实例.
	DismissModal(modalID string)
	// TopModal 返回显示中最顶层的模态，如果显示中没有模态，则返回 nil。
	TopModal() Modal
	// Reload 重新加载整个窗口显示.
	// 如果当前显示页面，则应调用页面的 HandleUserInteractions() 方法。
	// 如果显示模态框，则还应调用模态框的 Handle() 方法。
	Reload()
}
