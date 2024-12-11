package object

import "fyne.io/fyne/v2"

type TextEntry interface {
	TextInt() int
	TypedShortcut(shortcut fyne.Shortcut)
	TypedRune(r rune)
	SetText(text string)
	TextString() string
}
