# Kite CLI

A blazing fast command-line interface for Zerodha Kite. Execute trades, view portfolio, check margins - all from your terminal.

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash
```

### Manual Installation

If you prefer to install manually without the script:

```bash
# 1. Clone the repository
git clone https://github.com/kavinkumar999/kite-cli.git
cd kite-cli

# 2. Build the binary
go build -ldflags="-s -w" -o kite .

# 3. Create bin directory and move binary
mkdir -p ~/bin
mv kite ~/bin/

# 4. Add to PATH (add this to ~/.zshrc or ~/.bashrc)
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc

# 5. Reload shell
source ~/.zshrc

# 6. Verify installation
kite --help
```

### Shell Completion (Optional)

Enable tab autocompletion for commands and flags:

```bash
# For Zsh
mkdir -p ~/.zsh/completions
kite completion zsh > ~/.zsh/completions/_kite
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit; compinit' >> ~/.zshrc
source ~/.zshrc

# For Bash
kite completion bash > /etc/bash_completion.d/kite
```

## Setup

### 1. Get API Credentials

1. Go to [Kite Connect Developer Console](https://developers.kite.trade/)
2. Create a new app (costs ₹2000/month)
3. Note your **API Key** and **API Secret**

### 2. Configure CLI

**Option A: Interactive (stores in ~/.kite.yaml)**
```bash
kite auth
# Enter your API Key and API Secret when prompted
```

**Option B: Environment Variables (more secure)**
```bash
# Add to ~/.zshrc or ~/.bashrc
export KITE_API_KEY="your_api_key"
export KITE_API_SECRET="your_api_secret"
export KITE_ACCESS_TOKEN="your_access_token"
```

**Option C: Manual config file**
```bash
# Create ~/.kite.yaml manually
cat > ~/.kite.yaml << EOF
api_key: your_api_key
api_secret: your_api_secret
access_token: your_access_token
EOF

# Set secure permissions
chmod 600 ~/.kite.yaml
```

### 3. Login (Required Daily)

```bash
kite login
# 1. Open the URL in browser
# 2. Login with Zerodha credentials + OTP
# 3. Copy request_token from redirect URL
# 4. Paste it in terminal
```

> **Note**: Access token expires daily at ~6 AM IST. Run `kite login` each trading day.

## Usage

### Trading

```bash
# Buy stocks
kite buy ITC 10                    # Market order for 10 shares
kite buy ITC 10 -p 450             # Limit order at ₹450
kite buy ITC 10 --product MIS      # Intraday order
kite buy RELIANCE 5 -e BSE         # Buy from BSE

# Sell stocks
kite sell ITC 10                   # Market order
kite sell ITC 10 -p 460            # Limit order

# Test without placing real orders
kite buy ITC 10 --dry-run

# Cancel orders
kite cancel <order_id>
```

### Portfolio & Holdings

```bash
kite holdings                      # View demat holdings
kite portfolio                     # View open positions (alias: pos)
kite orders                        # View today's orders
kite margins                       # View available funds (alias: m)
```

### Market Data (Requires Quote API subscription)

```bash
kite quote ITC                     # Quick price check
kite quote RELIANCE -e BSE         # Quote from BSE
kite watchlist ITC RELIANCE TCS    # Multiple quotes
```

## Order Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--price` | `-p` | Limit price | Market order |
| `--trigger` | `-t` | Trigger price (SL orders) | - |
| `--product` | `-P` | CNC, MIS, NRML | CNC |
| `--exchange` | `-e` | NSE, BSE, NFO, MCX | NSE |
| `--validity` | `-v` | DAY, IOC | DAY |
| `--dry-run` | - | Test without placing order | false |

## Examples

```bash
# Quick intraday trade
kite buy SBIN 100 --product MIS

# Limit order with specific price
kite buy INFY 50 -p 1450

# Stop-loss order
kite sell TATASTEEL 25 -t 120 -p 119

# Test order (no real trade)
kite buy ITC 10 -p 400 --dry-run

# Check position and square off
kite portfolio
kite sell SBIN 100 --product MIS
```

## Configuration

### Config File Location

```
~/.kite.yaml
```

### Config Priority

1. **Environment variables** (highest priority)
2. **~/.kite.yaml** (fallback)

### Environment Variables

| Variable | Description |
|----------|-------------|
| `KITE_API_KEY` | Your Kite Connect API key |
| `KITE_API_SECRET` | Your Kite Connect API secret |
| `KITE_ACCESS_TOKEN` | Session access token |

## Update

To update to the latest version:

```bash
# Option 1: Re-run install script
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash

# Option 2: Manual update
cd /path/to/kite-cli
git pull origin main
go build -ldflags="-s -w" -o kite .
cp kite ~/bin/
```

## Uninstall

To completely remove Kite CLI:

```bash
# 1. Remove binary
rm ~/bin/kite

# 2. Remove config file (contains your credentials)
rm ~/.kite.yaml

# 3. Remove shell completions
rm ~/.zsh/completions/_kite 2>/dev/null

# 4. Remove PATH entry from ~/.zshrc (optional)
# Edit ~/.zshrc and remove the line: export PATH="$HOME/bin:$PATH"

# 5. Remove cloned repository (if exists)
rm -rf /path/to/kite-cli
```

### One-liner Uninstall

```bash
rm -f ~/bin/kite ~/.kite.yaml ~/.zsh/completions/_kite
```

> **Warning**: This will delete your saved credentials. Back up `~/.kite.yaml` if needed.

## Security

- Config file permissions are set to `600` (owner read/write only)
- Credentials can be stored in environment variables instead of file
- Never commit `~/.kite.yaml` to git

## Troubleshooting

### "api_key not configured"
Run `kite auth` or set `KITE_API_KEY` environment variable.

### "Invalid access token"
Access token expired. Run `kite login` to get a new one.

### "Insufficient permission"
Your Kite Connect app doesn't have Quote API access. Contact Zerodha.

### Command not found: kite
Add `~/bin` to your PATH:
```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## License

MIT
