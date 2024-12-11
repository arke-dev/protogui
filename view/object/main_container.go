package object

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainContainer struct {
	mainContainer fyne.CanvasObject
	jsonContainer *JSONContainer
	mapEntries    MapEntries
	formContainer *FormContainer
	positionY     float32
	window        fyne.Window
	lastObject    CanvasAdder
	vsplitOffset  float64
}

func NewMainContainer(window fyne.Window, positionY float32, withInputJSON bool) *MainContainer {
	if positionY <= 0 {
		positionY = 2
	}
	jsonContainer := NewJSONContainer(withInputJSON)

	mapEntries := make(MapEntries)
	mapEntries[LabelName("input")] = jsonContainer.Input()
	mapEntries[LabelName("result")] = jsonContainer.Result()

	return &MainContainer{
		jsonContainer: jsonContainer,
		positionY:     positionY,
		window:        window,
		mapEntries:    mapEntries,
		formContainer: NewFormContainer(),
	}
}

func (c *MainContainer) SetPosY(y float32) *MainContainer {
	c.positionY = y
	return c
}

func (c *MainContainer) AddContainer(container fyne.CanvasObject) *MainContainer {
	c.formContainer.AddContainer(container)
	return c
}

func (c *MainContainer) AddEntry(labelName LabelName, width float32, typeObj TypeObject, objs ...fyne.CanvasObject) *MainContainer {
	entry := NewEntry(string(labelName), width, c.positionY, typeObj, objs...)
	c.formContainer.AddContainer(entry.Container())
	c.mapEntries[labelName] = entry.EntryWidget()
	c.lastObject = entry
	return c
}

func (c *MainContainer) AddCompletion(labelName LabelName, width float32, objs ...fyne.CanvasObject) *MainContainer {
	entry := NewCompletion(string(labelName), width, c.positionY, objs...)
	c.formContainer.AddContainer(entry.Container())
	c.mapEntries[labelName] = entry
	c.lastObject = entry
	return c
}

func (c *MainContainer) AddButton(labelName LabelName, width float32, position float32, action func(MapEntries) func(), objs ...fyne.CanvasObject) *MainContainer {
	button := NewButton(string(labelName), width, position, action(c.mapEntries), objs...)
	c.lastObject = button
	c.formContainer.AddContainer(button.Container())
	return c
}

func (c *MainContainer) AddButtonSide(labelName LabelName, width float32, action func(MapEntries) func(), objs ...fyne.CanvasObject) *MainContainer {
	button := NewButton(string(labelName), width, 2, action(c.mapEntries), objs...)
	c.lastObject.AddCanvasObject(button.Button())
	return c
}

func (c *MainContainer) AddButtonIcon(labelName LabelName, icon fyne.Resource, size float32, position float32, action func(MapEntries) func(), objs ...fyne.CanvasObject) *MainContainer {
	button := NewButtonIcon(string(labelName), icon, size, position, action(c.mapEntries), objs...)
	c.lastObject = button
	c.formContainer.AddContainer(button.Container())
	return c
}

func (c *MainContainer) AddButtonIconSide(labelName LabelName, icon fyne.Resource, size float32, action func(MapEntries) func(), objs ...fyne.CanvasObject) *MainContainer {
	button := NewButtonIcon(string(labelName), icon, size, 2, action(c.mapEntries), objs...)
	c.lastObject.AddCanvasObject(button.Button())
	return c
}

func (c *MainContainer) AddProtobufEntry() *MainContainer {
	entry := NewProtobufEntry(c.window)
	c.formContainer.AddContainer(entry.Container())
	c.mapEntries[LabelName("protobuf")] = entry.EntryWidget()
	return c
}

func (c *MainContainer) VSplitOffset(offset float64) *MainContainer {
	c.vsplitOffset = offset
	return c
}

func (c *MainContainer) Apply() fyne.CanvasObject {
	mainContainer := container.NewVSplit(c.formContainer.MainContainer(), c.jsonContainer.Container())
	mainContainer.SetOffset(c.vsplitOffset)
	c.mainContainer = mainContainer
	return mainContainer
}
