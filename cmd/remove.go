package cmd

import (
	"fmt"
	"os"

	"github.com/kavinkumar999/kite-cli/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove an account",
	Long: `Remove a configured Kite account.

Examples:
  kite remove bob    # Remove account 'bob'`,
	Args: cobra.ExactArgs(1),
	Run:  runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) {
	alias := args[0]

	if err := config.RemoveAccount(alias); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Account '%s' removed\n", alias)
}
