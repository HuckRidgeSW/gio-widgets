package material

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
		th       *material.Theme
		eventKey int

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
		Font      font.Font
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

var tabbarKeyset = hrw.BuildKeyset(
	hrw.BuildKeygroup("Short-Shift", "{", "}", "[", "]"),
	hrw.BuildKeygroup("Short", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
)

func (tb *Tabbar) Layout(gtx C, wtb *hrw.Tabbar) D {
	tb.processEvents(gtx, wtb)

	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops).Pop()
	key.InputOp{
		Tag:  &tb.eventKey,
		Keys: key.Set(tabbarKeyset),
	}.Add(gtx.Ops)

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
					Min: image.Point{X: gtx.Dp(90), Y: gtx.Dp(48)},
					Max: image.Point{X: gtx.Dp(360), Y: gtx.Dp(48)},
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
								// log.Printf("tab %d: size: %v, baseline: %v, dp size - 2: %v, dp size / 2: %v", i,
								// 	lblDims.Size.Y,
								// 	lblDims.Baseline,
								// 	unit.Dp(lblDims.Size.Y-2),
								// 	unit.Dp(lblDims.Size.Y/2))
								// Not sure this is right.
								ib.Size = unit.Dp(lblDims.Size.Y / 2)
								// ib.Size = unit.Dp(lblDims.Size.Y - lblDims.Baseline)
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
				if tab == wtb.Active() {
					paint.FillShape(gtx.Ops,
						color.NRGBA{A: 0xff, B: 0xff},
						clip.Rect{
							Min: image.Point{X: 0, Y: dims.Size.Y - gtx.Dp(2)},
							Max: image.Point{X: dims.Size.X, Y: dims.Size.Y},
						}.Op(),
					)
				}

				return dims
			})
		}),
		layout.Rigid(func(gtx C) D {
			return wtb.Active().Layout(gtx)
		}),
	)
}

func (tb *Tabbar) processEvents(gtx C, wtb *hrw.Tabbar) {
	for _, e := range gtx.Events(&tb.eventKey) {
		switch ke := e.(type) {
		case key.Event:
			if ke.State != key.Press {
				break
			}
			switch ke.Name {
			// On macOS, cmd-shift-[ aka cmd-{ is reported as "{"; on Windows
			// it's "[".  Both report Shortcut+Shift, though.  (And similarly
			// below for }/].
			case "{", "[":
				// tl.log.Printf("toplevel process key: { %v", ke.Modifiers)
				if ke.Modifiers == key.ModShortcut|key.ModShift {
					// e.Consume()
					wtb.Prev()
				}
			case "}", "]":
				// tl.log.Printf("toplevel process key: } %v", ke.Modifiers)
				if ke.Modifiers == key.ModShortcut|key.ModShift {
					// e.Consume()
					wtb.Next()
				}
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				if ke.Modifiers == key.ModShortcut {
					// e.Consume()
					tabNum := int(ke.Name[0] - '1')
					// Ctrl/Cmd-9 always goes to the last tab.
					if ke.Name == "9" {
						wtb.Activate(wtb.Tabs[len(wtb.Tabs)-1].W)
					} else if tabNum < len(wtb.Tabs) {
						wtb.Activate(wtb.Tabs[tabNum].W)
					}
				}
			default:
				// tl.log.Printf("Toplevel.processEvents: Unknown key: %+v", ke)
			}
		default:
			// tl.log.Printf("Toplevel.processEvents: Unknown event type: %T %+v", ke, ke)
		}
	}
}
