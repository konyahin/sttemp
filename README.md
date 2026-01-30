# sttemp

A simple CLI tool for processing text templates with variable substitution from environment variables or user input.

Variables are resolved in the following order:
1. Environment variables. If a variable exists in the environment, its value is used.
2. Interactive prompt. If not found in environment, you'll be prompted to enter a value.

## Installation
You can build it from source, or use `go install`

```sh
go install github.com/konyahin/sttemp@latest
```

## Usage

```sh
sttemp [options] [template-name ...]
```

### Shell integration
If you want to get autocomplete, copy appropriate line into your shell settings.

```sh
# zsh
autoload -Uz compinit
compinit
compdef '_values "sttemp options" $(sttemp -l 2>/dev/null)' sttemp

# bash
complete -W "$(sttemp -l 2>/dev/null)" sttemp

# ksh
set -A complete_sttemp -- $(sttemp -l)
```

`fzf` integration can look like this

```sh
fst () {
    templates="$(sttemp -l)"
    selected=$(echo "$templates" | fzf)
    [ -n "$selected" ] && sttemp "$selected"
}
```

### Options
- `-C <path>` custom template directory (default: `~/.local/share/sttemp`)
- `-o <file>` output to file instead of stdout
- `-d` use template's subdirectory name as output filename
- `-h` show short help
- `--no-input` use only environment variables (do not ask user for substitution value); exit with error if some variable is missing
- `--edit` edit selected template in your console `$EDITOR`
- `-l` list all templates names

### Examples
```sh
sttemp                                  # list all templates
sttemp greeting                         # process template `greeting`, output to stdout
sttemp -o out.txt greeting              # save to file `out.txt`
sttemp -d mit                           # save as `LICENSE`, if `mit` in `LICENSE` subfolder (see files structure below)
sttemp --edit mit                       # open file with `mit` template in `$EDITOR`
export NAME="Alice" && sttemp greeting  # use environment variables
```

## Template Syntax

Use `{VARIABLE}` for placeholders. Variables are resolved from environment or prompted interactively. To include literal `{VARIABLE}` text in your template without substitution, escape it with a backslash as this `\{VARIABLE}`.

### Template example
```
Hello, {FIRST NAME}!
Server: {HOST}:{PORT}
Escape literals: \{NOT_A_VAR}
```

## Templates Organization
Store templates in subdirectories for auto-naming with `-d`:
```
~/.local/share/sttemp/
├── greeting
└── LICENSE/
    ├── mit   # `sttemp -d mit` creates file "LICENSE"
    └── GPLv3 # `sttemp -d GPLv3` also creates file "LICENSE"
```
