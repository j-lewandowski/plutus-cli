BINARY_NAME=plutus_dev
DEV_DB=./dev_plutus.sqlite

VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

PKG := plutus-cli/cmd

LDFLAGS := -X '$(PKG).Version=$(VERSION)' \
           -X '$(PKG).Commit=$(COMMIT)' \
           -X '$(PKG).BuildTime=$(BUILD_TIME)'

ifeq ($(OS),Windows_NT)
    RM = del /Q /F
    EXT = .exe
else
    RM = rm -f
    EXT =
endif

.PHONY: all build run clean test-update dashboard

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
	@PLUTUS_DEBUG_UPDATE=true go run -ldflags "-X '$(PKG).Version=v0.0.1'" ./cmd/plutus status

dashboard:
	@PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus dashboard

%:
	@: