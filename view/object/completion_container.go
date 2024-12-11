package object

import (
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	widgetx "fyne.io/x/fyne/widget"
)

type CompletionContainer struct {
	completion      *widgetx.CompletionEntry
	label           *widget.Label
	objs            []fyne.CanvasObject
	container       *fyne.Container
	curPosX         float32
	options         []string
	mux             sync.Mutex
	maxOptionsFound int
}

func NewCompletion(labelName string, width float32, startPosition float32, objs ...fyne.CanvasObject) *CompletionContainer {
	label := widget.NewLabelWithStyle(labelName, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	label.Resize(fyne.NewSize(20, 30))

	completionEntry := widgetx.NewCompletionEntry([]string{})
	label.Move(fyne.NewPos(startPosition, 2))
	completionEntry.Move(fyne.NewPos(startPosition+label.Size().Width+10, 2))
	completionEntry.Resize(fyne.NewSize(width, 37))

	completionContainer := &CompletionContainer{
		completion:      completionEntry,
		label:           label,
		objs:            objs,
		maxOptionsFound: 10,
	}

	completionEntry.OnChanged = func(s string) {
		if len(s) < 2 {
			completionEntry.HideCompletion()
			return
		}

		completionContainer.mux.Lock()
		defer completionContainer.mux.Unlock()

		result := make([]string, 0)
		found := 0
		for _, opt := range completionContainer.options {
			optLower := strings.ToLower(opt)
			sLower := strings.ToLower(s)
			if strings.Contains(optLower, sLower) {
				result = append(result, opt)
				found++
			}

			if found >= 10 {
				break
			}
		}

		completionEntry.SetOptions(result)
		if len(result) > 0 {
			completionEntry.ShowCompletion()
		}
	}

	completionContainer.container = container.NewWithoutLayout(label, completionEntry)

	for i := range objs {
		completionContainer.AddCanvasObject(objs[i])
	}

	completionContainer.Resize(fyne.NewSize(completionContainer.curPosX+10, 40))

	return completionContainer
}

func (o *CompletionContainer) Resize(size fyne.Size) {
	o.container.Resize(size)
}

func (o *CompletionContainer) AddCanvasObject(obj fyne.CanvasObject) {
	if o.curPosX == 0 {
		o.curPosX = o.completion.Size().Width + o.completion.Position().X + 10
	}

	obj.Move(fyne.NewPos(o.curPosX, 2))
	o.objs = append(o.objs, obj)
	o.curPosX = o.curPosX + obj.Size().Width + 10
	o.container.Add(obj)
}

func (o *CompletionContainer) SetAllOptions(opts []string) {
	o.options = opts
}

func (o *CompletionContainer) Container() *fyne.Container {
	return o.container
}

func (o *CompletionContainer) Completion() *widgetx.CompletionEntry {
	return o.completion
}
