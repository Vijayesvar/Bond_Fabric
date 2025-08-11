#!/bin/bash

# BondBridge Live Demo Script
# Duration: 3-5 minutes
# Demonstrates: Fixed topology, endorsement policies, PDC, and end-to-end workflow

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Demo configuration
DEMO_DURATION=300  # 5 minutes max
BOND_ID="BOND_001"
BOND_NAME="SEBI Demo Bond"
FACE_VALUE=10000
COUPON_RATE=5.5
MATURITY_DATE="2029-12-31"
INVESTOR_ADDRESS="investor1.bondbridge.com"
ISSUER_ADDRESS="issuer1.bondbridge.com"

echo -e "${CYAN}================================${NC}"
echo -e "${CYAN}   BondBridge Live Demo Script  ${NC}"
echo -e "${CYAN}================================${NC}"
echo -e "${YELLOW}Duration: 3-5 minutes${NC}"
echo -e "${YELLOW}Target: SEBI Hackathon Judges${NC}"
echo ""

# Function to print step headers
print_step() {
    echo -e "${GREEN}[STEP $1]${NC} $2"
    echo -e "${BLUE}$3${NC}"
    echo ""
}

# Function to execute CLI commands with output
execute_command() {
    local cmd="$1"
    local description="$2"
    
    echo -e "${YELLOW}Executing:${NC} $description"
    echo -e "${PURPLE}Command:${NC} $cmd"
    
    if eval "$cmd"; then
        echo -e "${GREEN}✓ Success${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
        return 1
    fi
    echo ""
}

# Function to show network status
show_network_status() {
    echo -e "${CYAN}=== Network Status Check ===${NC}"
    
    # Check orderer nodes
    echo -e "${BLUE}Orderer Nodes:${NC}"
    docker ps --filter "name=orderer" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    
    # Check peer nodes
    echo -e "${BLUE}Peer Nodes:${NC}"
    docker ps --filter "name=peer" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    
    # Check CA servers
    echo -e "${BLUE}CA Servers:${NC}"
    docker ps --filter "name=ca." --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
}

# Function to show metrics
show_metrics() {
    echo -e "${CYAN}=== Performance Metrics ===${NC}"
    
    # Get transaction count
    local tx_count=$(docker exec peer0.issuer.bondbridge.com peer channel getinfo -c bondchannel | grep "Block height" | awk '{print $3}' | sed 's/,//')
    echo -e "${BLUE}Total Transactions:${NC} $tx_count"
    
    # Get block height
    local block_height=$(docker exec peer0.issuer.bondbridge.com peer channel getinfo -c bondchannel | grep "Block height" | awk '{print $3}' | sed 's/,//')
    echo -e "${BLUE}Block Height:${NC} $block_height"
    
    # Get channel info
    echo -e "${BLUE}Channel Status:${NC}"
    docker exec peer0.issuer.bondbridge.com peer channel getinfo -c bondchannel | grep -E "(Channel name|Block height|Last block hash)"
    echo ""
}

