package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kavinkumar999/kite-cli/internal/config"
	"github.com/spf13/cobra"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Set up authentication with Kite Connect",
	Long: `Configure your Kite Connect API credentials for an account.

You need to:
1. Create an app at https://developers.kite.trade/
2. Get your API key and API secret
3. Run 'kite auth' and follow the prompts

You can add multiple accounts with different names.`,
	Run: runAuth,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login and get access token",
	Long: `Complete the login flow to get an access token for the current account.

1. Opens the Kite login URL
2. After login, you'll get a request token in the redirect URL
3. Enter that token to complete authentication

Use 'kite use <alias>' to switch to a different account before logging in.`,
	Run: runLogin,
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(loginCmd)
}

func runAuth(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter account name: ")
	alias, _ := reader.ReadString('\n')
	alias = strings.TrimSpace(alias)

	if alias == "" {
		fmt.Fprintln(os.Stderr, "Error: Account name cannot be empty")
		os.Exit(1)
	}

	// Check if account already exists
	if config.AccountExists(alias) {
		fmt.Printf("Account '%s' already exists. Do you want to overwrite? (y/N): ", alias)
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))
		if confirm != "y" && confirm != "yes" {
			fmt.Println("Cancelled.")
			return
		}
	}

	fmt.Print("Enter your API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	fmt.Print("Enter your API Secret: ")
	apiSecret, _ := reader.ReadString('\n')
	apiSecret = strings.TrimSpace(apiSecret)

	account := &config.Account{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	if err := config.SaveAccount(alias, account); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Account '%s' saved to %s\n", alias, config.GetConfigPath())
	fmt.Printf("\nNext step: Run 'kite login' to get your access token\n")
}

func runLogin(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cfg.APIKey == "" || cfg.APISecret == "" {
		fmt.Fprintln(os.Stderr, "API credentials not configured. Run 'kite auth' first.")
		os.Exit(1)
	}

	alias := cfg.Alias
	if alias == "" {
		alias = "default"
	}

	fmt.Printf("Logging in as: %s\n\n", alias)

	kc := kiteconnect.New(cfg.APIKey)
	loginURL := kc.GetLoginURL()

	fmt.Println("Open this URL in your browser to login:")
	fmt.Printf("\n  %s\n\n", loginURL)
	fmt.Println("After login, you'll be redirected to your redirect URL with a 'request_token' parameter.")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the request_token from the redirect URL: ")
	requestToken, _ := reader.ReadString('\n')
	requestToken = strings.TrimSpace(requestToken)

	data, err := kc.GenerateSession(requestToken, cfg.APISecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating session: %v\n", err)
		os.Exit(1)
	}

	// Update access token for the current account
	if err := config.UpdateAccessToken(alias, data.AccessToken, data.UserName); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving access token: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Login successful! Welcome %s\n", data.UserName)
	fmt.Println("You can now use kite commands like 'kite buy', 'kite holdings', etc.")
}
