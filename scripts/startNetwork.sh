#!/bin/bash

set -e

echo "Starting BondBridge Hyperledger Fabric Network..."

# Set environment variables
export FABRIC_CFG_PATH=${PWD}/network
export PATH=${PWD}/bin:${PWD}:$PATH

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

print_header() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}   $1${NC}"
    echo -e "${CYAN}================================${NC}"
}

# Check if required tools are installed
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    if ! command -v curl &> /dev/null; then
        print_error "curl is not installed. Please install curl first."
        exit 1
    fi
    
    print_status "Prerequisites check passed."
}

# Download Hyperledger Fabric binaries
download_fabric() {
    print_status "Downloading Hyperledger Fabric binaries..."
    
    if [ ! -d "bin" ]; then
        mkdir -p bin
        curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.0 1.5.0 0.4.22
        print_status "Fabric binaries downloaded successfully."
    else
        print_status "Fabric binaries already exist."
    fi
}

# Generate crypto materials
generate_crypto() {
    print_status "Generating crypto materials for 5 organizations..."
    
    if [ ! -d "organizations" ]; then
        mkdir -p organizations
        ./bin/cryptogen generate --config=./network/crypto-config.yaml --output=./organizations
        print_status "Crypto materials generated successfully."
    else
        print_status "Crypto materials already exist."
    fi
}

# Generate genesis block
generate_genesis() {
    print_status "Generating genesis block for 3-node Raft orderer cluster..."
    
    if [ ! -d "system-genesis-block" ]; then
        mkdir -p system-genesis-block
        ./bin/configtxgen -profile BondBridgeOrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block
        print_status "Genesis block generated successfully."
    else
        print_status "Genesis block already exists."
    fi
}

# Generate channel configuration
generate_channel_config() {
    print_status "Generating channel configuration for bondchannel..."
    
    if [ ! -f "channel-artifacts/bondchannel.tx" ]; then
        mkdir -p channel-artifacts
        ./bin/configtxgen -profile BondChannel -outputCreateChannelTx ./channel-artifacts/bondchannel.tx -channelID bondchannel
        print_status "Channel configuration generated successfully."
    else
        print_status "Channel configuration already exists."
    fi
}

# Start the network
start_network() {
    print_header "Starting BondBridge Network"
    
    print_status "Starting 3-node Raft orderer cluster..."
    docker-compose -f network/docker-compose.yaml up -d orderer0.bondbridge.com orderer1.bondbridge.com orderer2.bondbridge.com
    
    print_status "Waiting for orderer nodes to be ready..."
    sleep 10
    
    print_status "Starting Certificate Authority servers..."
    docker-compose -f network/docker-compose.yaml up -d ca.issuer.bondbridge.com ca.investor.bondbridge.com ca.regulator.bondbridge.com ca.marketmaker.bondbridge.com ca.custodian.bondbridge.com
    
    print_status "Waiting for CA servers to be ready..."
    sleep 15
    
    print_status "Starting peer nodes for all 5 organizations..."
    docker-compose -f network/docker-compose.yaml up -d peer0.issuer.bondbridge.com peer0.investor.bondbridge.com peer0.regulator.bondbridge.com peer0.marketmaker.bondbridge.com peer0.custodian.bondbridge.com
    
    print_status "Waiting for peer nodes to be ready..."
    sleep 20
    
    print_status "Network startup completed."
}

# Health check function
check_network_health() {
    print_header "Network Health Check"
    
    local all_healthy=true
    
    # Check orderer nodes
    print_status "Checking orderer nodes..."
    for orderer in orderer0 orderer1 orderer2; do
        if docker ps | grep -q "${orderer}.bondbridge.com"; then
            local status=$(docker ps --filter "name=${orderer}.bondbridge.com" --format "{{.Status}}")
            print_status "✓ ${orderer}: $status"
        else
            print_error "✗ ${orderer}: Not running"
            all_healthy=false
        fi
    done
    
    # Check peer nodes
    print_status "Checking peer nodes..."
    for org in issuer investor regulator marketmaker custodian; do
        if docker ps | grep -q "peer0.${org}.bondbridge.com"; then
            local status=$(docker ps --filter "name=peer0.${org}.bondbridge.com" --format "{{.Status}}")
            print_status "✓ peer0.${org}: $status"
        else
            print_error "✗ peer0.${org}: Not running"
            all_healthy=false
        fi
    done
    
    # Check CA servers
    print_status "Checking CA servers..."
    for org in issuer investor regulator marketmaker custodian; do
        if docker ps | grep -q "ca.${org}.bondbridge.com"; then
            local status=$(docker ps --filter "name=ca.${org}.bondbridge.com" --format "{{.Status}}")
            print_status "✓ ca.${org}: $status"
        else
            print_error "✗ ca.${org}: Not running"
            all_healthy=false
        fi
    done
    
    if [ "$all_healthy" = true ]; then
        print_status "✓ All network components are healthy!"
        return 0
    else
        print_error "✗ Some network components are not healthy."
        return 1
    fi
}

# Network status function
show_network_status() {
    print_header "Network Status"
    
    echo -e "${BLUE}Orderer Nodes:${NC}"
    docker ps --filter "name=orderer" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    
    echo -e "${BLUE}Peer Nodes:${NC}"
    docker ps --filter "name=peer" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    
    echo -e "${BLUE}CA Servers:${NC}"
    docker ps --filter "name=ca." --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    
    echo -e "${BLUE}Network Summary:${NC}"
    echo -e "• Orderer Nodes: 3 (Raft consensus)"
    echo -e "• Peer Nodes: 5 (1 per organization)"
    echo -e "• Organizations: Issuer, Investor, Regulator, Market Maker, Custodian"
    echo -e "• CA Servers: 5 (1 per organization)"
    echo ""
}

# Main execution
main() {
    print_header "BondBridge Network Startup"
    echo -e "${YELLOW}This will start a 5-organization network with 3-node Raft orderer cluster${NC}"
    echo ""
    
    check_prerequisites
    download_fabric
    generate_crypto
    generate_genesis
    generate_channel_config
    
    print_status "Starting network components..."
    start_network
    
    print_status "Performing health checks..."
    if check_network_health; then
        print_status "Network is ready for use!"
        show_network_status
        
        echo -e "${GREEN}Next steps:${NC}"
        echo -e "1. Deploy chaincodes: ./scripts/deployChaincode.sh"
        echo -e "2. Run demo: ./scripts/demo-script.sh"
        echo -e "3. Check logs: docker-compose -f network/docker-compose.yaml logs -f"
    else
        print_error "Network startup failed. Please check logs and try again."
        exit 1
    fi
}

# Error handling
trap 'print_error "Network startup failed!"; exit 1' ERR

# Execute main function
main "$@"

