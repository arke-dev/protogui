package object

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Button struct {
	container *fyne.Container
	button    *widget.Button
	curPosX   float32
	objs      []fyne.CanvasObject
}

func NewButton(labelName string, width float32, startPosition float32, action func(), objs ...fyne.CanvasObject) *Button {
	buttonWidget := widget.NewButton(labelName, action)
	buttonWidget.Move(fyne.NewPos(startPosition, 2))
	buttonWidget.Resize(fyne.NewSize(width, 40))

	button := &Button{
		button: buttonWidget,
		objs:   objs,
	}

	button.container = container.NewWithoutLayout(buttonWidget)

	for i := range objs {
		button.AddCanvasObject(objs[i])
	}

	button.Resize(fyne.NewSize(button.curPosX+10, 40))

	return button
}

func NewButtonIcon(labelName string, icon fyne.Resource, size float32, startPosition float32, action func(), objs ...fyne.CanvasObject) *Button {
	buttonWidget := widget.NewButtonWithIcon(labelName, icon, action)
	buttonWidget.Move(fyne.NewPos(startPosition, 2))
	buttonWidget.Resize(fyne.NewSquareSize(size))

	button := &Button{
		button: buttonWidget,
		objs:   objs,
	}

	button.container = container.NewWithoutLayout(buttonWidget)

	for i := range objs {
		button.AddCanvasObject(objs[i])
	}

	button.Resize(fyne.NewSize(button.curPosX+10, 40))

	return button
}

func (o *Button) Move(pos fyne.Position) {
	o.container.Move(pos)
}

func (o *Button) Resize(size fyne.Size) {
	o.container.Resize(size)
}

func (o *Button) AddCanvasObject(obj fyne.CanvasObject) {
	if o.curPosX == 0 {
		o.curPosX = o.button.Size().Width + o.button.Position().X + 10
	}

	obj.Move(fyne.NewPos(o.curPosX, 2))
	o.objs = append(o.objs, obj)
	o.curPosX = o.curPosX + obj.Size().Width + 10
	o.container.Add(obj)
}

func (o *Button) Container() *fyne.Container {
	return o.container
}

func (o *Button) Button() *widget.Button {
	return o.button
}
