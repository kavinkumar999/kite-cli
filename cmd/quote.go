package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kavinkumar999/kite-cli/internal/client"
	"github.com/spf13/cobra"
)

var quoteExchange string

var quoteCmd = &cobra.Command{
	Use:     "quote <symbol>",
	Aliases: []string{"q", "price"},
	Short:   "Get quick quote for a symbol",
	Long: `Get the current price and details for a symbol.

Examples:
  kite quote ITC
  kite quote RELIANCE -e BSE`,
	Args: cobra.ExactArgs(1),
	Run:  runQuote,
}

func init() {
	rootCmd.AddCommand(quoteCmd)
	quoteCmd.Flags().StringVarP(&quoteExchange, "exchange", "e", "NSE", "Exchange: NSE, BSE, NFO, MCX")
}

func runQuote(cmd *cobra.Command, args []string) {
	symbol := strings.ToUpper(args[0])
	instrument := quoteExchange + ":" + symbol

	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	quotes, err := c.Kite().GetQuote(instrument)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching quote: %v\n", err)
		os.Exit(1)
	}

	q, ok := quotes[instrument]
	if !ok {
		fmt.Fprintf(os.Stderr, "Symbol not found: %s\n", instrument)
		os.Exit(1)
	}

	changeSign := ""
	if q.NetChange > 0 {
		changeSign = "+"
	}
	changePercent := q.NetChange / q.OHLC.Open * 100

	fmt.Printf("%s (%s)\n", symbol, quoteExchange)
	fmt.Println(strings.Repeat("─", 30))
	fmt.Printf("  LTP:     ₹%.2f (%s%.2f / %s%.2f%%)\n", q.LastPrice, changeSign, q.NetChange, changeSign, changePercent)
	fmt.Printf("  Open:    ₹%.2f\n", q.OHLC.Open)
	fmt.Printf("  High:    ₹%.2f\n", q.OHLC.High)
	fmt.Printf("  Low:     ₹%.2f\n", q.OHLC.Low)
	fmt.Printf("  Close:   ₹%.2f (prev)\n", q.OHLC.Close)
	fmt.Printf("  Volume:  %d\n", q.Volume)

	if q.Depth.Buy[0].Price > 0 {
		fmt.Println()
		fmt.Printf("  Bid:     ₹%.2f x %d\n", q.Depth.Buy[0].Price, q.Depth.Buy[0].Quantity)
	}
	if q.Depth.Sell[0].Price > 0 {
		fmt.Printf("  Ask:     ₹%.2f x %d\n", q.Depth.Sell[0].Price, q.Depth.Sell[0].Quantity)
	}
}
