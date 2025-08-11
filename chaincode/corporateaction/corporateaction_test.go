package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStub is a mock implementation of the chaincode stub
type MockStub struct {
	mock.Mock
	state map[string][]byte
}

func (m *MockStub) GetState(key string) ([]byte, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStub) PutState(key string, value []byte) error {
	args := m.Called(key, value)
	m.state[key] = value
	return args.Error(0)
}

func (m *MockStub) GetStateByRange(startKey, endKey string) (contractapi.StateQueryIteratorInterface, error) {
	args := m.Called(startKey, endKey)
	return args.Get(0).(contractapi.StateQueryIteratorInterface), args.Error(1)
}

func (m *MockStub) GetTxID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockStub) SetEvent(name string, payload []byte) error {
	args := m.Called(name, payload)
	return args.Error(0)
}

// MockContext is a mock implementation of the transaction context
type MockContext struct {
	mock.Mock
	stub *MockStub
}

func (m *MockContext) GetStub() contractapi.TransactionContextInterface {
	return m
}

func (m *MockContext) GetState(key string) ([]byte, error) {
	return m.stub.GetState(key)
}

func (m *MockContext) PutState(key string, value []byte) error {
	return m.stub.PutState(key, value)
}

func (m *MockContext) GetStateByRange(startKey, endKey string) (contractapi.StateQueryIteratorInterface, error) {
	return m.stub.GetStateByRange(startKey, endKey)
}

func (m *MockContext) GetTxID() string {
	return m.stub.GetTxID()
}

func (m *MockContext) SetEvent(name string, payload []byte) error {
	return m.stub.SetEvent(name, payload)
}

// MockIterator is a mock implementation of the state query iterator
type MockIterator struct {
	mock.Mock
	results [][]byte
	index   int
}

func (m *MockIterator) HasNext() bool {
	return m.index < len(m.results)
}

func (m *MockIterator) Next() (*contractapi.QueryResult, error) {
	if m.index >= len(m.results) {
		return nil, nil
	}
	
	result := &contractapi.QueryResult{
		Key:   fmt.Sprintf("key_%d", m.index),
		Value: m.results[m.index],
	}
	m.index++
	return result, nil
}

func (m *MockIterator) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCorporateAction_Init(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	err := ca.Init(ctx)
	assert.NoError(t, err)
}

