package sandglass

import (
	"fmt"
	"strconv"
)

type configString string

const (
	DefaultChar  = 'X'
	DefaultSize  = 15
	DefaultColor = 0

	charPrefix  = "char:"
	sizePrefix  = "size:"
	colorPrefix = "color"
	lenOfPrefix = len(charPrefix)
)

func WithChar(char rune) configString {
	return configString(charPrefix + string(char))
}

func WithSize(size int) configString {
	return configString(sizePrefix + strconv.Itoa(size))
}

func WithColor(color int) configString {
	return configString(colorPrefix + strconv.Itoa(color))
}

func Print(configs ...configString) {
	char, size, color := getConfig(configs...)
	printPicture(string(char), size, color)
}

func getConfig(configs ...configString) (rune, int, int) {
	char, size, color := DefaultChar, DefaultSize, DefaultColor
	for _, config := range configs {
		prefix := config[:lenOfPrefix]
		switch prefix {
		case charPrefix:
			char = rune(config[lenOfPrefix])
		case sizePrefix:
			size, _ = strconv.Atoi(string(config[lenOfPrefix:]))
		case colorPrefix:
			color, _ = strconv.Atoi(string(config[lenOfPrefix:]))
		}
	}

	return char, size, color
}

func printPicture(char string, size int, color int) {
	fmt.Printf("\033[%dm", color)
	printSymbolNTimes(char, size)
	fmt.Println()

	for i := 1; i < size / 2; i++ {
		printLevelOfPyramid(char, size, i)
	}

	if size % 2 != 0 {
		printSymbolNTimes(" ", size / 2)
		fmt.Println(char)
	}

	for i := size / 2 - 1; i > 0; i-- {
		printLevelOfPyramid(char, size, i)
	}

	printSymbolNTimes(char, size)
	fmt.Println("\033[0m")
}

func printLevelOfPyramid(symbol string, size int, level int) {
	printSymbolNTimes(" ", level)
	fmt.Print(symbol)
	printSymbolNTimes(" ", size - (level + 1) * 2)
	fmt.Print(symbol)
	fmt.Println()
}

func printSymbolNTimes(symbol string, count int) {
	for i := 0; i < count; i++ {
		fmt.Print(symbol)
	}
}
