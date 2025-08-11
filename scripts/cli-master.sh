#!/bin/bash

# Master CLI Script for BondBridge Chaincodes
# This script provides a unified interface for all three chaincodes

set -e

# Configuration
CHANNEL_NAME="mychannel"
BONDTOKEN_CHAINCODE="bondtoken"
COMPLIANCE_CHAINCODE="compliance"
CORPORATEACTION_CHAINCODE="corporateaction"
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
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to display banner
show_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    BondBridge Master CLI                     ║"
    echo "║              Hyperledger Fabric Chaincodes                   ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Function to display main menu
show_main_menu() {
    echo -e "${BLUE}Main Menu:${NC}"
    echo "1. BondToken Operations"
    echo "2. Compliance Operations"
    echo "3. CorporateAction Operations"
    echo "4. Run All Tests"
    echo "5. Network Status"
    echo "6. Exit"
    echo ""
    echo -n "Select an option (1-6): "
}

# Function to display BondToken menu
show_bondtoken_menu() {
    echo -e "${GREEN}BondToken Operations:${NC}"
    echo "1. Create Bond"
    echo "2. Transfer Bond"
    echo "3. Get Bond Details"
    echo "4. Get All Bonds"
    echo "5. Get Bonds by Owner"
    echo "6. Update Bond Status"
    echo "7. Calculate Yield"
    echo "8. Back to Main Menu"
    echo ""
    echo -n "Select an option (1-8): "
}

# Function to display Compliance menu
show_compliance_menu() {
    echo -e "${PURPLE}Compliance Operations:${NC}"
    echo "1. Create KYC Record"
    echo "2. Approve KYC"
    echo "3. Reject KYC"
    echo "4. Create AML Check"
    echo "5. Update AML Check"
    echo "6. Check Compliance Status"
    echo "7. Get KYC Record"
    echo "8. Get AML Check"
    echo "9. Get All KYC Records"
    echo "10. Get All AML Checks"
    echo "11. Back to Main Menu"
    echo ""
    echo -n "Select an option (1-11): "
}

# Function to display CorporateAction menu
show_corporateaction_menu() {
    echo -e "${YELLOW}CorporateAction Operations:${NC}"
    echo "1. Create Coupon Payment"
    echo "2. Process Coupon Payment"
    echo "3. Create Redemption"
    echo "4. Process Redemption"
    echo "5. Get Coupon Details"
    echo "6. Get Redemption Details"
    echo "7. Get Coupons by Bond"
    echo "8. Get Redemptions by Bond"
    echo "9. Get Pending Coupons"
    echo "10. Get Pending Redemptions"
    echo "11. Calculate Coupon Amount"
    echo "12. Back to Main Menu"
    echo ""
    echo -n "Select an option (1-12): "
}

