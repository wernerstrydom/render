// Package funcs provides custom template functions for text manipulation.
package funcs

import (
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// Map returns a template.FuncMap with all custom functions.
func Map() template.FuncMap {
	return template.FuncMap{
		// Casing functions
		"lower":          strings.ToLower,
		"upper":          strings.ToUpper,
		"title":          toTitle,
		"camelCase":      toCamelCase,
		"pascalCase":     toPascalCase,
		"snakeCase":      toSnakeCase,
		"kebabCase":      toKebabCase,
		"upperSnakeCase": toUpperSnakeCase,
		"upperKebabCase": toUpperKebabCase,

		// Trimming functions
		"trim":       strings.TrimSpace,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"trimLeft":   strings.TrimLeft,
		"trimRight":  strings.TrimRight,
		"trimChars":  strings.Trim,

		// String manipulation
		"replace":   replace,
		"replaceN":  replaceN,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"repeat":    strings.Repeat,
		"reverse":   reverseString,
		"substr":    substring,
		"truncate":  truncate,
		"padLeft":   padLeft,
		"padRight":  padRight,
		"center":    center,
		"wrap":      wordWrap,
		"indent":    indent,
		"nindent":   nindent,

		// Splitting and joining
		"split":   strings.Split,
		"splitN":  strings.SplitN,
		"join":    join,
		"lines":   lines,
		"first":   first,
		"last":    last,
		"rest":    rest,
		"initial": initial,
		"nth":     nth,

		// Concatenation
		"concat": concat,
		"cat":    cat,

		// Conversion functions
		"toString":     toString,
		"toInt":        toInt,
		"toInt64":      toInt64,
		"toFloat":      toFloat,
		"toBool":       toBool,
		"toJson":       toJSON,
		"toPrettyJson": toPrettyJSON,
		"fromJson":     fromJSON,

		// Unicode normalization
		"nfc":   normalizeNFC,
		"nfd":   normalizeNFD,
		"nfkc":  normalizeNFKC,
		"nfkd":  normalizeNFKD,
		"ascii": toASCII,
		"slug":  toSlug,

		// Formatting
		"quote":  quote,
		"squote": squote,
		"printf": fmt.Sprintf,

		// Comparison and logic
		"eq":       eq,
		"ne":       ne,
		"lt":       lt,
		"le":       le,
		"gt":       gt,
		"ge":       ge,
		"and":      and,
		"or":       or,
		"not":      not,
		"default":  defaultVal,
		"empty":    empty,
		"coalesce": coalesce,
		"ternary":  ternary,

		// Collection functions
		"list":      list,
		"dict":      dict,
		"keys":      keys,
		"values":    values,
		"hasKey":    hasKey,
		"get":       get,
		"set":       set,
		"unset":     unset,
		"merge":     merge,
		"append":    appendList,
		"prepend":   prependList,
		"uniq":      uniq,
		"sortAlpha": sortAlpha,
		"len":       length,

		// Math functions
		"add":   add,
		"sub":   sub,
		"mul":   mul,
		"div":   div,
		"mod":   mod,
		"max":   max,
		"min":   min,
		"floor": floor,
		"ceil":  ceil,
		"round": round,

		// Regex functions
		"regexMatch":   regexMatch,
		"regexFind":    regexFind,
		"regexFindAll": regexFindAll,
		"regexReplace": regexReplace,
		"regexSplit":   regexSplit,

		// Counting
		"count":      count,
		"countWords": countWords,
		"countLines": countLines,
	}
}

// Casing functions

func toTitle(s string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(s))
}

func toCamelCase(s string) string {
	words := splitIntoWords(s)
	if len(words) == 0 {
		return ""
	}
	result := strings.ToLower(words[0])
	for _, word := range words[1:] {
		result += capitalizeFirst(word)
	}
	return result
}

func toPascalCase(s string) string {
	words := splitIntoWords(s)
	var result string
	for _, word := range words {
		result += capitalizeFirst(word)
	}
	return result
}

