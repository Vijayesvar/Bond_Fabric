# BondBridge CLI Tools

This directory contains command-line interface tools for interacting with the BondBridge Hyperledger Fabric chaincodes.

## Overview

The CLI tools provide a comprehensive interface for:
- **BondToken Chaincode**: Bond issuance, transfer, and management
- **Compliance Chaincode**: KYC/AML operations and compliance checks
- **CorporateAction Chaincode**: Coupon payments and bond redemptions

## Available Scripts

### 1. Master CLI (`cli-master.sh`)
A unified, menu-driven interface for all chaincode operations.

**Features:**
- Interactive menu system
- Color-coded output
- Comprehensive error handling
- Network status checking
- Integrated test runner

**Usage:**
```bash
./scripts/cli-master.sh
```

### 2. Individual Chaincode CLIs

#### BondToken CLI (`cli-bondtoken.sh`)
Direct command-line interface for bond operations.

**Commands:**
```bash
# Create a new bond
./scripts/cli-bondtoken.sh create-bond <id> <name> <currency> <face_value> <coupon_rate> <issue_date> <maturity_date> <status>

# Transfer bonds between owners
./scripts/cli-bondtoken.sh transfer-bond <bond_id> <from_owner> <to_owner>

# Query bond information
./scripts/cli-bondtoken.sh get-bond <bond_id>
./scripts/cli-bondtoken.sh get-all-bonds
./scripts/cli-bondtoken.sh get-bonds-by-owner <owner>

# Update bond status
./scripts/cli-bondtoken.sh update-status <bond_id> <new_status>

# Calculate yield
./scripts/cli-bondtoken.sh calculate-yield <bond_id> <current_price>
```

**Examples:**
```bash
# Create a US Treasury bond
./scripts/cli-bondtoken.sh create-bond BOND_001 "US Treasury Bond" USD 1000 5.0 2024-01-01 2029-01-01 ACTIVE

# Transfer bonds
./scripts/cli-bondtoken.sh transfer-bond BOND_001 alice bob

# Query bond details
./scripts/cli-bondtoken.sh get-bond BOND_001
```

#### Compliance CLI (`cli-compliance.sh`)
Interface for KYC and AML operations.

**Commands:**
```bash
# KYC Operations
./scripts/cli-compliance.sh create-kyc <address> <full_name> <dob> <nationality> <id_type> <id_number>
./scripts/cli-compliance.sh approve-kyc <address> <approved_by> <risk_level>
./scripts/cli-compliance.sh reject-kyc <address> <rejected_by> <reason>

# AML Operations
./scripts/cli-compliance.sh create-aml <address> <check_type> <risk_score> <details>
./scripts/cli-compliance.sh update-aml <address> <check_type> <status> <risk_score> <details>

# Query Operations
./scripts/cli-compliance.sh check-compliance <address>
./scripts/cli-compliance.sh get-kyc <address>
./scripts/cli-compliance.sh get-aml <address> <check_type>
./scripts/cli-compliance.sh get-all-kyc
./scripts/cli-compliance.sh get-all-aml <address>
```

**Examples:**
```bash
# Create KYC record
./scripts/cli-compliance.sh create-kyc alice "Alice Johnson" "1990-01-01" "US" "PASSPORT" "US123456"

# Approve KYC
./scripts/cli-compliance.sh approve-kyc alice admin1 LOW

# Create AML check
./scripts/cli-compliance.sh create-aml alice SANCTIONS 5 "No sanctions found"

# Check compliance
./scripts/cli-compliance.sh check-compliance alice
```

#### CorporateAction CLI (`cli-corporateaction.sh`)
Interface for coupon payments and bond redemptions.

**Commands:**
```bash
# Coupon Operations
./scripts/cli-corporateaction.sh create-coupon <bond_id> <payment_date> <amount>
./scripts/cli-corporateaction.sh process-coupon <coupon_id>

# Redemption Operations
./scripts/cli-corporateaction.sh create-redemption <bond_id> <redemption_date> <amount>
./scripts/cli-corporateaction.sh process-redemption <redemption_id>

# Query Operations
./scripts/cli-corporateaction.sh get-coupon <coupon_id>
./scripts/cli-corporateaction.sh get-redemption <redemption_id>
./scripts/cli-corporateaction.sh get-coupons-by-bond <bond_id>
./scripts/cli-corporateaction.sh get-redemptions-by-bond <bond_id>
./scripts/cli-corporateaction.sh get-pending-coupons
./scripts/cli-corporateaction.sh get-pending-redemptions

# Utility Operations
./scripts/cli-corporateaction.sh calculate-coupon <bond_id> <face_value> <coupon_rate>
```

**Examples:**
```bash
# Create coupon payment
./scripts/cli-corporateaction.sh create-coupon BOND_001 2024-06-01 50.00

# Process coupon payment
./scripts/cli-corporateaction.sh process-coupon COUPON_001

# Create redemption
./scripts/cli-corporateaction.sh create-redemption BOND_001 2029-01-01 1000.00

# Calculate coupon amount
./scripts/cli-corporateaction.sh calculate-coupon BOND_001 1000.00 5.0
```

## Prerequisites

### 1. Hyperledger Fabric Environment
- Fabric network must be running
- Peer CLI must be available
- Channel must be created and joined
- Chaincodes must be installed and instantiated