# Function to check if peer CLI is available
check_peer_cli() {
    if ! command -v peer &> /dev/null; then
        echo -e "${RED}Error: peer CLI not found${NC}"
        echo "Please ensure you are in the Fabric CLI environment"
        echo "Run: docker exec -it cli bash"
        return 1
    fi
    return 0
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

# Function to handle BondToken operations
handle_bondtoken_operations() {
    while true; do
        show_bondtoken_menu
        read -r choice

        case $choice in
            1)
                echo -n "Enter Bond ID: "
                read -r bond_id
                echo -n "Enter Issuer Name: "
                read -r issuer_name
                echo -n "Enter Currency: "
                read -r currency
                echo -n "Enter Face Value: "
                read -r face_value
                echo -n "Enter Coupon Rate (%): "
                read -r coupon_rate
                echo -n "Enter Issue Date (YYYY-MM-DD): "
                read -r issue_date
                echo -n "Enter Maturity Date (YYYY-MM-DD): "
                read -r maturity_date
                echo -n "Enter Status (ACTIVE/MATURED/DEFAULTED): "
                read -r status
                
                echo -e "${YELLOW}Creating bond: $bond_id${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $BONDTOKEN_CHAINCODE \
                    -c "{\"Args\":[\"IssueBond\",\"$bond_id\",\"ISSUER001\",\"$issuer_name\",\"$currency\",\"ISIN001\",\"AAA\",\"NONE\",\"$face_value\",\"$coupon_rate\",\"1000\",\"$maturity_date\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ Bond created successfully${NC}"
                ;;
            2)
                echo -n "Enter Bond ID: "
                read -r bond_id
                echo -n "Enter From Owner: "
                read -r from_owner
                echo -n "Enter To Owner: "
                read -r to_owner
                echo -n "Enter Quantity: "
                read -r quantity
                
                echo -e "${YELLOW}Transferring bond $bond_id${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $BONDTOKEN_CHAINCODE \
                    -c "{\"Args\":[\"Transfer\",\"$from_owner\",\"$to_owner\",\"$bond_id\",\"$quantity\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ Bond transferred successfully${NC}"
                ;;
            3)
                echo -n "Enter Bond ID: "
                read -r bond_id
                
                echo -e "${YELLOW}Querying bond: $bond_id${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $BONDTOKEN_CHAINCODE \
                    -c "{\"Args\":[\"GetBond\",\"$bond_id\"]}"
                ;;
            4)
                echo -e "${YELLOW}Querying all bonds${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $BONDTOKEN_CHAINCODE \
                    -c "{\"Args\":[\"GetAllBonds\"]}"
                ;;
            5)
                echo -n "Enter Owner Address: "
                read -r owner
                
                echo -e "${YELLOW}Querying bonds for owner: $owner${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $BONDTOKEN_CHAINCODE \
                    -c "{\"Args\":[\"GetBondHolders\"]}"
                ;;
            6)
                echo -n "Enter Bond ID: "
                read -r bond_id
                echo -n "Enter New Status: "
                read -r new_status
                
                echo -e "${YELLOW}Updating bond status${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $BONDTOKEN_CHAINCODE \
                    -c "{\"Args\":[\"UpdateBondStatus\",\"$bond_id\",\"$new_status\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ Bond status updated successfully${NC}"
                ;;
            7)
                echo -n "Enter Bond ID: "
                read -r bond_id
                echo -n "Enter Current Price: "
                read -r current_price
                
                echo -e "${YELLOW}Calculating yield${NC}"
                # This would call a yield calculation function
                echo "Yield calculation feature to be implemented"
                ;;
            8)
                break
                ;;
            *)
                echo -e "${RED}Invalid option. Please try again.${NC}"
                ;;
        esac
        
        echo ""
        echo "Press Enter to continue..."
        read -r
        clear
        show_banner
    done
}

# Function to handle Compliance operations
handle_compliance_operations() {
    while true; do
        show_compliance_menu
        read -r choice

        case $choice in
            1)
                echo -n "Enter Address: "
                read -r address
                echo -n "Enter Full Name: "
                read -r full_name
                echo -n "Enter Date of Birth (YYYY-MM-DD): "
                read -r dob
                echo -n "Enter Nationality: "
                read -r nationality
                echo -n "Enter ID Type: "
                read -r id_type
                echo -n "Enter ID Number: "
                read -r id_number
                
                echo -e "${YELLOW}Creating KYC record for: $address${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"CreateKYC\",\"$address\",\"$full_name\",\"$dob\",\"$nationality\",\"$id_type\",\"$id_number\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ KYC record created successfully${NC}"
                ;;
            2)
                echo -n "Enter Address: "
                read -r address
                echo -n "Enter Approved By: "
                read -r approved_by
                echo -n "Enter Risk Level (LOW/MEDIUM/HIGH): "
                read -r risk_level
                
                echo -e "${YELLOW}Approving KYC for: $address${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"ApproveKYC\",\"$address\",\"$approved_by\",\"$risk_level\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ KYC approved successfully${NC}"
                ;;
            3)
                echo -n "Enter Address: "
                read -r address
                echo -n "Enter Rejected By: "
                read -r rejected_by
                echo -n "Enter Reason: "
                read -r reason
                
                echo -e "${YELLOW}Rejecting KYC for: $address${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"RejectKYC\",\"$address\",\"$rejected_by\",\"$reason\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ KYC rejected successfully${NC}"
                ;;
            4)
                echo -n "Enter Address: "
                read -r address
                echo -n "Enter Check Type (SANCTIONS/PEP/ADVERSE_MEDIA): "
                read -r check_type
                echo -n "Enter Risk Score (1-10): "
                read -r risk_score
                echo -n "Enter Details: "
                read -r details
                
                echo -e "${YELLOW}Creating AML check for: $address${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"CreateAMLCheck\",\"$address\",\"$check_type\",\"$risk_score\",\"$details\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ AML check created successfully${NC}"
                ;;
            5)
                echo -n "Enter Address: "
                read -r address
                echo -n "Enter Check Type: "
                read -r check_type
                echo -n "Enter Status (PASSED/FAILED/PENDING): "
                read -r status
                echo -n "Enter Risk Score (1-10): "
                read -r risk_score
                echo -n "Enter Details: "
                read -r details
                
                echo -e "${YELLOW}Updating AML check for: $address${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"UpdateAMLCheck\",\"$address\",\"$check_type\",\"$status\",\"$risk_score\",\"$details\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ AML check updated successfully${NC}"
                ;;
            6)
                echo -n "Enter Address: "
                read -r address
                
                echo -e "${YELLOW}Checking compliance for: $address${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"CheckCompliance\",\"$address\"]}"
                ;;
            7)
                echo -n "Enter Address: "
                read -r address
                
                echo -e "${YELLOW}Querying KYC record for: $address${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"GetKYC\",\"$address\"]}"
                ;;
            8)
                echo -n "Enter Address: "
                read -r address
                echo -n "Enter Check Type: "
                read -r check_type
                
                local check_key="${address}_${check_type}"
                echo -e "${YELLOW}Querying AML check for: $address${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"GetAMLCheck\",\"$check_key\"]}"
                ;;
            9)
                echo -e "${YELLOW}Querying all KYC records${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"GetAllKYC\"]}"
                ;;
            10)
                echo -n "Enter Address: "
                read -r address
                
                echo -e "${YELLOW}Querying all AML checks for: $address${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $COMPLIANCE_CHAINCODE \
                    -c "{\"Args\":[\"GetAllAMLChecks\",\"$address\"]}"
                ;;
            11)
                break
                ;;
            *)
                echo -e "${RED}Invalid option. Please try again.${NC}"
                ;;
        esac
        
        echo ""
        echo "Press Enter to continue..."
        read -r
        clear
        show_banner
    done
}

# Function to handle CorporateAction operations
handle_corporateaction_operations() {
    while true; do
        show_corporateaction_menu
        read -r choice

        case $choice in
            1)
                echo -n "Enter Bond ID: "
                read -r bond_id
                echo -n "Enter Payment Date (YYYY-MM-DD): "
                read -r payment_date
                echo -n "Enter Amount: "
                read -r amount
                
                echo -e "${YELLOW}Creating coupon payment for bond: $bond_id${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"CreateCouponPayment\",\"$bond_id\",\"$payment_date\",\"$amount\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ Coupon payment created successfully${NC}"
                ;;
            2)
                echo -n "Enter Coupon ID: "
                read -r coupon_id
                
                echo -e "${YELLOW}Processing coupon payment: $coupon_id${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"ProcessCouponPayment\",\"$coupon_id\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ Coupon payment processed successfully${NC}"
                ;;
            3)
                echo -n "Enter Bond ID: "
                read -r bond_id
                echo -n "Enter Redemption Date (YYYY-MM-DD): "
                read -r redemption_date
                echo -n "Enter Amount: "
                read -r amount
                
                echo -e "${YELLOW}Creating redemption for bond: $bond_id${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"CreateRedemption\",\"$bond_id\",\"$redemption_date\",\"$amount\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ Redemption created successfully${NC}"
                ;;
            4)
                echo -n "Enter Redemption ID: "
                read -r redemption_id
                
                echo -e "${YELLOW}Processing redemption: $redemption_id${NC}"
                peer chaincode invoke \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"ProcessRedemption\",\"$redemption_id\"]}" \
                    --tls \
                    --cafile $ORDERER_CA
                echo -e "${GREEN}✓ Redemption processed successfully${NC}"
                ;;
            5)
                echo -n "Enter Coupon ID: "
                read -r coupon_id
                
                echo -e "${YELLOW}Querying coupon payment: $coupon_id${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"GetCouponPayment\",\"$coupon_id\"]}"
                ;;
            6)
                echo -n "Enter Redemption ID: "
                read -r redemption_id
                
                echo -e "${YELLOW}Querying redemption: $redemption_id${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"GetRedemption\",\"$redemption_id\"]}"
                ;;
            7)
                echo -n "Enter Bond ID: "
                read -r bond_id
                
                echo -e "${YELLOW}Querying coupon payments for bond: $bond_id${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"GetCouponPaymentsByBond\",\"$bond_id\"]}"
                ;;
            8)
                echo -n "Enter Bond ID: "
                read -r bond_id
                
                echo -e "${YELLOW}Querying redemptions for bond: $bond_id${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"GetRedemptionsByBond\",\"$bond_id\"]}"
                ;;
            9)
                echo -e "${YELLOW}Querying pending coupon payments${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"GetPendingCouponPayments\"]}"
                ;;
            10)
                echo -e "${YELLOW}Querying pending redemptions${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"GetPendingRedemptions\"]}"
                ;;
            11)
                echo -n "Enter Bond ID: "
                read -r bond_id
                echo -n "Enter Face Value: "
                read -r face_value
                echo -n "Enter Coupon Rate (%): "
                read -r coupon_rate
                
                echo -e "${YELLOW}Calculating coupon amount for bond: $bond_id${NC}"
                peer chaincode query \
                    -C $CHANNEL_NAME \
                    -n $CORPORATEACTION_CHAINCODE \
                    -c "{\"Args\":[\"CalculateCouponAmount\",\"$bond_id\",\"$face_value\",\"$coupon_rate\"]}"
                ;;
            12)
                break
                ;;
            *)
                echo -e "${RED}Invalid option. Please try again.${NC}"
                ;;
        esac
        
        echo ""
        echo "Press Enter to continue..."
        read -r
        clear
        show_banner
    done
}

