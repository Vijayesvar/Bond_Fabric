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

func (m *MockStub) DelState(key string) error {
	args := m.Called(key)
	delete(m.state, key)
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

func (m *MockContext) DelState(key string) error {
	return m.stub.DelState(key)
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
		return nil, fmt.Errorf("no more results")
	}
	
	result := &contractapi.QueryResult{
		Value: m.results[m.index],
	}
	m.index++
	return result, nil
}

func (m *MockIterator) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestBondToken_Init(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	err := bt.Init(ctx)
	assert.NoError(t, err)
}

func TestBondToken_CreateBond(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock the stub methods
	ctx.stub.On("GetState", "BOND_001").Return(nil, nil)
	ctx.stub.On("PutState", mock.Anything, mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "BondEvent", mock.Anything).Return(nil)
	
	err := bt.CreateBond(ctx, "BOND_001", "Test Bond", "USD", 1000.0, 5.0, "2024-01-01", "2029-01-01", "ACTIVE")
	assert.NoError(t, err)
	
	// Verify the bond was created
	ctx.stub.AssertExpectations(t)
}

func TestBondToken_CreateBond_AlreadyExists(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock existing bond
	existingBond := Bond{
		ID:          "BOND_001",
		Name:        "Existing Bond",
		Currency:    "USD",
		FaceValue:   1000.0,
		CouponRate:  5.0,
		IssueDate:   time.Now(),
		MaturityDate: time.Now().AddDate(5, 0, 0),
		Status:      "ACTIVE",
	}
	
	existingBondJSON, _ := json.Marshal(existingBond)
	ctx.stub.On("GetState", "BOND_001").Return(existingBondJSON, nil)
	
	err := bt.CreateBond(ctx, "BOND_001", "Test Bond", "USD", 1000.0, 5.0, "2024-01-01", "2029-01-01", "ACTIVE")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestBondToken_TransferBond(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a bond first
	bond := Bond{
		ID:          "BOND_001",
		Name:        "Test Bond",
		Currency:    "USD",
		FaceValue:   1000.0,
		CouponRate:  5.0,
		IssueDate:   time.Now(),
		MaturityDate: time.Now().AddDate(5, 0, 0),
		Status:      "ACTIVE",
		Owner:       "alice",
	}
	
	bondJSON, _ := json.Marshal(bond)
	ctx.stub.On("GetState", "BOND_001").Return(bondJSON, nil)
	ctx.stub.On("PutState", mock.Anything, mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "BondEvent", mock.Anything).Return(nil)
	
	err := bt.TransferBond(ctx, "BOND_001", "alice", "bob")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestBondToken_GetBond(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a bond
	bond := Bond{
		ID:          "BOND_001",
		Name:        "Test Bond",
		Currency:    "USD",
		FaceValue:   1000.0,
		CouponRate:  5.0,
		IssueDate:   time.Now(),
		MaturityDate: time.Now().AddDate(5, 0, 0),
		Status:      "ACTIVE",
		Owner:       "alice",
	}
	
	bondJSON, _ := json.Marshal(bond)
	ctx.stub.On("GetState", "BOND_001").Return(bondJSON, nil)
	
	retrievedBond, err := bt.GetBond(ctx, "BOND_001")
	assert.NoError(t, err)
	assert.Equal(t, bond.ID, retrievedBond.ID)
	assert.Equal(t, bond.Name, retrievedBond.Name)
	assert.Equal(t, bond.Owner, retrievedBond.Owner)
}

func TestBondToken_GetBond_NotFound(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	ctx.stub.On("GetState", "BOND_001").Return(nil, nil)
	
	_, err := bt.GetBond(ctx, "BOND_001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestBondToken_GetAllBonds(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with bond results
	bond1 := Bond{ID: "BOND_001", Name: "Bond 1"}
	bond2 := Bond{ID: "BOND_002", Name: "Bond 2"}
	
	bond1JSON, _ := json.Marshal(bond1)
	bond2JSON, _ := json.Marshal(bond2)
	
	mockIterator := &MockIterator{results: [][]byte{bond1JSON, bond2JSON}}
	
	ctx.stub.On("GetStateByRange", "", "").Return(mockIterator, nil)
	
	bonds, err := bt.GetAllBonds(ctx)
	assert.NoError(t, err)
	assert.Len(t, bonds, 2)
	assert.Equal(t, "BOND_001", bonds[0].ID)
	assert.Equal(t, "BOND_002", bonds[1].ID)
}

func TestBondToken_UpdateBondStatus(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a bond
	bond := Bond{
		ID:          "BOND_001",
		Name:        "Test Bond",
		Currency:    "USD",
		FaceValue:   1000.0,
		CouponRate:  5.0,
		IssueDate:   time.Now(),
		MaturityDate: time.Now().AddDate(5, 0, 0),
		Status:      "ACTIVE",
		Owner:       "alice",
	}
	
	bondJSON, _ := json.Marshal(bond)
	ctx.stub.On("GetState", "BOND_001").Return(bondJSON, nil)
	ctx.stub.On("PutState", mock.Anything, mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "BondEvent", mock.Anything).Return(nil)
	
	err := bt.UpdateBondStatus(ctx, "BOND_001", "MATURED")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestBondToken_CalculateYield(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a bond
	bond := Bond{
		ID:          "BOND_001",
		Name:        "Test Bond",
		Currency:    "USD",
		FaceValue:   1000.0,
		CouponRate:  5.0,
		IssueDate:   time.Now(),
		MaturityDate: time.Now().AddDate(5, 0, 0),
		Status:      "ACTIVE",
		Owner:       "alice",
	}
	
	bondJSON, _ := json.Marshal(bond)
	ctx.stub.On("GetState", "BOND_001").Return(bondJSON, nil)
	
	yield, err := bt.CalculateYield(ctx, "BOND_001", 950.0)
	assert.NoError(t, err)
	assert.Greater(t, yield, 0.0)
}

func TestBondToken_GetBondsByOwner(t *testing.T) {
	bt := &BondToken{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with bond results
	bond1 := Bond{ID: "BOND_001", Name: "Bond 1", Owner: "alice"}
	bond2 := Bond{ID: "BOND_002", Name: "Bond 2", Owner: "alice"}
	
	bond1JSON, _ := json.Marshal(bond1)
	bond2JSON, _ := json.Marshal(bond2)
	
	mockIterator := &MockIterator{results: [][]byte{bond1JSON, bond2JSON}}
	
	ctx.stub.On("GetStateByRange", "", "").Return(mockIterator, nil)
	
	bonds, err := bt.GetBondsByOwner(ctx, "alice")
	assert.NoError(t, err)
	assert.Len(t, bonds, 2)
	assert.Equal(t, "alice", bonds[0].Owner)
	assert.Equal(t, "alice", bonds[1].Owner)
}

