# Executable file name
BINARY_NAME := mhc

# Installation path - Linux systems typically use /usr/bin
INSTALL_PATH := /usr/bin

# Default target: build the project
all: build

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go
	@echo "Build completed: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	@echo "Cleanup completed"

# Install the program to system directory
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	# Check for sufficient permissions, use sudo if needed
	if [ ! -w $(INSTALL_PATH) ]; then \
		echo "Administrator privileges required, attempting installation with sudo..."; \
		sudo cp $(BINARY_NAME) $(INSTALL_PATH)/; \
	else \
		cp $(BINARY_NAME) $(INSTALL_PATH)/; \
	fi
	@echo "Installation completed, you can use the '$(BINARY_NAME)' command directly"

# Uninstall the program
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	if [ ! -w $(INSTALL_PATH) ]; then \
		sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "Uninstallation completed"

# Run the program
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Phony targets declaration
.PHONY: all build clean install uninstall run
