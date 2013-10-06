package mudlib

import (
	"fmt"
)

type color int

const (
	colorBlack color = iota
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

func (c color) toFg() int { return 30 + int(c) }

// TODO: background colors
// func (c color) toBg() int { return 40 + int(c) }

func setFg(c color, s string) string {
	return fmt.Sprintf("\033[0;%dm%s\033[0m", c.toFg(), s)
}

func setFgBold(c color, s string) string {
	return fmt.Sprintf("\033[1;%dm%s\033[0m", c.toFg(), s)
}
