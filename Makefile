.PHONY: build run run-ro run-sync run-ro-sync clean

build:
	@echo "Building binary..."
	@go build -o bin/$(BINARY_NAME) -v
	@echo "Done."

run:
	@echo "Running..."
	@go run main.go $(ARGS)

run-ro:
	@echo "Running read-only..."
	@go run main.go $(ARGS) -db-read-only

run-sync:
	@echo "Running with sync..."
	@go run main.go $(ARGS) -db-sync-type=sync

run-ro-sync:
	@echo "Running read-only with sync..."
	@go run main.go $(ARGS) -db-read-only -db-sync-type=sync

clean:
	@echo "Cleaning..."
	@rm -rf bin
	@echo "Done."