package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"golang.design/x/clipboard"

	"github.com/arke-dev/protogui/infra"
	"github.com/arke-dev/protogui/service"
	"github.com/arke-dev/protogui/view/consumerscreen"
	"github.com/arke-dev/protogui/view/grpcscreen"
	"github.com/arke-dev/protogui/view/protocscreen"

	"github.com/arke-dev/protogui/view/tabs"
	customTheme "github.com/arke-dev/protogui/view/theme"
)

func main() {
	os.Setenv("FYNE_SCALE", "0.8")
	_ = clipboard.Init()

	rabbitconn := infra.NewRabbitMQ("guest", "guest", "localhost", "", 5672)
	defer rabbitconn.Close()

	grpcConn := infra.NewGRPC()
	defer grpcConn.Close()

	protogui := app.NewWithID("protogui")
	protogui.Settings().SetTheme(customTheme.CustomTheme{})
	window := protogui.NewWindow("Fyne Desktop Client")
	protoC := service.NewProtoCompile()

	consumersvc := service.NewConsumer(infra.NewConsumerMQ(rabbitconn), protoC)
	grpcSvc := service.NewGRPC(protoC, grpcConn)

	protocScreen := protocscreen.New(protoC, window)

	consumerScreen := consumerscreen.NewConsumer(consumersvc, window)
	grpcScreen := grpcscreen.NewGRPC(grpcSvc, protoC, window)

	tabs := tabs.NewTabs(protocScreen, window, consumerScreen, grpcScreen)

	window.SetContent(tabs.MainContainer())
	window.Resize(fyne.NewSize(1200, 900))
	window.CenterOnScreen()
	window.ShowAndRun()
}
