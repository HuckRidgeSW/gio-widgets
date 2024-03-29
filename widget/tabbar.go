package widget

import (
	"strings"
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions

	Tabbar struct {
		Tabs      []*Tab
		byAddress map[interface{}]*Tab
		active    *Tab
		events    []TabEvent
		mux       sync.RWMutex
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
	tb.mux.RLock()
	defer tb.mux.RUnlock()

	var e []TabEvent
	for _, tab := range tb.Tabs {
		// Don't have to check tab.Closeable, because if it's false, there won't
		// be a CloseButton to get an event.
		if tab.CloseButton.Clicked(gtx) {
			e = append(e, TabEvent{
				Type: TabEventClose,
				Tab:  tab.W,
			})
		}
		if tab.button.Clicked(gtx) {
			e = append(e, TabEvent{
				Type: TabEventActivate,
				Tab:  tab.W,
			})
		}
	}

	return e
}

func (tb *Tabbar) Prev() {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	for i, tab := range tb.Tabs {
		if tab == tb.active {
			if i == 0 {
				tb.activate(tb.Tabs[len(tb.Tabs)-1].W)
			} else {
				tb.activate(tb.Tabs[i-1].W)
			}
			return
		}
	}
}

func (tb *Tabbar) Next() {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	for i, tab := range tb.Tabs {
		if tab == tb.active {
			if i < len(tb.Tabs)-1 {
				tb.activate(tb.Tabs[i+1].W)
			} else {
				tb.activate(tb.Tabs[0].W)
			}
			return
		}
	}
}

// Close closes the indicated tab.  You cannot close the active tab.
func (tb *Tabbar) Close(index int) {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	if index >= len(tb.Tabs) || !tb.Tabs[index].Closeable {
		return
	}

	tab := tb.Tabs[index]

	if tab == tb.active {
		return
	}

	copy(tb.Tabs[index:], tb.Tabs[index+1:])
	tb.Tabs = tb.Tabs[:len(tb.Tabs)-1]
	delete(tb.byAddress, tab.W)
}

func (tb *Tabbar) Activate(key interface{}) {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	tb.activate(key)
}

// activate a tab. Caller should have tb.mux locked.
func (tb *Tabbar) activate(key interface{}) {
	if tab, ok := tb.byAddress[key]; ok {
		if tb.active != nil {
			if act, ok := tb.active.W.(Deactivater); ok {
				act.Deactivate()
			}
		}
		tb.active = tab
		tab.BecameActive = true
		if act, ok := key.(Activater); ok {
			act.Activate()
		}
	}
}

func (tb *Tabbar) Append(t *Tab) {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	tb.insert(len(tb.Tabs), t)
}

func (tb *Tabbar) InsertAfter(after, new *Tab) {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	for i, tab := range tb.Tabs {
		if tab == after {
			tb.insert(i+1, new)
			return
		}
	}
}

func (tb *Tabbar) Insert(index int, t *Tab) {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	tb.insert(index, t)
}

// Should have tb.mux Locked
func (tb *Tabbar) insert(index int, t *Tab) {
	if index > len(tb.Tabs) {
		index = len(tb.Tabs)
	}
	tb.Tabs = append(tb.Tabs, nil)
	copy(tb.Tabs[index+1:], tb.Tabs[index:])
	tb.Tabs[index] = t
	tb.byAddress[t.W] = t
}

func (tb *Tabbar) IndexOf(key interface{}) int {
	tb.mux.RLock()
	defer tb.mux.RUnlock()

	for i, tab := range tb.Tabs {
		if tab.W == key {
			return i
		}
	}
	return -1 // could definitely happen
}

func (tb *Tabbar) ActiveIndex() int {
	tb.mux.RLock()
	defer tb.mux.RUnlock()

	for i, tab := range tb.Tabs {
		if tab == tb.active {
			return i
		}
	}
	return -1 // shouldn't happen
}

func (tb *Tabbar) Active() *Tab {
	tb.mux.Lock()
	defer tb.mux.Unlock()

	return tb.active
}

func NewTab(label string, w Layouter, closeable bool) *Tab {
	return &Tab{Label: label, W: w, Closeable: closeable}
}

func (t *Tab) LayoutButton(gtx C, w layout.Widget) D {
	return t.button.Layout(gtx, w)
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

// BuildKeyset joins groups together with "|".
func BuildKeyset(groups ...string) string {
	return strings.Join(groups, "|")
}

// BuildKeygroup builds a group of keys, e.g. Short-[C,V,X].  The prefix is
// optional.  If there's only one key in the list, does not create a [] group.
//
// This function is handy because you can mention individual letters ("C",
// "V", "X"), which makes it easier to search for them later.
func BuildKeygroup(prefix string, keys ...string) string {
	var group string
	if len(keys) == 1 {
		group = keys[0]
	} else {
		group = "[" + strings.Join(keys, ",") + "]"
	}
	if prefix == "" {
		return group
	}
	return prefix + "-" + group
}
