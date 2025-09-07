package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/lmittmann/tint"
)

// sets up the slog handler for the application
//
// TODO colorize errors
// TODO colorize a few things in display.go if possible
func slogHandler(level string) slog.Handler {
	lvl := slog.LevelInfo
	if level != "" {
		switch level {
		case "debug":
			lvl = slog.LevelDebug
		case "info":
			lvl = slog.LevelInfo
		case "warn":
			lvl = slog.LevelWarn
		case "error":
			lvl = slog.LevelError
		default:
			lvl = slog.LevelInfo
		}
	}

	prettySource := func(_ []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.SourceKey {
			source, ok := attr.Value.Any().(*slog.Source)
			if ok {
				dir, file := filepath.Split(source.File)
				return slog.String("src", fmt.Sprintf("%s:%d", filepath.Join(filepath.Base(dir), file), source.Line))
			}
		}
		if attr.Value.Kind() == slog.KindAny {
			if _, ok := attr.Value.Any().(error); ok {
				return tint.Attr(9, attr)
			}
		}
		return attr
	}

	return tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:   true,
		Level:       lvl,
		TimeFormat:  time.DateTime,
		ReplaceAttr: prettySource,
	})
}

// readFromFile reads the content of a specified file and returns its contents as a string.
func readFromFile(fileName string) (string, error) {
	slog.Debug("Reading from file", "filename", fileName)

	byteArr, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(byteArr), nil
}

// parseFunds reads the content of the specified file and parses it into a slice of Funds.
func parseFunds(filePath string) ([]Fund, error) {
	slog.Debug("Parsing funds", "filepath", filePath)

	fileContent, err := readFromFile(filePath)
	if err != nil {
		return nil, err
	}

	funds := Funds{}
	_, err = toml.Decode(fileContent, &funds)
	if err != nil {
		return nil, err
	}

	// TODO there ought to be a more correct way to parse the file directly into an array of Funds
	return funds.Funds, nil
}

func repeat(s string, count int) string {
	result := ""
	for range count {
		result += s
	}
	return result
}
