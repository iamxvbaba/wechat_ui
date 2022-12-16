package v

import (
	"gioui.org/unit"
	"gioui.org/widget/material"
	"wechat_ui/ui/assets"
	"wechat_ui/ui/values"
)

type Label struct {
	material.LabelStyle
}

func H1(txt string) Label {
	return labelWithDefaultColor(Label{material.H1(assets.Theme, txt)})
}

func H2(txt string) Label {
	return labelWithDefaultColor(Label{material.H2(assets.Theme, txt)})
}

func H3(txt string) Label {
	return labelWithDefaultColor(Label{material.H2(assets.Theme, txt)})
}

func H4(txt string) Label {
	return labelWithDefaultColor(Label{material.H4(assets.Theme, txt)})
}

func H5(txt string) Label {
	return labelWithDefaultColor(Label{material.H5(assets.Theme, txt)})
}

func H6(txt string) Label {
	return labelWithDefaultColor(Label{material.H6(assets.Theme, txt)})
}

func Body1(txt string) Label {
	return labelWithDefaultColor(Label{material.Body1(assets.Theme, txt)})
}

func Body2(txt string) Label {
	return labelWithDefaultColor(Label{material.Body2(assets.Theme, txt)})
}

func Caption(txt string) Label {
	return labelWithDefaultColor(Label{material.Caption(assets.Theme, txt)})
}

func ErrorLabel(txt string) Label {
	label := Caption(txt)
	label.Color = values.Danger
	return label
}

func NewLabel(size unit.Sp, txt string) Label {
	return labelWithDefaultColor(Label{material.Label(assets.Theme, size, txt)})
}

func labelWithDefaultColor(l Label) Label {
	l.Color = values.DeepBlue
	return l
}
