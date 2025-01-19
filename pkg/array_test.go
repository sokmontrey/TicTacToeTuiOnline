package pkg

import "testing"

func TestArray_RemoveByIndex(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	got := SliceRemoveByIndex(slice, 2)
	want := []int{1, 2, 4, 5}
	if len(got) != len(want) {
		t.Errorf("got len(got) = %d, wanted len(got) = %d", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got got[%d] = %d, wanted got[%d] = %d", i, got[i], i, want[i])
		}
	}
}

func TestArray_RemoveByValue(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	got := SliceRemoveByValue(slice, 3)
	want := []int{1, 2, 4, 5}
	if len(got) != len(want) {
		t.Errorf("got len(got) = %d, wanted len(got) = %d", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got got[%d] = %d, wanted got[%d] = %d", i, got[i], i, want[i])
		}
	}
}

func TestArray_RemoveByValueNotFound(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	got := SliceRemoveByValue(slice, 6)
	want := []int{1, 2, 3, 4, 5}
	if len(got) != len(want) {
		t.Errorf("got len(got) = %d, wanted len(got) = %d", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got got[%d] = %d, wanted got[%d] = %d", i, got[i], i, want[i])
		}
	}
}
