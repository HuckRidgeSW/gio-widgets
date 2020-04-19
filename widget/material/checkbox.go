package material

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	hrwidget "github.com/huckridgesw/gio-widgets/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	checkBoxCheckedIcon   = mustIcon(material.NewIcon(icons.ToggleCheckBox))
	checkBoxUncheckedIcon = mustIcon(material.NewIcon(icons.ToggleCheckBoxOutlineBlank))
)

type CheckBox struct {
	checkable
}

func NewCheckBox(t *material.Theme, label string) CheckBox {
	return CheckBox{
		checkable{
			Label:              label,
			Color:              t.Color.Text,
			IconColor:          t.Color.Primary,
			TextSize:           t.TextSize.Scale(14.0 / 16.0),
			Size:               unit.Dp(26),
			shaper:             t.Shaper,
			checkedStateIcon:   checkBoxCheckedIcon,
			uncheckedStateIcon: checkBoxUncheckedIcon,
		},
	}
}

func (c CheckBox) Layout(gtx *layout.Context, checkBox *hrwidget.CheckBox) {
	c.layout(gtx, checkBox.Checked(gtx))
	checkBox.Layout(gtx)
}
