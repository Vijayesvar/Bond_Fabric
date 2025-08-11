#!/bin/bash

# Compliance Chaincode CLI Script
# This script provides a command-line interface for interacting with the Compliance chaincode

set -e

# Configuration
CHANNEL_NAME="mychannel"
CHAINCODE_NAME="compliance"
CHAINCODE_VERSION="1.0"
PEER_ADDRESS="localhost:7051"
ORDERER_ADDRESS="localhost:7050"
MSP_ID="Org1MSP"
MSP_PATH="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display usage
show_usage() {
    echo -e "${BLUE}Compliance Chaincode CLI${NC}"
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  create-kyc <address> <full_name> <dob> <nationality> <id_type> <id_number>"
    echo "  approve-kyc <address> <approved_by> <risk_level>"
    echo "  reject-kyc <address> <rejected_by> <reason>"
    echo "  create-aml <address> <check_type> <risk_score> <details>"
    echo "  update-aml <address> <check_type> <status> <risk_score> <details>"
    echo "  check-compliance <address>"
    echo "  get-kyc <address>"
    echo "  get-aml <address> <check_type>"
    echo "  get-all-kyc"
    echo "  get-all-aml <address>"
    echo "  help"
    echo ""
    echo "Examples:"
    echo "  $0 create-kyc alice 'Alice Johnson' '1990-01-01' 'US' 'PASSPORT' 'US123456'"
    echo "  $0 approve-kyc alice admin1 LOW"
    echo "  $0 create-aml alice SANCTIONS 5 'No sanctions found'"
    echo "  $0 check-compliance alice"
}

# Function to check if peer CLI is available
check_peer_cli() {
    if ! command -v peer &> /dev/null; then
        echo -e "${RED}Error: peer CLI not found${NC}"
        echo "Please ensure you are in the Fabric CLI environment"
        echo "Run: docker exec -it cli bash"
        exit 1
    fi
}

# Function to check if we're in the right environment
check_environment() {
    if [ ! -d "$MSP_PATH" ]; then
        echo -e "${YELLOW}Warning: MSP path not found, using default${NC}"
        MSP_PATH=""
    fi

    # Set environment variables
    export CORE_PEER_LOCALMSPID=$MSP_ID
    if [ -n "$MSP_PATH" ]; then
        export CORE_PEER_MSPCONFIGPATH=$MSP_PATH
    fi
    export CORE_PEER_ADDRESS=$PEER_ADDRESS
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
}

# Function to create KYC record
create_kyc() {
    local address=$1
    local full_name=$2
    local dob=$3
    local nationality=$4
    local id_type=$5
    local id_number=$6

    echo -e "${YELLOW}Creating KYC record for: $address${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CreateKYC\",\"$address\",\"$full_name\",\"$dob\",\"$nationality\",\"$id_type\",\"$id_number\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ KYC record created successfully for $address${NC}"
}

# Function to approve KYC
approve_kyc() {
    local address=$1
    local approved_by=$2
    local risk_level=$3

    echo -e "${YELLOW}Approving KYC for: $address by $approved_by with risk level: $risk_level${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"ApproveKYC\",\"$address\",\"$approved_by\",\"$risk_level\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ KYC approved successfully for $address${NC}"
}

# Function to reject KYC
reject_kyc() {
    local address=$1
    local rejected_by=$2
    local reason=$3

    echo -e "${YELLOW}Rejecting KYC for: $address by $rejected_by with reason: $reason${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"RejectKYC\",\"$address\",\"$rejected_by\",\"$reason\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ KYC rejected successfully for $address${NC}"
}

# Function to create AML check
create_aml() {
    local address=$1
    local check_type=$2
    local risk_score=$3
    local details=$4

    echo -e "${YELLOW}Creating AML check for: $address, type: $check_type${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CreateAMLCheck\",\"$address\",\"$check_type\",\"$risk_score\",\"$details\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ AML check created successfully for $address${NC}"
}

# Function to update AML check
update_aml() {
    local address=$1
    local check_type=$2
    local status=$3
    local risk_score=$4
    local details=$5

    echo -e "${YELLOW}Updating AML check for: $address, type: $check_type${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"UpdateAMLCheck\",\"$address\",\"$check_type\",\"$status\",\"$risk_score\",\"$details\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ AML check updated successfully for $address${NC}"
}

# Function to check compliance
check_compliance() {
    local address=$1

    echo -e "${YELLOW}Checking compliance for: $address${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CheckCompliance\",\"$address\"]}"
}

# Function to get KYC record
get_kyc() {
    local address=$1

    echo -e "${YELLOW}Querying KYC record for: $address${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetKYC\",\"$address\"]}"
}

# Function to get AML check
get_aml() {
    local address=$1
    local check_type=$2

    echo -e "${YELLOW}Querying AML check for: $address, type: $check_type${NC}"

    local check_key="${address}_${check_type}"
    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetAMLCheck\",\"$check_key\"]}"
}

# Function to get all KYC records
get_all_kyc() {
    echo -e "${YELLOW}Querying all KYC records${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetAllKYC\"]}"
}

# Function to get all AML checks for an address
get_all_aml() {
    local address=$1

    echo -e "${YELLOW}Querying all AML checks for: $address${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetAllAMLChecks\",\"$address\"]}"
}

# Function to handle errors
handle_error() {
    echo -e "${RED}Error: $1${NC}"
    exit 1
}

# Main execution
main() {
    # Check prerequisites
    check_peer_cli
    check_environment

    # Parse command
    case "$1" in
        "create-kyc")
            if [ $# -ne 7 ]; then
                handle_error "create-kyc requires 6 arguments"
            fi
            create_kyc "$2" "$3" "$4" "$5" "$6" "$7"
            ;;
        "approve-kyc")
            if [ $# -ne 4 ]; then
                handle_error "approve-kyc requires 3 arguments"
            fi
            approve_kyc "$2" "$3" "$4"
            ;;
        "reject-kyc")
            if [ $# -ne 4 ]; then
                handle_error "reject-kyc requires 3 arguments"
            fi
            reject_kyc "$2" "$3" "$4"
            ;;
        "create-aml")
            if [ $# -ne 5 ]; then
                handle_error "create-aml requires 4 arguments"
            fi
            create_aml "$2" "$3" "$4" "$5"
            ;;
        "update-aml")
            if [ $# -ne 6 ]; then
                handle_error "update-aml requires 5 arguments"
            fi
            update_aml "$2" "$3" "$4" "$5" "$6"
            ;;
        "check-compliance")
            if [ $# -ne 2 ]; then
                handle_error "check-compliance requires 1 argument"
            fi
            check_compliance "$2"
            ;;
        "get-kyc")
            if [ $# -ne 2 ]; then
                handle_error "get-kyc requires 1 argument"
            fi
            get_kyc "$2"
            ;;
        "get-aml")
            if [ $# -ne 3 ]; then
                handle_error "get-aml requires 2 arguments"
            fi
            get_aml "$2" "$3"
            ;;
        "get-all-kyc")
            get_all_kyc
            ;;
        "get-all-aml")
            if [ $# -ne 2 ]; then
                handle_error "get-all-aml requires 1 argument"
            fi
            get_all_aml "$2"
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        "")
            show_usage
            ;;
        *)
            handle_error "Unknown command: $1. Use 'help' for usage information."
            ;;
    esac
}

# Run main function
main "$@"
