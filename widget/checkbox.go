package widget

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

func IsChecked(c widget.CheckBox) bool {
	return c.Checked(layout.NewContext(nil))
}
