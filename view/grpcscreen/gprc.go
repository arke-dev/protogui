package grpcscreen

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/lawmatsuyama/protogui/models"
	"github.com/lawmatsuyama/protogui/service"
	"github.com/lawmatsuyama/protogui/view/object"
)

type GRPCScreen struct {
	service service.GRPC
	protoC  service.ProtoCompiler
	window  fyne.Window
}

func NewGRPC(service service.GRPC, protoC service.ProtoCompiler, window fyne.Window) *GRPCScreen {
	return &GRPCScreen{
		service: service,
		window:  window,
		protoC:  protoC,
	}
}

func (p *GRPCScreen) NewCanvasObject() fyne.CanvasObject {
	return object.NewMainContainer(p.window, 100, true).
		AddProtobufEntry().
		AddCompletion("method", 500).
		AddButtonIconSide("", theme.ViewRefreshIcon(), 37, p.actionLoad).
		AddEntry("server address", 500, object.String).
		AddButton("call", 100, 2, p.actionCallButton).
		AddButtonSide("template", 100, p.actionTemplateButton).
		VSplitOffset(0.27).
		Apply()
}

func (p *GRPCScreen) actionCallButton(m object.MapEntries) func() {
	protobufEntry := m.GetProtobufEntry()
	methodsCompletion := m.GetCompletionByLabelName("method")
	addressEntry := m.GetEntryByLabelName("server address")
	inputTextEntry := m.GetInputTextEntry()
	resultTextEntry := m.GetResultTextEntry()

	return func() {
		res, err := p.service.Invoke(
			context.Background(),
			models.GRPCRequest{
				Path:           protobufEntry.TextString(),
				Method:         methodsCompletion.Completion().Text,
				RequestJsonMsg: inputTextEntry.Text,
				Address:        addressEntry.Text,
			},
		)

		if err != nil {
			resultTextEntry.SetText(err.Error())
			return
		}

		resultTextEntry.SetText(res)
	}
}

func (p *GRPCScreen) actionTemplateButton(m object.MapEntries) func() {
	protobufEntry := m.GetProtobufEntry()
	methodsCompletion := m.GetCompletionByLabelName("method")
	inputTextEntry := m.GetInputTextEntry()

	return func() {
		template, err := p.protoC.TemplateJSONFromMethod(protobufEntry.TextString(), methodsCompletion.Completion().Text)
		if err == nil {
			inputTextEntry.SetText(template)
		}
	}
}

func (p *GRPCScreen) actionLoad(m object.MapEntries) func() {
	protobufEntry := m.GetProtobufEntry()
	resultTextEntry := m.GetResultTextEntry()
	methodsCompletion := m.GetCompletionByLabelName("method")

	return func() {
		methods, err := p.protoC.GetRegisteredMethods(protobufEntry.TextString())
		if err != nil {
			resultTextEntry.SetText(err.Error())
			return
		}

		methodsCompletion.SetAllOptions(methods)
	}
}
