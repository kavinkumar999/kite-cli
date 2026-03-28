---
layout: home

hero:
  name: Kite CLI
  text: Manage family trading accounts from your terminal
  tagline: Switch between multiple Zerodha accounts instantly. Execute trades for yourself, spouse, parents - all from one terminal.
  actions:
    - theme: brand
      text: Get Started
      link: '#installation'
    - theme: alt
      text: View on GitHub
      link: https://github.com/kavinkumar999/kite-cli

features:
  - icon: 👨‍👩‍👧‍👦
    title: Family Account Management
    details: Add unlimited family accounts - spouse, parents, children. Each account with its own alias for easy identification.
  - icon: 🔄
    title: Instant Account Switching
    details: Switch between accounts with a single command. No logout/login needed. Execute trades for any family member in seconds.
  - icon: ⚡
    title: Lightning Fast Trading
    details: Place orders instantly without opening a browser. Buy, sell, check holdings - all from your terminal.
  - icon: 🛡️
    title: Secure & Local
    details: All credentials stored locally with strict file permissions. No cloud storage. Your data stays on your machine.
---

## Why Kite CLI?

Managing trading accounts for your entire family? Tired of logging in and out of different Zerodha accounts? Kite CLI lets you:

- **Add all family accounts once** - Set up accounts for yourself, spouse, parents, kids
- **Switch instantly** - One command to switch between any account
- **Trade for anyone** - Execute orders for any family member without browser hassle
- **View all portfolios** - Check holdings and positions across accounts quickly

```bash
# Morning routine - login all family accounts
kite use self && kite login
kite use spouse && kite login  
kite use dad && kite login

# Execute trades for different family members
kite use self && kite buy RELIANCE 10
kite use spouse && kite buy TCS 5
kite use dad && kite sell INFY 20
```

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash
```

### Manual Installation

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

::: code-group

```bash [Zsh]
mkdir -p ~/.zsh/completions
kite completion zsh > ~/.zsh/completions/_kite
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit; compinit' >> ~/.zshrc
source ~/.zshrc
```

```bash [Bash]
kite completion bash > /etc/bash_completion.d/kite
```

:::

## Setup

### 1. Get API Credentials

Each family member needs their own Kite Connect app:

1. Go to [Kite Connect Developer Console](https://developers.kite.trade/)
2. Login with the Zerodha account
3. Create a new app
4. Note the **API Key** and **API Secret**

::: tip Family Setup
Repeat this for each family member's Zerodha account. Each person needs their own API credentials.
:::

### 2. Add Family Accounts

Add each family member with a memorable alias:

```bash
# Add your account
kite auth
# Enter account name: self
# Enter your API Key: xxx
# Enter your API Secret: xxx

# Add spouse's account
kite auth
# Enter account name: spouse
# Enter your API Key: yyy
# Enter your API Secret: yyy

# Add parent's account
kite auth
# Enter account name: dad
# Enter your API Key: zzz
# Enter your API Secret: zzz
```

### 3. Daily Login

Login is required once per day for each account (token expires at ~6 AM IST):

```bash
# Login to each account
kite use self && kite login
kite use spouse && kite login
kite use dad && kite login
```

::: warning Daily Login Required
Access token expires daily at ~6 AM IST. Run `kite login` for each account every trading day.
:::

## Managing Family Accounts

### View All Accounts

```bash
kite ls
# Connected accounts:
#
#   ● self (Kavin Kumar) - current
#     spouse (Priya Kumar)
#     dad (Raj Kumar)
#     mom (Lakshmi Kumar)
#
# Use 'kite use <alias>' to switch accounts
```

### Switch Between Accounts

```bash
kite use spouse
# ✓ Switched to account: spouse (Priya Kumar)

kite use dad
# ✓ Switched to account: dad (Raj Kumar)
```

### Remove an Account

```bash
kite remove mom
# ✓ Account 'mom' removed
```

## Trading

### Execute Orders

```bash
# Switch to account and trade
kite use self
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

### Quick Family Trading

```bash
# Buy same stock for multiple family members
kite use self && kite buy RELIANCE 10
kite use spouse && kite buy RELIANCE 10
kite use dad && kite buy RELIANCE 5

# Check everyone's holdings
kite use self && kite holdings
kite use spouse && kite holdings
kite use dad && kite holdings
```

## Portfolio & Holdings

```bash
kite holdings                      # View demat holdings
kite portfolio                     # View open positions (alias: pos)
kite orders                        # View today's orders
kite margins                       # View available funds (alias: m)
```

## Market Data

::: info Quote API Subscription Required
Market data commands require Quote API access from Zerodha.
:::

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

## Configuration

### Config File Location

```
~/.kite.yaml
```

### Config File Format

```yaml
current_account: self
accounts:
  self:
    api_key: your_api_key
    api_secret: your_api_secret
    access_token: your_access_token
    user_name: Kavin Kumar
  spouse:
    api_key: spouse_api_key
    api_secret: spouse_api_secret
    access_token: spouse_access_token
    user_name: Priya Kumar
  dad:
    api_key: dad_api_key
    api_secret: dad_api_secret
    access_token: dad_access_token
    user_name: Raj Kumar
```

### Environment Variables

Environment variables take highest priority and override the config file:

| Variable | Description |
|----------|-------------|
| `KITE_API_KEY` | Your Kite Connect API key |
| `KITE_API_SECRET` | Your Kite Connect API secret |
| `KITE_ACCESS_TOKEN` | Session access token |

## Update

```bash
# Quick update
kite update

# Check for updates only
kite update --check

# Manual update
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash
```

## Uninstall

### Quick Uninstall

```bash
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/uninstall.sh | bash
```

### Manual Uninstall

```bash
rm -f ~/bin/kite ~/.kite.yaml ~/.zsh/completions/_kite
```

::: danger Warning
This will delete your saved credentials for all family accounts. Back up `~/.kite.yaml` if needed.
:::

## Security

- Config file permissions are set to `600` (owner read/write only)
- All credentials stored locally - no cloud sync
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

[MIT License](https://github.com/kavinkumar999/kite-cli/blob/main/LICENSE)
