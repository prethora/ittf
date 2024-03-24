# ittf (Invoice Text To Filename)

## Usage

```
ittf processes an input file based on given rules in an attempt to extract a formatted file name.

Outputs the file name, or exits with status 1 if no rules were matched.

Usage:
  ittf -i <input file> -r <rules file> [flags]
  ittf [command]

Available Commands:
  help        Help about any command
  validate    Validate a rules file
  version     Print the version number

Flags:
  -d, --date-output string   Date output format (default "20060102")
  -h, --help                 help for ittf
  -i, --input string         Path to the input file (required)
  -r, --rules string         Path to the rules file (required)
```

## Rules file format

A rules file is expected the be a YAML file in the following format:

```yaml
- vendor: 
  - vendor regexp 1
  - vendor regexp 2
  - ...
  date: date regexp
  dateFormat: a golang time format
  fileName: (date) - (0).pdf            # this is optional, see the section below to understand what is used by default.
- ...
- ...
```

## RegExp Index

Any regexp in the rules file can be optionally prefixed with a zero-based index in the format: `{index}/{regexp}`. For example:

```yaml
- vendor:
  - 1/Google (Cloud|Sevices|API)
  - ...
```

This does 2 things:

* Sets a requirement for the regexp "Google (Cloud|Sevices|API)" to match not only once, but twice (1 being the second element in a zero-based index)
* Extracts the second occurence as the value to be used if the match is referenced in the fileName format.

This can also be used for the date regexp, and can be useful if the invoice contains several dates in the same format and you want to be able to control which one is extracted.

If the index is omitted (meaning you just type in a normal regexp), the index defaults to 0 (requiring only one match).

## Date Format

The dateFormat field for a rule is required to tell ittf how to parse the matched date string into a recognized time. In golang, a time format is defined as how you would represent a specific chosen time. The time in question is:

`Mon Jan 02 15:04:05 MST 2006`

So if you want to specify a format as `MM/DD/YYYY` you actually have to specify it as `01/02/2006` - in other words, you need to specify it as how you would display January 2, 2006 in the format you want.

So for each rule, you need to set a dateFormat according to the date you are extracting.

## File Name Format

The fileName field for a rule is optional and defaults to 