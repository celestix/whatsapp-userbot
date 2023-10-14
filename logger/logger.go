package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Level string

var (
	LevelMain     Level = "MAIN"
	LevelInfo     Level = "INFO"
	LevelError    Level = "ERROR"
	LevelCritical Level = "CRITICAL"

	LevelToNum = map[Level]int{
		LevelInfo:     0,
		LevelMain:     1,
		LevelError:    2,
		LevelCritical: 3,
	}
)

type color string

var (
	colorReset color = "\033[0m"

	colorRed    color = "\033[31m"
	colorGreen  color = "\033[32m"
	colorYellow color = "\033[33m"
	colorBlue   color = "\033[34m"
	colorPurple color = "\033[35m"
	colorCyan   color = "\033[36m"
	colorWhite  color = "\033[37m"

	levelToColor = map[Level]color{
		LevelInfo:     colorBlue,
		LevelMain:     colorGreen,
		LevelError:    colorPurple,
		LevelCritical: colorRed,
	}
)

type Logger struct {
	Name      string
	Level     Level
	Separator string
	minima    int
	*log.Logger
}

func NewLogger(minLevel Level) *Logger {
	return &Logger{
		Logger:    log.New(os.Stdout, "", 0),
		Level:     LevelMain,
		minima:    LevelToNum[minLevel],
		Separator: " ",
	}
}

func (l *Logger) ChangeLevel(level Level) *Logger {
	l.Level = level
	return l
}

func (l *Logger) ChangeMinima(minLevel Level) *Logger {
	l.minima = LevelToNum[minLevel]
	return l
}

func (l *Logger) Create(name string) *Logger {
	return &Logger{
		Level:     l.Level,
		minima:    l.minima,
		Logger:    l.Logger,
		Separator: " ",
		Name:      fmt.Sprintf("%s[%s]", l.Name, strings.ToUpper(name)),
	}

}

func (l *Logger) shouldDo() bool {
	if LevelToNum[l.Level] < l.minima {
		return false
	}
	return true
}

func (l *Logger) rawNewline() string {
	if l.Name != "" {
		return fmt.Sprintf("[WAUB][%s]%s ", l.Level, l.Name)
	}
	return fmt.Sprintf("[WAUB][%s] ", l.Level)
}

func (l *Logger) Println(v ...any) {
	if !l.shouldDo() {
		return
	}
	l.Print(string(levelToColor[l.Level]), l.rawNewline(), l.separateText(v...), string(colorReset))
}

func (l *Logger) Fatalln(v ...any) {
	l.Fatal(string(colorRed), l.rawNewline(), l.separateText(v...), string(colorReset))
}

func (l *Logger) separateText(v ...any) string {
	var t string
	for n, i := range v {
		if n != 0 {
			t += fmt.Sprintf(" %v", i)
		} else {
			t += fmt.Sprintf("%v", i)
		}
	}
	return t
}
