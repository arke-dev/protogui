package consumerscreen

import (
	"context"

	"github.com/arke-dev/protogui/models"
	"github.com/arke-dev/protogui/service"
	"github.com/arke-dev/protogui/view/object"

	"fyne.io/fyne/v2"
)

type ConsumerScreen struct {
	service service.Consumer
	window  fyne.Window
}

func NewConsumer(service service.Consumer, window fyne.Window) *ConsumerScreen {
	return &ConsumerScreen{service: service, window: window}
}

func (p *ConsumerScreen) NewCanvasObject() fyne.CanvasObject {
	return object.NewMainContainer(p.window, 100.00, false).
		AddProtobufEntry().
		AddEntry("queue", 400, object.String).
		AddEntry("quantity", 100, object.Int).
		AddButton("get", 100, 2, p.actionGetButton).
		VSplitOffset(0.30).
		Apply()
}

func (p *ConsumerScreen) actionGetButton(mapEntries object.MapEntries) func() {
	protobufEntry := mapEntries.GetProtobufEntry()
	queueEntry := mapEntries.GetEntryByLabelName("queue")
	quantityEntry := mapEntries.GetEntryByLabelName("quantity")
	resultTextEntry := mapEntries.GetResultTextEntry()

	return func() {
		m, err := p.service.GetMessages(context.Background(), &models.GetMessagesRequest{
			Path:     protobufEntry.TextString(),
			Queue:    queueEntry.TextString(),
			Quantity: quantityEntry.TextInt(),
			Mode:     models.Nack,
		})
		if err != nil {
			resultTextEntry.SetText(err.Error())
			return
		}

		resultTextEntry.SetText(m)
	}
}
