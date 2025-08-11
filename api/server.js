const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');
const swaggerJsdoc = require('swagger-jsdoc');
const swaggerUi = require('swagger-ui-express');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 3001;

// Import routes
const bondRoutes = require('./routes/bonds');
const complianceRoutes = require('./routes/compliance');
const corporateActionRoutes = require('./routes/corporateActions');
const authRoutes = require('./routes/auth');
const userRoutes = require('./routes/users');

// Import blockchain service
const blockchainService = require('./services/blockchainService');

// Middleware
app.use(helmet());
app.use(cors());
app.use(morgan('combined'));
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Swagger configuration
const swaggerOptions = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'BondBridge API',
      version: '1.0.0',
      description: 'API for BondBridge Corporate Bond Tokenization Platform',
    },
    servers: [
      {
        url: `http://localhost:${PORT}`,
        description: 'Development server',
      },
    ],
  },
  apis: ['./routes/*.js'],
};

const swaggerSpec = swaggerJsdoc(swaggerOptions);
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec));

// Health check endpoint
app.get('/health', (req, res) => {
  res.status(200).json({
    status: 'OK',
    timestamp: new Date().toISOString(),
    service: 'BondBridge API',
    version: '1.0.0'
  });
});

// Blockchain status endpoint
app.get('/blockchain/status', async (req, res) => {
  try {
    const status = await blockchainService.getNetworkStatus();
    res.json(status);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// API Routes
app.use('/api/bonds', bondRoutes);
app.use('/api/compliance', complianceRoutes);
app.use('/api/corporate-actions', corporateActionRoutes);
app.use('/api/auth', authRoutes);
app.use('/api/users', userRoutes);

// Error handling middleware
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    error: 'Something went wrong!',
    message: process.env.NODE_ENV === 'development' ? err.message : 'Internal server error'
  });
});

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({ error: 'Route not found' });
});

// Start server
app.listen(PORT, async () => {
  console.log(`🚀 BondBridge API Server running on port ${PORT}`);
  console.log(`📚 API Documentation available at http://localhost:${PORT}/api-docs`);
  
  try {
    // Initialize blockchain connection
    await blockchainService.initialize();
    console.log('✅ Blockchain connection established');
  } catch (error) {
    console.error('❌ Failed to connect to blockchain:', error.message);
  }
});

module.exports = app;