### 2. Environment Variables
The scripts automatically set the following environment variables:
```bash
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

### 3. Network Configuration
- **Channel Name**: `mychannel`
- **Peer Address**: `localhost:7051`
- **Orderer Address**: `localhost:7050`
- **MSP ID**: `Org1MSP`

## Installation and Setup

### 1. Make Scripts Executable
```bash
chmod +x scripts/*.sh
```

### 2. Verify Fabric Environment
```bash
# Check if peer CLI is available
peer version

# Check if you're in the right environment
docker exec -it cli bash
```

### 3. Test Basic Functionality
```bash
# Run the test suite
./scripts/testChaincode.sh

# Try the master CLI
./scripts/cli-master.sh
```

## Usage Workflows

### 1. Complete Bond Lifecycle
```bash
# 1. Create KYC for issuer
./scripts/cli-compliance.sh create-kyc issuer001 "ABC Corp" "1980-01-01" "US" "EIN" "12-3456789"
./scripts/cli-compliance.sh approve-kyc issuer001 admin1 LOW

# 2. Issue bond
./scripts/cli-bondtoken.sh create-bond BOND_001 "ABC Corp Bond" USD 1000 5.0 2024-01-01 2029-01-01 ACTIVE

# 3. Create KYC for investor
./scripts/cli-compliance.sh create-kyc investor001 "John Doe" "1985-01-01" "US" "SSN" "123-45-6789"
./scripts/cli-compliance.sh approve-kyc investor001 admin1 LOW

# 4. Transfer bonds to investor
./scripts/cli-bondtoken.sh transfer-bond BOND_001 issuer001 investor001

# 5. Create coupon payment
./scripts/cli-corporateaction.sh create-coupon BOND_001 2024-06-01 50.00

# 6. Process coupon payment
./scripts/cli-corporateaction.sh process-coupon COUPON_001
```

### 2. Compliance Monitoring
```bash
# 1. Create comprehensive KYC
./scripts/cli-compliance.sh create-kyc user001 "Jane Smith" "1990-01-01" "UK" "PASSPORT" "GB123456"

# 2. Perform AML checks
./scripts/cli-compliance.sh create-aml user001 SANCTIONS 3 "Minor risk - enhanced monitoring required"
./scripts/cli-compliance.sh create-aml user001 PEP 7 "Politically exposed person - high risk"

# 3. Check overall compliance
./scripts/cli-compliance.sh check-compliance user001

# 4. Review all records
./scripts/cli-compliance.sh get-kyc user001
./scripts/cli-compliance.sh get-all-aml user001
```

## Error Handling

### Common Issues and Solutions

#### 1. Peer CLI Not Found
**Error**: `Error: peer CLI not found`
**Solution**: Ensure you're in the Fabric CLI environment
```bash
docker exec -it cli bash
```

#### 2. Channel Not Accessible
**Error**: `Channel not accessible`
**Solution**: Check if channel exists and is joined
```bash
peer channel list
peer channel join -b mychannel.block
```

#### 3. Chaincode Not Instantiated
**Error**: `Cannot query instantiated chaincodes`
**Solution**: Install and instantiate chaincodes
```bash
peer chaincode install -n bondtoken -v 1.0 -p github.com/bondtoken
peer chaincode instantiate -C mychannel -n bondtoken -v 1.0 -c '{"Args":["Init"]}'
```

#### 4. MSP Path Issues
**Warning**: `Warning: MSP path not found, using default`
**Solution**: Verify MSP configuration or use default paths

## Testing

### 1. Run Unit Tests
```bash
./scripts/testChaincode.sh
```

### 2. Test Individual Functions
```bash
# Test BondToken
./scripts/cli-bondtoken.sh help

# Test Compliance
./scripts/cli-compliance.sh help

# Test CorporateAction
./scripts/cli-corporateaction.sh help
```

### 3. Integration Testing
```bash
# Use the master CLI for comprehensive testing
./scripts/cli-master.sh
```

## Security Considerations

### 1. Access Control
- Scripts require appropriate MSP credentials
- Ensure proper channel policies are configured
- Use TLS for all communications

### 2. Input Validation
- All user inputs are validated before processing
- Sanitize inputs to prevent injection attacks
- Use parameterized queries

### 3. Audit Trail
- All operations are logged on the blockchain
- Maintain compliance with regulatory requirements
- Enable event monitoring

## Troubleshooting

### 1. Debug Mode
Enable verbose output for debugging:
```bash
# Set debug environment variable
export FABRIC_LOGGING_SPEC=DEBUG

# Run with verbose peer commands
peer chaincode query -C mychannel -n bondtoken -c '{"Args":["GetBond","BOND_001"]}' --tls --cafile $ORDERER_CA
```

### 2. Network Diagnostics
```bash
# Check network status
./scripts/cli-master.sh
# Select option 5: Network Status

# Check peer status
peer node status

# Check channel info
peer channel getinfo -c mychannel
```

### 3. Log Analysis
```bash
# Check peer logs
docker logs peer0.org1.example.com

# Check orderer logs
docker logs orderer.example.com
```

## Support and Maintenance

### 1. Regular Updates
- Keep scripts synchronized with chaincode versions
- Update configuration for network changes
- Maintain compatibility with Fabric versions

### 2. Monitoring
- Monitor script execution logs
- Track chaincode performance
- Monitor network health

### 3. Documentation
- Update README for new features
- Document configuration changes
- Maintain usage examples

## Contributing

### 1. Code Standards
- Follow bash scripting best practices
- Use consistent error handling
- Maintain backward compatibility

### 2. Testing Requirements
- All new features must include tests
- Maintain existing test coverage
- Validate against multiple Fabric versions

### 3. Documentation Updates
- Update README for new features
- Include usage examples
- Document configuration changes

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Contact

For questions and support:
- Create an issue in the project repository
- Contact the development team
- Check the project documentation