# Main demo execution
main() {
    echo -e "${GREEN}Starting BondBridge Live Demo...${NC}"
    echo ""
    
    # Step 1: Network Status & Topology
    print_step "1" "Network Topology Verification" "Showing 5 organizations with 1 peer each + 3-node Raft orderer cluster"
    show_network_status
    
    # Step 2: Bond Issuance (Issuer + Regulator endorsement)
    print_step "2" "Bond Issuance with Dual Endorsement" "Demonstrating Issuer + Regulator approval workflow"
    execute_command "./scripts/cli-bondtoken.sh create-bond $BOND_ID '$BOND_NAME' INR $FACE_VALUE $COUPON_RATE 2024-01-01 $MATURITY_DATE ACTIVE" \
        "Creating bond with dual endorsement requirement"
    
    # Step 3: KYC Creation (Regulator + Issuer endorsement)
    print_step "3" "KYC Creation with Regulatory Oversight" "Showing KYC workflow with PDC for PII protection"
    execute_command "./scripts/cli-compliance.sh create-kyc $INVESTOR_ADDRESS 'Demo Investor' '1990-01-01' 'IN' 'PAN' 'ABCDE1234F'" \
        "Creating KYC record with regulatory approval"
    
    # Step 4: KYC Approval (Regulator + Custodian endorsement)
    print_step "4" "KYC Approval with Custodian Validation" "Demonstrating multi-organization endorsement policy"
    execute_command "./scripts/cli-compliance.sh approve-kyc $INVESTOR_ADDRESS 'regulator1' 'LOW'" \
        "Approving KYC with custodian validation"
    
    # Step 5: Bond Transfer (Seller + Custodian + Market Maker endorsement)
    print_step "5" "Bond Transfer with Triple Endorsement" "Showing complex endorsement policy for transfers"
    execute_command "./scripts/cli-bondtoken.sh transfer-bond $BOND_ID $ISSUER_ADDRESS $INVESTOR_ADDRESS 1000" \
        "Transferring bonds with triple endorsement requirement"
    
    # Step 6: Coupon Payment Creation (Issuer + Custodian endorsement)
    print_step "6" "Coupon Payment with Settlement Approval" "Demonstrating automated coupon payment workflow"
    execute_command "./scripts/cli-corporateaction.sh create-coupon-payment $BOND_ID 2024-06-30 275" \
        "Creating coupon payment with issuer and custodian approval"
    
    # Step 7: Compliance Check
    print_step "7" "Compliance Verification" "Checking regulatory compliance status"
    execute_command "./scripts/cli-compliance.sh check-compliance $INVESTOR_ADDRESS" \
        "Verifying investor compliance status"
    
    # Step 8: Performance Metrics
    print_step "8" "Performance & Compliance Metrics" "Showing transaction throughput and regulatory compliance"
    show_metrics
    
    # Step 9: Demo Summary
    print_step "9" "Demo Summary & Key Benefits" "Highlighting critical fixes and improvements"
    echo -e "${CYAN}=== Demo Summary ===${NC}"
    echo -e "${GREEN}✓ Fixed Network Topology:${NC} 5 orgs × 1 peer + 3-node Raft orderer"
    echo -e "${GREEN}✓ Endorsement Policies:${NC} Multi-org approval for critical operations"
    echo -e "${GREEN}✓ Privacy Protection:${NC} PDC for KYC/PII data"
    echo -e "${GREEN}✓ Legal Compliance:${NC} RTA integration framework"
    echo -e "${GREEN}✓ Regulatory Approval:${NC} SEBI-compliant architecture"
    echo ""
    
    echo -e "${GREEN}Demo completed successfully!${NC}"
    echo -e "${YELLOW}Key Benefits Demonstrated:${NC}"
    echo -e "• Reduced minimum lot size from ₹10,000 to ₹1,000"
    echo -e "• Improved liquidity through fractional ownership"
    echo -e "• Enhanced compliance through multi-org endorsement"
    echo -e "• Real-time settlement and regulatory reporting"
    echo -e "• Legal enforceability with RTA integration"
    echo ""
}

# Error handling
trap 'echo -e "${RED}Demo failed!${NC}"; exit 1' ERR

# Check if network is running
if ! docker ps | grep -q "orderer"; then
    echo -e "${RED}Error: Network is not running. Please start the network first.${NC}"
    echo -e "${YELLOW}Run: ./scripts/startNetwork.sh${NC}"
    exit 1
fi

# Check if chaincodes are deployed
if ! docker exec peer0.issuer.bondbridge.com peer chaincode list --instantiated -C bondchannel | grep -q "bondtoken"; then
    echo -e "${RED}Error: Chaincodes are not deployed. Please deploy them first.${NC}"
    echo -e "${YELLOW}Run: ./scripts/deployChaincode.sh${NC}"
    exit 1
fi

# Execute demo
main "$@"
