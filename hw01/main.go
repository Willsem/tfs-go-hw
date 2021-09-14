package main

import (
	"flag"

	"github.com/willsem/tfs-go-hw/hw01/sandglass"
)

var char string
var size int
var color int

func init() {
	flag.StringVar(&char, "char", string(sandglass.DefaultChar), "symbol for output sandglass")
	flag.IntVar(&size, "size", sandglass.DefaultSize, "size of sandglass")
	flag.IntVar(&color, "color", sandglass.DefaultColor, "color of output")
}

func main() {
	flag.Parse()
	sandglass.Print(
		sandglass.WithChar([]rune(char)[0]),
		sandglass.WithSize(size),
		sandglass.WithColor(color),
	)
}
