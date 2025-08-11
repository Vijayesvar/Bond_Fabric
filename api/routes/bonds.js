const express = require('express');
const router = express.Router();
const blockchainService = require('../services/blockchainService');
const { validateBondData, validateTransferData } = require('../middleware/validation');
const auth = require('../middleware/auth');

/**
 * @swagger
 * components:
 *   schemas:
 *     Bond:
 *       type: object
 *       required:
 *         - id
 *         - issuerID
 *         - issuerName
 *         - faceValue
 *         - couponRate
 *         - totalSupply
 *         - maturityDate
 *         - currency
 *         - isin
 *         - rating
 *       properties:
 *         id:
 *           type: string
 *           description: Unique bond identifier
 *         issuerID:
 *           type: string
 *           description: Issuer organization ID
 *         issuerName:
 *           type: string
 *           description: Name of the issuing organization
 *         faceValue:
 *           type: number
 *           description: Face value of the bond
 *         couponRate:
 *           type: number
 *           description: Annual coupon rate percentage
 *         totalSupply:
 *           type: integer
 *           description: Total number of tokens issued
 *         maturityDate:
 *           type: string
 *           format: date
 *           description: Bond maturity date
 *         currency:
 *           type: string
 *           description: Currency of the bond
 *         isin:
 *           type: string
 *           description: International Securities Identification Number
 *         rating:
 *           type: string
 *           description: Credit rating of the bond
 *         collateral:
 *           type: string
 *           description: Collateral backing the bond
 */

/**
 * @swagger
 * /api/bonds:
 *   get:
 *     summary: Get all bonds
 *     tags: [Bonds]
 *     responses:
 *       200:
 *         description: List of all bonds
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/Bond'
 */
router.get('/', async (req, res) => {
  try {
    const bonds = await blockchainService.getAllBonds();
    res.json(bonds);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/bonds/{id}:
 *   get:
 *     summary: Get bond by ID
 *     tags: [Bonds]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *         description: Bond ID
 *     responses:
 *       200:
 *         description: Bond details
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Bond'
 *       404:
 *         description: Bond not found
 */
router.get('/:id', async (req, res) => {
  try {
    const bond = await blockchainService.getBond(req.params.id);
    if (!bond) {
      return res.status(404).json({ error: 'Bond not found' });
    }
    res.json(bond);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/bonds:
 *   post:
 *     summary: Issue a new bond
 *     tags: [Bonds]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/Bond'
 *     responses:
 *       201:
 *         description: Bond issued successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 txId:
 *                   type: string
 *                 bond:
 *                   $ref: '#/components/schemas/Bond'
 *       400:
 *         description: Invalid bond data
 *       401:
 *         description: Unauthorized
 */
router.post('/', auth, validateBondData, async (req, res) => {
  try {
    const result = await blockchainService.issueBond(req.body);
    
    // Get the created bond
    const bond = await blockchainService.getBond(req.body.id);
    
    res.status(201).json({
      success: true,
      txId: result.txId,
      bond
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/bonds/{id}/transfer:
 *   post:
 *     summary: Transfer bond tokens
 *     tags: [Bonds]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *         description: Bond ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - from
 *               - to
 *               - quantity
 *             properties:
 *               from:
 *                 type: string
 *                 description: Sender address
 *               to:
 *                 type: string
 *                 description: Recipient address
 *               quantity:
 *                 type: integer
 *                 description: Number of tokens to transfer
 *     responses:
 *       200:
 *         description: Transfer successful
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 txId:
 *                   type: string
 *       400:
 *         description: Invalid transfer data
 *       401:
 *         description: Unauthorized
 */
router.post('/:id/transfer', auth, validateTransferData, async (req, res) => {
  try {
    const { from, to, quantity } = req.body;
    const bondId = req.params.id;
    
    const result = await blockchainService.transferTokens(from, to, bondId, quantity);
    
    res.json({
      success: true,
      txId: result.txId,
      message: `Successfully transferred ${quantity} tokens of bond ${bondId} from ${from} to ${to}`
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/bonds/{id}/balance/{address}:
 *   get:
 *     summary: Get token balance for a specific address
 *     tags: [Bonds]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *         description: Bond ID
 *       - in: path
 *         name: address
 *         required: true
 *         schema:
 *           type: string
 *         description: Address to check balance for
 *     responses:
 *       200:
 *         description: Balance information
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 bondId:
 *                   type: string
 *                 address:
 *                   type: string
 *                 balance:
 *                   type: integer
 *       404:
 *         description: Bond not found
 */
router.get('/:id/balance/:address', async (req, res) => {
  try {
    const { id, address } = req.params;
    
    // Check if bond exists
    const bond = await blockchainService.getBond(id);
    if (!bond) {
      return res.status(404).json({ error: 'Bond not found' });
    }
    
    const balance = await blockchainService.getBalance(address, id);
    
    res.json({
      bondId: id,
      address,
      balance,
      bond: bond
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/bonds/{id}/holders:
 *   get:
 *     summary: Get all token holders for a bond
 *     tags: [Bonds]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *         description: Bond ID
 *     responses:
 *       200:
 *         description: List of token holders
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 type: object
 *                 properties:
 *                   address:
 *                     type: string
 *                   bondId:
 *                     type: string
 *                   quantity:
 *                     type: integer
 *                   lastUpdated:
 *                     type: string
 *                     format: date-time
 *       404:
 *         description: Bond not found
 */
router.get('/:id/holders', async (req, res) => {
  try {
    const { id } = req.params;
    
    // Check if bond exists
    const bond = await blockchainService.getBond(id);
    if (!bond) {
      return res.status(404).json({ error: 'Bond not found' });
    }
    
    // Note: This would need to be implemented in the blockchain service
    // For now, we'll return a placeholder
    res.json({
      bondId: id,
      holders: [],
      message: 'Holders endpoint not yet implemented'
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

module.exports = router;
