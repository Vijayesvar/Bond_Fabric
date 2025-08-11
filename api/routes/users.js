const express = require('express');
const router = express.Router();
const auth = require('../middleware/auth');

// Mock user database (in production, use a real database)
const users = [
  {
    id: 1,
    email: 'admin@bondbridge.com',
    role: 'REGULATOR',
    organization: 'BondBridge Regulatory Authority',
    blockchainAddress: '0x1234567890abcdef',
    createdAt: new Date(),
    updatedAt: new Date()
  }
];

/**
 * @swagger
 * /api/users:
 *   get:
 *     summary: Get all users
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: List of all users
 *       401:
 *         description: Unauthorized
 */
router.get('/', auth, (req, res) => {
  try {
    // Remove sensitive information
    const safeUsers = users.map(user => {
      const { password, ...safeUser } = user;
      return safeUser;
    });
    
    res.json(safeUsers);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/users/{id}:
 *   get:
 *     summary: Get user by ID
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *         description: User ID
 *     responses:
 *       200:
 *         description: User details
 *       404:
 *         description: User not found
 *       401:
 *         description: Unauthorized
 */
router.get('/:id', auth, (req, res) => {
  try {
    const user = users.find(u => u.id === parseInt(req.params.id));
    
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    
    // Remove sensitive information
    const { password, ...safeUser } = user;
    res.json(safeUser);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

/**
 * @swagger
 * /api/users/profile:
 *   get:
 *     summary: Get current user profile
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Current user profile
 *       401:
 *         description: Unauthorized
 */
router.get('/profile', auth, (req, res) => {
  try {
    const user = users.find(u => u.id === req.user.userId);
    
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    
    // Remove sensitive information
    const { password, ...safeUser } = user;
    res.json(safeUser);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

module.exports = router;
