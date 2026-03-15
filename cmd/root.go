package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kite",
	Short: "Kite CLI - Fast trading from your terminal",
	Long: `Kite CLI is a blazing fast command-line interface for Zerodha Kite.

Execute trades, view portfolio, check margins - all from your terminal.

Examples:
  kite buy ITC 10              # Buy 10 shares of ITC at market price
  kite buy ITC 10 --price 450  # Buy 10 shares of ITC at ₹450
  kite sell RELIANCE 5         # Sell 5 shares of RELIANCE
  kite holdings                # View your holdings
  kite portfolio               # View your portfolio
  kite margins                 # View available margins
  kite watchlist               # View your watchlist`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kite.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kite")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// Config file found and loaded
	}
}
