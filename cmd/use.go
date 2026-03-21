package cmd

import (
	"fmt"
	"os"

	"github.com/kavinkumar999/kite-cli/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <alias>",
	Short: "Switch to a different account",
	Long: `Switch to a different Kite account.

Examples:
  kite use alice  # Switch to account 'alice'
  kite use bob    # Switch to account 'bob'`,
	Args: cobra.ExactArgs(1),
	Run:  runUse,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(cmd *cobra.Command, args []string) {
	alias := args[0]

	if err := config.SetCurrentAccount(alias); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	account, _ := config.GetAccount(alias)
	if account != nil && account.UserName != "" {
		fmt.Printf("✓ Switched to account: %s (%s)\n", alias, account.UserName)
	} else {
		fmt.Printf("✓ Switched to account: %s\n", alias)
	}
}
