package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kavinkumar999/kite-cli/internal/client"
	"github.com/spf13/cobra"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var (
	sellPrice     float64
	sellTrigger   float64
	sellOrderType string
	sellProduct   string
	sellExchange  string
	sellValidity  string
	sellDryRun    bool
)

var sellCmd = &cobra.Command{
	Use:   "sell <symbol> <quantity>",
	Short: "Sell stocks",
	Long: `Place a sell order for the specified stock.

Examples:
  kite sell ITC 10                    # Market order to sell 10 shares
  kite sell ITC 10 --price 450        # Limit order at ₹450
  kite sell RELIANCE 5 --product MIS  # Square off intraday position`,
	Args: cobra.ExactArgs(2),
	Run:  runSell,
}

func init() {
	rootCmd.AddCommand(sellCmd)

	sellCmd.Flags().Float64VarP(&sellPrice, "price", "p", 0, "Limit price (omit for market order)")
	sellCmd.Flags().Float64VarP(&sellTrigger, "trigger", "t", 0, "Trigger price for SL orders")
	sellCmd.Flags().StringVarP(&sellOrderType, "type", "T", "", "Order type: MARKET, LIMIT, SL, SL-M")
	sellCmd.Flags().StringVarP(&sellProduct, "product", "P", "CNC", "Product type: CNC, MIS, NRML")
	sellCmd.Flags().StringVarP(&sellExchange, "exchange", "e", "NSE", "Exchange: NSE, BSE, NFO, MCX")
	sellCmd.Flags().StringVarP(&sellValidity, "validity", "v", "DAY", "Validity: DAY, IOC, TTL")
	sellCmd.Flags().BoolVar(&sellDryRun, "dry-run", false, "Simulate order without placing (for testing)")
}

func runSell(cmd *cobra.Command, args []string) {
	symbol := strings.ToUpper(args[0])
	quantity, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid quantity: %s\n", args[1])
		os.Exit(1)
	}

	orderType := determineOrderType(sellPrice, sellTrigger, sellOrderType)

	if sellDryRun {
		fmt.Println("🧪 DRY RUN (no order placed)")
		fmt.Printf("   Action:   SELL\n")
		fmt.Printf("   Symbol:   %s:%s\n", sellExchange, symbol)
		fmt.Printf("   Quantity: %d\n", quantity)
		fmt.Printf("   Type:     %s\n", orderType)
		fmt.Printf("   Product:  %s\n", sellProduct)
		if sellPrice > 0 {
			fmt.Printf("   Price:    ₹%.2f\n", sellPrice)
		}
		if sellTrigger > 0 {
			fmt.Printf("   Trigger:  ₹%.2f\n", sellTrigger)
		}
		return
	}

	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	orderParams := kiteconnect.OrderParams{
		Exchange:        sellExchange,
		Tradingsymbol:   symbol,
		Quantity:        quantity,
		TransactionType: kiteconnect.TransactionTypeSell,
		Product:         sellProduct,
		OrderType:       orderType,
		Validity:        sellValidity,
	}

	if sellPrice > 0 {
		orderParams.Price = sellPrice
	}
	if sellTrigger > 0 {
		orderParams.TriggerPrice = sellTrigger
	}

	resp, err := c.Kite().PlaceOrder(kiteconnect.VarietyRegular, orderParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Order failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ SELL order placed: %s x %d\n", symbol, quantity)
	fmt.Printf("  Order ID: %s\n", resp.OrderID)
	if sellPrice > 0 {
		fmt.Printf("  Price: ₹%.2f (%s)\n", sellPrice, orderType)
	} else {
		fmt.Printf("  Type: %s\n", orderType)
	}
}
