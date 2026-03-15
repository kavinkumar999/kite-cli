package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kavinkumar999/kite-cli/internal/client"
	"github.com/spf13/cobra"
)

var marginsCmd = &cobra.Command{
	Use:     "margins",
	Aliases: []string{"m", "margin", "funds"},
	Short:   "View available margins/funds",
	Long:    `Display your available margins for equity and commodity segments.`,
	Run:     runMargins,
}

func init() {
	rootCmd.AddCommand(marginsCmd)
}

func runMargins(cmd *cobra.Command, args []string) {
	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	margins, err := c.Kite().GetUserMargins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching margins: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println("═══ EQUITY ═══")
	fmt.Fprintln(w, "Available\tUsed\tOpening Balance")
	fmt.Fprintf(w, "₹%.2f\t₹%.2f\t₹%.2f\n",
		margins.Equity.Available.Cash,
		margins.Equity.Used.Debits,
		margins.Equity.Available.OpeningBalance,
	)
	w.Flush()

	fmt.Println()
	fmt.Println("─── Details ───")
	fmt.Printf("  Cash:           ₹%.2f\n", margins.Equity.Available.Cash)
	fmt.Printf("  Collateral:     ₹%.2f\n", margins.Equity.Available.Collateral)
	fmt.Printf("  Intraday Payin: ₹%.2f\n", margins.Equity.Available.IntradayPayin)

	if margins.Commodity.Available.Cash > 0 || margins.Commodity.Used.Debits > 0 {
		fmt.Println()
		fmt.Println("═══ COMMODITY ═══")
		fmt.Printf("  Available: ₹%.2f\n", margins.Commodity.Available.Cash)
		fmt.Printf("  Used:      ₹%.2f\n", margins.Commodity.Used.Debits)
	}
}
