package object

import "fyne.io/fyne/v2"

func NewProtobufEntry(window fyne.Window) *EntryContainer {
	protobufEntry := NewEntry("protobuf", 700, 100, String)
	dialog := NewProtoDialogButton(protobufEntry, window)
	protobufEntry.AddCanvasObject(dialog.Button())
	return protobufEntry
}
