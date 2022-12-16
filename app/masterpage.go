package app

// MasterPage  是一个可以显示子页面的页面.
// 它是 GenericPageModal 的扩展，提供对用于显示 MasterPage 的 Window 或 PageNavigator 的访问.
// MasterPage 的 ParentNavigator 通常在 MasterPage 被 WindowNavigator 或 PageNavigator 推入显示窗口时设置.
// MasterPage 必须由要显示子页面的页面嵌入. 这些页面必须满足 MasterPage 尚未满足的 Page 接口的其他方法.
type MasterPage struct {
	*GenericPageModal
	subPages *PageStack
}

// NewMasterPage returns an instance of MasterPage.
func NewMasterPage(id string) *MasterPage {
	return &MasterPage{
		GenericPageModal: NewGenericPageModal(id),
		subPages:         NewPageStack(id),
	}
}

// CurrentPage 返回位于堆栈顶部的页面。如果堆栈为空，则返回 nil.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CurrentPage() Page {
	return masterPage.subPages.Top()
}

//CurrentPageID 返回当前页面的 ID，如果没有显示页面，则返回空字符串
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CurrentPageID() string {
	if currentPage := masterPage.CurrentPage(); currentPage != nil {
		return currentPage.ID()
	}
	return ""
}

// Display 使指定的页面显示在父窗口或页面上。同一页面的所有其他实例将被关闭并从后台堆栈中删除
// Part of the PageNavigator interface.
func (masterPage *MasterPage) Display(newPage Page) {
	pushed := masterPage.subPages.Push(newPage, masterPage)
	if pushed {
		masterPage.ParentWindow().Reload()
	}
}

// CloseCurrentPage 关闭堆栈顶部的页面并准备好显示下一页.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CloseCurrentPage() {
	popped := masterPage.subPages.Pop()
	if popped {
		masterPage.ParentWindow().Reload()
	}
}

// ClosePagesAfter 关闭堆栈顶部的所有页面，直到（并排除）具有指定 ID 的页面。
// 如果没有找到具有提供的 ID 的页面，则不会弹出任何页面。弹出其他页面后会显示指定ID的页面.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) ClosePagesAfter(keepPageID string) {
	popped := masterPage.subPages.PopAfter(func(page Page) bool {
		return page.ID() == keepPageID
	})
	if popped {
		masterPage.ParentWindow().Reload()
	}
}

// ClearStackAndDisplay 关闭堆栈中的所有页面并显示指定的页面。
// Part of the PageNavigator interface.
func (masterPage *MasterPage) ClearStackAndDisplay(newPage Page) {
	newPage.OnAttachedToNavigator(masterPage)
	masterPage.subPages.Reset(newPage)
	masterPage.ParentWindow().Reload()
}

// CloseAllPages 关闭堆栈中的所有页面.
// Part of the PageNavigator interface.
func (masterPage *MasterPage) CloseAllPages() {
	masterPage.subPages.Reset()
	masterPage.ParentWindow().Reload()
}
