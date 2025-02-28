#!/bin/bash
# build.sh - Build script for govm (Go Virtual Machine)

set -e

# Default values
BUILD_MODE="dev"
OUTPUT_DIR="bin"
BINARY_NAME="govm"
PLATFORMS=""
VERBOSE=0
CLEAN=0

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
    echo "  -m, --mode MODE       Build mode: dev, release (default: dev)"
    echo "  -o, --output DIR      Output directory (default: bin)"
    echo "  -n, --name NAME       Binary name (default: govm)"
    echo "  -p, --platforms LIST  Comma-separated list of platforms to build for"
    echo "                        Format: os/arch (e.g., linux/amd64,darwin/amd64)"
    echo "  -c, --clean           Clean before building"
    echo "  -v, --verbose         Verbose output"
    echo "  -h, --help            Show this help message"
    echo
    echo "Examples:"
    echo "  $0 --mode release --platforms linux/amd64,darwin/amd64,windows/amd64"
    echo "  $0 --clean --verbose"
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--mode)
            BUILD_MODE="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -n|--name)
            BINARY_NAME="$2"
            shift 2
            ;;
        -p|--platforms)
            PLATFORMS="$2"
            shift 2
            ;;
        -c|--clean)
            CLEAN=1
            shift
            ;;
        -v|--verbose)
            VERBOSE=1
            shift
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

# Validate build mode
if [[ "$BUILD_MODE" != "dev" && "$BUILD_MODE" != "release" ]]; then
    echo -e "${RED}Error:${NC} Invalid build mode: $BUILD_MODE"
    echo "Valid modes: dev, release"
    exit 1
fi

# Set build flags based on mode
if [[ "$BUILD_MODE" == "dev" ]]; then
    BUILD_FLAGS=""
    echo -e "${YELLOW}Building in development mode${NC}"
else
    BUILD_FLAGS="-ldflags=\"-s -w\""
    echo -e "${YELLOW}Building in release mode${NC}"
fi

# Clean if requested
if [[ $CLEAN -eq 1 ]]; then
    echo -e "${BLUE}Cleaning...${NC}"
    rm -rf "$OUTPUT_DIR"
    go clean
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build for current platform if no platforms specified
if [[ -z "$PLATFORMS" ]]; then
    echo -e "${GREEN}Building $BINARY_NAME for current platform...${NC}"
    
    BUILD_CMD="go build"
    if [[ $VERBOSE -eq 1 ]]; then
        BUILD_CMD="$BUILD_CMD -v"
    fi
    
    if [[ "$BUILD_MODE" == "release" ]]; then
        eval "$BUILD_CMD -ldflags=\"-s -w\" -o \"$OUTPUT_DIR/$BINARY_NAME\" ."
    else
        eval "$BUILD_CMD -o \"$OUTPUT_DIR/$BINARY_NAME\" ."
    fi
    
    echo -e "${GREEN}Build complete: $OUTPUT_DIR/$BINARY_NAME${NC}"
else
    # Build for specified platforms
    echo -e "${GREEN}Building $BINARY_NAME for multiple platforms...${NC}"
    
    IFS=',' read -ra PLATFORM_LIST <<< "$PLATFORMS"
    for platform in "${PLATFORM_LIST[@]}"; do
        IFS='/' read -ra parts <<< "$platform"
        if [[ ${#parts[@]} -ne 2 ]]; then
            echo -e "${RED}Invalid platform format: $platform. Expected format: os/arch${NC}"
            continue
        fi
        
        OS="${parts[0]}"
        ARCH="${parts[1]}"
        
        echo -e "${BLUE}Building for $OS/$ARCH...${NC}"
        
        # Set output filename based on platform
        if [[ "$OS" == "windows" ]]; then
            output_file="$OUTPUT_DIR/${BINARY_NAME}-${OS}-${ARCH}.exe"
        else
            output_file="$OUTPUT_DIR/${BINARY_NAME}-${OS}-${ARCH}"
        fi
        
        BUILD_CMD="GOOS=$OS GOARCH=$ARCH go build"
        if [[ $VERBOSE -eq 1 ]]; then
            BUILD_CMD="$BUILD_CMD -v"
        fi
        
        if [[ "$BUILD_MODE" == "release" ]]; then
            eval "$BUILD_CMD -ldflags=\"-s -w\" -o \"$output_file\" ."
        else
            eval "$BUILD_CMD -o \"$output_file\" ."
        fi
        
        if [[ $? -eq 0 ]]; then
            echo -e "${GREEN}Successfully built: $output_file${NC}"
        else
            echo -e "${RED}Failed to build for $OS/$ARCH${NC}"
        fi
    done
fi

echo -e "${GREEN}Build process completed!${NC}"
