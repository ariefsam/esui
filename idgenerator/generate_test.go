package idgenerator

import "testing"

func TestGenerate(t *testing.T) {
	id := Generate()
	if len(id) == 0 {
		t.Error("Expected id to have a length greater than 0")
	}
}