func TestCorporateAction_CreateCouponPayment(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock the stub methods
	ctx.stub.On("PutState", mock.Anything, mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "CorporateActionEvent", mock.Anything).Return(nil)
	
	err := ca.CreateCouponPayment(ctx, "BOND_001", "2024-06-01", 50.0)
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCorporateAction_CreateCouponPayment_InvalidDate(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	err := ca.CreateCouponPayment(ctx, "BOND_001", "invalid-date", 50.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payment date format")
}

func TestCorporateAction_ProcessCouponPayment(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a coupon payment first
	couponPayment := CouponPayment{
		ID:          "COUPON_BOND_001_20240601",
		BondID:      "BOND_001",
		PaymentDate: time.Now(),
		Amount:      50.0,
		Status:      "PENDING",
	}
	
	couponJSON, _ := json.Marshal(couponPayment)
	ctx.stub.On("GetState", "COUPON_BOND_001_20240601").Return(couponJSON, nil)
	ctx.stub.On("PutState", "COUPON_BOND_001_20240601", mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "CorporateActionEvent", mock.Anything).Return(nil)
	
	err := ca.ProcessCouponPayment(ctx, "COUPON_BOND_001_20240601")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCorporateAction_ProcessCouponPayment_NotPending(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a coupon payment with non-pending status
	couponPayment := CouponPayment{
		ID:          "COUPON_BOND_001_20240601",
		BondID:      "BOND_001",
		PaymentDate: time.Now(),
		Amount:      50.0,
		Status:      "PAID",
	}
	
	couponJSON, _ := json.Marshal(couponPayment)
	ctx.stub.On("GetState", "COUPON_BOND_001_20240601").Return(couponJSON, nil)
	
	err := ca.ProcessCouponPayment(ctx, "COUPON_BOND_001_20240601")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not pending")
}

func TestCorporateAction_CreateRedemption(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock the stub methods
	ctx.stub.On("PutState", mock.Anything, mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "CorporateActionEvent", mock.Anything).Return(nil)
	
	err := ca.CreateRedemption(ctx, "BOND_001", "2029-01-01", 1000.0)
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCorporateAction_CreateRedemption_InvalidDate(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	err := ca.CreateRedemption(ctx, "BOND_001", "invalid-date", 1000.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid redemption date format")
}

func TestCorporateAction_ProcessRedemption(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a redemption first
	redemption := Redemption{
		ID:             "REDEMPTION_BOND_001_20290101",
		BondID:         "BOND_001",
		RedemptionDate: time.Now(),
		Amount:         1000.0,
		Status:         "PENDING",
	}
	
	redemptionJSON, _ := json.Marshal(redemption)
	ctx.stub.On("GetState", "REDEMPTION_BOND_001_20290101").Return(redemptionJSON, nil)
	ctx.stub.On("PutState", "REDEMPTION_BOND_001_20290101", mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "CorporateActionEvent", mock.Anything).Return(nil)
	
	err := ca.ProcessRedemption(ctx, "REDEMPTION_BOND_001_20290101")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCorporateAction_ProcessRedemption_NotPending(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a redemption with non-pending status
	redemption := Redemption{
		ID:             "REDEMPTION_BOND_001_20290101",
		BondID:         "BOND_001",
		RedemptionDate: time.Now(),
		Amount:         1000.0,
		Status:         "COMPLETED",
	}
	
	redemptionJSON, _ := json.Marshal(redemption)
	ctx.stub.On("GetState", "REDEMPTION_BOND_001_20290101").Return(redemptionJSON, nil)
	
	err := ca.ProcessRedemption(ctx, "REDEMPTION_BOND_001_20290101")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not pending")
}

func TestCorporateAction_GetCouponPayment(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a coupon payment
	couponPayment := CouponPayment{
		ID:          "COUPON_BOND_001_20240601",
		BondID:      "BOND_001",
		PaymentDate: time.Now(),
		Amount:      50.0,
		Status:      "PENDING",
	}
	
	couponJSON, _ := json.Marshal(couponPayment)
	ctx.stub.On("GetState", "COUPON_BOND_001_20240601").Return(couponJSON, nil)
	
	retrievedCoupon, err := ca.GetCouponPayment(ctx, "COUPON_BOND_001_20240601")
	assert.NoError(t, err)
	assert.Equal(t, couponPayment.ID, retrievedCoupon.ID)
	assert.Equal(t, couponPayment.BondID, retrievedCoupon.BondID)
	assert.Equal(t, couponPayment.Amount, retrievedCoupon.Amount)
}

func TestCorporateAction_GetCouponPayment_NotFound(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	ctx.stub.On("GetState", "COUPON_BOND_001_20240601").Return(nil, nil)
	
	_, err := ca.GetCouponPayment(ctx, "COUPON_BOND_001_20240601")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestCorporateAction_GetRedemption(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a redemption
	redemption := Redemption{
		ID:             "REDEMPTION_BOND_001_20290101",
		BondID:         "BOND_001",
		RedemptionDate: time.Now(),
		Amount:         1000.0,
		Status:         "PENDING",
	}
	
	redemptionJSON, _ := json.Marshal(redemption)
	ctx.stub.On("GetState", "REDEMPTION_BOND_001_20290101").Return(redemptionJSON, nil)
	
	retrievedRedemption, err := ca.GetRedemption(ctx, "REDEMPTION_BOND_001_20290101")
	assert.NoError(t, err)
	assert.Equal(t, redemption.ID, retrievedRedemption.ID)
	assert.Equal(t, redemption.BondID, retrievedRedemption.BondID)
	assert.Equal(t, redemption.Amount, retrievedRedemption.Amount)
}

func TestCorporateAction_GetRedemption_NotFound(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	ctx.stub.On("GetState", "REDEMPTION_BOND_001_20290101").Return(nil, nil)
	
	_, err := ca.GetRedemption(ctx, "REDEMPTION_BOND_001_20290101")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestCorporateAction_GetCouponPaymentsByBond(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with coupon payment results
	coupon1 := CouponPayment{ID: "COUPON_BOND_001_20240601", BondID: "BOND_001"}
	coupon2 := CouponPayment{ID: "COUPON_BOND_001_20241201", BondID: "BOND_001"}
	
	coupon1JSON, _ := json.Marshal(coupon1)
	coupon2JSON, _ := json.Marshal(coupon2)
	
	mockIterator := &MockIterator{results: [][]byte{coupon1JSON, coupon2JSON}}
	
	ctx.stub.On("GetStateByRange", "", "").Return(mockIterator, nil)
	
	coupons, err := ca.GetCouponPaymentsByBond(ctx, "BOND_001")
	assert.NoError(t, err)
	assert.Len(t, coupons, 2)
	assert.Equal(t, "BOND_001", coupons[0].BondID)
	assert.Equal(t, "BOND_001", coupons[1].BondID)
}

func TestCorporateAction_GetRedemptionsByBond(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with redemption results
	redemption1 := Redemption{ID: "REDEMPTION_BOND_001_20290101", BondID: "BOND_001"}
	redemption2 := Redemption{ID: "REDEMPTION_BOND_001_20290701", BondID: "BOND_001"}
	
	redemption1JSON, _ := json.Marshal(redemption1)
	redemption2JSON, _ := json.Marshal(redemption2)
	
	mockIterator := &MockIterator{results: [][]byte{redemption1JSON, redemption2JSON}}
	
	ctx.stub.On("GetStateByRange", "", "").Return(mockIterator, nil)
	
	redemptions, err := ca.GetRedemptionsByBond(ctx, "BOND_001")
	assert.NoError(t, err)
	assert.Len(t, redemptions, 2)
	assert.Equal(t, "BOND_001", redemptions[0].BondID)
	assert.Equal(t, "BOND_001", redemptions[1].BondID)
}

func TestCorporateAction_GetPendingCouponPayments(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with pending coupon payment results
	coupon1 := CouponPayment{ID: "COUPON_BOND_001_20240601", BondID: "BOND_001", Status: "PENDING"}
	coupon2 := CouponPayment{ID: "COUPON_BOND_002_20240601", BondID: "BOND_002", Status: "PENDING"}
	
	coupon1JSON, _ := json.Marshal(coupon1)
	coupon2JSON, _ := json.Marshal(coupon2)
	
	mockIterator := &MockIterator{results: [][]byte{coupon1JSON, coupon2JSON}}
	
	ctx.stub.On("GetStateByRange", "", "").Return(mockIterator, nil)
	
	pendingPayments, err := ca.GetPendingCouponPayments(ctx)
	assert.NoError(t, err)
	assert.Len(t, pendingPayments, 2)
	assert.Equal(t, "PENDING", pendingPayments[0].Status)
	assert.Equal(t, "PENDING", pendingPayments[1].Status)
}

func TestCorporateAction_GetPendingRedemptions(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with pending redemption results
	redemption1 := Redemption{ID: "REDEMPTION_BOND_001_20290101", BondID: "BOND_001", Status: "PENDING"}
	redemption2 := Redemption{ID: "REDEMPTION_BOND_002_20290101", BondID: "BOND_002", Status: "PENDING"}
	
	redemption1JSON, _ := json.Marshal(redemption1)
	redemption2JSON, _ := json.Marshal(redemption2)
	
	mockIterator := &MockIterator{results: [][]byte{redemption1JSON, redemption2JSON}}
	
	ctx.stub.On("GetStateByRange", "", "").Return(mockIterator, nil)
	
	pendingRedemptions, err := ca.GetPendingRedemptions(ctx)
	assert.NoError(t, err)
	assert.Len(t, pendingRedemptions, 2)
	assert.Equal(t, "PENDING", pendingRedemptions[0].Status)
	assert.Equal(t, "PENDING", pendingRedemptions[1].Status)
}

func TestCorporateAction_CalculateCouponAmount(t *testing.T) {
	ca := &CorporateAction{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	amount, err := ca.CalculateCouponAmount(ctx, "BOND_001", 1000.0, 5.0)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, amount)
	
	// Test with different values
	amount, err = ca.CalculateCouponAmount(ctx, "BOND_002", 5000.0, 3.5)
	assert.NoError(t, err)
	assert.Equal(t, 175.0, amount)
}

