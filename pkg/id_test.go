package pkg

import "testing"

func TestGenerateId(t *testing.T) {
	want := 4
	id := GenerateId(want)
	got := len(id)
	if got != want {
		t.Errorf("got len(id) = %d, wanted len(id) = %d ", got, want)
	} else {
		println("Generated id: ", id)
	}
}
