package core

import (
	"testing"
)

func TestCore(t *testing.T) {
	result := Core("works")
	if result != "Core works" {
		t.Error("Expected Core to append 'works'")
	}
}
