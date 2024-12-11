package object

import "fyne.io/fyne/v2"

type CanvasAdder interface {
	AddCanvasObject(obj fyne.CanvasObject)
}
