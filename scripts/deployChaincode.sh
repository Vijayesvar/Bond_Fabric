#!/bin/bash

set -e

echo "Deploying BondBridge Chaincode..."

# Set environment variables
export FABRIC_CFG_PATH=${PWD}/network
export PATH=${PWD}/bin:${PWD}:$PATH

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if network is running
check_network() {
    print_status "Checking if network is running..."
    
    if ! docker ps | grep -q "orderer.bondbridge.com"; then
        print_error "Network is not running. Please start the network first: ./scripts/startNetwork.sh"
        exit 1
    fi
    
    print_status "Network is running."
}

# Package chaincode
package_chaincode() {
    print_status "Packaging chaincode..."
    
    if [ ! -d "chaincode" ]; then
        print_error "Chaincode directory not found. Please create chaincode first."
        exit 1
    fi
    
    cd chaincode
    
    # Package BondToken chaincode
    if [ -d "bondtoken" ]; then
        print_status "Packaging BondToken chaincode..."
        peer lifecycle chaincode package bondtoken.tar.gz --path ./bondtoken --lang golang --label bondtoken_1.0
    fi
    
    # Package Compliance chaincode
    if [ -d "compliance" ]; then
        print_status "Packaging Compliance chaincode..."
        peer lifecycle chaincode package compliance.tar.gz --path ./compliance --lang golang --label compliance_1.0
    fi
    
    # Package CorporateAction chaincode
    if [ -d "corporateaction" ]; then
        print_status "Packaging CorporateAction chaincode..."
        peer lifecycle chaincode package corporateaction.tar.gz --path ./corporateaction --lang golang --label corporateaction_1.0
    fi
    
    cd ..
}

# Install chaincode on issuer peer
install_on_issuer() {
    print_status "Installing chaincode on issuer peer..."
    
    export CORE_PEER_TLS_ENABLED=false
    export CORE_PEER_LOCALMSPID=IssuerMSP
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/users/Admin@issuer.bondbridge.com/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/peers/peer0.issuer.bondbridge.com/tls/ca.crt
    export CORE_PEER_ADDRESS=localhost:7051
    
    # Install BondToken chaincode
    if [ -f "chaincode/bondtoken.tar.gz" ]; then
        peer lifecycle chaincode install chaincode/bondtoken.tar.gz
        print_status "BondToken chaincode installed on issuer peer."
    fi
    
    # Install Compliance chaincode
    if [ -f "chaincode/compliance.tar.gz" ]; then
        peer lifecycle chaincode install chaincode/compliance.tar.gz
        print_status "Compliance chaincode installed on issuer peer."
    fi
    
    # Install CorporateAction chaincode
    if [ -f "chaincode/corporateaction.tar.gz" ]; then
        peer lifecycle chaincode install chaincode/corporateaction.tar.gz
        print_status "CorporateAction chaincode installed on issuer peer."
    fi
}

# Install chaincode on investor peer
install_on_investor() {
    print_status "Installing chaincode on investor peer..."
    
    export CORE_PEER_TLS_ENABLED=false
    export CORE_PEER_LOCALMSPID=InvestorMSP
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/investor.bondbridge.com/users/Admin@investor.bondbridge.com/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/investor.bondbridge.com/peers/peer0.investor.bondbridge.com/tls/ca.crt
    export CORE_PEER_ADDRESS=localhost:8051
    
    # Install BondToken chaincode
    if [ -f "chaincode/bondtoken.tar.gz" ]; then
        peer lifecycle chaincode install chaincode/bondtoken.tar.gz
        print_status "BondToken chaincode installed on investor peer."
    fi
    
    # Install Compliance chaincode
    if [ -f "chaincode/compliance.tar.gz" ]; then
        peer lifecycle chaincode install chaincode/compliance.tar.gz
        print_status "Compliance chaincode installed on investor peer."
    fi
    
    # Install CorporateAction chaincode
    if [ -f "chaincode/corporateaction.tar.gz" ]; then
        peer lifecycle chaincode install chaincode/corporateaction.tar.gz
        print_status "CorporateAction chaincode installed on investor peer."
    fi
}

