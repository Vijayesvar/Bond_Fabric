#!/bin/bash

# Test Runner Script for Bond Tokenization Chaincodes
# This script runs unit tests for all chaincodes and provides a summary

set -e

echo "=========================================="
echo "Bond Tokenization Chaincode Test Runner"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run tests for a specific chaincode
run_chaincode_tests() {
    local chaincode_name=$1
    local chaincode_path=$2
    
    echo -e "\n${YELLOW}Testing $chaincode_name...${NC}"
    echo "----------------------------------------"
    
    if [ ! -d "$chaincode_path" ]; then
        echo -e "${RED}Error: $chaincode_path not found${NC}"
        return 1
    fi
    
    cd "$chaincode_path"
    
    # Check if go.mod exists
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}Error: go.mod not found in $chaincode_path${NC}"
        cd - > /dev/null
        return 1
    fi
    
    # Download dependencies
    echo "Downloading dependencies..."
    go mod download
    
    # Run tests with verbose output
    echo "Running tests..."
    if go test -v ./...; then
        echo -e "${GREEN}✓ All tests passed for $chaincode_name${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ Some tests failed for $chaincode_name${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    cd - > /dev/null
}

# Function to check if Go is installed
check_go_installation() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
        echo "Please install Go 1.19 or later and try again"
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    local major_version=$(echo $go_version | cut -d. -f1)
    local minor_version=$(echo $go_version | cut -d. -f2)
    
    if [ "$major_version" -lt 1 ] || ([ "$major_version" -eq 1 ] && [ "$minor_version" -lt 19 ]); then
        echo -e "${RED}Error: Go version $go_version is not supported${NC}"
        echo "Please install Go 1.19 or later"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Go $go_version detected${NC}"
}

# Function to check if testify is available
check_testify() {
    echo "Checking testify dependency..."
    if ! go list -m github.com/stretchr/testify &> /dev/null; then
        echo -e "${YELLOW}Warning: testify not found, installing...${NC}"
        go get github.com/stretchr/testify
    fi
}

# Function to display test summary
display_summary() {
    echo -e "\n=========================================="
    echo "Test Summary"
    echo "=========================================="
    echo "Total Chaincodes Tested: $TOTAL_TESTS"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All chaincode tests passed!${NC}"
        exit 0
    else
        echo -e "\n${RED}❌ Some chaincode tests failed${NC}"
        exit 1
    fi
}

# Main execution
main() {
    # Check prerequisites
    check_go_installation
    check_testify
    
    echo -e "\n${GREEN}Starting chaincode tests...${NC}"
    
    # Test each chaincode
    run_chaincode_tests "BondToken" "chaincode/bondtoken"
    run_chaincode_tests "Compliance" "chaincode/compliance"
    run_chaincode_tests "CorporateAction" "chaincode/corporateaction"
    
    # Display summary
    display_summary
}

# Run main function
main "$@"

