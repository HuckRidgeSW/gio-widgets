// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	hrw "github.com/huckridgesw/gio-widgets/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
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
			Active    color.RGBA
			Inactive  color.RGBA
			Container color.RGBA
			Divider   color.RGBA
		}
		Font      text.Font
		IconType  IconType
		Buttons   layout.List
		CloseIcon *material.Icon
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
// func (th *material.Theme) Tabbar() *Tabbar {
func NewTabbar(th *material.Theme) *Tabbar {
	tb := Tabbar{
		th:         th,
		Scrollable: false,
		Font:       text.Font{
			// Size: th.TextSize.Scale(14.0 / 16.0),
		},
		CloseIcon: mustIcon(material.NewIcon(icons.NavigationClose)),
	}
	tb.Color.Active = th.Color.Text
	tb.Color.Inactive = rgb(0xffffff) // FIXME
	tb.Color.Container = th.Color.Primary
	tb.Color.Divider = rgb(0xffffff) // FIXME
	return &tb
}

func mustIcon(ic *material.Icon, err error) *material.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func (tb *Tabbar) Layout(gtx *layout.Context, wtb *hrw.Tabbar) {
	wtb.ProcessEvents(gtx)

	for i, tab := range wtb.Tabs {
		if tab.BecameActive {
			tab.BecameActive = false
			tb.Buttons.ScrollTo(i)
			break
		}
	}
	layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.
		Layout(gtx,
			layout.Rigid(func() {
				tb.Buttons.Layout(gtx, len(wtb.Tabs), func(i int) {
					tab := wtb.Tabs[i]
					// From https://material.io/components/tabs/#specs
					gtx.Constraints = layout.Constraints{
						Width:  layout.Constraint{Min: gtx.Px(unit.Dp(90)), Max: gtx.Px(unit.Dp(360))},
						Height: layout.Constraint{Min: gtx.Px(unit.Dp(48)), Max: gtx.Px(unit.Dp(48))},
					}
					var m op.MacroOp
					m.Record(gtx.Ops)
					layout.Inset{
						Top: unit.Dp(12), Bottom: unit.Dp(12),
						Left: unit.Dp(16), Right: unit.Dp(16),
					}.Layout(gtx, func() {
						// log.Printf("Inset cons: %+v", gtx.Constraints)
						// gtx.Constraints.Width.Min = 0
						// gtx.Constraints.Height.Min = 0

						// var m2 op.MacroOp
						// m2.Record(gtx.Ops)
						// tb.th.Body2(tab.GetLabel()).Layout(gtx)
						// m2.Stop()
						// lblDims := gtx.Dimensions
						// log.Printf("lblDims: %+v", lblDims)

						// if tab.Closeable {
						// 	layout.W.Layout(gtx, func() {
						// 		ib := tb.th.IconButton(tb.CloseIcon)
						// 		ib.Size = unit.Px(float32(lblDims.Size.Y - 10))
						// 		ib.Padding = unit.Dp(0)
						// 		ib.Layout(gtx, &tab.CloseButton)
						// 	})
						// }
						layout.Center.Layout(gtx, func() {
							// log.Printf("Center cons: %+v", gtx.Constraints)

							// Get the dimensions of the label
							var m2 op.MacroOp
							m2.Record(gtx.Ops)
							tb.th.Body2(tab.GetLabel()).Layout(gtx)
							m2.Stop()
							lblDims := gtx.Dimensions
							// log.Printf("lblDims: %+v", lblDims)

							layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func() {
									if tab.Closeable {
										ib := tb.th.IconButton(tb.CloseIcon)
										ib.Size = unit.Px(float32(lblDims.Size.Y - 2))
										ib.Padding = unit.Dp(0)
										ib.Layout(gtx, &tab.CloseButton)
									}
								}),
								layout.Rigid(func() {
									// Draw the label
									m2.Add()
									gtx.Dimensions = lblDims
								}))
						})
					})
					m.Stop()
					dims := gtx.Dimensions
					pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
					tab.LayoutButton(gtx)
					m.Add()

					// Underline the active item
					if tab == wtb.Active {
						paint.ColorOp{Color: color.RGBA{
							A: 0xff, B: 0xff,
						}}.Add(gtx.Ops)
						paint.PaintOp{
							Rect: f32.Rectangle{
								Min: f32.Point{Y: float32(dims.Size.Y - gtx.Px(unit.Dp(2)))},
								Max: f32.Point{X: float32(dims.Size.X), Y: float32(dims.Size.Y)},
							},
						}.Add(gtx.Ops)
					}
				})
			}),
			layout.Rigid(func() {
				wtb.Active.Layout(gtx)
			}),
		)
}
