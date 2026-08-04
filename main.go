package main

import "github.com/spf13/cobra"

func main() {
	rootCmd := &cobra.Command{
		Use:   "changes",
		Short: "changes manages changesets and generates CHANGELOG.md",
	}
	_ = rootCmd.Execute()
}
