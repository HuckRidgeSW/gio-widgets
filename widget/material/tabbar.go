package material

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	hrw "github.com/huckridgesw/gio-widgets/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type (
	// See https://material.io/components/tabs
	Tabbar struct {
		th *material.Theme

		// Currently ignored: All tabbars are scrollable.
		//
		// Mostly because I don't know how to make them not scrollable.
		Scrollable bool

		ClusteredFixed bool
		Alignment      layout.Alignment // only for ClusteredFixed
		Color          struct {
			// Active: use Theme.Color.Text
			// Container: use Theme.Color.Primary
			Active    color.NRGBA
			Inactive  color.NRGBA
			Container color.NRGBA
			Divider   color.NRGBA
		}
		Font      text.Font
		IconType  IconType
		Buttons   layout.List
		CloseIcon *widget.Icon
	}

	IconType int
)

const (
	IconNone IconType = iota
	IconLeader
	IconTop
)

// Tabbar creates a new tab bar.  You should store it and reuse it for each
// layout, because the embedded layout.List has some state you need to keep
// around.
func NewTabbar(th *material.Theme) *Tabbar {
	tb := Tabbar{
		th:        th,
		CloseIcon: mustIcon(widget.NewIcon(icons.NavigationClose)),
	}
	tb.Color.Active = th.Palette.Fg
	tb.Color.Inactive = nrgb(0xffffff) // FIXME
	tb.Color.Container = th.Palette.Fg
	tb.Color.Divider = nrgb(0xffffff) // FIXME
	return &tb
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

// nrgb cribbed from rgb in
// git.sr.ht/eliasnaur/gioui.org/widget/material/theme.go
func nrgb(c uint32) color.NRGBA {
	return nargb(0xff000000 | c)
}

// nargb cribbed from argb in
// git.sr.ht/eliasnaur/gioui.org/widget/material/theme.go.
func nargb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func (tb *Tabbar) Layout(gtx C, wtb *hrw.Tabbar) D {
	for i, tab := range wtb.Tabs {
		if tab.BecameActive {
			tab.BecameActive = false
			tb.Buttons.ScrollTo(i)
			break
		}
	}
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			// log.Printf("tabbar buttons")
			return tb.Buttons.Layout(gtx, len(wtb.Tabs), func(gtx C, i int) D {
				tab := wtb.Tabs[i]
				// From https://material.io/components/tabs/#specs
				gtx.Constraints = layout.Constraints{
					Min: image.Point{X: gtx.Px(unit.Dp(90)), Y: gtx.Px(unit.Dp(48))},
					Max: image.Point{X: gtx.Px(unit.Dp(360)), Y: gtx.Px(unit.Dp(48))},
				}
				buttonMacro := op.Record(gtx.Ops)
				dims := layout.Inset{
					Top: unit.Dp(12), Bottom: unit.Dp(12),
					Left: unit.Dp(16), Right: unit.Dp(16),
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						// log.Printf("Center cons: %+v", gtx.Constraints)

						// Get the dimensions of the label
						labelMacro := op.Record(gtx.Ops)
						lblDims := material.Body2(tb.th, tab.GetLabel()).Layout(gtx)
						labelCall := labelMacro.Stop()
						// log.Printf("lblDims: %+v", lblDims)

						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if !tab.Closeable {
									return D{}
								}
								ib := material.IconButton(tb.th, &tab.CloseButton, tb.CloseIcon, "")
								ib.Size = unit.Px(float32(lblDims.Size.Y - 2))
								ib.Inset = layout.Inset{}
								return ib.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								// Draw the label
								labelCall.Add(gtx.Ops)
								return lblDims
							}))
					})
				})
				buttonCall := buttonMacro.Stop()
				gtx.Constraints = layout.Exact(dims.Size)
				tab.LayoutButton(gtx, func(gtx C) D {
					buttonCall.Add(gtx.Ops)
					return dims
				})

				// Underline the active item
				if tab == wtb.Active {
					paint.FillShape(gtx.Ops,
						color.NRGBA{A: 0xff, B: 0xff},
						clip.Rect{
							Min: image.Point{X: 0, Y: dims.Size.Y - gtx.Px(unit.Dp(2))},
							Max: image.Point{X: dims.Size.X, Y: dims.Size.Y},
						}.Op(),
					)
				}

				return dims
			})
		}),
		layout.Rigid(func(gtx C) D {
			return wtb.Active.Layout(gtx)
		}),
	)
}
