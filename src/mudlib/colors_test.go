package mudlib

import (
  "testing"
)

func TestColorsFg(t *testing.T) {
  cases := []struct {
    c color
    bold bool
    want string
  }{
    {
      c: colorGreen,
      bold: false,
      want: "\033[0;32mtest\033[0m",
    },
    {
      c: colorRed,
      bold: true,
      want: "\033[1;31mtest\033[0m",
    },
  }

  for _, tt := range cases {
    got := ""
    if tt.bold {
      got = setFgBold(tt.c, "test")
    } else {
      got = setFg(tt.c, "test")
    }
    if got != tt.want {
      t.Errorf("got %q, want %q", got, tt.want)
    }
  }
}