# Function to run all tests
run_all_tests() {
    echo -e "${YELLOW}Running all chaincode tests...${NC}"
    
    if [ -f "scripts/testChaincode.sh" ]; then
        bash scripts/testChaincode.sh
    else
        echo -e "${RED}Test script not found${NC}"
    fi
}

# Function to check network status
check_network_status() {
    echo -e "${YELLOW}Checking network status...${NC}"
    
    # Check if peer CLI is available
    if ! check_peer_cli; then
        echo -e "${RED}Network not accessible - peer CLI not found${NC}"
        return
    fi
    
    # Check channel info
    echo -e "${CYAN}Channel Information:${NC}"
    peer channel getinfo -c $CHANNEL_NAME 2>/dev/null || echo "Channel not accessible"
    
    # Check chaincode list
    echo -e "${CYAN}Installed Chaincodes:${NC}"
    peer chaincode list --installed 2>/dev/null || echo "Cannot query installed chaincodes"
    
    echo -e "${CYAN}Instantiated Chaincodes:${NC}"
    peer chaincode list --instantiated -C $CHANNEL_NAME 2>/dev/null || echo "Cannot query instantiated chaincodes"
}

# Main execution
main() {
    clear
    show_banner
    
    # Check prerequisites
    if ! check_peer_cli; then
        echo -e "${YELLOW}Warning: Running in demo mode (peer CLI not available)${NC}"
        echo "Some functions may not work without Fabric CLI environment"
        echo ""
    fi
    
    # Set up environment
    check_environment
    
    # Main menu loop
    while true; do
        show_main_menu
        read -r choice
        
        case $choice in
            1)
                clear
                show_banner
                handle_bondtoken_operations
                clear
                show_banner
                ;;
            2)
                clear
                show_banner
                handle_compliance_operations
                clear
                show_banner
                ;;
            3)
                clear
                show_banner
                handle_corporateaction_operations
                clear
                show_banner
                ;;
            4)
                run_all_tests
                echo ""
                echo "Press Enter to continue..."
                read -r
                clear
                show_banner
                ;;
            5)
                check_network_status
                echo ""
                echo "Press Enter to continue..."
                read -r
                clear
                show_banner
                ;;
            6)
                echo -e "${GREEN}Thank you for using BondBridge Master CLI!${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}Invalid option. Please try again.${NC}"
                echo ""
                ;;
        esac
    done
}

# Run main function
main "$@"
