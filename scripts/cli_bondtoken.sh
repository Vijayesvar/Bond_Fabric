#!/bin/bash

# BondToken Chaincode CLI Script
# This script provides a command-line interface for interacting with the BondToken chaincode

set -e

# Configuration
CHANNEL_NAME="mychannel"
CHAINCODE_NAME="bondtoken"
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
    echo -e "${BLUE}BondToken Chaincode CLI${NC}"
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  create-bond <id> <name> <currency> <face_value> <coupon_rate> <issue_date> <maturity_date> <status>"
    echo "  transfer-bond <bond_id> <from_owner> <to_owner>"
    echo "  get-bond <bond_id>"
    echo "  get-all-bonds"
    echo "  get-bonds-by-owner <owner>"
    echo "  update-status <bond_id> <new_status>"
    echo "  calculate-yield <bond_id> <current_price>"
    echo "  help"
    echo ""
    echo "Examples:"
    echo "  $0 create-bond BOND_001 'US Treasury Bond' USD 1000 5.0 2024-01-01 2029-01-01 ACTIVE"
    echo "  $0 transfer-bond BOND_001 alice bob"
    echo "  $0 get-bond BOND_001"
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

# Function to create a bond
create_bond() {
    local bond_id=$1
    local name=$2
    local currency=$3
    local face_value=$4
    local coupon_rate=$5
    local issue_date=$6
    local maturity_date=$7
    local status=$8
    
    echo -e "${YELLOW}Creating bond: $bond_id${NC}"
    
    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CreateBond\",\"$bond_id\",\"$name\",\"$currency\",\"$face_value\",\"$coupon_rate\",\"$issue_date\",\"$maturity_date\",\"$status\"]}" \
        --tls \
        --cafile $ORDERER_CA
    
    echo -e "${GREEN}✓ Bond $bond_id created successfully${NC}"
}

# Function to transfer a bond
transfer_bond() {
    local bond_id=$1
    local from_owner=$2
    local to_owner=$3
    
    echo -e "${YELLOW}Transferring bond $bond_id from $from_owner to $to_owner${NC}"
    
    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"TransferBond\",\"$bond_id\",\"$from_owner\",\"$to_owner\"]}" \
        --tls \
        --cafile $ORDERER_CA
    
    echo -e "${GREEN}✓ Bond $bond_id transferred successfully${NC}"
}

# Function to get a bond
get_bond() {
    local bond_id=$1
    
    echo -e "${YELLOW}Querying bond: $bond_id${NC}"
    
    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetBond\",\"$bond_id\"]}"
}

# Function to get all bonds
get_all_bonds() {
    echo -e "${YELLOW}Querying all bonds${NC}"
    
    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetAllBonds\"]}"
}

# Function to get bonds by owner
get_bonds_by_owner() {
    local owner=$1
    
    echo -e "${YELLOW}Querying bonds owned by: $owner${NC}"
    
    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetBondsByOwner\",\"$owner\"]}"
}

# Function to update bond status
update_bond_status() {
    local bond_id=$1
    local new_status=$2
    
    echo -e "${YELLOW}Updating bond $bond_id status to: $new_status${NC}"
    
    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"UpdateBondStatus\",\"$bond_id\",\"$new_status\"]}" \
        --tls \
        --cafile $ORDERER_CA
    
    echo -e "${GREEN}✓ Bond $bond_id status updated successfully${NC}"
}

# Function to calculate yield
calculate_yield() {
    local bond_id=$1
    local current_price=$2
    
    echo -e "${YELLOW}Calculating yield for bond $bond_id at price $current_price${NC}"
    
    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CalculateYield\",\"$bond_id\",\"$current_price\"]}"
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
        "create-bond")
            if [ $# -ne 9 ]; then
                handle_error "create-bond requires 8 arguments"
            fi
            create_bond "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9"
            ;;
        "transfer-bond")
            if [ $# -ne 4 ]; then
                handle_error "transfer-bond requires 3 arguments"
            fi
            transfer_bond "$2" "$3" "$4"
            ;;
        "get-bond")
            if [ $# -ne 2 ]; then
                handle_error "get-bond requires 1 argument"
            fi
            get_bond "$2"
            ;;
        "get-all-bonds")
            get_all_bonds
            ;;
        "get-bonds-by-owner")
            if [ $# -ne 2 ]; then
                handle_error "get-bonds-by-owner requires 1 argument"
            fi
            get_bonds_by_owner "$2"
            ;;
        "update-status")
            if [ $# -ne 3 ]; then
                handle_error "update-status requires 2 arguments"
            fi
            update_bond_status "$2" "$3"
            ;;
        "calculate-yield")
            if [ $# -ne 3 ]; then
                handle_error "calculate-yield requires 2 arguments"
            fi
            calculate_yield "$2" "$3"
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

