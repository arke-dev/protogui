package object

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TypeObject string

const (
	String TypeObject = "string"
	Int    TypeObject = "int"
)

type EntryWidget struct {
	tp TypeObject
	widget.Entry
}

func NewEntryWidget(tp TypeObject) *EntryWidget {
	entry := &EntryWidget{tp: tp}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *EntryWidget) TypedRune(r rune) {
	if e.tp == String {
		e.Entry.TypedRune(r)
		return
	}

	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
	}
}

func (e *EntryWidget) TypedShortcut(shortcut fyne.Shortcut) {
	if e.tp == String {
		e.Entry.TypedShortcut(shortcut)
	}

	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.Atoi(content); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *EntryWidget) TextInt() int {
	if e.tp != Int {
		return 0
	}

	n, _ := strconv.Atoi(e.Text)
	return n
}

func (e *EntryWidget) TextString() string {
	return e.Text
}
