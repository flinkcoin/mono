package main

import (
	"github.com/flinkcoin/mono/apps/broker/internal"
	"testing"
)

func TestHello(t *testing.T) {
	app := internal.Init()
	if result.Greeter.Message != "Hello world" {
		t.Error("Expected Hello to append 'world'")
	}
}
