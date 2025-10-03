package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func initLog() zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}
	consoleWriter.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = "\x1b[35mDEBG\x1b[0m"
			case "info":
				l = "\x1b[32mINFO\x1b[0m"
			case "warn":
				l = "\x1b[33mWARN\x1b[0m"
			case "error":
				l = "\x1b[31mERRO\x1b[0m"
			case "fatal":
				l = "\x1b[31mFATL\x1b[0m"
			case "panic":
				l = "\x1b[31mPANC\x1b[0m"
			default:
				l = strings.ToUpper(ll)
			}
		} else {
			l = strings.ToUpper(fmt.Sprintf("%s", i))
		}
		return fmt.Sprintf("| %-10s|", l)
	}
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	return logger
}
