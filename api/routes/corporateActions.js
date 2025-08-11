const express = require('express');
const router = express.Router();
const blockchainService = require('../services/blockchainService');
const auth = require('../middleware/auth');

/**
 * @swagger
 * components:
 *   schemas:
 *     CorporateAction:
 *       type: object
 *       required:
 *         - id
 *         - bondId
 *         - type
 *         - description
 *         - recordDate
 *         - paymentDate
 *         - amount
 *       properties:
 *         id:
 *           type: string
 *           description: Unique corporate action identifier
 *         bondId:
 *           type: string
 *           description: Associated bond ID
 *         type:
 *           type: string
 *           enum: [COUPON_PAYMENT, PRINCIPAL_PAYMENT, INTEREST_PAYMENT, DEFAULT, CALL, PUT]
 *           description: Type of corporate action
 *         description:
 *           type: string
 *           description: Description of the corporate action
 *         recordDate:
 *           type: string
 *           format: date
 *           description: Record date for the action
 *         paymentDate:
 *           type: string
 *           format: date
 *           description: Payment date for the action
 *         amount:
 *           type: number
 *           description: Amount associated with the action
 *         status:
 *           type: string
 *           enum: [PENDING, PROCESSED, FAILED]
 *           description: Processing status
 */

/**
 * @swagger
 * /api/corporate-actions:
 *   post:
 *     summary: Create a new corporate action
 *     tags: [Corporate Actions]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/CorporateAction'
 *     responses:
 *       201:
 *         description: Corporate action created successfully
 *       400:
 *         description: Invalid corporate action data
 *       401:
 *         description: Unauthorized
 */
router.post('/', auth, async (req, res) => {
  try {
    const result = await blockchainService.createCorporateAction(req.body);
    res.status(201).json({
      success: true,
      txId: result.txId,
      message: 'Corporate action created successfully'
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/corporate-actions/bond/{bondId}:
 *   get:
 *     summary: Get all corporate actions for a specific bond
 *     tags: [Corporate Actions]
 *     parameters:
 *       - in: path
 *         name: bondId
 *         required: true
 *         schema:
 *           type: string
 *         description: Bond ID to get corporate actions for
 *     responses:
 *       200:
 *         description: List of corporate actions
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/CorporateAction'
 *       404:
 *         description: Bond not found
 */
router.get('/bond/:bondId', async (req, res) => {
  try {
    const actions = await blockchainService.getCorporateActions(req.params.bondId);
    res.json(actions);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

module.exports = router;
