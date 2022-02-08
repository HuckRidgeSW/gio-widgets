module github.com/huckridgesw/gio-widgets

go 1.17

require (
	gioui.org v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20210722180016-6781d3edade3
)

require golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect

replace gioui.org => ../../../git.sr.ht/eliasnaur/gioui.org
