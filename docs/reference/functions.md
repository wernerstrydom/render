# Template Functions Reference

render provides 80+ template functions in addition to Go's built-in functions.

## Casing Functions

Transform string case.

| Function | Description | Example |
|----------|-------------|---------|
| `lower` | Lowercase | `{{ "Hello" \| lower }}` → `hello` |
| `upper` | Uppercase | `{{ "Hello" \| upper }}` → `HELLO` |
| `title` | Title Case | `{{ "hello world" \| title }}` → `Hello World` |
| `camelCase` | camelCase | `{{ "hello world" \| camelCase }}` → `helloWorld` |
| `pascalCase` | PascalCase | `{{ "hello world" \| pascalCase }}` → `HelloWorld` |
| `snakeCase` | snake_case | `{{ "HelloWorld" \| snakeCase }}` → `hello_world` |
| `kebabCase` | kebab-case | `{{ "HelloWorld" \| kebabCase }}` → `hello-world` |
| `upperSnakeCase` | UPPER_SNAKE | `{{ "helloWorld" \| upperSnakeCase }}` → `HELLO_WORLD` |
| `upperKebabCase` | UPPER-KEBAB | `{{ "helloWorld" \| upperKebabCase }}` → `HELLO-WORLD` |

## Trimming Functions

Remove characters from strings.

| Function | Signature | Description |
|----------|-----------|-------------|
| `trim` | `trim s` | Remove leading/trailing whitespace |
| `trimPrefix` | `trimPrefix prefix s` | Remove prefix |
| `trimSuffix` | `trimSuffix suffix s` | Remove suffix |
| `trimLeft` | `trimLeft cutset s` | Remove leading chars in cutset |
| `trimRight` | `trimRight cutset s` | Remove trailing chars in cutset |
| `trimChars` | `trimChars cutset s` | Remove leading and trailing chars in cutset |

Examples:
```
{{ "  hello  " | trim }}                    → "hello"
{{ "hello.txt" | trimSuffix ".txt" }}       → "hello"
{{ "###hello###" | trimChars "#" }}         → "hello"
```

## String Manipulation

| Function | Signature | Description |
|----------|-----------|-------------|
| `replace` | `replace old new s` | Replace all occurrences |
| `replaceN` | `replaceN old new n s` | Replace first n occurrences |
| `contains` | `contains substr s` | Check if contains substring |
| `hasPrefix` | `hasPrefix prefix s` | Check prefix |
| `hasSuffix` | `hasSuffix suffix s` | Check suffix |
| `repeat` | `repeat s n` | Repeat string n times |
| `reverse` | `reverse s` | Reverse string |
| `substr` | `substr s start end` | Substring (supports negative indices) |
| `truncate` | `truncate s length` | Truncate to length |
| `padLeft` | `padLeft s length pad` | Pad left to length |
| `padRight` | `padRight s length pad` | Pad right to length |
| `center` | `center s length pad` | Center with padding |
| `wrap` | `wrap s width` | Word wrap at width |
| `indent` | `indent spaces s` | Indent with spaces |
| `nindent` | `nindent spaces s` | Newline + indent |

Examples:
```
{{ "hello" | replace "l" "L" }}             → "heLLo"
{{ "hello" | substr 0 3 }}                  → "hel"
{{ "hi" | padLeft 5 "0" }}                  → "000hi"
{{ "yaml:\n  key: value" | indent 2 }}      → "  yaml:\n    key: value"
```

## Splitting and Joining

| Function | Signature | Description |
|----------|-----------|-------------|
| `split` | `split sep s` | Split string by separator |
| `splitN` | `splitN sep n s` | Split into at most n parts |
| `join` | `join sep items` | Join array with separator |
| `lines` | `lines s` | Split into lines |
| `first` | `first items` | First element |
| `last` | `last items` | Last element |
| `rest` | `rest items` | All but first |
| `initial` | `initial items` | All but last |
| `nth` | `nth n items` | Nth element (0-indexed) |

Examples:
```
{{ "a,b,c" | split "," }}                   → ["a", "b", "c"]
{{ list "a" "b" "c" | join "-" }}           → "a-b-c"
{{ list 1 2 3 | first }}                    → 1
{{ list 1 2 3 | rest }}                     → [2, 3]
```

## Concatenation

| Function | Signature | Description |
|----------|-----------|-------------|
| `concat` | `concat items...` | Concatenate strings |
| `cat` | `cat items...` | Concatenate with spaces |

Examples:
```
{{ concat "hello" "world" }}                → "helloworld"
{{ cat "hello" "world" }}                   → "hello world"
```

## Conversion Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `toString` | `toString v` | Convert to string |
| `toInt` | `toInt v` | Convert to int |
| `toInt64` | `toInt64 v` | Convert to int64 |
| `toFloat` | `toFloat v` | Convert to float64 |
| `toBool` | `toBool v` | Convert to bool |
| `toJson` | `toJson v` | Convert to JSON string |
| `toPrettyJson` | `toPrettyJson v` | Convert to pretty JSON |
| `fromJson` | `fromJson s` | Parse JSON string |

Examples:
```
{{ 42 | toString }}                         → "42"
{{ "123" | toInt }}                         → 123
{{ .config | toJson }}                      → '{"key":"value"}'
{{ .config | toPrettyJson }}                → formatted JSON
```

