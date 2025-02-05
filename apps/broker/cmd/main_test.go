package main

import (
	"github.com/flinkcoin/mono/apps/broker/app"
	"testing"
)

func TestHello(t *testing.T) {
	app := app.Init()
	if result.Greeter.Message != "Hello world" {
		t.Error("Expected Hello to append 'world'")
	}
}
