package utils

import (
	"log"
	"runtime"
)

func SayHi() {
	log.Println("SayHi")
}

func SayBye() {
	log.Println("SayBye")
}

const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

func SayColor() {
	LogColor(ColorBlack, "Black")
	LogColor(ColorRed, "Red")
	LogColor(ColorGreen, "Green")
	LogColor(ColorYellow, "Yellow")
	LogColor(ColorBlue, "Blue")
	LogColor(ColorMagenta, "Magenta")
	LogColor(ColorCyan, "Cyan")
	LogColor(ColorWhite, "White")
}

func LogColor(color int, str string) {
	if runtime.GOOS == "windows" {
		log.Println(str)
		return
	}

	switch color {
	case ColorBlack:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorBlack, str)
	case ColorRed:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorRed, str)
	case ColorGreen:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorGreen, str)
	case ColorYellow:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorYellow, str)
	case ColorBlue:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorBlue, str)
	case ColorMagenta:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorMagenta, str)
	case ColorCyan:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorCyan, str)
	case ColorWhite:
		log.Printf("\x1b[0;%dm%s\x1b[0m", ColorWhite, str)
	default:
		log.Println(str)
	}
}
