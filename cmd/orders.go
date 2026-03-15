package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kavinkumar999/kite-cli/internal/client"
	"github.com/spf13/cobra"
)

var ordersCmd = &cobra.Command{
	Use:     "orders",
	Aliases: []string{"o", "order"},
	Short:   "View your orders",
	Long:    `Display all orders placed today with their status.`,
	Run:     runOrders,
}

var cancelCmd = &cobra.Command{
	Use:   "cancel <order_id>",
	Short: "Cancel an order",
	Long:  `Cancel a pending order by its order ID.`,
	Args:  cobra.ExactArgs(1),
	Run:   runCancel,
}

func init() {
	rootCmd.AddCommand(ordersCmd)
	rootCmd.AddCommand(cancelCmd)
}

func runOrders(cmd *cobra.Command, args []string) {
	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	orders, err := c.Kite().GetOrders()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching orders: %v\n", err)
		os.Exit(1)
	}

	if len(orders) == 0 {
		fmt.Println("No orders found today.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ORDER ID\tTYPE\tSYMBOL\tQTY\tPRICE\tSTATUS\tTIME")
	fmt.Fprintln(w, "--------\t----\t------\t---\t-----\t------\t----")

	for _, o := range orders {
		price := "MARKET"
		if o.Price > 0 {
			price = fmt.Sprintf("₹%.2f", o.Price)
		}

		orderTime := ""
		if !o.OrderTimestamp.IsZero() {
			orderTime = o.OrderTimestamp.Format("15:04:05")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\t%s\n",
			o.OrderID,
			o.TransactionType,
			o.TradingSymbol,
			o.Quantity,
			price,
			o.Status,
			orderTime,
		)
	}

	w.Flush()
}

func runCancel(cmd *cobra.Command, args []string) {
	orderID := args[0]

	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	_, err = c.Kite().CancelOrder("regular", orderID, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cancelling order: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Order %s cancelled\n", orderID)
}
