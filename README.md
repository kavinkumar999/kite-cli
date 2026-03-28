# Kite CLI

Manage multiple Zerodha Kite accounts from your terminal. Perfect for families - switch between accounts instantly and execute trades for yourself, spouse, parents, all from one place.

## Features

- **Family Account Management** - Add unlimited accounts with custom aliases (self, spouse, dad, mom)
- **Instant Switching** - Switch between accounts with `kite use <alias>`
- **LLM & AI Ready** - Connect with Claude, ChatGPT, or any AI assistant to trade using natural language
- **Fast Trading** - Execute orders without opening a browser
- **Secure** - All credentials stored locally with strict file permissions

## Quick Start

```bash
# Install
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash

# Add family accounts
kite auth  # Add your account as 'self'
kite auth  # Add spouse's account as 'spouse'
kite auth  # Add parent's account as 'dad'

# Switch and trade
kite use self && kite buy RELIANCE 10
kite use spouse && kite buy TCS 5
```

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash
```

### Manual Installation

```bash
git clone https://github.com/kavinkumar999/kite-cli.git
cd kite-cli
go build -ldflags="-s -w" -o kite .
mkdir -p ~/bin && mv kite ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Shell Completion (Optional)

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

Each family member needs their own Kite Connect app:

1. Go to [Kite Connect Developer Console](https://developers.kite.trade/)
2. Login with the family member's Zerodha account
3. Create a new app
4. Note the **API Key** and **API Secret**

### 2. Add Family Accounts

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

```bash
kite use self && kite login
kite use spouse && kite login
kite use dad && kite login
```

> **Note**: Access token expires daily at ~6 AM IST. Run `kite login` for each account every trading day.

## Managing Accounts

### List All Accounts

```bash
kite ls
# Connected accounts:
#
#   ● self (Kavin Kumar) - current
#     spouse (Priya Kumar)
#     dad (Raj Kumar)
#
# Use 'kite use <alias>' to switch accounts
```

### Switch Accounts

```bash
kite use spouse
# ✓ Switched to account: spouse (Priya Kumar)
```

### Remove Account

```bash
kite remove dad
# ✓ Account 'dad' removed
```

## Trading

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

### Quick Family Trading

```bash
# Buy same stock for multiple family members
kite use self && kite buy RELIANCE 10
kite use spouse && kite buy RELIANCE 10
kite use dad && kite buy RELIANCE 5
```

## Portfolio & Holdings

```bash
kite holdings                      # View demat holdings
kite portfolio                     # View open positions (alias: pos)
kite orders                        # View today's orders
kite margins                       # View available funds (alias: m)
```

## Market Data

```bash
kite quote ITC                     # Quick price check
kite quote RELIANCE -e BSE         # Quote from BSE
kite watchlist ITC RELIANCE TCS    # Multiple quotes
```

> **Note**: Market data commands require Quote API subscription from Zerodha.

## LLM & AI Integration

Connect Kite CLI with any AI assistant that can execute shell commands. Trade using natural language!

### Supported Tools

- **Claude** - Claude Desktop or API with tool use
- **ChatGPT** - Code Interpreter or custom GPTs
- **Cursor/Cline/Aider** - Built-in terminal access
- **Open Interpreter** - Natural language to shell
- **LangChain/AutoGPT** - Custom agents with shell tools

### Example Conversation

```
You: "Buy 10 shares of Reliance for my dad's account"
AI: Executing: kite use dad && kite buy RELIANCE 10

You: "Check spouse's portfolio and show profit/loss"
AI: Executing: kite use spouse && kite holdings
```

### Safety: Use Dry Run

Always test with `--dry-run` when using AI assistants:

```bash
kite buy RELIANCE 10 --dry-run
# ✓ Dry run - order NOT placed
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

### Config File

Location: `~/.kite.yaml`

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
```

### Environment Variables (Override)

| Variable | Description |
|----------|-------------|
| `KITE_API_KEY` | Your Kite Connect API key |
| `KITE_API_SECRET` | Your Kite Connect API secret |
| `KITE_ACCESS_TOKEN` | Session access token |

## Update

```bash
kite update              # Quick update
kite update --check      # Check for updates only
```

## Uninstall

```bash
# Quick uninstall
curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/uninstall.sh | bash

# Manual uninstall
rm -f ~/bin/kite ~/.kite.yaml ~/.zsh/completions/_kite
```

> **Warning**: This will delete credentials for all family accounts. Back up `~/.kite.yaml` if needed.

## Security

- Config file permissions set to `600` (owner read/write only)
- All credentials stored locally - no cloud sync
- Never commit `~/.kite.yaml` to git

## Troubleshooting

| Error | Solution |
|-------|----------|
| "api_key not configured" | Run `kite auth` |
| "Invalid access token" | Run `kite login` |
| "Insufficient permission" | Contact Zerodha for Quote API access |
| "Command not found: kite" | Add `~/bin` to PATH: `echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc` |

## License

MIT