## Unicode Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `nfc` | `nfc s` | NFC normalization |
| `nfd` | `nfd s` | NFD normalization |
| `nfkc` | `nfkc s` | NFKC normalization |
| `nfkd` | `nfkd s` | NFKD normalization |
| `ascii` | `ascii s` | Convert to ASCII |
| `slug` | `slug s` | URL-safe slug |

Examples:
```
{{ "café" | ascii }}                        → "cafe"
{{ "Hello World!" | slug }}                 → "hello-world"
```

## Formatting

| Function | Signature | Description |
|----------|-----------|-------------|
| `quote` | `quote s` | Double-quote string |
| `squote` | `squote s` | Single-quote string |
| `printf` | `printf format args...` | Formatted string |

Examples:
```
{{ "hello" | quote }}                       → '"hello"'
{{ printf "%s-%d" "item" 42 }}              → "item-42"
```

## Comparison and Logic

| Function | Signature | Description |
|----------|-----------|-------------|
| `eq` | `eq a b` | Equal (deep comparison) |
| `ne` | `ne a b` | Not equal |
| `lt` | `lt a b` | Less than |
| `le` | `le a b` | Less than or equal |
| `gt` | `gt a b` | Greater than |
| `ge` | `ge a b` | Greater than or equal |
| `and` | `and a b` | Boolean AND |
| `or` | `or a b` | Boolean OR |
| `not` | `not a` | Boolean NOT |
| `default` | `default def val` | Default if empty |
| `empty` | `empty v` | Check if empty |
| `coalesce` | `coalesce items...` | First non-empty |
| `ternary` | `ternary cond true false` | Conditional value |

Examples:
```
{{ if eq .status "active" }}active{{ end }}
{{ .name | default "Anonymous" }}
{{ .value | empty }}                        → true/false
{{ ternary .enabled "on" "off" }}
```

## Collection Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `list` | `list items...` | Create a list |
| `dict` | `dict pairs...` | Create a dictionary |
| `keys` | `keys m` | Get map keys |
| `values` | `values m` | Get map values |
| `hasKey` | `hasKey m key` | Check if key exists |
| `get` | `get m key` | Get value by key |
| `set` | `set m key val` | Set value (returns new map) |
| `unset` | `unset m key` | Remove key (returns new map) |
| `merge` | `merge maps...` | Merge maps |
| `append` | `append items val` | Append to list |
| `prepend` | `prepend items val` | Prepend to list |
| `uniq` | `uniq items` | Remove duplicates |
| `sortAlpha` | `sortAlpha items` | Sort alphabetically |
| `len` | `len items` | Length |

Examples:
```
{{ list 1 2 3 }}                            → [1, 2, 3]
{{ dict "a" 1 "b" 2 }}                      → {"a": 1, "b": 2}
{{ .config | keys }}                        → list of keys
{{ hasKey .config "debug" }}                → true/false
{{ .items | sortAlpha }}                    → sorted list
```

## Math Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `add` | `add a b` | Addition |
| `sub` | `sub a b` | Subtraction |
| `mul` | `mul a b` | Multiplication |
| `div` | `div a b` | Division |
| `mod` | `mod a b` | Modulo |
| `max` | `max items...` | Maximum |
| `min` | `min items...` | Minimum |
| `floor` | `floor v` | Floor |
| `ceil` | `ceil v` | Ceiling |
| `round` | `round v` | Round |

Examples:
```
{{ add 1 2 }}                               → 3
{{ mul 3 4 }}                               → 12
{{ div 10 3 }}                              → 3
{{ max 1 5 3 }}                             → 5
{{ 3.7 | floor }}                           → 3
```

## Regex Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `regexMatch` | `regexMatch pattern s` | Check if matches |
| `regexFind` | `regexFind pattern s` | Find first match |
| `regexFindAll` | `regexFindAll pattern s n` | Find all matches |
| `regexReplace` | `regexReplace pattern repl s` | Replace matches |
| `regexSplit` | `regexSplit pattern s n` | Split by pattern |

Examples:
```
{{ regexMatch "^[a-z]+$" "hello" }}         → true
{{ regexFind "[0-9]+" "abc123def" }}        → "123"
{{ regexReplace "[0-9]" "X" "a1b2c3" }}     → "aXbXcX"
```

## Counting Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `count` | `count substr s` | Count occurrences |
| `countWords` | `countWords s` | Count words |
| `countLines` | `countLines s` | Count lines |

Examples:
```
{{ "hello" | count "l" }}                   → 2
{{ "hello world" | countWords }}            → 2
{{ "a\nb\nc" | countLines }}                → 3
```

## Pipeline Usage

Functions can be chained using pipes:

```
{{ .name | lower | replace " " "_" | truncate 20 }}
{{ .items | sortAlpha | join ", " }}
{{ .config | toPrettyJson | indent 2 }}
```

## Error Handling

Some functions return errors for invalid input:

```
{{ div 10 0 }}                              → error: division by zero
{{ "abc" | toInt }}                         → error: cannot convert "abc" to int
{{ regexMatch "[invalid" "test" }}          → error: invalid regex pattern
```

Use conditionals or `default` to handle potential errors:

```
{{ if .value }}{{ .value | toInt }}{{ else }}0{{ end }}
```
