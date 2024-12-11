package object

import (
	"bytes"
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

// JSONEntry is a widget entry for json input and json result.
// TODO: ExtendBaseWidget is not working when try to extend MultiLineEntry so should think in another solution in the future.
// type JSONEntry struct {
// 	*widget.Entry
// }

type JSONContainer struct {
	input     *widget.Entry
	result    *widget.Entry
	clip      *widget.Button
	container *fyne.Container
}

func NewJSONContainer(withInput bool) *JSONContainer {
	inputJSON := newInputJSON()
	resultJSON := newResultJSON()

	jsonContainer := &JSONContainer{
		input:  inputJSON,
		result: resultJSON,
	}

	clipContainer := jsonContainer.addClip(resultJSON)

	var containerJSON fyne.CanvasObject
	if withInput {
		containerJSON = container.NewHSplit(inputJSON, resultJSON)
	} else {
		containerJSON = container.NewStack(resultJSON)
	}

	jsonContainer.container = container.NewBorder(nil, clipContainer, nil, nil, containerJSON)

	return jsonContainer
}

func newInputJSON() *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapBreak
	entry.OnChanged = func(s string) {
		ok := json.Valid([]byte(s))
		if !ok {
			return
		}

		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, []byte(s), "", "\t")
		if err != nil {
			return
		}
		entry.SetText(prettyJSON.String())
	}
	return entry
}

func newResultJSON() *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapBreak
	entry.Disable()
	return entry
}

func (j *JSONContainer) addClip(result *widget.Entry) *fyne.Container {
	clip := widget.NewButtonWithIcon("copy", theme.ContentCopyIcon(), func() {
		_ = clipboard.Write(clipboard.FmtText, []byte(result.Text))
	})

	clip.Resize(fyne.NewSize(100, 50))
	clip.Move(fyne.NewPos(900, 2))
	j.clip = clip

	return container.NewWithoutLayout(clip)
}

func (i *JSONContainer) Input() *widget.Entry {
	return i.input
}

func (i *JSONContainer) Result() *widget.Entry {
	return i.result
}

func (i *JSONContainer) Clip() *widget.Button {
	return i.clip
}

func (i *JSONContainer) Container() fyne.CanvasObject {
	return i.container
}
