package main

import (
	"github.com/flinkcoin/mono/apps/broker/internal"
)

func main() {

	/*broker, err :=*/
	app := internal.Init()
	app.Harbor.Init()

	//	host.Init()
}
