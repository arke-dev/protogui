package object

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type EntryContainer struct {
	container *fyne.Container
	entry     *EntryWidget
	label     *widget.Label
	objects   []fyne.CanvasObject
	curPosX   float32
}

func NewEntry(labelName string, width float32, startPosition float32, typeObj TypeObject, objs ...fyne.CanvasObject) *EntryContainer {
	label := widget.NewLabelWithStyle(labelName, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	label.Resize(fyne.NewSize(20, 30))

	entryWidget := NewEntryWidget(typeObj)
	entryWidget.Resize(fyne.NewSize(width, 37))

	entryWidget.TypedKey(&fyne.KeyEvent{})

	label.Move(fyne.NewPos(startPosition, 2))
	entryWidget.Move(fyne.NewPos(startPosition+label.Size().Width+10, 2))

	entry := &EntryContainer{
		entry:   entryWidget,
		label:   label,
		objects: objs,
	}

	entry.container = container.NewWithoutLayout(label, entryWidget)

	for i := range objs {
		entry.AddCanvasObject(objs[i])
	}

	entry.Resize(fyne.NewSize(entry.curPosX+10, 40))

	return entry
}

func (o *EntryContainer) SetText(text string) {
	o.entry.SetText(text)
}

func (o *EntryContainer) Move(pos fyne.Position) {
	o.container.Move(pos)
}

func (o *EntryContainer) Resize(size fyne.Size) {
	o.container.Resize(size)
}

func (o *EntryContainer) AddCanvasObject(obj fyne.CanvasObject) {
	if o.curPosX == 0 {
		o.curPosX = o.entry.Size().Width + o.entry.Position().X + 10
	}

	obj.Move(fyne.NewPos(o.curPosX, 2))
	o.objects = append(o.objects, obj)
	o.curPosX = o.curPosX + obj.Size().Width + 10
	o.container.Add(obj)
}

func (o *EntryContainer) Container() *fyne.Container {
	return o.container
}

func (o *EntryContainer) EntryWidget() *EntryWidget {
	return o.entry
}

// func (o *Entry) TextInt() (int, error) {
// 	text := o.entry.Text
// 	return strconv.Atoi(text)
// }
