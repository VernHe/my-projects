package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(wordCmd) // word子命令
	rootCmd.AddCommand(timeCmd) // time子命令
}
