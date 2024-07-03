package utils

import "unicode"

func reformString(name string) string {
	var reformedString []rune
	for _, ch := range name {
		if unicode.IsSymbol(ch) || unicode.IsPunct(ch) {
			continue
		}
		if unicode.IsSpace(ch) {
			reformedString = append(reformedString, '_')
			continue
		}
		reformedString = append(reformedString, ch)
	}
	return string(reformedString)
}
