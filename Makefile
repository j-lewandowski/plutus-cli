BINARY_NAME=plutus

DEV_DB=./dev_plutus.sqlite

ifeq ($(OS),Windows_NT)
    RM = del /Q /F
    EXT = .exe
else
    RM = rm -f
    EXT =
endif

.PHONY: all build run clean

all: build

build:
	@echo "🔨 Building production version..."
	go build -o $(BINARY_NAME)$(EXT) ./cmd/plutus
	@echo "✅ Build complete: $(BINARY_NAME)$(EXT)"

run:
	@echo "🔧 Running command in DEV mode..."
	PLUTUS_DB=$(DEV_DB) go run ./cmd/plutus $(args)

clean:
	@echo "🧹 Cleaning up..."
	$(RM) $(BINARY_NAME)$(EXT)
	$(RM) $(DEV_DB)
	@echo "✨ Cleaned."