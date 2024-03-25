# ittf (Invoice Text To Filename)

## Installation

```bash
go install github.com/prethora/ittf@latest
```

## Usage

```
ittf processes an input file based on given rules in an attempt to extract a formatted file 
name.

Outputs the file name, or exits with a non-zero exit code otherwise. See README.md file for 
error codes.

Usage:
  ittf -f <file> -r <rules file> [flags]
  ittf [command]

Available Commands:
  help        Help about any command
  validate    Validate a rules file

Flags:
  -d, --date-output string   Date output format (default "20060102")
  -f, --file string          Path to the input file (required)
  -h, --help                 help for ittf
  -r, --rules string         Path to the rules file (required)
  -s, --subgroup-match       Use single sub-group match extraction
  -v, --version              version for ittf
```

### Error Codes

* 0 - (success) - matched a rule, and has output a filename
* 1 - Error: Input file does not exist or is not readable 
* 2 - Error: Rules file does not exist or is not readable 
* 3 - Error: File matched one or more rules, but could not parse date
* 4 - Error: File did not match any rules
* 5 - Error: Provided date output format is invalid

## Rules file format

A rules file is expected to be a YAML file in the following format:

```yaml
- vendor: 
    - vendor query 1
    - vendor query 2
    - ...
  date: date query
  dateFormat: a golang time format
  fileName: (date) - (1).pdf        # this is optional, and if omitted defaults to this value
- ...
- ...
```

## Rules File Validation

You can use the following command to validate a rules file:

```bash
ittf validate <rules file>
```

This will output any errors if present, or a success message with a rule count. Useful while creating/editing a rules file.

## Query Format

A query in its simplest form is just a regexp string. However, a query can also be a map with the following fields:

```yaml
match: regexp         # required
index: 1              # optional - defaults to 1
after: query          # optional
before: query         # optional
```

This is a powerful setup which gives you a lot of flexibility on how to define what you are looking for in an invoice.

### index

The `index` field is a one-based positive or negative index which both sets a requirement for how many times the regexp should match (at least) and which of the matches to extract. For example:

* `index: 3` means at least 3 matches should exist, and the third one should be extracted
* `index: -1` means at least 1 match should exist, and the last one should be extracted
* `index: -2` means at least 2 matches should exist, and the one before last should be extracted

The `index` field is optional, and defaults to `1`

### after

The `after` field is optional and can itself be a query (meaning an after query can be a regexp string, but also a map with an optional index and even its own before/after queries etc). The point of the `after` field is to find a string within the input string which acts as an opening boundary. In other words, the query's regexp will only apply to the text that comes AFTER the boundary. If omitted, the opening boundary is naturally the beginning of the input string. If the `after` query is set and fails to match, the whole query fails.

### before

The `before` field is optional and can itself be a query (meaning a before query can be a regexp string, but also a map with an optional index and even its own before/after query etc). The point of the `before` field is to find a string within the input string which acts as a closing boundary. In other words, the query's regexp will only apply to the text that comes BEFORE the boundary. If omitted, the closing boundary is naturally the end of the input string. If the `before` query is set and fails to match, the whole query fails. 

### Boundaries

The `before` field applies to the whole input string and is not affected by the `after` field and vice versa. Both boundary fields can be used conjointly to apply the query's regexp to an inner region of the input string. If both boundaries are used together and the `before` matched string comes before or overlaps with the `after` matched string, the query fails.

## RegExp Single Subgroup Matching

The string matched by a vendor or date query is extracted and can be used in the file name output (through configuring the fileName field). By default, the string used is the string matched by the query's regexp as a whole. If you want to use a specific part of the match instead, you can add the -s flag to the command and use a single subgroup in the regexp.

For example:

```yaml
- vendor:
    - for (Google Cloud) Services
    - ...
  ...
  fileName: (date) - (1).pdf  
```

The first vendor query regexp here will match "for Google Cloud Services" in the invoice text input but the file name output might be `20240425 - Google Cloud.pdf`. This gives you some flexibility in situations where you want to more liberally use the regexp to identify the invoice source, but also have fine-grained control over what you use in the file name output.

Don't forget though, this feature is not enabled by default, you need to add the -s flag when running the command to use it. Also, when using this feature, if you need to use brackets `()` to logically group things make sure to use `(?:)` instead to not interfere with the subgroup matching. 

**This feature will only work if a single subgroup match is included in the regexp.**

## Date Format

The `dateFormat` field for a rule is required to tell ittf how to parse the matched date string into a recognized time. In golang, a time format is defined as how you would represent a specific language-defined time. The time in question is:

`Mon Jan 02 15:04:05 MST 2006`

So if you want to specify a format as `MM/DD/YYYY` you actually have to specify it as `01/02/2006` - in other words, you need to specify it as how you would display January 2, 2006 in the format you want.

So for each rule, you need to set a `dateFormat` according to the date you are extracting.

Golang supports 2-digit years, so if you need to support such an archaic invoice, instead of using `2006` in your `dateFormat`, you would simply use `06`. 

An invalid date such a Feb 31, 2024 might match but it will cause the rule to fail.

## File Name Format

The fileName field for a rule is optional and defaults to `(date) - (1).pdf`. If you are explicitly specifying it however, the `(date)` placeholder is required. You can optionally use one-based index placeholders within brackets to insert the matched string for a vendor query. So `(3) - (date).pdf` might result in an output of `Google Cloud - 20240325.pdf` if the third vendor query matched to `Google Cloud`.

## Output Date Format

The date that is inserted into the file name is by default in the format `20060102`. This can however be customized if necessary using the -d flag. For example:

```bash
ittf -f <file> -r <rules file> -d "01-02-2006"
```

This might result in an output such as `03-25-2024 - Google Cloud.pdf`

## Development

Clone the repo:

```bash
git clone github.com/prethora/ittf
cd ittf
```

Run the project:

```bash
go run .
```

Build the project:

```bash
go build
```

Install from source:

```bash
go install
```