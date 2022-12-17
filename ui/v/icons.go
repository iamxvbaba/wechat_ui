package v

import "wechat_ui/ui/assets"

var (
	ContactIcon         = NewImage(assets.IconList["contact"])
	ContactIconInactive = NewImage(assets.IconList["contact_inactive"])
	MsgIcon             = NewImage(assets.IconList["msg"])
	MsgIconInactive     = NewImage(assets.IconList["msg_inactive"])

	SpIconInactive = NewImage(assets.IconList["sp_inactive"])
	PhoneInactive  = NewImage(assets.IconList["phone_inactive"])
	MoreInactive   = NewImage(assets.IconList["more_inactive"])
)
