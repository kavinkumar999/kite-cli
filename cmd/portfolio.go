package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kavinkumar999/kite-cli/internal/client"
	"github.com/spf13/cobra"
)

var portfolioCmd = &cobra.Command{
	Use:     "portfolio",
	Aliases: []string{"pos", "positions"},
	Short:   "View your positions (intraday & overnight)",
	Long:    `Display all open positions with current P&L.`,
	Run:     runPortfolio,
}

func init() {
	rootCmd.AddCommand(portfolioCmd)
}

func runPortfolio(cmd *cobra.Command, args []string) {
	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	positions, err := c.Kite().GetPositions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching positions: %v\n", err)
		os.Exit(1)
	}

	if len(positions.Net) == 0 {
		fmt.Println("No open positions.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SYMBOL\tPRODUCT\tQTY\tAVG PRICE\tLTP\tP&L")
	fmt.Fprintln(w, "------\t-------\t---\t---------\t---\t---")

	var totalPnl float64

	for _, p := range positions.Net {
		if p.Quantity == 0 {
			continue
		}

		pnlSign := ""
		if p.PnL > 0 {
			pnlSign = "+"
		}
		totalPnl += p.PnL

		fmt.Fprintf(w, "%s\t%s\t%d\t₹%.2f\t₹%.2f\t%s₹%.2f\n",
			p.Tradingsymbol,
			p.Product,
			p.Quantity,
			p.AveragePrice,
			p.LastPrice,
			pnlSign, p.PnL,
		)
	}

	w.Flush()

	pnlSign := ""
	if totalPnl > 0 {
		pnlSign = "+"
	}
	fmt.Printf("\nTotal P&L: %s₹%.2f\n", pnlSign, totalPnl)
}
