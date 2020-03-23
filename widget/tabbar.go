package widget

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type (
	Tabbar struct {
		Tabs      []*Tab
		byAddress map[interface{}]*Tab
		Active    *Tab
	}

	Tab struct {
		Label        string
		W            Layouter
		Closeable    bool
		CloseButton  widget.Button
		BecameActive bool
		button       widget.Button
	}

	Labeler interface {
		Label() string
	}

	Layouter interface {
		Layout(*layout.Context)
	}

	Activater interface {
		Activate()
	}

	Deactivater interface {
		Deactivate()
	}

	Closer interface {
		Close()
	}
)

func NewTabbar(tabs ...*Tab) *Tabbar {
	tb := Tabbar{
		Tabs:      tabs,
		byAddress: map[interface{}]*Tab{},
	}
	for _, tab := range tabs {
		tb.byAddress[tab.W] = tab
	}
	return &tb
}

func (tb *Tabbar) ProcessEvents(gtx *layout.Context) {
	for i, tab := range tb.Tabs {
		if tab.CloseButton.Clicked(gtx) {
			tb.Close(i)
		}
		if tab.button.Clicked(gtx) {
			tb.Activate(tab.W)
		}
	}
}

func (tb *Tabbar) Prev() {
	for i, tab := range tb.Tabs {
		if tab == tb.Active {
			if i == 0 {
				tb.Activate(tb.Tabs[len(tb.Tabs)-1].W)
			} else {
				tb.Activate(tb.Tabs[i-1].W)
			}
			return
		}
	}
}

func (tb *Tabbar) Next() {
	for i, tab := range tb.Tabs {
		if tab == tb.Active {
			if i < len(tb.Tabs)-1 {
				tb.Activate(tb.Tabs[i+1].W)
			} else {
				tb.Activate(tb.Tabs[0].W)
			}
			return
		}
	}
}

// Close closes the indicated tab.  If that tab is active, activates the one
// to its right, or if there isn't one, the (new) last tab.
func (tb *Tabbar) Close(index int) {
	if index >= len(tb.Tabs) || !tb.Tabs[index].Closeable {
		return
	}
	tab := tb.Tabs[index]
	copy(tb.Tabs[index:], tb.Tabs[index+1:])
	tb.Tabs = tb.Tabs[:len(tb.Tabs)-1]
	delete(tb.byAddress, tab.W)

	// If this is the active tab, activate the tab to the right, or if there
	// isn't one, the new last tab.
	if tab == tb.Active {
		if index >= len(tb.Tabs) {
			index = len(tb.Tabs) - 1
		}
		tb.Activate(tb.Tabs[index].W)
	}

	// Send close message to the closed tab.  (Deactivate is sent by the above
	// tb.Activate() call.)
	if closer, ok := tab.W.(Closer); ok {
		closer.Close()
	}
}

func (tb *Tabbar) Activate(key interface{}) {
	if tab, ok := tb.byAddress[key]; ok {
		if tb.Active != nil {
			if act, ok := tb.Active.W.(Deactivater); ok {
				act.Deactivate()
			}
		}
		tb.Active = tab
		tab.BecameActive = true
		if act, ok := key.(Activater); ok {
			act.Activate()
		}
	}
}

func (tb *Tabbar) Append(t *Tab) {
	tb.Insert(len(tb.Tabs), t)
}

func (tb *Tabbar) InsertAfter(after, new *Tab) {
	for i, tab := range tb.Tabs {
		if tab == after {
			tb.Insert(i+1, new)
			return
		}
	}
}

func (tb *Tabbar) Insert(index int, t *Tab) {
	if index > len(tb.Tabs) {
		index = len(tb.Tabs)
	}
	tb.Tabs = append(tb.Tabs, nil)
	copy(tb.Tabs[index+1:], tb.Tabs[index:])
	tb.Tabs[index] = t
	tb.byAddress[t.W] = t
}

func NewTab(label string, w Layouter, closeable bool) *Tab {
	return &Tab{Label: label, W: w, Closeable: closeable}
}

func (t *Tab) LayoutButton(gtx *layout.Context) {
	t.button.Layout(gtx)
}

func (t *Tab) Layout(gtx *layout.Context) {
	t.W.Layout(gtx)
}

func (t *Tab) GetLabel() string {
	if labeler, ok := t.W.(Labeler); ok {
		return labeler.Label()
	}
	return t.Label
}
