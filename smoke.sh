#!/bin/bash

# Smoke Test Script for Bragdoc CLI
# This script performs basic smoke tests to verify core functionality
# Usage: ./smoke.sh [target_os]
# Example: ./smoke.sh darwin/amd64

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Get target OS from environment variable or parameter
TARGET_OS="${TARGET_OS:-${1:-$(go env GOOS)/$(go env GOARCH)}}"
BINARY_NAME="bragdoc"

# Parse OS and ARCH
IFS='/' read -r GOOS GOARCH <<< "$TARGET_OS"
export GOOS GOARCH

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Bragdoc CLI Smoke Test Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Target OS: ${YELLOW}${TARGET_OS}${NC}"
echo -e "Binary: ${YELLOW}${BINARY_NAME}${NC}"
echo ""

# Function to print test result
print_result() {
    local test_name=$1
    local result=$2
    local message=$3
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [ "$result" -eq 0 ]; then
        echo -e "${GREEN}✓${NC} ${test_name}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗${NC} ${test_name}"
        if [ -n "$message" ]; then
            echo -e "  ${RED}Error: ${message}${NC}"
        fi
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Function to run a test
run_test() {
    local test_name=$1
    shift
    local output
    local exit_code
    
    output=$("$@" 2>&1) || exit_code=$?
    exit_code=${exit_code:-0}
    
    if [ $exit_code -eq 0 ]; then
        print_result "$test_name" 0
        return 0
    else
        print_result "$test_name" 1 "$output"
        return 1
    fi
}

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}Cleaning up test environment...${NC}"
    rm -rf ./.bragdoc
    rm -f ./${BINARY_NAME}
    echo -e "${GREEN}Cleanup complete${NC}"
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Clean up any existing test environment and binary
rm -rf ./.bragdoc
rm -f ./${BINARY_NAME}

# Step 1: Build the binary
echo -e "${BLUE}Step 1: Building binary for ${TARGET_OS}...${NC}"
if CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BINARY_NAME} ./cmd/cli/main.go 2>&1; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi
echo ""

# Make binary executable
chmod +x ./${BINARY_NAME}

# Set HOME to current directory for testing (will use ./.bragdoc)
export HOME=$(pwd)
mkdir -p $HOME/.bragdoc

# Step 2: Test help command (should work without initialization)
echo -e "${BLUE}Step 2: Testing help commands...${NC}"
run_test "bragdoc --help" ./${BINARY_NAME} --help
run_test "bragdoc version" ./${BINARY_NAME} version
run_test "bragdoc brag --help" ./${BINARY_NAME} brag --help
echo ""

# Step 3: Test that commands fail without initialization
echo -e "${BLUE}Step 3: Testing initialization requirement...${NC}"
if ./${BINARY_NAME} brag list 2>&1 | grep -q "not initialized"; then
    print_result "brag list fails without init" 0
else
    print_result "brag list fails without init" 1 "Should require initialization"
fi

if ./${BINARY_NAME} brag add --title "Test" --description "Should fail" 2>&1 | grep -q "not initialized"; then
    print_result "brag add fails without init" 0
else
    print_result "brag add fails without init" 1 "Should require initialization"
fi
echo ""

# Step 4: Initialize bragdoc
echo -e "${BLUE}Step 4: Initializing bragdoc...${NC}"
if ./${BINARY_NAME} init --name "Smoke Test User" --email "smoke@test.com" --locale "en-US" 2>&1 | grep -q "initialized successfully"; then
    print_result "bragdoc init" 0
else
    print_result "bragdoc init" 1 "Initialization failed"
    exit 1
fi
echo ""

# Step 5: Test brag add command
echo -e "${BLUE}Step 5: Testing brag add...${NC}"
if ./${BINARY_NAME} brag add \
    --title "Test Achievement" \
    --description "This is a test achievement for smoke testing" \
    --category "achievement" \
    --tags "test,smoke,automation" 2>&1 | grep -q "created successfully"; then
    print_result "brag add with tags" 0
else
    print_result "brag add with tags" 1 "Failed to create brag"
fi

if ./${BINARY_NAME} brag add \
    --title "Leadership Project" \
    --description "Led a team to deliver a critical project on time" \
    --category "leadership" \
    --tags "leadership,team" 2>&1 | grep -q "created successfully"; then
    print_result "brag add leadership" 0
