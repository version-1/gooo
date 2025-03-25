package strings

import (
	"regexp"
	"strings"
)

// ToCase converts a string to a specific case
func ToCase(s string, f func(byte, int) string) string {
	ns := ""
	for i := 0; i < len(s); i++ {
		ns += f(s[i], i)
	}

	return ns
}

func IsDelimiter(s byte) bool {
	return string(s) == "-" || string(s) == " " || string(s) == "." || string(s) == "_"
}

func IsUpperCase(s string) bool {
	flag := false
	for i := 0; i < len(s); i++ {
		if !IsDelimiter(s[i]) && (byte('A') > s[i] || s[i] > byte('Z')) {
			return false
		}
		flag = true
	}

	return flag
}

func IsUpperCaseChar(s byte) bool {
	return byte('A') <= s && s <= byte('Z')
}

func ToSnakeCase(s string) string {
	return ToCase(s,
		func(c byte, i int) string {
			if IsUpperCaseChar(c) {
				if i == 0 || IsDelimiter(s[i-1]) || IsUpperCaseChar(s[i-1]) {
					return strings.ToLower(string(c))
				}

				return "_" + strings.ToLower(string(c))
			}

			if IsDelimiter(c) {
				return "_"
			}

			return string(c)
		})
}

func ToCamelCase(s string) string {
	return ToCase(s,
		func(c byte, i int) string {
			if i == 0 {
				if IsDelimiter(c) {
					return ""
				}

				return strings.ToUpper(string(c))
			}

			if IsDelimiter(c) {
				return ""
			}

			if IsUpperCaseChar(c) && !IsUpperCaseChar(s[i-1]) {
				return strings.ToUpper(string(c))
			}

			if IsDelimiter(s[i-1]) {
				return strings.ToUpper(string(c))
			}

			return strings.ToLower(string(c))
		})
}

func ToPascalCase(s string) string {
	return ToCase(s,
		func(s byte, i int) string {
			if i == 0 {
				return strings.ToUpper(string(s))
			}

			return string(s)
		})
}

func ToKebabCase(s string) string {
	return ToCase(s,
		func(s byte, i int) string {
			if 'A' <= s && s <= 'Z' {
				if i == 0 {
					return string(s + 32)
				}

				return "_" + string(s+32)
			}

			if string(s) == "-" || string(s) == " " || string(s) == "." {
				return "_"
			}

			return string(s)
		})
}

func ToPlural(word string) string {
	w := strings.ToLower(word)
	irregular := map[string]string{
		"man":    "men",
		"woman":  "women",
		"child":  "children",
		"tooth":  "teeth",
		"foot":   "feet",
		"mouse":  "mice",
		"person": "people",
	}

	if plural, ok := irregular[w]; ok {
		return plural
	}

	// 末尾が "s", "x", "z", "ch", "sh" で終わる単語の場合
	if matched, _ := regexp.MatchString("(s|x|z|ch|sh)$", w); matched {
		return word + "es"
	}

	// 末尾が "f" または "fe" で終わる単語の場合
	if strings.HasSuffix(w, "fe") {
		return word[:len(w)-2] + "ves"
	} else if strings.HasSuffix(w, "f") {
		return word[:len(w)-1] + "ves"
	}

	// 子音 + yで終わる単語の場合
	if matched, _ := regexp.MatchString("[^aeiou]y$", w); matched {
		return w[:len(w)-1] + "ies"
	}

	// 通常の規則 (末尾に "s" を追加)
	return w + "s"
}
