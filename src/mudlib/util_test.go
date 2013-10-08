package mudlib

import (
	"reflect"
	"testing"
)

func TestRemoveStringFromList(t *testing.T) {
	cases := []struct {
		list      []string
		s         string
		wantList  []string
		wantError bool
	}{
		{list: []string{"a", "b", "c"}, s: "b",
			wantList: []string{"a", "c"}, wantError: false},
		{list: []string{"a", "c"}, s: "b",
			wantList: []string{"a", "c"}, wantError: true},
	}

	for _, tt := range cases {
		gotError := removeStringFromList(tt.s, &tt.list)
		if tt.wantError && gotError == nil {
			t.Errorf("want error, got no error\n")
		}
		if !tt.wantError && gotError != nil {
			t.Errorf("want no error, got error\n")
		}
		if !reflect.DeepEqual(tt.wantList, tt.list) {
			t.Errorf("want %v, got %v\n", tt.wantList, tt.list)
		}
	}
}
