const Joi = require('joi');

const validateBondData = (req, res, next) => {
  const schema = Joi.object({
    id: Joi.string().required(),
    issuerID: Joi.string().required(),
    issuerName: Joi.string().required(),
    faceValue: Joi.number().positive().required(),
    couponRate: Joi.number().min(0).max(100).required(),
    totalSupply: Joi.number().integer().positive().required(),
    maturityDate: Joi.string().pattern(/^\d{4}-\d{2}-\d{2}$/).required(),
    currency: Joi.string().required(),
    isin: Joi.string().required(),
    rating: Joi.string().required(),
    collateral: Joi.string().required()
  });

  const { error } = schema.validate(req.body);
  if (error) {
    return res.status(400).json({ error: error.details[0].message });
  }

  next();
};

const validateTransferData = (req, res, next) => {
  const schema = Joi.object({
    from: Joi.string().required(),
    to: Joi.string().required(),
    quantity: Joi.number().integer().positive().required()
  });

  const { error } = schema.validate(req.body);
  if (error) {
    return res.status(400).json({ error: error.details[0].message });
  }

  next();
};

const validateLoginData = (req, res, next) => {
  const schema = Joi.object({
    email: Joi.string().email().required(),
    password: Joi.string().min(6).required()
  });

  const { error } = schema.validate(req.body);
  if (error) {
    return res.status(400).json({ error: error.details[0].message });
  }

  next();
};

const validateRegisterData = (req, res, next) => {
  const schema = Joi.object({
    email: Joi.string().email().required(),
    password: Joi.string().min(6).required(),
    role: Joi.string().valid('ISSUER', 'INVESTOR', 'REGULATOR', 'MARKET_MAKER', 'CUSTODIAN').required(),
    organization: Joi.string().required(),
    blockchainAddress: Joi.string().optional()
  });

  const { error } = schema.validate(req.body);
  if (error) {
    return res.status(400).json({ error: error.details[0].message });
  }

  next();
};

module.exports = {
  validateBondData,
  validateTransferData,
  validateLoginData,
  validateRegisterData
};
