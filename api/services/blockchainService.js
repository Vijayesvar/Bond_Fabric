const { Wallets, Gateway } = require('fabric-network');
const { FabricCAServices } = require('fabric-ca-client');
const path = require('path');
const fs = require('fs');

class BlockchainService {
  constructor() {
    this.gateway = null;
    this.network = null;
    this.contracts = {};
    this.wallet = null;
    this.connectionProfile = null;
    this.organizations = ['issuer', 'investor', 'regulator', 'marketmaker', 'custodian'];
  }

  async initialize() {
    try {
      // Load connection profile
      this.connectionProfile = JSON.parse(
        fs.readFileSync(path.join(__dirname, '../config/connection-profile.json'), 'utf8')
      );

      // Create wallet
      const walletPath = path.join(__dirname, '../wallet');
      this.wallet = await Wallets.newFileSystemWallet(walletPath);

      // Initialize contracts
      await this.initializeContracts();
      
      console.log('Blockchain service initialized successfully');
    } catch (error) {
      console.error('Failed to initialize blockchain service:', error);
      throw error;
    }
  }

  async initializeContracts() {
    try {
      // Initialize gateway
      this.gateway = new Gateway();
      
      // Connect to network
      await this.gateway.connect(this.connectionProfile, {
        wallet: this.wallet,
        identity: 'admin',
        discovery: { enabled: true, asLocalhost: true }
      });

      // Get network and channel
      this.network = await this.gateway.getNetwork('bondchannel');
      
      // Initialize contracts
      this.contracts.bondToken = await this.network.getContract('bondtoken');
      this.contracts.compliance = await this.network.getContract('compliance');
      this.contracts.corporateAction = await this.network.getContract('corporateaction');
      
      console.log('Contracts initialized successfully');
    } catch (error) {
      console.error('Failed to initialize contracts:', error);
      throw error;
    }
  }

  async getNetworkStatus() {
    try {
      if (!this.network) {
        return { status: 'disconnected', message: 'Network not initialized' };
      }

      const channel = this.network.getChannel();
      const info = await channel.queryInfo();
      
      return {
        status: 'connected',
        channel: 'bondchannel',
        height: info.height,
        currentBlockHash: info.currentBlockHash.toString('hex'),
        previousBlockHash: info.previousBlockHash.toString('hex'),
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      return { status: 'error', message: error.message };
    }
  }

  // Bond Token Contract Methods
  async issueBond(bondData) {
    try {
      const result = await this.contracts.bondToken.submitTransaction(
        'IssueBond',
        bondData.id,
        bondData.issuerID,
        bondData.issuerName,
        bondData.currency,
        bondData.isin,
        bondData.rating,
        bondData.collateral,
        bondData.faceValue.toString(),
        bondData.couponRate.toString(),
        bondData.totalSupply.toString(),
        bondData.maturityDate
      );
      
      return { success: true, txId: result.toString() };
    } catch (error) {
      throw new Error(`Failed to issue bond: ${error.message}`);
    }
  }

  async getBond(bondId) {
    try {
      const result = await this.contracts.bondToken.evaluateTransaction('GetBond', bondId);
      return JSON.parse(result.toString());
    } catch (error) {
      throw new Error(`Failed to get bond: ${error.message}`);
    }
  }

  async getAllBonds() {
    try {
      const result = await this.contracts.bondToken.evaluateTransaction('GetAllBonds');
      return JSON.parse(result.toString());
    } catch (error) {
      throw new Error(`Failed to get all bonds: ${error.message}`);
    }
  }

  async transferTokens(from, to, bondId, quantity) {
    try {
      const result = await this.contracts.bondToken.submitTransaction(
        'Transfer',
        from,
        to,
        bondId,
        quantity.toString()
      );
      
      return { success: true, txId: result.toString() };
    } catch (error) {
      throw new Error(`Failed to transfer tokens: ${error.message}`);
    }
  }

  async getBalance(address, bondId) {
    try {
      const result = await this.contracts.bondToken.evaluateTransaction('GetBalance', address, bondId);
      return parseInt(result.toString());
    } catch (error) {
      throw new Error(`Failed to get balance: ${error.message}`);
    }
  }

  // Compliance Contract Methods
  async createKYC(kycData) {
    try {
      const result = await this.contracts.compliance.submitTransaction(
        'CreateKYC',
        kycData.address,
        kycData.fullName,
        kycData.dateOfBirth,
        kycData.nationality,
        kycData.idType,
        kycData.idNumber
      );
      
      return { success: true, txId: result.toString() };
    } catch (error) {
      throw new Error(`Failed to create KYC: ${error.message}`);
    }
  }

  async approveKYC(address, approvedBy, riskLevel) {
    try {
      const result = await this.contracts.compliance.submitTransaction(
        'ApproveKYC',
        address,
        approvedBy,
        riskLevel
      );
      
      return { success: true, txId: result.toString() };
    } catch (error) {
      throw new Error(`Failed to approve KYC: ${error.message}`);
    }
  }

  async getKYC(address) {
    try {
      const result = await this.contracts.compliance.evaluateTransaction('GetKYC', address);
      return JSON.parse(result.toString());
    } catch (error) {
      throw new Error(`Failed to get KYC: ${error.message}`);
    }
  }

  async checkCompliance(address) {
    try {
      const result = await this.contracts.compliance.evaluateTransaction('CheckCompliance', address);
      const [isCompliant, details] = JSON.parse(result.toString());
      return { isCompliant, details };
    } catch (error) {
      throw new Error(`Failed to check compliance: ${error.message}`);
    }
  }

  // Corporate Action Contract Methods
  async createCorporateAction(actionData) {
    try {
      const result = await this.contracts.corporateAction.submitTransaction(
        'CreateCorporateAction',
        actionData.id,
        actionData.bondId,
        actionData.type,
        actionData.description,
        actionData.recordDate,
        actionData.paymentDate,
        actionData.amount.toString()
      );
      
      return { success: true, txId: result.toString() };
    } catch (error) {
      throw new Error(`Failed to create corporate action: ${error.message}`);
    }
  }

  async getCorporateActions(bondId) {
    try {
      const result = await this.contracts.corporateAction.evaluateTransaction('GetCorporateActions', bondId);
      return JSON.parse(result.toString());
    } catch (error) {
      throw new Error(`Failed to get corporate actions: ${error.message}`);
    }
  }

  // Utility Methods
  async disconnect() {
    if (this.gateway) {
      this.gateway.disconnect();
    }
  }

  async createIdentity(org, userId) {
    try {
      const caClient = new FabricCAServices(
        this.connectionProfile.certificateAuthorities[`ca.${org}.bondbridge.com`].url
      );

      const enrollment = await caClient.enroll({
        enrollmentID: userId,
        enrollmentSecret: 'password'
      });

      const identity = {
        credentials: {
          certificate: enrollment.certificate,
          privateKey: enrollment.key.toBytes()
        },
        mspId: `${org}MSP`,
        type: 'X.509'
      };

      await this.wallet.put(userId, identity);
      return { success: true, userId };
    } catch (error) {
      throw new Error(`Failed to create identity: ${error.message}`);
    }
  }
}

module.exports = new BlockchainService();
