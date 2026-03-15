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
	buyPrice      float64
	buyTrigger    float64
	buyOrderType  string
	buyProduct    string
	buyExchange   string
	buyValidity   string
	buyDryRun     bool
)

var buyCmd = &cobra.Command{
	Use:   "buy <symbol> <quantity>",
	Short: "Buy stocks",
	Long: `Place a buy order for the specified stock.

Examples:
  kite buy ITC 10                    # Market order for 10 shares of ITC (NSE)
  kite buy ITC 10 --price 450        # Limit order at ₹450
  kite buy ITC 10 -e BSE             # Buy from BSE instead of NSE
  kite buy RELIANCE 5 --product MIS  # Intraday order
  kite buy NIFTY24MARFUT 50          # Buy futures`,
	Args: cobra.ExactArgs(2),
	Run:  runBuy,
}

func init() {
	rootCmd.AddCommand(buyCmd)

	buyCmd.Flags().Float64VarP(&buyPrice, "price", "p", 0, "Limit price (omit for market order)")
	buyCmd.Flags().Float64VarP(&buyTrigger, "trigger", "t", 0, "Trigger price for SL orders")
	buyCmd.Flags().StringVarP(&buyOrderType, "type", "T", "", "Order type: MARKET, LIMIT, SL, SL-M (auto-detected)")
	buyCmd.Flags().StringVarP(&buyProduct, "product", "P", "CNC", "Product type: CNC, MIS, NRML")
	buyCmd.Flags().StringVarP(&buyExchange, "exchange", "e", "NSE", "Exchange: NSE, BSE, NFO, MCX")
	buyCmd.Flags().StringVarP(&buyValidity, "validity", "v", "DAY", "Validity: DAY, IOC, TTL")
	buyCmd.Flags().BoolVar(&buyDryRun, "dry-run", false, "Simulate order without placing (for testing)")
}

func runBuy(cmd *cobra.Command, args []string) {
	symbol := strings.ToUpper(args[0])
	quantity, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid quantity: %s\n", args[1])
		os.Exit(1)
	}

	orderType := determineOrderType(buyPrice, buyTrigger, buyOrderType)

	if buyDryRun {
		fmt.Println("🧪 DRY RUN (no order placed)")
		fmt.Printf("   Action:   BUY\n")
		fmt.Printf("   Symbol:   %s:%s\n", buyExchange, symbol)
		fmt.Printf("   Quantity: %d\n", quantity)
		fmt.Printf("   Type:     %s\n", orderType)
		fmt.Printf("   Product:  %s\n", buyProduct)
		if buyPrice > 0 {
			fmt.Printf("   Price:    ₹%.2f\n", buyPrice)
		}
		if buyTrigger > 0 {
			fmt.Printf("   Trigger:  ₹%.2f\n", buyTrigger)
		}
		return
	}

	c, err := client.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	orderParams := kiteconnect.OrderParams{
		Exchange:        buyExchange,
		Tradingsymbol:   symbol,
		Quantity:        quantity,
		TransactionType: kiteconnect.TransactionTypeBuy,
		Product:         buyProduct,
		OrderType:       orderType,
		Validity:        buyValidity,
	}

	if buyPrice > 0 {
		orderParams.Price = buyPrice
	}
	if buyTrigger > 0 {
		orderParams.TriggerPrice = buyTrigger
	}

	resp, err := c.Kite().PlaceOrder(kiteconnect.VarietyRegular, orderParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Order failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ BUY order placed: %s x %d\n", symbol, quantity)
	fmt.Printf("  Order ID: %s\n", resp.OrderID)
	if buyPrice > 0 {
		fmt.Printf("  Price: ₹%.2f (%s)\n", buyPrice, orderType)
	} else {
		fmt.Printf("  Type: %s\n", orderType)
	}
}

func determineOrderType(price, trigger float64, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if trigger > 0 && price > 0 {
		return kiteconnect.OrderTypeSL
	}
	if trigger > 0 {
		return kiteconnect.OrderTypeSLM
	}
	if price > 0 {
		return kiteconnect.OrderTypeLimit
	}
	return kiteconnect.OrderTypeMarket
}
