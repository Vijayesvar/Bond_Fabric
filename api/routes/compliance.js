const express = require('express');
const router = express.Router();
const blockchainService = require('../services/blockchainService');
const auth = require('../middleware/auth');

/**
 * @swagger
 * components:
 *   schemas:
 *     KYC:
 *       type: object
 *       required:
 *         - address
 *         - fullName
 *         - dateOfBirth
 *         - nationality
 *         - idType
 *         - idNumber
 *       properties:
 *         address:
 *           type: string
 *           description: User's blockchain address
 *         fullName:
 *           type: string
 *           description: User's full name
 *         dateOfBirth:
 *           type: string
 *           description: User's date of birth
 *         nationality:
 *           type: string
 *           description: User's nationality
 *         idType:
 *           type: string
 *           description: Type of identification document
 *         idNumber:
 *           type: string
 *           description: Identification number
 *         status:
 *           type: string
 *           enum: [PENDING, APPROVED, REJECTED]
 *           description: KYC approval status
 *         riskLevel:
 *           type: string
 *           enum: [LOW, MEDIUM, HIGH]
 *           description: Risk assessment level
 */

/**
 * @swagger
 * /api/compliance/kyc:
 *   post:
 *     summary: Create a new KYC record
 *     tags: [Compliance]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/KYC'
 *     responses:
 *       201:
 *         description: KYC record created successfully
 *       400:
 *         description: Invalid KYC data
 *       401:
 *         description: Unauthorized
 */
router.post('/kyc', auth, async (req, res) => {
  try {
    const result = await blockchainService.createKYC(req.body);
    res.status(201).json({
      success: true,
      txId: result.txId,
      message: 'KYC record created successfully'
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/compliance/kyc/{address}:
 *   get:
 *     summary: Get KYC record by address
 *     tags: [Compliance]
 *     parameters:
 *       - in: path
 *         name: address
 *         required: true
 *         schema:
 *           type: string
 *         description: User's blockchain address
 *     responses:
 *       200:
 *         description: KYC record details
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/KYC'
 *       404:
 *         description: KYC record not found
 */
router.get('/kyc/:address', async (req, res) => {
  try {
    const kyc = await blockchainService.getKYC(req.params.address);
    if (!kyc) {
      return res.status(404).json({ error: 'KYC record not found' });
    }
    res.json(kyc);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/compliance/kyc/{address}/approve:
 *   post:
 *     summary: Approve KYC record
 *     tags: [Compliance]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: address
 *         required: true
 *         schema:
 *           type: string
 *         description: User's blockchain address
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - approvedBy
 *               - riskLevel
 *             properties:
 *               approvedBy:
 *                 type: string
 *                 description: ID of the approver
 *               riskLevel:
 *                 type: string
 *                 enum: [LOW, MEDIUM, HIGH]
 *                 description: Risk assessment level
 *     responses:
 *       200:
 *         description: KYC approved successfully
 *       400:
 *         description: Invalid approval data
 *       401:
 *         description: Unauthorized
 */
router.post('/kyc/:address/approve', auth, async (req, res) => {
  try {
    const { approvedBy, riskLevel } = req.body;
    const { address } = req.params;
    
    const result = await blockchainService.approveKYC(address, approvedBy, riskLevel);
    
    res.json({
      success: true,
      txId: result.txId,
      message: 'KYC approved successfully'
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/compliance/check/{address}:
 *   get:
 *     summary: Check compliance status for an address
 *     tags: [Compliance]
 *     parameters:
 *       - in: path
 *         name: address
 *         required: true
 *         schema:
 *           type: string
 *         description: Address to check compliance for
 *     responses:
 *       200:
 *         description: Compliance status
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 isCompliant:
 *                   type: boolean
 *                   description: Whether the address is compliant
 *                 details:
 *                   type: string
 *                   description: Compliance details
 *       404:
 *         description: Address not found
 */
router.get('/check/:address', async (req, res) => {
  try {
    const compliance = await blockchainService.checkCompliance(req.params.address);
    res.json(compliance);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

module.exports = router;
