package object

import (
	"fyne.io/fyne/v2/widget"
)

type LabelName string

type MapEntries map[LabelName]any

func (m MapEntries) GetEntryByLabelName(label string) *EntryWidget {
	return m[LabelName(label)].(*EntryWidget)
}

func (m MapEntries) GetCompletionByLabelName(label string) *CompletionContainer {
	return m[LabelName(label)].(*CompletionContainer)
}

func (m MapEntries) GetInputTextEntry() *widget.Entry {
	return m[LabelName("input")].(*widget.Entry)
}

func (m MapEntries) GetResultTextEntry() *widget.Entry {
	return m[LabelName("result")].(*widget.Entry)
}

func (m MapEntries) GetProtobufEntry() *EntryWidget {
	return m.GetEntryByLabelName("protobuf")
}