func toSnakeCase(s string) string {
	words := splitIntoWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "_")
}

func toKebabCase(s string) string {
	words := splitIntoWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "-")
}

func toUpperSnakeCase(s string) string {
	words := splitIntoWords(s)
	for i := range words {
		words[i] = strings.ToUpper(words[i])
	}
	return strings.Join(words, "_")
}

func toUpperKebabCase(s string) string {
	words := splitIntoWords(s)
	for i := range words {
		words[i] = strings.ToUpper(words[i])
	}
	return strings.Join(words, "-")
}

// splitIntoWords splits a string into words, handling various formats
func splitIntoWords(s string) []string {
	// First, normalize unicode
	s = normalizeNFC(s)

	// Handle existing separators (space, underscore, hyphen)
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	// Insert space before uppercase letters in camelCase/PascalCase
	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			// Insert space if previous is lowercase, or if we're at the start of a new word
			// in an acronym (e.g., "XMLParser" -> "XML Parser")
			if unicode.IsLower(prev) {
				result.WriteRune(' ')
			} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				result.WriteRune(' ')
			}
		}
		result.WriteRune(r)
	}

	// Split by spaces and filter empty strings
	parts := strings.Fields(result.String())
	var words []string
	for _, p := range parts {
		if p != "" {
			words = append(words, p)
		}
	}
	return words
}

func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// String manipulation functions

// replace replaces all occurrences of old with new in s.
// Arguments are ordered (old, new, s) for template pipeline compatibility.
func replace(old, new, s string) string {
	return strings.ReplaceAll(s, old, new)
}

// replaceN replaces n occurrences of old with new in s.
// Arguments are ordered (old, new, n, s) for template pipeline compatibility.
func replaceN(old, new string, n int, s string) string {
	return strings.Replace(s, old, new, n)
}

func reverseString(s string) string {
	runes := []rune(s)
	slices.Reverse(runes)
	return string(runes)
}

func substring(s string, start, end int) string {
	runes := []rune(s)
	if start < 0 {
		start = len(runes) + start
	}
	if end < 0 {
		end = len(runes) + end
	}
	if start < 0 {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}
	if start >= end {
		return ""
	}
	return string(runes[start:end])
}

func truncate(s string, length int) string {
	runes := []rune(s)
	if len(runes) <= length {
		return s
	}
	return string(runes[:length])
}

func padLeft(s string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	sRunes := []rune(s)
	sLen := len(sRunes)
	// If string is already >= target length, return as-is
	if sLen >= length {
		return s
	}
	// Calculate padding needed
	padRunes := []rune(pad)
	padLen := len(padRunes)
	needed := length - sLen
	// Build result efficiently
	result := make([]rune, length)
	// Fill padding from the start
	for i := 0; i < needed; i++ {
		result[i] = padRunes[i%padLen]
	}
	// Copy original string
	copy(result[needed:], sRunes)
	return string(result)
}

func padRight(s string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	runes := []rune(s)
	// If string is already >= target length, return as-is
	if len(runes) >= length {
		return s
	}
	// Pad to target length
	padRunes := []rune(pad)
	for len(runes) < length {
		runes = append(runes, padRunes...)
	}
	// Trim excess padding from the right
	return string(runes[:length])
}

func center(s string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	sLen := len([]rune(s))
	if sLen >= length {
		return s
	}
	leftPad := (length - sLen) / 2
	rightPad := length - sLen - leftPad
	return strings.Repeat(pad, leftPad) + s + strings.Repeat(pad, rightPad)
}

func wordWrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var result strings.Builder
	currentLine := 0
	words := strings.Fields(s)
	for _, word := range words {
		wordLen := len([]rune(word))
		if currentLine+wordLen > width && currentLine > 0 {
			result.WriteString("\n")
			currentLine = 0
		}
		if currentLine > 0 {
			result.WriteString(" ")
			currentLine++
		}
		result.WriteString(word)
		currentLine += wordLen
	}
	return result.String()
}

func indent(spaces int, s string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
}

