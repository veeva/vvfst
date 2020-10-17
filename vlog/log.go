/*
This code serves as an example and is not meant for production use.

Copyright 2020 Veeva Systems Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions
and limitations under the License.
*/
package vlog

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

var logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
var consoleNoColor = true

func InitLog(nocolor bool) {
	output := zerolog.ConsoleWriter{
		Out:         os.Stdout,
		TimeFormat:  zerolog.TimeFormatUnix,
		FormatLevel: consoleDefaultFormatLevel(nocolor),
	}

	consoleNoColor = nocolor
	logger = zerolog.New(output).With().Timestamp().Logger().Level(zerolog.InfoLevel)
}

func Trace(msg string) {
	logger.Trace().Msg(msg)
}

func Tracef(msg string, v ...interface{}) {
	logger.Trace().Msgf(msg, v...)
}

func Debug(msg string) {
	logger.Debug().Msg(msg)
}

func Debugf(msg string, v ...interface{}) {
	logger.Debug().Msgf(msg, v...)
}

func Info(msg string) {
	logger.Info().Msg(msg)
}

func Infof(msg string, v ...interface{}) {
	logger.Info().Msgf(msg, v...)
}

func Warn(msg string) {
	logger.Warn().Msg(msg)
}

func Warnf(msg string, v ...interface{}) {
	logger.Warn().Msgf(msg, v...)
}

func Error(msg string) {
	logger.Error().Msg(msg)
}

func Errorf(msg string, v ...interface{}) {
	logger.Error().Msgf(msg, v...)
}

func Fatal(msg string) {
	logger.Fatal().Msg(msg)
}

func Fatalf(msg string, v ...interface{}) {
	logger.Fatal().Msgf(msg, v...)
}

func Panic(msg string) {
	logger.Panic().Msg(msg)
}

func Panicf(msg string, v ...interface{}) {
	logger.Panic().Msgf(msg, v...)
}

func NoFormatLog(msg string) {
	fmt.Println(msg)
}

func NoFormatLogf(msg string, v ...interface{}) {
	fmt.Printf(msg+"\n", v...)
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

	ColorBold     = 1
	ColorDarkGray = 90
)

func consoleDefaultFormatLevel(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				l = colorize("TRACE", ColorMagenta, noColor)
			case "debug":
				l = colorize("DEBUG", ColorYellow, noColor)
			case "info":
				l = colorize("INFO ", ColorGreen, noColor)
			case "warn":
				l = colorize("WARN ", ColorRed, noColor)
			case "error":
				l = colorize(colorize("ERROR", ColorRed, noColor), ColorBold, noColor)
			case "fatal":
				l = colorize(colorize("FATAL", ColorRed, noColor), ColorBold, noColor)
			case "panic":
				l = colorize(colorize("PANIC", ColorRed, noColor), ColorBold, noColor)
			default:
				l = colorize("?????", ColorBold, noColor)
			}
		} else {
			if i == nil {
				l = colorize("?????", ColorBold, noColor)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}
}

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func Colorize(s interface{}, c int) string {
	if consoleNoColor {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
