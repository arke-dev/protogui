package object

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ProtoDialogButton struct {
	button *widget.Button
}

func NewProtoDialogButton(entry *EntryContainer, window fyne.Window) *ProtoDialogButton {
	button := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		fDial := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
			if lu != nil {
				entry.SetText(lu.Path())
			}
		}, window)
		fDial.Resize(fyne.NewSize(800, 800))
		fDial.Show()
	})
	button.Resize(fyne.NewSquareSize(37))

	return &ProtoDialogButton{
		button: button,
	}
}

func (b *ProtoDialogButton) Button() *widget.Button {
	return b.button
}