else
    print_result "brag add leadership" 1 "Failed to create leadership brag"
fi

if ./${BINARY_NAME} brag add \
    --title "Innovation Initiative" \
    --description "Implemented innovative solution that improved efficiency" \
    --category "innovation" 2>&1 | grep -q "created successfully"; then
    print_result "brag add innovation" 0
else
    print_result "brag add innovation" 1 "Failed to create innovation brag"
fi
echo ""

# Step 6: Test brag list command
echo -e "${BLUE}Step 6: Testing brag list...${NC}"
if ./${BINARY_NAME} brag list 2>&1 | grep -q "Test Achievement"; then
    print_result "brag list (table format)" 0
else
    print_result "brag list (table format)" 1 "Failed to list brags"
fi

if ./${BINARY_NAME} brag list --format json 2>&1 | grep -q '"Title"'; then
    print_result "brag list (json format)" 0
else
    print_result "brag list (json format)" 1 "Failed to list brags in JSON"
fi

if ./${BINARY_NAME} brag list --format yaml 2>&1 | grep -q "title:"; then
    print_result "brag list (yaml format)" 0
else
    print_result "brag list (yaml format)" 1 "Failed to list brags in YAML"
fi
echo ""

# Step 7: Test brag filtering
echo -e "${BLUE}Step 7: Testing brag filtering...${NC}"
if ./${BINARY_NAME} brag list --category leadership 2>&1 | grep -q "Leadership Project"; then
    print_result "filter by category" 0
else
    print_result "filter by category" 1 "Failed to filter by category"
fi

if ./${BINARY_NAME} brag list --tags test 2>&1 | grep -q "Test Achievement"; then
    print_result "filter by tags" 0
else
    print_result "filter by tags" 1 "Failed to filter by tags"
fi
echo ""

# Step 8: Test brag show command
echo -e "${BLUE}Step 8: Testing brag show...${NC}"
if ./${BINARY_NAME} brag show 1 2>&1 | grep -q "Test Achievement"; then
    print_result "show single brag" 0
else
    print_result "show single brag" 1 "Failed to show brag"
fi

if ./${BINARY_NAME} brag show 1,2 2>&1 | grep -q "Test Achievement"; then
    print_result "show multiple brags" 0
else
    print_result "show multiple brags" 1 "Failed to show multiple brags"
fi

if ./${BINARY_NAME} brag show 1-3 2>&1 | grep -q "Test Achievement"; then
    print_result "show range of brags" 0
else
    print_result "show range of brags" 1 "Failed to show range of brags"
fi
echo ""

# Step 9: Test brag edit command
echo -e "${BLUE}Step 9: Testing brag edit...${NC}"
if ./${BINARY_NAME} brag edit 1 --title "Updated Test Achievement" 2>&1 | grep -q "updated successfully"; then
    print_result "edit brag title" 0
else
    print_result "edit brag title" 1 "Failed to edit brag"
fi

if ./${BINARY_NAME} brag show 1 2>&1 | grep -q "Updated Test Achievement"; then
    print_result "verify edit persisted" 0
else
    print_result "verify edit persisted" 1 "Edit did not persist"
fi
echo ""

# Step 10: Test brag remove command
echo -e "${BLUE}Step 10: Testing brag remove...${NC}"
if ./${BINARY_NAME} brag remove 3 --force 2>&1 | grep -q "Successfully removed"; then
    print_result "remove brag with --force" 0
else
    print_result "remove brag with --force" 1 "Failed to remove brag"
fi

if ! ./${BINARY_NAME} brag show 3 2>&1 | grep -q "Innovation Initiative"; then
    print_result "verify removal" 0
else
    print_result "verify removal" 1 "Brag was not removed"
fi
echo ""

# Step 11: Test tag management
echo -e "${BLUE}Step 11: Testing tag management...${NC}"
if ./${BINARY_NAME} tag list 2>&1 | grep -q "test\|smoke\|automation\|leadership\|team"; then
    print_result "tag list (auto-created tags)" 0
else
    print_result "tag list (auto-created tags)" 1 "Failed to list auto-created tags"
fi

if ./${BINARY_NAME} tag add --name "golang" 2>&1 | grep -q "created successfully"; then
    print_result "tag add" 0
