#!/bin/bash
# release.sh - Release script for smg (Stack Machine Go)

set -e

# Default values
VERSION=""
OUTPUT_DIR="dist"
PLATFORMS="linux/amd64,darwin/amd64,windows/amd64"
SKIP_TESTS=0
SKIP_TAG=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print usage information
usage() {
    echo -e "${BLUE}Usage:${NC} $0 [options] -v VERSION"
    echo
    echo "Options:"
    echo "  -v, --version VERSION  Version to release (required)"
    echo "  -o, --output DIR       Output directory (default: dist)"
    echo "  -p, --platforms LIST   Comma-separated list of platforms to build for"
    echo "                         Format: os/arch (default: linux/amd64,darwin/amd64,windows/amd64)"
    echo "  -s, --skip-tests       Skip running tests"
    echo "  -t, --skip-tag         Skip creating git tag"
    echo "  -h, --help             Show this help message"
    echo
    echo "Examples:"
    echo "  $0 --version 1.0.0"
    echo "  $0 --version 1.0.0 --platforms linux/amd64,darwin/amd64 --skip-tests"
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -p|--platforms)
            PLATFORMS="$2"
            shift 2
            ;;
        -s|--skip-tests)
            SKIP_TESTS=1
            shift
            ;;
        -t|--skip-tag)
            SKIP_TAG=1
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

# Check if version is provided
if [[ -z "$VERSION" ]]; then
    echo -e "${RED}Error:${NC} Version is required"
    usage
fi

# Validate version format (semver)
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?(\+[a-zA-Z0-9.]+)?$ ]]; then
    echo -e "${RED}Error:${NC} Invalid version format: $VERSION"
    echo "Version should follow semantic versioning (e.g., 1.0.0, 1.0.0-alpha, 1.0.0+build.1)"
    exit 1
fi

echo -e "${GREEN}Preparing release v$VERSION${NC}"

# Run tests if not skipped
if [[ $SKIP_TESTS -eq 0 ]]; then
    echo -e "${BLUE}Running tests...${NC}"
    ./scripts/test.sh
    if [[ $? -ne 0 ]]; then
        echo -e "${RED}Tests failed. Aborting release.${NC}"
        exit 1
    fi
    echo -e "${GREEN}Tests passed!${NC}"
else
    echo -e "${YELLOW}Skipping tests${NC}"
fi

# Create release directory
mkdir -p "$OUTPUT_DIR"

# Build release binaries
echo -e "${BLUE}Building release binaries...${NC}"
./scripts/build.sh --mode release --output "$OUTPUT_DIR" --platforms "$PLATFORMS"

# Create release archive
echo -e "${BLUE}Creating release archive...${NC}"
tar -czf "$OUTPUT_DIR/smg-$VERSION.tar.gz" -C "$OUTPUT_DIR" .

# Create git tag if not skipped
if [[ $SKIP_TAG -eq 0 ]]; then
    echo -e "${BLUE}Creating git tag v$VERSION...${NC}"
    git tag -a "v$VERSION" -m "Release v$VERSION"
    echo -e "${GREEN}Tag created: v$VERSION${NC}"
    echo -e "${YELLOW}Remember to push the tag with: git push origin v$VERSION${NC}"
else
    echo -e "${YELLOW}Skipping git tag creation${NC}"
fi

echo -e "${GREEN}Release v$VERSION created successfully!${NC}"
echo -e "${BLUE}Release files are available in: $OUTPUT_DIR${NC}" 