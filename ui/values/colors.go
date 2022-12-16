package values

import "image/color"

var (
	Primary          = rgb(0x2970ff)
	Primary50        = rgb(0xE3F2FF)
	PrimaryHighlight = rgb(0x1B41B3)

	PageNavText = rgb(0x091440)
	Text        = rgb(0x091440)
	InvText     = rgb(0xffffff) // inverted default color #ffffff
	GrayText1   = rgb(0x3d5873) // darker shade #3D5873
	GrayText2   = rgb(0x596D81) // lighter shade of GrayText1 #596D81
	GrayText3   = rgb(0x8997a5) // lighter shade of GrayText2 #8997A5 (hint)
	GrayText4   = rgb(0xc4cbd2) // lighter shade of GrayText3 ##C4CBD2
	GreenText   = rgb(0x41BE53) // green text #41BE53

	// background colors
	Background       = argb(0x22444444)
	Black            = rgb(0x000000)
	BlueProgressTint = rgb(0x73d7ff)
	Danger           = rgb(0xed6d47)
	DeepBlue         = rgb(0x091440)
	NavyBlue         = rgb(0x1F45B0)
	LightBlue        = rgb(0xe4f6ff)
	LightBlue2       = rgb(0x75D8FF)
	LightBlue3       = rgb(0xBCE8FF)
	LightBlue4       = rgb(0xBBDEFF)
	LightBlue5       = rgb(0x70CBFF)
	LightBlue6       = rgb(0x4B91D8)
	Gray1            = rgb(0x3d5873)
	Gray2            = rgb(0xe6eaed)
	Gray3            = rgb(0xc4cbd2)
	Gray4            = rgb(0xf3f5f6)
	Gray5            = rgb(0xf5f5f5)
	Green50          = rgb(0xE8F7EA)
	Green500         = rgb(0x41BE53)
	Orange           = rgb(0xD34A21)
	Orange2          = rgb(0xF8E8E7)
	Orange3          = rgb(0xF8CABC)
	OrangeRipple     = rgb(0xD32F2F)
	Success          = rgb(0x41bf53)
	Success2         = rgb(0xE1F8EF)
	Surface          = rgb(0xffffff)
	Turquoise100     = rgb(0xB6EED7)
	Turquoise300     = rgb(0x2DD8A3)
	Turquoise700     = rgb(0x00A05F)
	Turquoise800     = rgb(0x008F52)
	Yellow           = rgb(0xffc84e)
	White            = rgb(0xffffff)

	SurfaceHighlight color.NRGBA

	DarkGray = rgb(0xFF2E2E2E)
)

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
