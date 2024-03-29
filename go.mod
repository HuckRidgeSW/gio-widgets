module github.com/huckridgesw/gio-widgets

go 1.21

replace gioui.org => ../../../git.sr.ht/eliasnaur/gioui.org

replace gioui.org/x => ../../../github.com/gioui/gio-x

require (
	github.com/go-text/typesetting v0.0.0-20230803102845-24e03d8b5372 // indirect
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)

require (
	gioui.org v0.0.0-00010101000000-000000000000
	golang.org/x/exp/shiny v0.0.0-20220827204233-334a2380cb91
)
