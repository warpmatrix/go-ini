package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"unicode"
)

var (
	// DefaultSection is the name of default section. You can use this constant or the string literal.
	// In most of cases, an empty string is all you need to access the section.
	DefaultSection = "DEFAULT"

	// KeyValueDelim is a string which contents all possiable delimiter symbols
	KeyValueDelim = "="

	// CommentSym based on running os to set ';' or '#'
	CommentSym = "#"

	// LineBreak is the delimiter to determine or compose a new line.
	// This variable will be changed to "\r\n" automatically on Windows at package init time.
	LineBreak = "\n"
)

func parse(reader *bufio.Reader) (*Config, error) {
	var cfg Config
	cfg.init()
	isEOF := false
	secName := DefaultSection
	for !isEOF {
		line, err := getLine(reader, &isEOF)
		if err != nil {
			return nil, err
		}
		if len(line) == 0 {
			continue
		}
		// Comments
		if line[0] == CommentSym[0] {
			continue
		}
		// Section
		if line[0] == '[' {
			secName, err = parseSecName(line, &cfg)
			if err != nil {
				return nil, err
			}
			continue
		}
		// new Default Section
		if len(cfg.SecList) == 0 {
			err = cfg.newSection(secName)
			if err != nil {
				return nil, err
			}
		}
		// Keyname
		keyName, offset, err := parseKeyName(string(line))
		if err != nil {
			return nil, err
		}
		// Value
		value, err := parseValue(string(line[offset:]))
		if err != nil {
			return nil, err
		}
		// parse Key-Value
		err = cfg.Sections[secName].newKeyVal(keyName, value)
		if err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

func init() {
	switch runtime.GOOS {
	case "windows":
		CommentSym = ";"
		LineBreak = "\r\n"
	case "linux":
		CommentSym = "#"
		LineBreak = "\n"
	}
}

func getLine(buf *bufio.Reader, isEOF *bool) ([]byte, error) {
	line, err := buf.ReadBytes('\n')
	if err == io.EOF {
		*isEOF = true
		err = nil
	} else if err != nil {
		return nil, err
	}
	line = bytes.TrimLeftFunc(line, unicode.IsSpace)
	return line, err

}

func parseSecName(line []byte, cfg *Config) (string, error) {
	closeIdx := bytes.LastIndexByte(line, ']')
	if closeIdx == -1 {
		return "", fmt.Errorf("unclosed section: %s", line)
	}
	secName := string(line[1:closeIdx])
	err := cfg.newSection(secName)
	if err != nil {
		return "", err
	}
	return secName, nil
}

func parseKeyName(line string) (string, int, error) {
	endIdx := strings.IndexAny(line, KeyValueDelim)
	if endIdx < 0 {
		return "", -1, fmt.Errorf("delimiter(%s) not found", KeyValueDelim)
	}
	return strings.TrimSpace(line[0:endIdx]), endIdx + 1, nil
}

func parseValue(line string) (string, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", nil
	}
	// TODO: Check continuation lines when desired

	// Check inline comment
	i := strings.IndexAny(line, CommentSym)
	if i > -1 {
		line = strings.TrimSpace(line[:i])
	}
	return line, nil
}
