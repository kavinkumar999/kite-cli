package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kavinkumar999/kite-cli/internal/client"
	"github.com/spf13/cobra"
)

var watchlistCmd = &cobra.Command{
	Use:     "watchlist",
	Aliases: []string{"w", "watch"},
	Short:   "View your watchlist with live quotes",
	Long: `Display instruments from your watchlist with current prices.

Note: Kite API doesn't have a direct watchlist endpoint.
This command shows quotes for instruments you specify.

Examples:
  kite watchlist ITC RELIANCE TCS
  kite watchlist NIFTY50 BANKNIFTY`,
	Args: cobra.MinimumNArgs(1),
	Run:  runWatchlist,
}

func init() {
	rootCmd.AddCommand(watchlistCmd)
}

func runWatchlist(cmd *cobra.Command, args []string) {
	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	instruments := make([]string, len(args))
	for i, arg := range args {
		instruments[i] = "NSE:" + arg
	}

	quotes, err := c.Kite().GetQuote(instruments...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching quotes: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SYMBOL\tLTP\tCHANGE\tCHANGE %\tOPEN\tHIGH\tLOW\tVOLUME")
	fmt.Fprintln(w, "------\t---\t------\t--------\t----\t----\t---\t------")

	for instrument, q := range quotes {
		changeSign := ""
		if q.NetChange > 0 {
			changeSign = "+"
		}

		fmt.Fprintf(w, "%s\t₹%.2f\t%s%.2f\t%s%.2f%%\t₹%.2f\t₹%.2f\t₹%.2f\t%d\n",
			instrument,
			q.LastPrice,
			changeSign, q.NetChange,
			changeSign, q.NetChange/q.OHLC.Open*100,
			q.OHLC.Open,
			q.OHLC.High,
			q.OHLC.Low,
			q.Volume,
		)
	}

	w.Flush()
}
