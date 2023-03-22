package main

import "testing"

func TestDummy(t *testing.T) {
	want := "aboba"

	t.Run("dummy", func(t *testing.T) {
		got := "aboba"
		if want != got {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
