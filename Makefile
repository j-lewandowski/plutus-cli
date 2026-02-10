BINARY_NAME=plutus_dev
DEV_DB=./dev_plutus.sqlite

VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X 'main.Version=$(VERSION)' \
           -X 'main.Commit=$(COMMIT)' \
           -X 'main.BuildTime=$(BUILD_TIME)'

ifeq ($(OS),Windows_NT)
    RM = del /Q /F
    EXT = .exe
else
    RM = rm -f
    EXT =
endif

.PHONY: all build run clean test-update dashboard db-reset db-seed db-shell db-dump

all: build

build:
	@echo "🔨 Building version $(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)$(EXT) ./cmd/plutus
	@echo "✅ Build complete: $(BINARY_NAME)$(EXT)"

run:
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus $(filter-out $@,$(MAKECMDGOALS))

clean:
	@$(RM) $(BINARY_NAME)$(EXT)
	@$(RM) $(DEV_DB)

test-update:
	@rm -f ~/.plutus_update_check
	@PLUTUS_DEBUG_UPDATE=true go run -ldflags "-X 'main.Version=v0.0.1'" ./cmd/plutus status

dashboard: build
	@PLUTUS_DB=$(DEV_DB) ./$(BINARY_NAME)$(EXT) dashboard


db-reset:
	@echo "🗑️  Resetting dev database..."
	@$(RM) $(DEV_DB)
	@echo "✅ Dev database deleted. It will be recreated on next run."

db-seed: db-reset
	@echo "🌱 Seeding dev database..."
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus add 100 01.06.2025
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus add 200 15.07.2025
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus add 150 01.09.2025
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus add 300 15.10.2025
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus add 123 12.12.2025
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus add 250 15.01.2026
	@echo "📡 Syncing data..."
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus sync
	@echo "✅ Seed complete! Run 'make dashboard' or 'make run status' to verify."

db-shell:
	@sqlite3 $(DEV_DB)

db-dump:
	@echo "=== deposits ==="
	@sqlite3 $(DEV_DB) "SELECT id, deposit_date, deposit_amount_in_eurocents/100.0 || ' €' as amount FROM deposit ORDER BY deposit_date;"
	@echo ""
	@echo "=== index_price (last 5) ==="
	@sqlite3 $(DEV_DB) "SELECT date, price_in_eurocents/100.0 || ' €' as price, CASE is_real WHEN 1 THEN 'real' ELSE 'filled' END as type FROM index_price ORDER BY date DESC LIMIT 5;"
	@echo ""
	@echo "=== eur_exchange_rate (last 5) ==="
	@sqlite3 $(DEV_DB) "SELECT date, price_pln_in_grosz/100.0 || ' PLN' as rate FROM eur_exchange_rate ORDER BY date DESC LIMIT 5;"

%:
	@: