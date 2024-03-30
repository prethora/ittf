package aliases

var Aliases = map[string]Alias{
	"DD/M/YYYY": {
		Match:  `\d+/\w+/\d+`,
		Format: `02/Jan/2006`,
	},
	"MMMM DD, YYYY": {
		Match:  `\w+\s+\d+,\s+\d+`,
		Format: `January 2, 2006`,
	},
	"DD M YYYY": {
		Match:  `\d+\s+\w+\s+\d+`,
		Format: `02 Jan 2006`,
	},
	"DDth MMMM YYYY": {
		Match:  `\d+\w+\s+\w+\s+\d+`,
		Format: `02 January 2006`,
	},
	"Dth MMMM YY": {
		Match:  `\d+\w+\s+\w+\s+\d+`,
		Format: `2 January 06`,
	},
	"DD M YY": {
		Match:  `\d+\s+\w+\s+\d+`,
		Format: `02 Jan 06`,
	},
	"D/M/YYYY": {
		Match:  `\d+/\w+/\d+`,
		Format: `2/Jan/2006`,
	},
	"Dth M YY": {
		Match:  `\d+\w+\s+\w+\s+\d+`,
		Format: `2 Jan 06`,
	},
	"YY-MM-DD": {
		Match:  `\d+-\d+-\d+`,
		Format: `06-01-02`,
	},
	"YYYY-MM-DD": {
		Match:  `\d+-\d+-\d+`,
		Format: `2006-01-02`,
	},
}
