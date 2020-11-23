package widget

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions

	Tabbar struct {
		Tabs      []*Tab
		byAddress map[interface{}]*Tab
		Active    *Tab
		events    []TabEvent
	}

	Tab struct {
		Label        string
		W            Layouter
		Closeable    bool
		CloseButton  widget.Clickable
		BecameActive bool
		button       widget.Clickable
	}

	TabEvent struct {
		Type TabEventType
		Tab  Layouter
	}

	TabEventType int

	Labeler interface {
		Label() string
	}

	Layouter interface {
		Layout(C) D
	}

	Activater interface {
		Activate()
	}

	Deactivater interface {
		Deactivate()
	}
)

const (
	TabEventClose TabEventType = iota
	TabEventActivate
	// MoveLeft
	// MoveRight
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

func (tb *Tabbar) Events(gtx layout.Context) []TabEvent {
	var e []TabEvent
	for _, tab := range tb.Tabs {
		// Don't have to check tab.Closeable, because if it's false, there won't
		// be a CloseButton to get an event.
		if tab.CloseButton.Clicked() {
			e = append(e, TabEvent{
				Type: TabEventClose,
				Tab:  tab.W,
			})
		}
		if tab.button.Clicked() {
			e = append(e, TabEvent{
				Type: TabEventActivate,
				Tab:  tab.W,
			})
		}
	}

	return e
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

// Close closes the indicated tab.  You cannot close the active tab.
func (tb *Tabbar) Close(index int) {
	if index >= len(tb.Tabs) || !tb.Tabs[index].Closeable {
		return
	}

	tab := tb.Tabs[index]

	if tab == tb.Active {
		return
	}

	copy(tb.Tabs[index:], tb.Tabs[index+1:])
	tb.Tabs = tb.Tabs[:len(tb.Tabs)-1]
	delete(tb.byAddress, tab.W)
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

func (tb *Tabbar) IndexOf(key interface{}) int {
	for i, tab := range tb.Tabs {
		if tab.W == key {
			return i
		}
	}
	return -1 // could definitely happen
}

func (tb *Tabbar) ActiveIndex() int {
	for i, tab := range tb.Tabs {
		if tab == tb.Active {
			return i
		}
	}
	return -1 // shouldn't happen
}

func NewTab(label string, w Layouter, closeable bool) *Tab {
	return &Tab{Label: label, W: w, Closeable: closeable}
}

func (t *Tab) LayoutButton(gtx C) D {
	return t.button.Layout(gtx)
}

func (t *Tab) Layout(gtx C) D {
	return t.W.Layout(gtx)
}

func (t *Tab) GetLabel() string {
	if labeler, ok := t.W.(Labeler); ok {
		return labeler.Label()
	}
	return t.Label
}
