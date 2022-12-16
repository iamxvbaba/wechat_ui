package app

// GenericPageModal 实现了 ID() 和 OnAttachedToNavigator() 方法
// 大多数页面和模态框都需要。它还定义了 ParentNavigator() 和
// ParentWindow() 辅助方法，使页面能够访问导航器
// 显示页面和根 WindowNavigator。实际的页面和模态框可以嵌入这个结构并根据需要实现其他方法
type GenericPageModal struct {
	id        string
	parentNav PageNavigator
}

// NewGenericPageModal returns an instance of a GenericPageModal.
func NewGenericPageModal(id string) *GenericPageModal {
	return &GenericPageModal{
		id: id,
	}
}

// ID is a unique string that identifies this page or modal and may be used to
// differentiate this page or modal from other pages or modals.
// Part of the Page and Modal interfaces.
func (pageModal *GenericPageModal) ID() string {
	return pageModal.id
}

// OnAttachedToNavigator 导航发生时调用,即当此页面或模式被推入窗口的显示时。
// navigator 参数是用于显示此页面或模态的 PageNavigator 或 WindowNavigator 对象。
//OnAttachedToNavigator 在 OnResume（用于模态）和 OnNavigatedTo（用于页面）之前调用。 Page 和 Modal 界面的一部分
func (pageModal *GenericPageModal) OnAttachedToNavigator(parentNav PageNavigator) {
	pageModal.parentNav = parentNav
}

// ParentNavigator 是一个帮助方法，它返回将该内容推送到显示中的 Navigator，
//它可以是 WindowNavigator 或任何其他实现 PageNavigator 接口的页面（例如 MasterPage）。对于模态，这始终是 WindowNavigator.
func (pageModal *GenericPageModal) ParentNavigator() PageNavigator {
	return pageModal.parentNav
}

//ParentWindow 是一个帮助方法，
//如果它是 WindowNavigator，则返回显示此页面或模态框的 Navigator，否则它递归检查父导航器以查找并返回 WindowNavigator.
func (pageModal *GenericPageModal) ParentWindow() WindowNavigator {
	parentNav := pageModal.ParentNavigator()
	for {
		if parentNav == nil {
			return nil
		}
		if windowNav, isWindowNav := parentNav.(WindowNavigator); isWindowNav {
			return windowNav
		}
		if navigatedPageModal, ok := parentNav.(interface{ ParentNavigator() PageNavigator }); ok {
			parentNav = navigatedPageModal.ParentNavigator()
		}
	}
}
