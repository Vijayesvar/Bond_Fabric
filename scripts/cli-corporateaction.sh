#!/bin/bash

# CorporateAction Chaincode CLI Script
# This script provides a command-line interface for interacting with the CorporateAction chaincode

set -e

# Configuration
CHANNEL_NAME="mychannel"
CHAINCODE_NAME="corporateaction"
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
    echo -e "${BLUE}CorporateAction Chaincode CLI${NC}"
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  create-coupon <bond_id> <payment_date> <amount>"
    echo "  process-coupon <coupon_id>"
    echo "  create-redemption <bond_id> <redemption_date> <amount>"
    echo "  process-redemption <redemption_id>"
    echo "  get-coupon <coupon_id>"
    echo "  get-redemption <redemption_id>"
    echo "  get-coupons-by-bond <bond_id>"
    echo "  get-redemptions-by-bond <bond_id>"
    echo "  get-pending-coupons"
    echo "  get-pending-redemptions"
    echo "  calculate-coupon <bond_id> <face_value> <coupon_rate>"
    echo "  help"
    echo ""
    echo "Examples:"
    echo "  $0 create-coupon BOND_001 2024-06-01 50.00"
    echo "  $0 process-coupon COUPON_001"
    echo "  $0 create-redemption BOND_001 2029-01-01 1000.00"
    echo "  $0 calculate-coupon BOND_001 1000.00 5.0"
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

# Function to create coupon payment
create_coupon() {
    local bond_id=$1
    local payment_date=$2
    local amount=$3

    echo -e "${YELLOW}Creating coupon payment for bond: $bond_id${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CreateCouponPayment\",\"$bond_id\",\"$payment_date\",\"$amount\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ Coupon payment created successfully for bond $bond_id${NC}"
}

# Function to process coupon payment
process_coupon() {
    local coupon_id=$1

    echo -e "${YELLOW}Processing coupon payment: $coupon_id${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"ProcessCouponPayment\",\"$coupon_id\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ Coupon payment processed successfully: $coupon_id${NC}"
}

# Function to create redemption
create_redemption() {
    local bond_id=$1
    local redemption_date=$2
    local amount=$3

    echo -e "${YELLOW}Creating redemption for bond: $bond_id${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CreateRedemption\",\"$bond_id\",\"$redemption_date\",\"$amount\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ Redemption created successfully for bond $bond_id${NC}"
}

# Function to process redemption
process_redemption() {
    local redemption_id=$1

    echo -e "${YELLOW}Processing redemption: $redemption_id${NC}"

    peer chaincode invoke \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"ProcessRedemption\",\"$redemption_id\"]}" \
        --tls \
        --cafile $ORDERER_CA

    echo -e "${GREEN}✓ Redemption processed successfully: $redemption_id${NC}"
}

# Function to get coupon payment
get_coupon() {
    local coupon_id=$1

    echo -e "${YELLOW}Querying coupon payment: $coupon_id${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetCouponPayment\",\"$coupon_id\"]}"
}

# Function to get redemption
get_redemption() {
    local redemption_id=$1

    echo -e "${YELLOW}Querying redemption: $redemption_id${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetRedemption\",\"$redemption_id\"]}"
}

# Function to get coupon payments by bond
get_coupons_by_bond() {
    local bond_id=$1

    echo -e "${YELLOW}Querying coupon payments for bond: $bond_id${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetCouponPaymentsByBond\",\"$bond_id\"]}"
}

# Function to get redemptions by bond
get_redemptions_by_bond() {
    local bond_id=$1

    echo -e "${YELLOW}Querying redemptions for bond: $bond_id${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetRedemptionsByBond\",\"$bond_id\"]}"
}

# Function to get pending coupon payments
get_pending_coupons() {
    echo -e "${YELLOW}Querying pending coupon payments${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetPendingCouponPayments\"]}"
}

# Function to get pending redemptions
get_pending_redemptions() {
    echo -e "${YELLOW}Querying pending redemptions${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"GetPendingRedemptions\"]}"
}

# Function to calculate coupon amount
calculate_coupon() {
    local bond_id=$1
    local face_value=$2
    local coupon_rate=$3

    echo -e "${YELLOW}Calculating coupon amount for bond: $bond_id${NC}"

    peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c "{\"Args\":[\"CalculateCouponAmount\",\"$bond_id\",\"$face_value\",\"$coupon_rate\"]}"
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
        "create-coupon")
            if [ $# -ne 4 ]; then
                handle_error "create-coupon requires 3 arguments"
            fi
            create_coupon "$2" "$3" "$4"
            ;;
        "process-coupon")
            if [ $# -ne 2 ]; then
                handle_error "process-coupon requires 1 argument"
            fi
            process_coupon "$2"
            ;;
        "create-redemption")
            if [ $# -ne 4 ]; then
                handle_error "create-redemption requires 3 arguments"
            fi
            create_redemption "$2" "$3" "$4"
            ;;
        "process-redemption")
            if [ $# -ne 2 ]; then
                handle_error "process-redemption requires 1 argument"
            fi
            process_redemption "$2"
            ;;
        "get-coupon")
            if [ $# -ne 2 ]; then
                handle_error "get-coupon requires 1 argument"
            fi
            get_coupon "$2"
            ;;
        "get-redemption")
            if [ $# -ne 2 ]; then
                handle_error "get-redemption requires 1 argument"
            fi
            get_redemption "$2"
            ;;
        "get-coupons-by-bond")
            if [ $# -ne 2 ]; then
                handle_error "get-coupons-by-bond requires 1 argument"
            fi
            get_coupons_by_bond "$2"
            ;;
        "get-redemptions-by-bond")
            if [ $# -ne 2 ]; then
                handle_error "get-redemptions-by-bond requires 1 argument"
            fi
            get_redemptions_by_bond "$2"
            ;;
        "get-pending-coupons")
            get_pending_coupons
            ;;
        "get-pending-redemptions")
            get_pending_redemptions
            ;;
        "calculate-coupon")
            if [ $# -ne 4 ]; then
                handle_error "calculate-coupon requires 3 arguments"
            fi
            calculate_coupon "$2" "$3" "$4"
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
