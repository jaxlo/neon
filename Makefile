# Compiler and flags
CC = gcc
CFLAGS = -Wall -Wextra -std=c99 -O2
LIBS = -lraylib -lGL -lm -lpthread -ldl -lrt -lX11

# Source files
SOURCES = src/main.c \
    src/player/*.c \
    src/renderer/*.c \
	src/world/*.c \
	src/ui/*ui.c

# Output
TARGET = a.out
BUILD_DIR = build

# Default target
all: $(TARGET)

# Build the game
$(TARGET): $(SOURCES)
	@mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) $(SOURCES) -o $(BUILD_DIR)/$(TARGET) $(LIBS)

# Development target (build and run)
dev: $(TARGET)
	./$(BUILD_DIR)/$(TARGET)

# Debug build
debug: CFLAGS += -g -DDEBUG
debug: $(TARGET)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Help
help:
	@echo "Available targets:"
	@echo "  all        - Build the game"
	@echo "  dev        - Build and run the game"
	@echo "  debug      - Build with debug symbols"
	@echo "  clean      - Remove build artifacts"
	@echo "  help       - Show this help"

.PHONY: all dev debug clean help