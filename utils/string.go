package utils

import (
	"regexp"
	"strings"
)

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func ToPascal(str string) string {
	return camelCaseGenerator(str, true)
}

func ToCamel(str string) string {
	if str == "" {
		return str
	}

	if firstChar := rune(str[0]); firstChar >= 'A' && firstChar <= 'Z' {
		str = strings.ToLower(string(firstChar)) + str[1:]
	}

	return camelCaseGenerator(str, false)
}

func ToSnake(str string) string {
	return snakeCaseGenerator(str, '_', false)
}

func ToScreamingSnake(str string) string {
	return snakeCaseGenerator(str, '_', true)
}

func ToKebab(str string) string {
	return snakeCaseGenerator(str, '-', false)
}

func ToScreamingKebab(str string) string {
	return snakeCaseGenerator(str, '-', true)
}

func snakeCaseGenerator(str string, delimiter uint8, screaming bool) string {
	str = addWordBoundariesToNumbers(str)
	str = strings.Trim(str, " ")

	result := ""

	for i, char := range str {
		nextCaseIsChanged := false

		if i+1 < len(str) {
			next := str[i+1]
			if (char >= 'A' && char <= 'Z' && next >= 'a' && next <= 'z') || (char >= 'a' && char <= 'z' && next >= 'A' && next <= 'Z') {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && result[len(result)-1] != delimiter && nextCaseIsChanged {
			if char >= 'A' && char <= 'Z' {
				result += string(delimiter) + string(char)
			} else if char >= 'a' && char <= 'z' {
				result += string(char) + string(delimiter)
			}
		} else if char == ' ' || char == '_' || char == '-' {
			result += string(delimiter)
		} else {
			result = result + string(char)
		}
	}

	if screaming {
		result = strings.ToUpper(result)
	} else {
		result = strings.ToLower(result)
	}
	return result
}

func camelCaseGenerator(str string, upperInit bool) string {
	str = addWordBoundariesToNumbers(str)
	str = strings.Trim(str, " ")
	result := ""

	capNext := upperInit

	for _, char := range str {
		if char >= 'A' && char <= 'Z' {
			result += string(char)
		} else if char >= '0' && char <= '9' {
			result += string(char)
		} else if char >= 'a' && char <= 'z' {
			if capNext {
				result += strings.ToUpper(string(char))
			} else {
				result += string(char)
			}
		}

		if char == '_' || char == ' ' || char == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}

	return result
}

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}