func nindent(spaces int, s string) string {
	return "\n" + indent(spaces, s)
}

// Splitting and joining functions

func join(sep string, items any) string {
	switch v := items.(type) {
	case []string:
		return strings.Join(v, sep)
	case []any:
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = toString(item)
		}
		return strings.Join(strs, sep)
	default:
		return toString(items)
	}
}

func lines(s string) []string {
	return strings.Split(s, "\n")
}

func first(items any) any {
	switch v := items.(type) {
	case []any:
		if len(v) > 0 {
			return v[0]
		}
	case []string:
		if len(v) > 0 {
			return v[0]
		}
	case string:
		if len(v) > 0 {
			return string([]rune(v)[0])
		}
	}
	return nil
}

func last(items any) any {
	switch v := items.(type) {
	case []any:
		if len(v) > 0 {
			return v[len(v)-1]
		}
	case []string:
		if len(v) > 0 {
			return v[len(v)-1]
		}
	case string:
		runes := []rune(v)
		if len(runes) > 0 {
			return string(runes[len(runes)-1])
		}
	}
	return nil
}

func rest(items any) any {
	switch v := items.(type) {
	case []any:
		if len(v) > 1 {
			return v[1:]
		}
		return []any{}
	case []string:
		if len(v) > 1 {
			return v[1:]
		}
		return []string{}
	case string:
		runes := []rune(v)
		if len(runes) > 1 {
			return string(runes[1:])
		}
		return ""
	}
	return nil
}

func initial(items any) any {
	switch v := items.(type) {
	case []any:
		if len(v) > 1 {
			return v[:len(v)-1]
		}
		return []any{}
	case []string:
		if len(v) > 1 {
			return v[:len(v)-1]
		}
		return []string{}
	case string:
		runes := []rune(v)
		if len(runes) > 1 {
			return string(runes[:len(runes)-1])
		}
		return ""
	}
	return nil
}

func nth(n int, items any) any {
	switch v := items.(type) {
	case []any:
		if n >= 0 && n < len(v) {
			return v[n]
		}
	case []string:
		if n >= 0 && n < len(v) {
			return v[n]
		}
	case string:
		runes := []rune(v)
		if n >= 0 && n < len(runes) {
			return string(runes[n])
		}
	}
	return nil
}

// Concatenation functions

func concat(items ...string) string {
	return strings.Join(items, "")
}

func cat(items ...any) string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = toString(item)
	}
	return strings.Join(strs, " ")
}

// Conversion functions

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case nil:
		return ""
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toInt(v any) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int64:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %q to int: %w", val, err)
		}
		return i, nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

func toInt64(v any) (int64, error) {
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int64:
		return val, nil
	case float64:
		return int64(val), nil
	case string:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %q to int64: %w", val, err)
		}
		return i, nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

func toFloat(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %q to float: %w", val, err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float", v)
	}
}

func toBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("cannot convert %q to bool: %w", val, err)
		}
		return b, nil
	case int:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case float64:
		return val != 0, nil
	default:
		return v != nil, nil
	}
}

// Internal helper functions for math operations (panic on invalid input,
// which is caught by the template engine and returned as an error)
func mustFloat(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			panic(fmt.Sprintf("cannot convert %q to float: %v", val, err))
		}
		return f
	default:
		panic(fmt.Sprintf("cannot convert %T to float", v))
	}
}

func mustInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			panic(fmt.Sprintf("cannot convert %q to int: %v", val, err))
		}
		return i
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		panic(fmt.Sprintf("cannot convert %T to int", v))
	}
}

func mustBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			panic(fmt.Sprintf("cannot convert %q to bool: %v", val, err))
		}
		return b
	case int:
		return val != 0
	case int64:
		return val != 0
	case float64:
		return val != 0
	default:
		panic(fmt.Sprintf("cannot convert %T to bool", v))
	}
}

func toJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(b), nil
}

func toPrettyJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(b), nil
}

func fromJSON(s string) (any, error) {
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return v, nil
}

