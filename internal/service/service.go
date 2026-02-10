package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

var errUnknownSimbols = errors.New("unknown simbols")

func CodeDetection(str string) (string, error) {
	var result string

	textFunc := func(r rune) bool {
		return (r >= 'А' && r <= 'я') || (r >= '0' && r <= '9') || (r == '"') || (r == '/') ||
			(r == '(') || (r == ')') || (r == '\'') || (r == '?') || (r == ':') || (r == ',')
	}
	morseFunc := func(r rune) bool {
		return r == '.' || r == '-'
	}

	if strings.ContainsFunc(str, textFunc) {
		result = morse.ToMorse(str)
	} else if strings.ContainsFunc(str, morseFunc) {
		result = morse.ToText(str)
	} else {
		return "", errUnknownSimbols
	}

	return result, nil
}
