# Kite CLI

A blazing fast command-line interface for Zerodha Kite. Execute trades, view portfolio, check margins - all from your terminal.

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash
```

## Setup

### 1. Get API Credentials

1. Go to [Kite Connect Developer Console](https://developers.kite.trade/)
2. Create a new app
3. Note your **API Key** and **API Secret**

### 2. Configure CLI

```bash
kite auth
# Enter your API Key and API Secret when prompted
```

### 3. Login

```bash
kite login
# Opens login URL, complete login, and enter the request token
```

## Usage

### Trading

```bash
# Buy stocks
kite buy ITC 10                    # Market order for 10 shares
kite buy ITC 10 --price 450        # Limit order at ₹450
kite buy ITC 10 --product MIS      # Intraday order
kite buy RELIANCE 5 -e BSE         # Buy from BSE

# Sell stocks
kite sell ITC 10                   # Market order
kite sell ITC 10 --price 460       # Limit order

# Cancel orders
kite cancel <order_id>
```

### Portfolio & Holdings

```bash
kite holdings                      # View demat holdings
kite portfolio                     # View open positions
kite orders                        # View today's orders
```

### Market Data

```bash
kite quote ITC                     # Quick price check
kite quote RELIANCE -e BSE         # Quote from BSE
kite watchlist ITC RELIANCE TCS    # Multiple quotes
```

### Funds

```bash
kite margins                       # View available funds
```

## Order Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--price` | `-p` | Limit price |
| `--trigger` | `-t` | Trigger price (SL orders) |
| `--product` | `-P` | CNC (delivery), MIS (intraday), NRML (F&O) |
| `--exchange` | `-e` | NSE, BSE, NFO, MCX |
| `--validity` | `-v` | DAY, IOC |

## Examples

```bash
# Quick intraday trade
kite buy SBIN 100 --product MIS

# Limit order with specific price
kite buy INFY 50 --price 1450

# Stop-loss order
kite sell TATASTEEL 25 --trigger 120 --price 119

# Check position and square off
kite portfolio
kite sell SBIN 100 --product MIS
```

## Configuration

Config is stored in `~/.kite.yaml`:

```yaml
api_key: your_api_key
api_secret: your_api_secret
access_token: your_access_token
```

## Building from Source

```bash
go build -o kite .

# With optimizations
go build -ldflags="-s -w" -o kite .
```

## License

MIT