else
    print_result "tag add" 1 "Failed to create tag"
fi

# Add another tag specifically for testing removal (not associated with any brag)
if ./${BINARY_NAME} tag add --name "removeme" 2>&1 | grep -q "created successfully"; then
    print_result "tag add (for removal test)" 0
else
    print_result "tag add (for removal test)" 1 "Failed to create tag for removal test"
fi

if ./${BINARY_NAME} tag list --format json 2>&1 | grep -q '"Name"'; then
    print_result "tag list (json format)" 0
else
    print_result "tag list (json format)" 1 "Failed to list tags in JSON"
fi

if ./${BINARY_NAME} tag list --format yaml 2>&1 | grep -q "name:"; then
    print_result "tag list (yaml format)" 0
else
    print_result "tag list (yaml format)" 1 "Failed to list tags in YAML"
fi

# Test tag validation
if ./${BINARY_NAME} tag add --name "a" 2>&1 | grep -q "at least 2 characters"; then
    print_result "tag validation (too short)" 0
else
    print_result "tag validation (too short)" 1 "Should reject short tag name"
fi

if ./${BINARY_NAME} tag add --name "this-is-way-too-long-x" 2>&1 | grep -q "cannot exceed 20 characters"; then
    print_result "tag validation (too long)" 0
else
    print_result "tag validation (too long)" 1 "Should reject long tag name"
fi

if ./${BINARY_NAME} tag add --name "golang" 2>&1 | grep -q "already exists"; then
    print_result "tag validation (duplicate)" 0
else
    print_result "tag validation (duplicate)" 1 "Should reject duplicate tag"
fi

# Get a tag ID to test removal (use the "removeme" tag we just created)
TAG_LIST=$(./${BINARY_NAME} tag list 2>&1)
if echo "$TAG_LIST" | grep -q "removeme"; then
    # Extract the ID from the table format (first column)
    # Skip the header lines (ID, --) and get only numeric IDs
    TAG_ID=$(echo "$TAG_LIST" | grep "removeme" | grep -v "^ID" | grep -v "^--" | awk '{print $1}' | grep -E '^[0-9]+$')
    
    if [ -n "$TAG_ID" ] && [ "$TAG_ID" -gt 0 ] 2>/dev/null; then
        REMOVE_OUTPUT=$(./${BINARY_NAME} tag remove "$TAG_ID" --confirm 2>&1)
        if echo "$REMOVE_OUTPUT" | grep -q "removed successfully"; then
            print_result "tag remove" 0
        else
            print_result "tag remove" 1 "Failed to remove tag ID $TAG_ID"
        fi
        
        # Verify the tag was actually removed
        VERIFY_LIST=$(./${BINARY_NAME} tag list 2>&1)
        if ! echo "$VERIFY_LIST" | grep -q "removeme"; then
            print_result "verify tag removal" 0
        else
            print_result "verify tag removal" 1 "Tag was not removed"
        fi
    else
        print_result "tag remove" 1 "Invalid tag ID: '$TAG_ID'"
    fi
else
    print_result "tag remove" 1 "removeme tag not found in list"
fi
echo ""

# Step 12: Test error handling
echo -e "${BLUE}Step 12: Testing error handling...${NC}"
if ./${BINARY_NAME} brag add --title "Bad" --description "Too short" 2>&1 | grep -q "at least"; then
    print_result "validation error (short title)" 0
else
    print_result "validation error (short title)" 1 "Should reject short title"
fi

if ./${BINARY_NAME} brag add --title "Valid Title" --description "Short" 2>&1 | grep -q "at least"; then
    print_result "validation error (short description)" 0
else
    print_result "validation error (short description)" 1 "Should reject short description"
fi

if ./${BINARY_NAME} brag add --title "Valid Title" --description "This is a valid description" --category "invalid" 2>&1 | grep -q "invalid category"; then
    print_result "validation error (invalid category)" 0
else
    print_result "validation error (invalid category)" 1 "Should reject invalid category"
fi
echo ""

# Print summary
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Total tests run: ${YELLOW}${TESTS_RUN}${NC}"
echo -e "Tests passed: ${GREEN}${TESTS_PASSED}${NC}"
echo -e "Tests failed: ${RED}${TESTS_FAILED}${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All smoke tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some smoke tests failed${NC}"
    exit 1
fi
