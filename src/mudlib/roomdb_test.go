package mudlib

import (
	"testing"
)

func TestRoomDbGet(t *testing.T) {
	LoadRoomDb("testdata/rooms/", "start")
	r, err := rooms.get("start")
	if err != nil {
		t.Errorf("%+v", err)
	}
	r2, err := rooms.get("start")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if r != r2 {
		t.Errorf("Mismatch pointers %p vs %p", r, r2)
	}
}
