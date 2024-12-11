package object

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type FormContainer struct {
	containers    []fyne.CanvasObject
	mainContainer *fyne.Container
	curPosY       float32
}

func NewFormContainer(containers ...fyne.CanvasObject) *FormContainer {
	form := &FormContainer{mainContainer: container.NewWithoutLayout()}

	if len(containers) == 0 {
		return form
	}

	for i := range containers {
		form.AddContainer(containers[i])
	}

	return form
}

func (c *FormContainer) AddContainer(container fyne.CanvasObject) {
	if c.curPosY == 0 {
		c.curPosY = 2
	}

	container.Move(fyne.NewPos(2, c.curPosY))

	c.mainContainer.Add(container)
	c.containers = append(c.containers, container)
	c.curPosY = c.curPosY + container.Size().Height + 10
}

func (c *FormContainer) MainContainer() *fyne.Container {
	return c.mainContainer
}
