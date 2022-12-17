package v

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

type (
	// Line represents a rectangle widget with an initial thickness of 1
	Line struct {
		Height     int
		Width      int
		Color      color.NRGBA
		isVertical bool
	}
)

// SeparatorVertical returns a vertical line widget instance
func SeparatorVertical(height, width int, col color.NRGBA) Line {
	vLine := NewLine(height, width, col)
	vLine.isVertical = true
	return vLine
}

// NewLine returns a line widget instance
func NewLine(height, width int, col color.NRGBA) Line {
	if height == 0 {
		height = 1
	}

	//col := values.Primary
	//col.A = 150
	return Line{
		Height: height,
		Width:  width,
		Color:  col,
	}
}

func NewSeparator(col color.NRGBA) Line {
	l := NewLine(1, 0, col)
	l.Color = col
	return l
}

// Layout renders the line widget
func (l Line) Layout(gtx C) D {
	if l.Width == 0 {
		l.Width = gtx.Constraints.Max.X
	}

	if l.isVertical && l.Height == 0 {
		l.Height = gtx.Constraints.Max.Y
	}

	line := image.Rectangle{
		Max: image.Point{
			X: l.Width,
			Y: l.Height,
		},
	}
	defer clip.Rect(line).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, l.Color)

	return layout.Dimensions{Size: line.Max}
}
