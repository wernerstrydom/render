package funcs

import (
	"testing"
)

func TestCasingFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) string
		input    string
		expected string
	}{
		// camelCase
		{"camelCase from space", toCamelCase, "hello world", "helloWorld"},
		{"camelCase from snake", toCamelCase, "hello_world", "helloWorld"},
		{"camelCase from kebab", toCamelCase, "hello-world", "helloWorld"},
		{"camelCase from pascal", toCamelCase, "HelloWorld", "helloWorld"},
		{"camelCase single word", toCamelCase, "hello", "hello"},
		{"camelCase empty", toCamelCase, "", ""},
		{"camelCase with acronym", toCamelCase, "XMLParser", "xmlParser"},

		// pascalCase
		{"pascalCase from space", toPascalCase, "hello world", "HelloWorld"},
		{"pascalCase from snake", toPascalCase, "hello_world", "HelloWorld"},
		{"pascalCase from kebab", toPascalCase, "hello-world", "HelloWorld"},
		{"pascalCase from camel", toPascalCase, "helloWorld", "HelloWorld"},
		{"pascalCase single word", toPascalCase, "hello", "Hello"},
		{"pascalCase empty", toPascalCase, "", ""},

		// snakeCase
		{"snakeCase from space", toSnakeCase, "hello world", "hello_world"},
		{"snakeCase from camel", toSnakeCase, "helloWorld", "hello_world"},
		{"snakeCase from pascal", toSnakeCase, "HelloWorld", "hello_world"},
		{"snakeCase from kebab", toSnakeCase, "hello-world", "hello_world"},
		{"snakeCase single word", toSnakeCase, "hello", "hello"},
		{"snakeCase empty", toSnakeCase, "", ""},

		// kebabCase
		{"kebabCase from space", toKebabCase, "hello world", "hello-world"},
		{"kebabCase from camel", toKebabCase, "helloWorld", "hello-world"},
		{"kebabCase from pascal", toKebabCase, "HelloWorld", "hello-world"},
		{"kebabCase from snake", toKebabCase, "hello_world", "hello-world"},
		{"kebabCase single word", toKebabCase, "hello", "hello"},
		{"kebabCase empty", toKebabCase, "", ""},

		// upperSnakeCase
		{"upperSnakeCase from space", toUpperSnakeCase, "hello world", "HELLO_WORLD"},
		{"upperSnakeCase from camel", toUpperSnakeCase, "helloWorld", "HELLO_WORLD"},
		{"upperSnakeCase single word", toUpperSnakeCase, "hello", "HELLO"},

		// upperKebabCase
		{"upperKebabCase from space", toUpperKebabCase, "hello world", "HELLO-WORLD"},
		{"upperKebabCase from camel", toUpperKebabCase, "helloWorld", "HELLO-WORLD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.input)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStringManipulation(t *testing.T) {
	t.Run("reverseString", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"hello", "olleh"},
			{"", ""},
			{"a", "a"},
			{"ab", "ba"},
			{"日本語", "語本日"},
		}
		for _, tt := range tests {
			result := reverseString(tt.input)
			if result != tt.expected {
				t.Errorf("reverseString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("substring", func(t *testing.T) {
		tests := []struct {
			input    string
			start    int
			end      int
			expected string
		}{
			{"hello", 0, 5, "hello"},
			{"hello", 1, 4, "ell"},
			{"hello", 0, 0, ""},
			{"hello", -2, 5, "lo"},
			{"hello", 0, -1, "hell"},
			{"hello", 0, 100, "hello"},
			{"日本語", 0, 2, "日本"},
		}
		for _, tt := range tests {
			result := substring(tt.input, tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("substring(%q, %d, %d) = %q, want %q", tt.input, tt.start, tt.end, result, tt.expected)
			}
		}
	})

	t.Run("truncate", func(t *testing.T) {
		tests := []struct {
			input    string
			length   int
			expected string
		}{
			{"hello world", 5, "hello"},
			{"hello", 10, "hello"},
			{"hello", 0, ""},
			{"日本語", 2, "日本"},
		}
		for _, tt := range tests {
			result := truncate(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.length, result, tt.expected)
			}
		}
	})

	t.Run("padLeft", func(t *testing.T) {
		tests := []struct {
			input    string
			length   int
			pad      string
			expected string
		}{
			{"hello", 10, " ", "     hello"},
			{"hello", 10, "0", "00000hello"},
			{"hello", 3, " ", "hello"}, // string longer than target, return as-is
			{"hello", 5, " ", "hello"},
		}
		for _, tt := range tests {
			result := padLeft(tt.input, tt.length, tt.pad)
			if result != tt.expected {
				t.Errorf("padLeft(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		}
	})

	t.Run("padRight", func(t *testing.T) {
		tests := []struct {
			input    string
			length   int
			pad      string
			expected string
		}{
			{"hello", 10, " ", "hello     "},
			{"hello", 10, "0", "hello00000"},
			{"hello", 3, " ", "hello"}, // string longer than target, return as-is
			{"hello", 5, " ", "hello"},
		}
		for _, tt := range tests {
			result := padRight(tt.input, tt.length, tt.pad)
			if result != tt.expected {
				t.Errorf("padRight(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		}
	})

	t.Run("center", func(t *testing.T) {
		tests := []struct {
			input    string
			length   int
			pad      string
			expected string
		}{
			{"hello", 11, " ", "   hello   "},
			{"hello", 10, "-", "--hello---"},
			{"hello", 5, " ", "hello"},
			{"hello", 3, " ", "hello"},
		}
		for _, tt := range tests {
			result := center(tt.input, tt.length, tt.pad)
			if result != tt.expected {
				t.Errorf("center(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		}
	})

	t.Run("indent", func(t *testing.T) {
		tests := []struct {
			spaces   int
			input    string
			expected string
		}{
			{2, "hello", "  hello"},
			{2, "hello\nworld", "  hello\n  world"},
			{0, "hello", "hello"},
		}
		for _, tt := range tests {
			result := indent(tt.spaces, tt.input)
			if result != tt.expected {
				t.Errorf("indent(%d, %q) = %q, want %q", tt.spaces, tt.input, result, tt.expected)
			}
		}
	})

	t.Run("nindent", func(t *testing.T) {
		result := nindent(2, "hello")
		expected := "\n  hello"
		if result != expected {
			t.Errorf("nindent(2, %q) = %q, want %q", "hello", result, expected)
		}
	})
}

func TestSplittingJoining(t *testing.T) {
	t.Run("join", func(t *testing.T) {
		tests := []struct {
			sep      string
			items    any
			expected string
		}{
			{", ", []string{"a", "b", "c"}, "a, b, c"},
			{"-", []any{"a", "b", "c"}, "a-b-c"},
			{"", []string{"a", "b"}, "ab"},
		}
		for _, tt := range tests {
			result := join(tt.sep, tt.items)
			if result != tt.expected {
				t.Errorf("join(%q, %v) = %q, want %q", tt.sep, tt.items, result, tt.expected)
			}
		}
	})

	t.Run("first", func(t *testing.T) {
		if result := first([]any{"a", "b", "c"}); result != "a" {
			t.Errorf("first([]any) = %v, want 'a'", result)
		}
		if result := first([]string{"a", "b", "c"}); result != "a" {
			t.Errorf("first([]string) = %v, want 'a'", result)
		}
		if result := first("abc"); result != "a" {
			t.Errorf("first(string) = %v, want 'a'", result)
		}
		if result := first([]any{}); result != nil {
			t.Errorf("first(empty) = %v, want nil", result)
		}
	})

	t.Run("last", func(t *testing.T) {
		if result := last([]any{"a", "b", "c"}); result != "c" {
			t.Errorf("last([]any) = %v, want 'c'", result)
		}
		if result := last([]string{"a", "b", "c"}); result != "c" {
			t.Errorf("last([]string) = %v, want 'c'", result)
		}
		if result := last("abc"); result != "c" {
			t.Errorf("last(string) = %v, want 'c'", result)
		}
	})

	t.Run("rest", func(t *testing.T) {
		result := rest([]any{"a", "b", "c"}).([]any)
		if len(result) != 2 || result[0] != "b" || result[1] != "c" {
			t.Errorf("rest([]any) = %v, want [b c]", result)
		}
	})

	t.Run("initial", func(t *testing.T) {
		result := initial([]any{"a", "b", "c"}).([]any)
		if len(result) != 2 || result[0] != "a" || result[1] != "b" {
			t.Errorf("initial([]any) = %v, want [a b]", result)
		}
	})

	t.Run("nth", func(t *testing.T) {
		if result := nth(1, []any{"a", "b", "c"}); result != "b" {
			t.Errorf("nth(1, []any) = %v, want 'b'", result)
		}
		if result := nth(10, []any{"a", "b", "c"}); result != nil {
			t.Errorf("nth(10, []any) = %v, want nil", result)
		}
	})
}

func TestConversionFunctions(t *testing.T) {
	t.Run("toString", func(t *testing.T) {
		tests := []struct {
			input    any
			expected string
		}{
			{"hello", "hello"},
			{123, "123"},
			{12.5, "12.5"},
			{true, "true"},
			{nil, ""},
			{[]byte("bytes"), "bytes"},
		}
		for _, tt := range tests {
			result := toString(tt.input)
			if result != tt.expected {
				t.Errorf("toString(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("toInt", func(t *testing.T) {
		tests := []struct {
			input    any
			expected int
		}{
			{123, 123},
			{int64(123), 123},
			{float64(123.9), 123},
			{"123", 123},
			{true, 1},
			{false, 0},
		}
		for _, tt := range tests {
			result, err := toInt(tt.input)
			if err != nil {
				t.Errorf("toInt(%v) unexpected error: %v", tt.input, err)
				continue
			}
			if result != tt.expected {
				t.Errorf("toInt(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		}

		// Test error case
		_, err := toInt("not a number")
		if err == nil {
			t.Error("toInt('not a number') should return error")
		}
	})

	t.Run("toFloat", func(t *testing.T) {
		tests := []struct {
			input    any
			expected float64
		}{
			{123, 123.0},
			{12.5, 12.5},
			{"12.5", 12.5},
			{int64(123), 123.0},
		}
		for _, tt := range tests {
			result, err := toFloat(tt.input)
			if err != nil {
				t.Errorf("toFloat(%v) unexpected error: %v", tt.input, err)
				continue
			}
			if result != tt.expected {
				t.Errorf("toFloat(%v) = %f, want %f", tt.input, result, tt.expected)
			}
		}

		// Test error case
		_, err := toFloat("not a number")
		if err == nil {
			t.Error("toFloat('not a number') should return error")
		}
	})

	t.Run("toBool", func(t *testing.T) {
		tests := []struct {
			input    any
			expected bool
		}{
			{true, true},
			{false, false},
			{"true", true},
			{"false", false},
			{1, true},
			{0, false},
			{1.0, true},
			{0.0, false},
		}
		for _, tt := range tests {
			result, err := toBool(tt.input)
			if err != nil {
				t.Errorf("toBool(%v) unexpected error: %v", tt.input, err)
				continue
			}
			if result != tt.expected {
				t.Errorf("toBool(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		}

		// Test error case
		_, err := toBool("not a bool")
		if err == nil {
			t.Error("toBool('not a bool') should return error")
		}
	})

	t.Run("toJSON", func(t *testing.T) {
		result, err := toJSON(map[string]any{"key": "value"})
		if err != nil {
			t.Fatalf("toJSON() unexpected error: %v", err)
		}
		expected := `{"key":"value"}`
		if result != expected {
			t.Errorf("toJSON() = %q, want %q", result, expected)
		}
	})

	t.Run("fromJSON", func(t *testing.T) {
		result, err := fromJSON(`{"key":"value"}`)
		if err != nil {
			t.Fatalf("fromJSON() unexpected error: %v", err)
		}
		m := result.(map[string]any)
		if m["key"] != "value" {
			t.Errorf("fromJSON() key = %v, want 'value'", m["key"])
		}

		// Test error case
		_, err = fromJSON("invalid json")
		if err == nil {
			t.Error("fromJSON('invalid json') should return error")
		}
	})
}

func TestUnicodeFunctions(t *testing.T) {
	t.Run("toASCII", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"hello", "hello"},
			{"café", "cafe"},
			{"naïve", "naive"},
			{"日本語", ""},
		}
		for _, tt := range tests {
			result := toASCII(tt.input)
			if result != tt.expected {
				t.Errorf("toASCII(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("toSlug", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"Hello World", "hello-world"},
			{"Hello  World", "hello-world"},
			{"Hello---World", "hello-world"},
			{"Café Naïve", "cafe-naive"},
			{"  hello  ", "hello"},
		}
		for _, tt := range tests {
			result := toSlug(tt.input)
			if result != tt.expected {
				t.Errorf("toSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		}
	})
}

func TestLogicFunctions(t *testing.T) {
	t.Run("eq", func(t *testing.T) {
		if !eq("a", "a") {
			t.Error("eq('a', 'a') should be true")
		}
		if eq("a", "b") {
			t.Error("eq('a', 'b') should be false")
		}
		if !eq(1, 1) {
			t.Error("eq(1, 1) should be true")
		}
	})

	t.Run("ne", func(t *testing.T) {
		if !ne("a", "b") {
			t.Error("ne('a', 'b') should be true")
		}
		if ne("a", "a") {
			t.Error("ne('a', 'a') should be false")
		}
	})

	t.Run("lt/le/gt/ge", func(t *testing.T) {
		if !lt(1, 2) {
			t.Error("lt(1, 2) should be true")
		}
		if !le(1, 1) {
			t.Error("le(1, 1) should be true")
		}
		if !gt(2, 1) {
			t.Error("gt(2, 1) should be true")
		}
		if !ge(1, 1) {
			t.Error("ge(1, 1) should be true")
		}
	})

	t.Run("and/or/not", func(t *testing.T) {
		if !and(true, true) {
			t.Error("and(true, true) should be true")
		}
		if and(true, false) {
			t.Error("and(true, false) should be false")
		}
		if !or(true, false) {
			t.Error("or(true, false) should be true")
		}
		if or(false, false) {
			t.Error("or(false, false) should be false")
		}
		if !not(false) {
			t.Error("not(false) should be true")
		}
	})

	t.Run("default", func(t *testing.T) {
		if result := defaultVal("default", ""); result != "default" {
			t.Errorf("defaultVal with empty = %v, want 'default'", result)
		}
		if result := defaultVal("default", "value"); result != "value" {
			t.Errorf("defaultVal with value = %v, want 'value'", result)
		}
	})

	t.Run("empty", func(t *testing.T) {
		if !empty("") {
			t.Error("empty('') should be true")
		}
		if !empty(nil) {
			t.Error("empty(nil) should be true")
		}
		if !empty([]any{}) {
			t.Error("empty([]) should be true")
		}
		if empty("hello") {
			t.Error("empty('hello') should be false")
		}
	})

	t.Run("coalesce", func(t *testing.T) {
		if result := coalesce("", nil, "value", "other"); result != "value" {
			t.Errorf("coalesce() = %v, want 'value'", result)
		}
	})

	t.Run("ternary", func(t *testing.T) {
		if result := ternary(true, "yes", "no"); result != "yes" {
			t.Errorf("ternary(true) = %v, want 'yes'", result)
		}
		if result := ternary(false, "yes", "no"); result != "no" {
			t.Errorf("ternary(false) = %v, want 'no'", result)
		}
	})
}

func TestCollectionFunctions(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		result := list("a", "b", "c")
		if len(result) != 3 || result[0] != "a" {
			t.Errorf("list() = %v, want [a b c]", result)
		}
	})

	t.Run("dict", func(t *testing.T) {
		result := dict("key1", "value1", "key2", "value2")
		if result["key1"] != "value1" || result["key2"] != "value2" {
			t.Errorf("dict() = %v", result)
		}
	})

	t.Run("keys", func(t *testing.T) {
		m := map[string]any{"a": 1, "b": 2}
		result := keys(m)
		if len(result) != 2 {
			t.Errorf("keys() len = %d, want 2", len(result))
		}
	})

	t.Run("values", func(t *testing.T) {
		m := map[string]any{"a": 1, "b": 2}
		result := values(m)
		if len(result) != 2 {
			t.Errorf("values() len = %d, want 2", len(result))
		}
	})

	t.Run("hasKey", func(t *testing.T) {
		m := map[string]any{"key": "value"}
		if !hasKey(m, "key") {
			t.Error("hasKey() should be true")
		}
		if hasKey(m, "missing") {
			t.Error("hasKey() should be false for missing key")
		}
	})

	t.Run("get", func(t *testing.T) {
		m := map[string]any{"key": "value"}
		if result := get(m, "key"); result != "value" {
			t.Errorf("get() = %v, want 'value'", result)
		}
	})

	t.Run("set", func(t *testing.T) {
		m := map[string]any{"key": "value"}
		result := set(m, "new", "newvalue")
		if result["new"] != "newvalue" {
			t.Errorf("set() = %v", result)
		}
	})

	t.Run("merge", func(t *testing.T) {
		m1 := map[string]any{"a": 1}
		m2 := map[string]any{"b": 2}
		result := merge(m1, m2)
		if result["a"] != 1 || result["b"] != 2 {
			t.Errorf("merge() = %v", result)
		}
	})

	t.Run("append", func(t *testing.T) {
		result := appendList([]any{"a", "b"}, "c")
		if len(result) != 3 || result[2] != "c" {
			t.Errorf("append() = %v", result)
		}
	})

	t.Run("uniq", func(t *testing.T) {
		result := uniq([]any{"a", "b", "a", "c", "b"})
		if len(result) != 3 {
			t.Errorf("uniq() len = %d, want 3", len(result))
		}
	})

	t.Run("sortAlpha", func(t *testing.T) {
		result := sortAlpha([]any{"c", "a", "b"})
		if result[0] != "a" || result[1] != "b" || result[2] != "c" {
			t.Errorf("sortAlpha() = %v, want [a b c]", result)
		}
	})

	t.Run("length", func(t *testing.T) {
		if result := length("hello"); result != 5 {
			t.Errorf("length(string) = %d, want 5", result)
		}
		if result := length([]any{1, 2, 3}); result != 3 {
			t.Errorf("length([]any) = %d, want 3", result)
		}
		if result := length(map[string]any{"a": 1}); result != 1 {
			t.Errorf("length(map) = %d, want 1", result)
		}
	})
}

func TestMathFunctions(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		if result := add(1, 2); result != 3 {
			t.Errorf("add(1, 2) = %v, want 3", result)
		}
		// add returns int when result is whole number
		if result := add(1.5, 2.5); result != 4 {
			t.Errorf("add(1.5, 2.5) = %v, want 4", result)
		}
		// add returns float when result has decimal
		if result := add(1.5, 2.3); result != 3.8 {
			t.Errorf("add(1.5, 2.3) = %v, want 3.8", result)
		}
	})

	t.Run("sub", func(t *testing.T) {
		if result := sub(5, 3); result != 2 {
			t.Errorf("sub(5, 3) = %v, want 2", result)
		}
	})

	t.Run("mul", func(t *testing.T) {
		if result := mul(3, 4); result != 12 {
			t.Errorf("mul(3, 4) = %v, want 12", result)
		}
	})

	t.Run("div", func(t *testing.T) {
		result, err := div(10, 2)
		if err != nil {
			t.Errorf("div(10, 2) unexpected error: %v", err)
		}
		if result != 5 {
			t.Errorf("div(10, 2) = %v, want 5", result)
		}

		// Test division by zero returns error
		_, err = div(10, 0)
		if err == nil {
			t.Error("div(10, 0) should return error")
		}
	})

	t.Run("mod", func(t *testing.T) {
		result, err := mod(10, 3)
		if err != nil {
			t.Errorf("mod(10, 3) unexpected error: %v", err)
		}
		if result != 1 {
			t.Errorf("mod(10, 3) = %v, want 1", result)
		}

		// Test division by zero
		_, err = mod(10, 0)
		if err == nil {
			t.Error("mod(10, 0) should return error")
		}
	})

	t.Run("max", func(t *testing.T) {
		if result := max(1, 5, 3); result != 5 {
			t.Errorf("max(1, 5, 3) = %v, want 5", result)
		}
	})

	t.Run("min", func(t *testing.T) {
		if result := min(1, 5, 3); result != 1 {
			t.Errorf("min(1, 5, 3) = %v, want 1", result)
		}
	})

	t.Run("floor", func(t *testing.T) {
		if result := floor(3.7); result != 3 {
			t.Errorf("floor(3.7) = %v, want 3", result)
		}
	})

	t.Run("ceil", func(t *testing.T) {
		if result := ceil(3.2); result != 4 {
			t.Errorf("ceil(3.2) = %v, want 4", result)
		}
	})

	t.Run("round", func(t *testing.T) {
		if result := round(3.4); result != 3 {
			t.Errorf("round(3.4) = %v, want 3", result)
		}
		if result := round(3.6); result != 4 {
			t.Errorf("round(3.6) = %v, want 4", result)
		}
	})
}

func TestRegexFunctions(t *testing.T) {
	t.Run("regexMatch", func(t *testing.T) {
		result, err := regexMatch(`\d+`, "abc123def")
		if err != nil {
			t.Fatalf("regexMatch() unexpected error: %v", err)
		}
		if !result {
			t.Error("regexMatch should match digits")
		}

		result, err = regexMatch(`\d+`, "abcdef")
		if err != nil {
			t.Fatalf("regexMatch() unexpected error: %v", err)
		}
		if result {
			t.Error("regexMatch should not match")
		}

		// Test invalid regex
		_, err = regexMatch(`[invalid`, "test")
		if err == nil {
			t.Error("regexMatch() should return error for invalid regex")
		}
	})

	t.Run("regexFind", func(t *testing.T) {
		result, err := regexFind(`\d+`, "abc123def")
		if err != nil {
			t.Fatalf("regexFind() unexpected error: %v", err)
		}
		if result != "123" {
			t.Errorf("regexFind() = %q, want '123'", result)
		}

		// Test invalid regex
		_, err = regexFind(`[invalid`, "test")
		if err == nil {
			t.Error("regexFind() should return error for invalid regex")
		}
	})

	t.Run("regexFindAll", func(t *testing.T) {
		result, err := regexFindAll(`\d+`, "a1b2c3", -1)
		if err != nil {
			t.Fatalf("regexFindAll() unexpected error: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("regexFindAll() len = %d, want 3", len(result))
		}

		// Test invalid regex
		_, err = regexFindAll(`[invalid`, "test", -1)
		if err == nil {
			t.Error("regexFindAll() should return error for invalid regex")
		}
	})

	t.Run("regexReplace", func(t *testing.T) {
		result, err := regexReplace(`\d+`, "X", "a1b2c3")
		if err != nil {
			t.Fatalf("regexReplace() unexpected error: %v", err)
		}
		if result != "aXbXcX" {
			t.Errorf("regexReplace() = %q, want 'aXbXcX'", result)
		}

		// Test invalid regex
		_, err = regexReplace(`[invalid`, "X", "test")
		if err == nil {
			t.Error("regexReplace() should return error for invalid regex")
		}
	})

	t.Run("regexSplit", func(t *testing.T) {
		result, err := regexSplit(`\d+`, "a1b2c3", -1)
		if err != nil {
			t.Fatalf("regexSplit() unexpected error: %v", err)
		}
		if len(result) != 4 || result[0] != "a" {
			t.Errorf("regexSplit() = %v", result)
		}

		// Test invalid regex
		_, err = regexSplit(`[invalid`, "test", -1)
		if err == nil {
			t.Error("regexSplit() should return error for invalid regex")
		}
	})
}

func TestCountingFunctions(t *testing.T) {
	t.Run("count", func(t *testing.T) {
		if result := count("o", "hello world"); result != 2 {
			t.Errorf("count('o', 'hello world') = %d, want 2", result)
		}
	})

	t.Run("countWords", func(t *testing.T) {
		if result := countWords("hello world foo"); result != 3 {
			t.Errorf("countWords() = %d, want 3", result)
		}
	})

	t.Run("countLines", func(t *testing.T) {
		if result := countLines("line1\nline2\nline3"); result != 3 {
			t.Errorf("countLines() = %d, want 3", result)
		}
		if result := countLines(""); result != 0 {
			t.Errorf("countLines('') = %d, want 0", result)
		}
	})
}

func TestQuoteFunctions(t *testing.T) {
	t.Run("quote", func(t *testing.T) {
		if result := quote("hello"); result != `"hello"` {
			t.Errorf("quote() = %q, want '\"hello\"'", result)
		}
		if result := quote(`say "hi"`); result != `"say \"hi\""` {
			t.Errorf("quote() with quotes = %q", result)
		}
	})

	t.Run("squote", func(t *testing.T) {
		if result := squote("hello"); result != `'hello'` {
			t.Errorf("squote() = %q, want \"'hello'\"", result)
		}
	})
}

func TestFuncMap(t *testing.T) {
	funcMap := Map()

	// Verify that the map contains expected functions
	expectedFuncs := []string{
		"lower", "upper", "camelCase", "snakeCase", "kebabCase",
		"trim", "replace", "split", "join",
		"toInt", "toString", "toJson",
		"add", "sub", "mul", "div",
		"regexMatch", "regexReplace",
	}

	for _, name := range expectedFuncs {
		if _, ok := funcMap[name]; !ok {
			t.Errorf("FuncMap missing function: %s", name)
		}
	}
}

// TestPipelineArgOrder verifies that wrapper functions have correct argument
// ordering for Go template pipeline usage. In pipelines, the piped value
// becomes the LAST argument. These wrappers exist because the stdlib functions
// take the string-to-operate-on as the FIRST argument, which is the wrong
// position for pipelines.
//
// If these tests fail, it likely means a wrapper was replaced with a direct
// stdlib binding, which breaks all pipeline usage in user templates.
func TestPipelineArgOrder(t *testing.T) {
	t.Run("trimPrefix", func(t *testing.T) {
		// {{ "hello.txt" | trimPrefix "hello" }} should produce ".txt"
		result := trimPrefix("hello", "hello.txt")
		if result != ".txt" {
			t.Errorf("trimPrefix(\"hello\", \"hello.txt\") = %q, want \".txt\"", result)
		}
	})

	t.Run("trimSuffix", func(t *testing.T) {
		// {{ "hello.txt" | trimSuffix ".txt" }} should produce "hello"
		result := trimSuffix(".txt", "hello.txt")
		if result != "hello" {
			t.Errorf("trimSuffix(\".txt\", \"hello.txt\") = %q, want \"hello\"", result)
		}
	})

	t.Run("trimLeft", func(t *testing.T) {
		// {{ "###hello" | trimLeft "#" }} should produce "hello"
		result := trimLeft("#", "###hello")
		if result != "hello" {
			t.Errorf("trimLeft(\"#\", \"###hello\") = %q, want \"hello\"", result)
		}
	})

	t.Run("trimRight", func(t *testing.T) {
		// {{ "hello###" | trimRight "#" }} should produce "hello"
		result := trimRight("#", "hello###")
		if result != "hello" {
			t.Errorf("trimRight(\"#\", \"hello###\") = %q, want \"hello\"", result)
		}
	})

	t.Run("trimChars", func(t *testing.T) {
		// {{ "###hello###" | trimChars "#" }} should produce "hello"
		result := trimChars("#", "###hello###")
		if result != "hello" {
			t.Errorf("trimChars(\"#\", \"###hello###\") = %q, want \"hello\"", result)
		}
	})

	t.Run("contains", func(t *testing.T) {
		// {{ "hello world" | contains "world" }} should be true
		if !contains("world", "hello world") {
			t.Error("contains(\"world\", \"hello world\") should be true")
		}
		if contains("xyz", "hello world") {
			t.Error("contains(\"xyz\", \"hello world\") should be false")
		}
	})

	t.Run("hasPrefix", func(t *testing.T) {
		// {{ "hello world" | hasPrefix "hello" }} should be true
		if !hasPrefix("hello", "hello world") {
			t.Error("hasPrefix(\"hello\", \"hello world\") should be true")
		}
		if hasPrefix("world", "hello world") {
			t.Error("hasPrefix(\"world\", \"hello world\") should be false")
		}
	})

	t.Run("hasSuffix", func(t *testing.T) {
		// {{ "hello world" | hasSuffix "world" }} should be true
		if !hasSuffix("world", "hello world") {
			t.Error("hasSuffix(\"world\", \"hello world\") should be true")
		}
		if hasSuffix("hello", "hello world") {
			t.Error("hasSuffix(\"hello\", \"hello world\") should be false")
		}
	})

	t.Run("split", func(t *testing.T) {
		// {{ "a,b,c" | split "," }} should produce ["a", "b", "c"]
		result := split(",", "a,b,c")
		if len(result) != 3 || result[0] != "a" || result[1] != "b" || result[2] != "c" {
			t.Errorf("split(\",\", \"a,b,c\") = %v, want [a b c]", result)
		}
	})

	t.Run("splitN", func(t *testing.T) {
		// {{ "a,b,c" | splitN "," 2 }} should produce ["a", "b,c"]
		result := splitN(",", 2, "a,b,c")
		if len(result) != 2 || result[0] != "a" || result[1] != "b,c" {
			t.Errorf("splitN(\",\", 2, \"a,b,c\") = %v, want [a b,c]", result)
		}
	})
}
