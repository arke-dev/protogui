package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"golang.design/x/clipboard"

	"github.com/lawmatsuyama/protogui/infra"
	"github.com/lawmatsuyama/protogui/service"
	"github.com/lawmatsuyama/protogui/view/consumerscreen"
	"github.com/lawmatsuyama/protogui/view/grpcscreen"
	"github.com/lawmatsuyama/protogui/view/protocscreen"

	"github.com/lawmatsuyama/protogui/view/tabs"
	customTheme "github.com/lawmatsuyama/protogui/view/theme"
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
