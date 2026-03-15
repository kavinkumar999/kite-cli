package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kavinkumar999/kite-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Set up authentication with Kite Connect",
	Long: `Configure your Kite Connect API credentials.

You need to:
1. Create an app at https://developers.kite.trade/
2. Get your API key and API secret
3. Run 'kite auth' and follow the prompts`,
	Run: runAuth,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login and get access token",
	Long: `Complete the login flow to get an access token.

1. Opens the Kite login URL
2. After login, you'll get a request token in the redirect URL
3. Enter that token to complete authentication`,
	Run: runLogin,
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(loginCmd)
}

func runAuth(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	fmt.Print("Enter your API Secret: ")
	apiSecret, _ := reader.ReadString('\n')
	apiSecret = strings.TrimSpace(apiSecret)

	cfg := &config.Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Credentials saved to %s\n", config.GetConfigPath())
	fmt.Println("\nNext step: Run 'kite login' to get your access token")
}

func runLogin(cmd *cobra.Command, args []string) {
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		fmt.Fprintln(os.Stderr, "API credentials not configured. Run 'kite auth' first.")
		os.Exit(1)
	}

	kc := kiteconnect.New(apiKey)
	loginURL := kc.GetLoginURL()

	fmt.Println("Open this URL in your browser to login:")
	fmt.Printf("\n  %s\n\n", loginURL)
	fmt.Println("After login, you'll be redirected to your redirect URL with a 'request_token' parameter.")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the request_token from the redirect URL: ")
	requestToken, _ := reader.ReadString('\n')
	requestToken = strings.TrimSpace(requestToken)

	data, err := kc.GenerateSession(requestToken, apiSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating session: %v\n", err)
		os.Exit(1)
	}

	viper.Set("access_token", data.AccessToken)
	if err := viper.WriteConfig(); err != nil {
		home, _ := os.UserHomeDir()
		viper.WriteConfigAs(home + "/.kite.yaml")
	}

	fmt.Printf("\n✓ Login successful! Welcome %s\n", data.UserName)
	fmt.Println("You can now use kite commands like 'kite buy', 'kite holdings', etc.")
}
