BINARY_NAME := disk_not_sleep

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)

.PHONY: all build clean