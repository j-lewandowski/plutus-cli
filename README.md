# Plutus CLI

Plutus is a local-first CLI tool written in Go for tracking long-term investments and visualizing portfolio performance over time.

---

## Why Plutus?

I like to call myself an investor - although my portfolio occasionally disagrees.

Over time, I realized that tracking investments in spreadsheets quickly turns into a repetitive and boring process: downloading data, adjusting formulas, fixing broken charts, and repeating the same steps every time I wanted to see an updated report.

Plutus was created to automate this workflow. Instead of manually preparing data, I wanted a tool that would reliably fetch market information, persist it locally, and generate reports with a single command.

The name [*Plutus*](https://en.wikipedia.org/wiki/Plutus) comes from Greek mythology a subtle nod to my interest in long-term investing and a personal appreciation for Greek mythology.

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
- Charts generation