// Unicode normalization functions

func normalizeNFC(s string) string {
	return norm.NFC.String(s)
}

func normalizeNFD(s string) string {
	return norm.NFD.String(s)
}

func normalizeNFKC(s string) string {
	return norm.NFKC.String(s)
}

func normalizeNFKD(s string) string {
	return norm.NFKD.String(s)
}

func toASCII(s string) string {
	// Normalize to NFKD to decompose characters
	s = norm.NFKD.String(s)
	// Remove non-ASCII characters
	var result strings.Builder
	for _, r := range s {
		if r < 128 {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func toSlug(s string) string {
	// Normalize and convert to ASCII
	s = toASCII(s)
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace non-alphanumeric with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	// Trim hyphens from ends
	s = strings.Trim(s, "-")
	return s
}

// Formatting functions

func quote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

func squote(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `\'`) + `'`
}

// Comparison and logic functions

func eq(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func ne(a, b any) bool {
	return !reflect.DeepEqual(a, b)
}

func lt(a, b any) bool {
	return mustFloat(a) < mustFloat(b)
}

func le(a, b any) bool {
	return mustFloat(a) <= mustFloat(b)
}

func gt(a, b any) bool {
	return mustFloat(a) > mustFloat(b)
}

func ge(a, b any) bool {
	return mustFloat(a) >= mustFloat(b)
}

func and(a, b any) bool {
	return mustBool(a) && mustBool(b)
}

func or(a, b any) bool {
	return mustBool(a) || mustBool(b)
}

func not(a any) bool {
	return !mustBool(a)
}

func defaultVal(def, val any) any {
	if empty(val) {
		return def
	}
	return val
}

func empty(v any) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case []any:
		return len(val) == 0
	case []string:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	case bool:
		return !val
	case int:
		return val == 0
	case int64:
		return val == 0
	case float64:
		return val == 0
	default:
		return false
	}
}

func coalesce(items ...any) any {
	for _, item := range items {
		if !empty(item) {
			return item
		}
	}
	return nil
}

func ternary(cond bool, trueVal, falseVal any) any {
	if cond {
		return trueVal
	}
	return falseVal
}

// Collection functions

func list(items ...any) []any {
	return items
}

func dict(pairs ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(pairs); i += 2 {
		key := toString(pairs[i])
		m[key] = pairs[i+1]
	}
	return m
}

func keys(m any) []string {
	switch v := m.(type) {
	case map[string]any:
		return slices.Collect(maps.Keys(v))
	default:
		return nil
	}
}

func values(m any) []any {
	switch v := m.(type) {
	case map[string]any:
		return slices.Collect(maps.Values(v))
	default:
		return nil
	}
}

func hasKey(m any, key string) bool {
	switch v := m.(type) {
	case map[string]any:
		_, ok := v[key]
		return ok
	default:
		return false
	}
}

func get(m any, key string) any {
	switch v := m.(type) {
	case map[string]any:
		return v[key]
	default:
		return nil
	}
}

func set(m any, key string, val any) map[string]any {
	switch v := m.(type) {
	case map[string]any:
		// Create a copy to avoid modifying the input map
		result := maps.Clone(v)
		result[key] = val
		return result
	default:
		return map[string]any{key: val}
	}
}

func unset(m any, key string) map[string]any {
	switch v := m.(type) {
	case map[string]any:
		// Create a copy to avoid modifying the input map
		result := make(map[string]any, len(v))
		for k, val := range v {
			if k != key {
				result[k] = val
			}
		}
		return result
	default:
		return map[string]any{}
	}
}

func merge(inputMaps ...any) map[string]any {
	result := make(map[string]any)
	for _, m := range inputMaps {
		if v, ok := m.(map[string]any); ok {
			maps.Copy(result, v)
		}
	}
	return result
}

func appendList(items any, val any) []any {
	switch v := items.(type) {
	case []any:
		return append(v, val)
	default:
		return []any{items, val}
	}
}

func prependList(items any, val any) []any {
	switch v := items.(type) {
	case []any:
		return append([]any{val}, v...)
	default:
		return []any{val, items}
	}
}

func uniq(items any) []any {
	switch v := items.(type) {
	case []any:
		seen := make(map[string]bool)
		result := make([]any, 0)
		for _, item := range v {
			key := fmt.Sprintf("%v", item)
			if !seen[key] {
				seen[key] = true
				result = append(result, item)
			}
		}
		return result
	case []string:
		seen := make(map[string]bool)
		result := make([]any, 0)
		for _, item := range v {
			if !seen[item] {
				seen[item] = true
				result = append(result, item)
			}
		}
		return result
	default:
		return []any{items}
	}
}

func sortAlpha(items any) []string {
	var strs []string
	switch v := items.(type) {
	case []any:
		for _, item := range v {
			strs = append(strs, toString(item))
		}
	case []string:
		// Make a copy to avoid modifying the input
		strs = make([]string, len(v))
		copy(strs, v)
	default:
		return nil
	}
	slices.Sort(strs)
	return strs
}

func length(items any) int {
	switch v := items.(type) {
	case string:
		return len([]rune(v))
	case []any:
		return len(v)
	case []string:
		return len(v)
	case map[string]any:
		return len(v)
	default:
		return 0
	}
}

// Math functions

func add(a, b any) any {
	af, bf := mustFloat(a), mustFloat(b)
	result := af + bf
	if result == float64(int(result)) {
		return int(result)
	}
	return result
}

func sub(a, b any) any {
	af, bf := mustFloat(a), mustFloat(b)
	result := af - bf
	if result == float64(int(result)) {
		return int(result)
	}
	return result
}

func mul(a, b any) any {
	af, bf := mustFloat(a), mustFloat(b)
	result := af * bf
	if result == float64(int(result)) {
		return int(result)
	}
	return result
}

func div(a, b any) (any, error) {
	af, bf := mustFloat(a), mustFloat(b)
	if bf == 0 {
		return 0, fmt.Errorf("division by zero in div operation")
	}
	result := af / bf
	if result == float64(int(result)) {
		return int(result), nil
	}
	return result, nil
}

func mod(a, b any) (int, error) {
	bi := mustInt(b)
	if bi == 0 {
		return 0, fmt.Errorf("division by zero in mod operation")
	}
	return mustInt(a) % bi, nil
}

func max(items ...any) any {
	if len(items) == 0 {
		return nil
	}
	maxVal := mustFloat(items[0])
	for _, item := range items[1:] {
		if v := mustFloat(item); v > maxVal {
			maxVal = v
		}
	}
	if maxVal == float64(int(maxVal)) {
		return int(maxVal)
	}
	return maxVal
}

func min(items ...any) any {
	if len(items) == 0 {
		return nil
	}
	minVal := mustFloat(items[0])
	for _, item := range items[1:] {
		if v := mustFloat(item); v < minVal {
			minVal = v
		}
	}
	if minVal == float64(int(minVal)) {
		return int(minVal)
	}
	return minVal
}

func floor(v any) int {
	return int(math.Floor(mustFloat(v)))
}

func ceil(v any) int {
	return int(math.Ceil(mustFloat(v)))
}

func round(v any) int {
	return int(math.Round(mustFloat(v)))
}

// Regex functions
// These functions return (result, error) to properly surface regex compilation errors
// Go's text/template will catch the error and report it to the user

func regexMatch(pattern, s string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return re.MatchString(s), nil
}

func regexFind(pattern, s string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return re.FindString(s), nil
}

func regexFindAll(pattern, s string, n int) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return re.FindAllString(s, n), nil
}

func regexReplace(pattern, replacement, s string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return re.ReplaceAllString(s, replacement), nil
}

func regexSplit(pattern, s string, n int) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return re.Split(s, n), nil
}

// Counting functions

func count(substr, s string) int {
	return strings.Count(s, substr)
}

func countWords(s string) int {
	return len(strings.Fields(s))
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}
