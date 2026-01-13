# Plutus CLI

Plutus is a local-first CLI tool written in Go for tracking long-term investments and visualizing portfolio performance over time.

---

## Why Plutus?

I like to call myself an investor - although my portfolio occasionally disagrees.

Over time, I realized that tracking investments in spreadsheets quickly turns into a repetitive and boring process: downloading data, adjusting formulas, fixing broken charts, and repeating the same steps every time I wanted to see an updated report.

Plutus was created to automate this workflow. Instead of manually preparing data, I wanted a tool that would reliably fetch market information, persist it locally, and generate reports with a single command.

The name [_Plutus_](https://en.wikipedia.org/wiki/Plutus) comes from Greek mythology a subtle nod to my interest in long-term investing and a personal appreciation for Greek mythology.

---

## Why Go?

This project is also a hands-on learning experience with Go, chosen for its:

- simplicity and readability
- excellent support for concurrency
- single static binary with no external runtime
- great fit for CLI tools and backend services

Plutus is distributed as a single executable and does not require an interpreter or runtime environment like Node.js or Python.

---

## Features

- Concurrent market data fetching from multiple sources
- Storing market data in local-first SQLite database
- Charts generation [Comming soon!]

---

## Supported Assets

> **Note:** As of version **1.0.0**, Plutus only supports **Invesco S&P 500 UCITS ETF (P500.DE)**.
>
> There are plans to implement custom configuration in future releases, allowing users to track any index or asset they prefer.

---

## Installation

To install `plutus`, simply run the following command:

```bash
curl -fsSL https://raw.githubusercontent.com/j-lewandowski/plutus-cli/main/install.sh | bash
```

## Commands

Plutus provides several commands to manage your portfolio.

### `add`

Adds a new deposit to your portfolio.

```bash
plutus add <amount> [date]
```

- **amount**: The amount of money deposited (e.g., `100`, `125.50`).
- **date** (optional): The date of the deposit. Supported formats: `DD.MM.YYYY`, `DD-MM-YYYY`, `YYYY-MM-DD`.

### `sync`

Fetches up-to-date market data from external sources (NBP, Yahoo Finance) to update the CLI's local database.

```bash
plutus sync
```

### `status`

Displays the current portfolio value, total invested amount, and profit/loss performance.

```bash
plutus status
```

### `help`

Shows information about available commands.

```bash
plutus help
```
