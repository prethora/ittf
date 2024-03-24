package commands

import (
	"github.com/spf13/cobra"
)

var inputFilePath string
var rulesFilePath string
var dateOutputFormat string

var RootCmd = &cobra.Command{
	Use:   "ittf -i <input file> -r <rules file>",
	Short: "ittf processes an input file based on given rules in an attempt to extract a formatted file name.\n\nOutputs the file name, or exits with status 1 if no rules were matched.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		execRoot()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		execVersion()
	},
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

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.AddCommand(validateCmd)
	RootCmd.AddCommand(versionCmd)

	RootCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Path to the input file (required)")
	RootCmd.Flags().StringVarP(&rulesFilePath, "rules", "r", "", "Path to the rules file (required)")
	RootCmd.Flags().StringVarP(&dateOutputFormat, "date-output", "d", "20060102", "Date output format")
	RootCmd.MarkFlagRequired("input")
	RootCmd.MarkFlagRequired("rules")
}
