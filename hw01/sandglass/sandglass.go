package sandglass

import (
	"fmt"
)

type configureFunc func()

const (
	DefaultChar  = 'X'
	DefaultSize  = 15
	DefaultColor = 0
)

var (
	char     rune
	WithChar = func(c rune) configureFunc {
		return func() {
			char = c
		}
	}

	size     int
	WithSize = func(s int) configureFunc {
		return func() {
			size = s
		}
	}

	color     int
	WithColor = func(c int) configureFunc {
		return func() {
			color = c
		}
	}
)

func Print(functions ...configureFunc) {
	char, size, color = DefaultChar, DefaultSize, DefaultColor
	for _, function := range functions {
		function()
	}

	printPicture(string(char), size, color)
}

func printPicture(char string, size int, color int) {
	fmt.Printf("\033[%dm", color)
	printSymbolNTimes(char, size)
	fmt.Println()

	for i := 1; i < size/2; i++ {
		printLevelOfPyramid(char, size, i)
	}

	if size%2 != 0 {
		printSymbolNTimes(" ", size/2)
		fmt.Println(char)
	}

	for i := size/2 - 1; i > 0; i-- {
		printLevelOfPyramid(char, size, i)
	}

	printSymbolNTimes(char, size)
	fmt.Println("\033[0m")
}

func printLevelOfPyramid(symbol string, size int, level int) {
	printSymbolNTimes(" ", level)
	fmt.Print(symbol)
	printSymbolNTimes(" ", size-(level+1)*2)
	fmt.Print(symbol)
	fmt.Println()
}

func printSymbolNTimes(symbol string, count int) {
	for i := 0; i < count; i++ {
		fmt.Print(symbol)
	}
}
