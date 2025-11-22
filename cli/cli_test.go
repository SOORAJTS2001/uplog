package main

import "testing"

func TestDummy(t *testing.T) {
	if 1 != 1 {
		t.Fatalf("Dummy Go Test Success")
	}
}
