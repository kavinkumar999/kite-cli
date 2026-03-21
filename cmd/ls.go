package cmd

import (
	"fmt"
	"os"

	"github.com/kavinkumar999/kite-cli/internal/config"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all connected accounts",
	Long: `List all configured Kite accounts.

The current account is marked with (current).

Examples:
  kite ls    # List all accounts`,
	Run: runLs,
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func runLs(cmd *cobra.Command, args []string) {
	aliases, currentAccount, err := config.ListAccounts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(aliases) == 0 {
		fmt.Println("No accounts configured. Run 'kite auth' to add an account.")
		return
	}

	fmt.Println("Connected accounts:")
	fmt.Println()
	for _, alias := range aliases {
		account, _ := config.GetAccount(alias)
		if alias == currentAccount {
			if account != nil && account.UserName != "" {
				fmt.Printf("  ● %s (%s) - current\n", alias, account.UserName)
			} else {
				fmt.Printf("  ● %s - current\n", alias)
			}
		} else {
			if account != nil && account.UserName != "" {
				fmt.Printf("    %s (%s)\n", alias, account.UserName)
			} else {
				fmt.Printf("    %s\n", alias)
			}
		}
	}
	fmt.Println()
	fmt.Println("Use 'kite use <alias>' to switch accounts")
}