# Approve chaincode definitions
approve_chaincode() {
    print_status "Approving chaincode definitions..."
    
    export CORE_PEER_TLS_ENABLED=false
    export CORE_PEER_LOCALMSPID=IssuerMSP
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/users/Admin@issuer.bondbridge.com/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/peers/peer0.issuer.bondbridge.com/tls/ca.crt
    export CORE_PEER_ADDRESS=localhost:7051
    
    # Get package IDs
    BONDTOKEN_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "bondtoken_1.0" | awk '{print $3}' | sed 's/,//')
    COMPLIANCE_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "compliance_1.0" | awk '{print $3}' | sed 's/,//')
    CORPORATEACTION_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "corporateaction_1.0" | awk '{print $3}' | sed 's/,//')
    
    # Approve BondToken
    if [ ! -z "$BONDTOKEN_PACKAGE_ID" ]; then
        peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name bondtoken --version 1.0 --package-id $BONDTOKEN_PACKAGE_ID --sequence 1
        print_status "BondToken chaincode approved by issuer."
    fi
    
    # Approve Compliance
    if [ ! -z "$COMPLIANCE_PACKAGE_ID" ]; then
        peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name compliance --version 1.0 --package-id $COMPLIANCE_PACKAGE_ID --sequence 1
        print_status "Compliance chaincode approved by issuer."
    fi
    
    # Approve CorporateAction
    if [ ! -z "$CORPORATEACTION_PACKAGE_ID" ]; then
        peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name corporateaction --version 1.0 --package-id $CORPORATEACTION_PACKAGE_ID --sequence 1
        print_status "CorporateAction chaincode approved by issuer."
    fi
    
    # Approve by investor
    export CORE_PEER_LOCALMSPID=InvestorMSP
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/investor.bondbridge.com/users/Admin@investor.bondbridge.com/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/investor.bondbridge.com/peers/peer0.investor.bondbridge.com/tls/ca.crt
    export CORE_PEER_ADDRESS=localhost:8051
    
    if [ ! -z "$BONDTOKEN_PACKAGE_ID" ]; then
        peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name bondtoken --version 1.0 --package-id $BONDTOKEN_PACKAGE_ID --sequence 1
        print_status "BondToken chaincode approved by investor."
    fi
    
    if [ ! -z "$COMPLIANCE_PACKAGE_ID" ]; then
        peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name compliance --version 1.0 --package-id $COMPLIANCE_PACKAGE_ID --sequence 1
        print_status "Compliance chaincode approved by investor."
    fi
    
    if [ ! -z "$CORPORATEACTION_PACKAGE_ID" ]; then
        peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name corporateaction --version 1.0 --package-id $CORPORATEACTION_PACKAGE_ID --sequence 1
        print_status "CorporateAction chaincode approved by investor."
    fi
}

# Commit chaincode definitions
commit_chaincode() {
    print_status "Committing chaincode definitions..."
    
    export CORE_PEER_TLS_ENABLED=false
    export CORE_PEER_LOCALMSPID=IssuerMSP
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/users/Admin@issuer.bondbridge.com/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/peers/peer0.issuer.bondbridge.com/tls/ca.crt
    export CORE_PEER_ADDRESS=localhost:7051
    
    # Commit BondToken
    if [ -f "chaincode/bondtoken.tar.gz" ]; then
        peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name bondtoken --version 1.0 --sequence 1
        print_status "BondToken chaincode committed to bondchannel."
    fi
    
    # Commit Compliance
    if [ -f "chaincode/compliance.tar.gz" ]; then
        peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name compliance --version 1.0 --sequence 1
        print_status "Compliance chaincode committed to bondchannel."
    fi
    
    # Commit CorporateAction
    if [ -f "chaincode/corporateaction.tar.gz" ]; then
        peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com --channelID bondchannel --name corporateaction --version 1.0 --sequence 1
        print_status "CorporateAction chaincode committed to bondchannel."
    fi
}

# Test chaincode
test_chaincode() {
    print_status "Testing chaincode functionality..."
    
    export CORE_PEER_TLS_ENABLED=false
    export CORE_PEER_LOCALMSPID=IssuerMSP
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/users/Admin@issuer.bondbridge.com/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/issuer.bondbridge.com/peers/peer0.issuer.bondbridge.com/tls/ca.crt
    export CORE_PEER_ADDRESS=localhost:7051
    
    # Test BondToken initialization
    if [ -f "chaincode/bondtoken.tar.gz" ]; then
        peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com -C bondchannel -n bondtoken --isInit -c '{"Args":["Init"]}'
        print_status "BondToken chaincode initialized successfully."
    fi
    
    # Test Compliance initialization
    if [ -f "chaincode/compliance.tar.gz" ]; then
        peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com -C bondchannel -n compliance --isInit -c '{"Args":["Init"]}'
        print_status "Compliance chaincode initialized successfully."
    fi
    
    # Test CorporateAction initialization
    if [ -f "chaincode/corporateaction.tar.gz" ]; then
        peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.bondbridge.com -C bondchannel -n corporateaction --isInit -c '{"Args":["Init"]}'
        print_status "CorporateAction chaincode initialized successfully."
    fi
}

# Main execution
main() {
    print_status "Starting chaincode deployment..."
    
    check_network
    package_chaincode
    install_on_issuer
    install_on_investor
    approve_chaincode
    commit_chaincode
    test_chaincode
    
    print_status "Chaincode deployment completed successfully!"
    print_status "All smart contracts are now active on the bondchannel."
    
    echo ""
    print_status "Next steps:"
    echo "1. Start API server: cd api && npm start"
    echo "2. Start frontend: cd frontend && npm start"
    echo "3. Test the platform with sample bond issuance and trading"
}

# Run main function
main "$@"

