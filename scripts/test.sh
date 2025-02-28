#!/bin/bash
# test.sh - Test script for smg (Go Virtual Machine)

set -e

# Default values
VERBOSE=0
COVERAGE=0
COVERAGE_DIR="coverage"
PACKAGE="./..."
RACE=0
BENCH=0
BENCH_TIME="1s"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print usage information
usage() {
    echo -e "${BLUE}Usage:${NC} $0 [options]"
    echo
    echo "Options:"
    echo "  -v, --verbose         Verbose output"
    echo "  -c, --coverage        Generate coverage report"
    echo "  -d, --coverage-dir DIR Set coverage output directory (default: coverage)"
    echo "  -p, --package PKG     Package to test (default: ./...)"
    echo "  -r, --race            Enable race detection"
    echo "  -b, --bench           Run benchmarks"
    echo "  -t, --bench-time TIME Set benchmark time (default: 1s)"
    echo "  -h, --help            Show this help message"
    echo
    echo "Examples:"
    echo "  $0 --verbose --coverage"
    echo "  $0 --package ./parser/... --race"
    echo "  $0 --bench --bench-time 2s"
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=1
            shift
            ;;
        -c|--coverage)
            COVERAGE=1
            shift
            ;;
        -d|--coverage-dir)
            COVERAGE_DIR="$2"
            shift 2
            ;;
        -p|--package)
            PACKAGE="$2"
            shift 2
            ;;
        -r|--race)
            RACE=1
            shift
            ;;
        -b|--bench)
            BENCH=1
            shift
            ;;
        -t|--bench-time)
            BENCH_TIME="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo -e "${RED}Error:${NC} Unknown option $1"
            usage
            ;;
    esac
done

# Build test command
TEST_CMD="go test"

if [[ $VERBOSE -eq 1 ]]; then
    TEST_CMD="$TEST_CMD -v"
fi

if [[ $RACE -eq 1 ]]; then
    TEST_CMD="$TEST_CMD -race"
    echo -e "${YELLOW}Race detection enabled${NC}"
fi

if [[ $COVERAGE -eq 1 ]]; then
    mkdir -p "$COVERAGE_DIR"
    COVERAGE_FILE="$COVERAGE_DIR/coverage.out"
    TEST_CMD="$TEST_CMD -coverprofile=$COVERAGE_FILE"
    echo -e "${YELLOW}Coverage report will be generated at $COVERAGE_FILE${NC}"
fi

# Run tests or benchmarks
if [[ $BENCH -eq 1 ]]; then
    echo -e "${GREEN}Running benchmarks for $PACKAGE...${NC}"
    TEST_CMD="$TEST_CMD -bench=. -benchtime=$BENCH_TIME -run=^$ $PACKAGE"
else
    echo -e "${GREEN}Running tests for $PACKAGE...${NC}"
    TEST_CMD="$TEST_CMD $PACKAGE"
fi

echo -e "${BLUE}Executing: $TEST_CMD${NC}"
eval "$TEST_CMD"

# Generate HTML coverage report if requested
if [[ $COVERAGE -eq 1 ]]; then
    HTML_COVERAGE_FILE="$COVERAGE_DIR/coverage.html"
    echo -e "${GREEN}Generating HTML coverage report at $HTML_COVERAGE_FILE${NC}"
    go tool cover -html="$COVERAGE_FILE" -o "$HTML_COVERAGE_FILE"
    
    # Print coverage summary
    echo -e "${BLUE}Coverage summary:${NC}"
    go tool cover -func="$COVERAGE_FILE"
fi

echo -e "${GREEN}Test process completed!${NC}" 