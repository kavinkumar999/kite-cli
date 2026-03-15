package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kavinkumar999/kite-cli/internal/client"
	"github.com/spf13/cobra"
)

var holdingsCmd = &cobra.Command{
	Use:     "holdings",
	Aliases: []string{"h", "hold"},
	Short:   "View your holdings",
	Long:    `Display all stocks in your demat account with current P&L.`,
	Run:     runHoldings,
}

func init() {
	rootCmd.AddCommand(holdingsCmd)
}

func runHoldings(cmd *cobra.Command, args []string) {
	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	holdings, err := c.Kite().GetHoldings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching holdings: %v\n", err)
		os.Exit(1)
	}

	if len(holdings) == 0 {
		fmt.Println("No holdings found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SYMBOL\tQTY\tAVG PRICE\tLTP\tP&L\tP&L %")
	fmt.Fprintln(w, "------\t---\t---------\t---\t---\t-----")

	var totalInvested, totalCurrent float64

	for _, h := range holdings {
		invested := h.AveragePrice * float64(h.Quantity)
		current := h.LastPrice * float64(h.Quantity)
		pnl := current - invested
		pnlPercent := (pnl / invested) * 100

		totalInvested += invested
		totalCurrent += current

		pnlSign := ""
		if pnl > 0 {
			pnlSign = "+"
		}

		fmt.Fprintf(w, "%s\t%d\t₹%.2f\t₹%.2f\t%s₹%.2f\t%s%.2f%%\n",
			h.Tradingsymbol,
			h.Quantity,
			h.AveragePrice,
			h.LastPrice,
			pnlSign, pnl,
			pnlSign, pnlPercent,
		)
	}

	w.Flush()

	totalPnl := totalCurrent - totalInvested
	totalPnlPercent := (totalPnl / totalInvested) * 100
	pnlSign := ""
	if totalPnl > 0 {
		pnlSign = "+"
	}

	fmt.Println()
	fmt.Printf("Invested: ₹%.2f | Current: ₹%.2f | P&L: %s₹%.2f (%s%.2f%%)\n",
		totalInvested, totalCurrent, pnlSign, totalPnl, pnlSign, totalPnlPercent)
}
