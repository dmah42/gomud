package mudlib

import (
	"testing"
)

func TestIsValidNick(t *testing.T) {
	cases := []struct {
		nick      string
		wantValid bool
	}{
		{nick: "dma", wantValid: true},
		{nick: "d m", wantValid: false},
		{nick: "", wantValid: false},
	}

	for _, tt := range cases {
		gotValid := isValidNick(tt.nick)
		if tt.wantValid != gotValid {
			t.Errorf("%q want %v, got %v\n", tt.nick, tt.wantValid, gotValid)
		}
	}
}
