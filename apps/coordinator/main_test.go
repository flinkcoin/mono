package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	result, _ := InitializeEvent("Hello world")
	if result.Greeter.Message != "Hello world" {
		t.Error("Expected Hello to append 'world'")
	}
}
