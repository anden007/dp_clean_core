package misc

import (
	"fmt"
	"time"
)

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextPurple
	TextCyan
	TextWhite
)

func ColorText(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}

func ServiceLoadInfo(serviceName string, enable bool, loadTime time.Time) {
	if enable {
		fmt.Println(ColorText(TextGreen, fmt.Sprintf("[✔] %s Enable(%s)", serviceName, time.Since(loadTime))))
	} else {
		fmt.Println(ColorText(TextRed, fmt.Sprintf("[✘] %s Disable", serviceName)))
	}
}

func PrintErrorInfo(str string) {
	fmt.Println(ColorText(TextRed, str))
}

func PrintInfo(str string) {
	fmt.Println(ColorText(TextGreen, str))
}
