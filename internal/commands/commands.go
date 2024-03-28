package commands

import (
	"github.com/spf13/cobra"
)

var inputFilePath string
var rulesFilePath string
var dateOutputFormat string
var useSubgroupMatch bool

var RootCmd = &cobra.Command{
	Use:   "ittf -f <file> -r <rules file>",
	Short: "ittf processes an input file based on given rules in an attempt to extract a formatted file name.\n\nOutputs the file name, or exits with a non-zero exit code otherwise. See README.md file for error codes.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		execRoot()
	},
	Version: "1.0.3",
}

var validateCmd = &cobra.Command{
	Use:   "validate <rules file>",
	Short: "Validate a rules file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rulesFilePath = args[0]
		execValidate()
	},
}

var aliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Output a list of configured aliases",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		execAliases()
	},
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.AddCommand(validateCmd)
	RootCmd.AddCommand(aliasesCmd)

	RootCmd.Flags().StringVarP(&inputFilePath, "file", "f", "", "Path to the input file (required)")
	RootCmd.Flags().StringVarP(&rulesFilePath, "rules", "r", "", "Path to the rules file (required)")
	RootCmd.Flags().StringVarP(&dateOutputFormat, "date-output", "d", "20060102", "Date output format")
	RootCmd.Flags().BoolVarP(&useSubgroupMatch, "subgroup-match", "s", false, "Use single sub-group match extraction")
	RootCmd.MarkFlagRequired("input")
	RootCmd.MarkFlagRequired("rules")
}
