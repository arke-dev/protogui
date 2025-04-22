package protocscreen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/arke-dev/protogui/service"
	"github.com/arke-dev/protogui/view/object"
)

type ProtoCScreen struct {
	protoC service.ProtoCompiler
	window fyne.Window
}

func New(protoC service.ProtoCompiler, window fyne.Window) *ProtoCScreen {
	return &ProtoCScreen{
		protoC: protoC,
		window: window,
	}
}

func (p *ProtoCScreen) String() string {
	return "protoc"
}

func (p *ProtoCScreen) NewCanvasObject() fyne.CanvasObject {
	return object.NewMainContainer(p.window, 100.00, true).
		AddProtobufEntry().
		AddCompletion("typename", 500).
		AddButtonIconSide("", theme.ViewRefreshIcon(), 37, p.actionLoad).
		AddButton("decode", 100, 2, p.actionDecodeButton).
		AddButtonSide("encode", 100, p.actionEncodeButton).
		AddButtonSide("template", 100, p.actionTemplateButton).
		VSplitOffset(0.2).
		Apply()
}

func (p *ProtoCScreen) actionDecodeButton(m object.MapEntries) func() {
	protobufEntry := m.GetProtobufEntry()
	typenameEntry := m.GetCompletionByLabelName("typename")
	inputTextEntry := m.GetInputTextEntry()
	resultTextEntry := m.GetResultTextEntry()
	return func() {
		res, err := p.protoC.Decode(protobufEntry.TextString(), typenameEntry.Completion().Text, inputTextEntry.Text)
		if err != nil {
			resultTextEntry.SetText(err.Error())
			return
		}
		resultTextEntry.SetText(res)
	}
}

func (p *ProtoCScreen) actionEncodeButton(m object.MapEntries) func() {
	protobufEntry := m.GetProtobufEntry()
	typenameEntry := m.GetCompletionByLabelName("typename")
	inputTextEntry := m.GetInputTextEntry()
	resultTextEntry := m.GetResultTextEntry()

	return func() {
		res, err := p.protoC.Encode(protobufEntry.TextString(), typenameEntry.Completion().Text, inputTextEntry.Text)
		if err != nil {
			resultTextEntry.SetText(err.Error())
			return
		}
		resultTextEntry.SetText(res)
	}
}

func (p *ProtoCScreen) actionTemplateButton(m object.MapEntries) func() {
	protobufEntry := m.GetProtobufEntry()
	typenameEntry := m.GetCompletionByLabelName("typename")
	inputTextEntry := m.GetInputTextEntry()

	return func() {
		template, err := p.protoC.TemplateJSON(protobufEntry.TextString(), typenameEntry.Completion().Text)
		if err == nil {
			inputTextEntry.SetText(template)
		}
	}
}

func (p *ProtoCScreen) actionLoad(m object.MapEntries) func() {
	protobufEntry := m.GetProtobufEntry()
	resultTextEntry := m.GetResultTextEntry()
	typenameCompletion := m.GetCompletionByLabelName("typename")

	return func() {
		types, err := p.protoC.GetRegisteredTypes(protobufEntry.TextString())
		if err != nil {
			resultTextEntry.SetText(err.Error())
			return
		}

		typenameCompletion.SetAllOptions(types)
	}
}
