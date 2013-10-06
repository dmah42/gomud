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

func addColor(fg, bg color, s string) string {
	return fmt.Sprintf("\033[1;%d;%dm%s\033[0m", 30+int(fg), 40+int(bg), s)
}
