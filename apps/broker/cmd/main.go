package main

import (
	"github.com/flinkcoin/mono/apps/broker/app"
	"github.com/flinkcoin/mono/libs/core/pkg/core"
)

func main() {

	core.Log.Info("Starting broker service!")

	/*broker, err :=*/
	app := app.Init()
	app.Harbor.Init()

	//	host.Init()
}
