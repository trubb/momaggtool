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
func slogHandler(level string) slog.Handler {
	// TODO make the function care about level
	if level != "" {
		switch level {
		case "debug":
			_ = slog.LevelDebug
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
		return attr
	}

	return tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:   true,
		Level:       slog.LevelDebug,
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

	// slog.Debug("Parsed funds", "funds", fmt.Sprintf("%+v", funds))

	// TODO there ought to be a more correct way to parse the file directly into an array of Funds
	return funds.Funds, nil
}
