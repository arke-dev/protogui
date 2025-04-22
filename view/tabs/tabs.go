package tabs

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"
	"github.com/arke-dev/protogui/view/consumerscreen"
	"github.com/arke-dev/protogui/view/grpcscreen"

	"github.com/arke-dev/protogui/view/protocscreen"
)

type Screen interface {
	NewCanvasObject() fyne.CanvasObject
}
type Tabs struct {
	mux            sync.Mutex
	tabs           *container.AppTabs
	protocScreen   *protocscreen.ProtoCScreen
	window         fyne.Window
	closeButton    *widget.Button
	tabsContainer  *fyne.Container
	mainContainer  *fyne.Container
	consumerScreen *consumerscreen.ConsumerScreen
	grpcScreen     *grpcscreen.GRPCScreen
	buttonPosX     float32
}

func NewTabs(
	protocScreen *protocscreen.ProtoCScreen,
	window fyne.Window,
	consumerScreen *consumerscreen.ConsumerScreen,
	grpcsScreen *grpcscreen.GRPCScreen) *Tabs {
	t := &Tabs{
		tabs:           container.NewAppTabs(),
		protocScreen:   protocScreen,
		window:         window,
		consumerScreen: consumerScreen,
		grpcScreen:     grpcsScreen,
		buttonPosX:     2,
	}

	t.setup()
	return t
}

func (t *Tabs) setup() {
	close := widget.NewButton("close tab", func() { t.removeTab() })
	close.Resize(fyne.NewSize(100, 30))
	close.Move(fyne.NewPos(1050, 20))
	t.tabs.Move(fyne.NewPos(2, 20))
	t.tabs.Resize(fyne.NewSize(1000, 800))
	t.closeButton = close
	t.tabsContainer = container.NewWithoutLayout(t.closeButton, t.tabs)

	protocButton := t.NewButton("protoc", t.protocScreen)
	consumerButton := t.NewButton("consumer", t.consumerScreen)
	grpcButton := t.NewButton("grpc", t.grpcScreen)
	menu := container.NewWithoutLayout(protocButton, consumerButton, grpcButton)
	t.mainContainer = container.NewVBox(menu, t.tabsContainer)
}

func (t *Tabs) NewButton(label string, screen Screen) *widget.Button {
	button := widget.NewButton(label, func() {
		tab := container.NewTabItem(label, screen.NewCanvasObject())
		tab.Content.Resize(fyne.NewSize(20, 90))
		t.tabs.Append(tab)
		t.tabs.Select(tab)
	})

	button.Resize(fyne.NewSize(100, 50))
	button.Move(fyne.NewPos(t.buttonPosX, 2))
	t.buttonPosX = t.buttonPosX + button.Size().Width + 10
	return button
}

func (t *Tabs) MainContainer() *fyne.Container {
	return t.mainContainer
}

func (t *Tabs) removeTab() {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.tabs.RemoveIndex(t.tabs.SelectedIndex())
}
