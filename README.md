# BondBridge - Corporate Bond Tokenization Platform

A Hyperledger Fabric-based platform for corporate bond tokenization and trading, designed for the SEBI hackathon.

## Architecture Overview

BondBridge implements a permissioned blockchain network with the following components:

- **Network Topology**: 2 peer nodes + 3 orderers (Raft consensus)
- **Channels**: `bondchannel` for bond operations
- **Organizations**: Issuer, Regulator, Market-Maker, Custodian, Investor
- **Smart Contracts**: BondToken, Compliance, CorporateAction
- **APIs**: REST/gRPC services with Fabric SDK integration
- **Frontend**: React-based web interface

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Node.js 18+
- Go 1.19+
- Hyperledger Fabric binaries

### Setup
```bash
# Clone and setup
git clone <repository>
cd BondBridge

# Start the network
./scripts/startNetwork.sh

# Install and instantiate chaincode
./scripts/deployChaincode.sh

# Start the API server
cd api && npm install && npm start

# Start the frontend
cd frontend && npm install && npm start
```

## Project Structure

```
BondBridge/
├── network/           # Fabric network configuration
├── chaincode/         # Smart contracts
├── api/              # REST/gRPC API layer
├── frontend/         # React web interface
├── scripts/          # Deployment and utility scripts
└── docs/             # Documentation and runbooks
```



MIT License
