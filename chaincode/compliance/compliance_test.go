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

func TestCompliance_Init(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	err := c.Init(ctx)
	assert.NoError(t, err)
}

func TestCompliance_CreateKYC(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock the stub methods
	ctx.stub.On("GetState", "alice").Return(nil, nil)
	ctx.stub.On("PutState", "alice", mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "KYCEvent", mock.Anything).Return(nil)
	
	err := c.CreateKYC(ctx, "alice", "Alice Johnson", "1990-01-01", "US", "PASSPORT", "US123456")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCompliance_CreateKYC_AlreadyExists(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockContext{stub: &MockStub{state: make(map[string][]byte)}}}
	
	// Mock existing KYC
	existingKYC := KYCRecord{
		Address:     "alice",
		FullName:    "Alice Johnson",
		DateOfBirth: time.Now(),
		Nationality: "US",
		IDType:      "PASSPORT",
		IDNumber:    "US123456",
		Status:      "APPROVED",
	}
	
	existingKYCJSON, _ := json.Marshal(existingKYC)
	ctx.stub.On("GetState", "alice").Return(existingKYCJSON, nil)
	
	err := c.CreateKYC(ctx, "alice", "Alice Johnson", "1990-01-01", "US", "PASSPORT", "US123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestCompliance_ApproveKYC(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a KYC record first
	kyc := KYCRecord{
		Address:     "alice",
		FullName:    "Alice Johnson",
		DateOfBirth: time.Now(),
		Nationality: "US",
		IDType:      "PASSPORT",
		IDNumber:    "US123456",
		Status:      "PENDING",
	}
	
	kycJSON, _ := json.Marshal(kyc)
	ctx.stub.On("GetState", "alice").Return(kycJSON, nil)
	ctx.stub.On("PutState", "alice", mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "KYCEvent", mock.Anything).Return(nil)
	
	err := c.ApproveKYC(ctx, "alice", "admin", "LOW")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCompliance_RejectKYC(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a KYC record first
	kyc := KYCRecord{
		Address:     "alice",
		FullName:    "Alice Johnson",
		DateOfBirth: time.Now(),
		Nationality: "US",
		IDType:      "PASSPORT",
		IDNumber:    "US123456",
		Status:      "PENDING",
	}
	
	kycJSON, _ := json.Marshal(kyc)
	ctx.stub.On("GetState", "alice").Return(kycJSON, nil)
	ctx.stub.On("PutState", "alice", mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "KYCEvent", mock.Anything).Return(nil)
	
	err := c.RejectKYC(ctx, "alice", "admin", "Incomplete documentation")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCompliance_CreateAMLCheck(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock the stub methods
	ctx.stub.On("PutState", "alice_SANCTIONS", mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "AMLEvent", mock.Anything).Return(nil)
	
	err := c.CreateAMLCheck(ctx, "alice", "SANCTIONS", 75, "Sanctions check completed")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCompliance_UpdateAMLCheck(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create an AML check first
	amlCheck := AMLCheck{
		Address:    "alice",
		CheckType:  "SANCTIONS",
		Status:     "PENDING",
		RiskScore:  75,
		CheckDate:  time.Now(),
		ExpiryDate: time.Now().AddDate(0, 6, 0),
		Details:    "Sanctions check completed",
		CheckedBy:  "SYSTEM",
	}
	
	amlCheckJSON, _ := json.Marshal(amlCheck)
	ctx.stub.On("GetState", "alice_SANCTIONS").Return(amlCheckJSON, nil)
	ctx.stub.On("PutState", "alice_SANCTIONS", mock.Anything).Return(nil)
	ctx.stub.On("GetTxID").Return("tx123")
	ctx.stub.On("SetEvent", "AMLEvent", mock.Anything).Return(nil)
	
	err := c.UpdateAMLCheck(ctx, "alice", "SANCTIONS", "PASSED", 25, "Sanctions check passed")
	assert.NoError(t, err)
	
	ctx.stub.AssertExpectations(t)
}

func TestCompliance_CheckCompliance(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock KYC record
	kyc := KYCRecord{
		Address:     "alice",
		FullName:    "Alice Johnson",
		DateOfBirth: time.Now(),
		Nationality: "US",
		IDType:      "PASSPORT",
		IDNumber:    "US123456",
		Status:      "APPROVED",
	}
	
	kycJSON, _ := json.Marshal(kyc)
	ctx.stub.On("GetState", "alice").Return(kycJSON, nil)
	
	// Mock AML checks - no sanctions or PEP failures
	ctx.stub.On("GetState", "alice_SANCTIONS").Return(nil, nil)
	ctx.stub.On("GetState", "alice_PEP").Return(nil, nil)
	
	compliant, reason, err := c.CheckCompliance(ctx, "alice")
	assert.NoError(t, err)
	assert.True(t, compliant)
	assert.Equal(t, "Compliant", reason)
}

func TestCompliance_CheckCompliance_KYCNotApproved(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Mock KYC record with pending status
	kyc := KYCRecord{
		Address:     "alice",
		FullName:    "Alice Johnson",
		DateOfBirth: time.Now(),
		Nationality: "US",
		IDType:      "PASSPORT",
		IDNumber:    "US123456",
		Status:      "PENDING",
	}
	
	kycJSON, _ := json.Marshal(kyc)
	ctx.stub.On("GetState", "alice").Return(kycJSON, nil)
	
	compliant, reason, err := c.CheckCompliance(ctx, "alice")
	assert.NoError(t, err)
	assert.False(t, compliant)
	assert.Contains(t, reason, "KYC status: PENDING")
}

func TestCompliance_GetKYC(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create a KYC record
	kyc := KYCRecord{
		Address:     "alice",
		FullName:    "Alice Johnson",
		DateOfBirth: time.Now(),
		Nationality: "US",
		IDType:      "PASSPORT",
		IDNumber:    "US123456",
		Status:      "APPROVED",
	}
	
	kycJSON, _ := json.Marshal(kyc)
	ctx.stub.On("GetState", "alice").Return(kycJSON, nil)
	
	retrievedKYC, err := c.GetKYC(ctx, "alice")
	assert.NoError(t, err)
	assert.Equal(t, kyc.Address, retrievedKYC.Address)
	assert.Equal(t, kyc.FullName, retrievedKYC.FullName)
	assert.Equal(t, kyc.Status, retrievedKYC.Status)
}

func TestCompliance_GetKYC_NotFound(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	ctx.stub.On("GetState", "alice").Return(nil, nil)
	
	_, err := c.GetKYC(ctx, "alice")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestCompliance_GetAMLCheck(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create an AML check
	amlCheck := AMLCheck{
		Address:    "alice",
		CheckType:  "SANCTIONS",
		Status:     "PASSED",
		RiskScore:  25,
		CheckDate:  time.Now(),
		ExpiryDate: time.Now().AddDate(0, 6, 0),
		Details:    "Sanctions check passed",
		CheckedBy:  "SYSTEM",
	}
	
	amlCheckJSON, _ := json.Marshal(amlCheck)
	ctx.stub.On("GetState", "alice_SANCTIONS").Return(amlCheckJSON, nil)
	
	retrievedCheck, err := c.GetAMLCheck(ctx, "alice_SANCTIONS")
	assert.NoError(t, err)
	assert.Equal(t, amlCheck.Address, retrievedCheck.Address)
	assert.Equal(t, amlCheck.CheckType, retrievedCheck.CheckType)
	assert.Equal(t, amlCheck.Status, retrievedCheck.Status)
}

func TestCompliance_GetAMLCheck_NotFound(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	ctx.stub.On("GetState", "alice_SANCTIONS").Return(nil, nil)
	
	_, err := c.GetAMLCheck(ctx, "alice_SANCTIONS")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestCompliance_KYCExists(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Test existing KYC
	kyc := KYCRecord{Address: "alice"}
	kycJSON, _ := json.Marshal(kyc)
	ctx.stub.On("GetState", "alice").Return(kycJSON, nil)
	
	exists, err := c.KYCExists(ctx, "alice")
	assert.NoError(t, err)
	assert.True(t, exists)
	
	// Test non-existing KYC
	ctx.stub.On("GetState", "bob").Return(nil, nil)
	
	exists, err = c.KYCExists(ctx, "bob")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCompliance_GetAllKYC(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with KYC results
	kyc1 := KYCRecord{Address: "alice", FullName: "Alice Johnson"}
	kyc2 := KYCRecord{Address: "bob", FullName: "Bob Smith"}
	
	kyc1JSON, _ := json.Marshal(kyc1)
	kyc2JSON, _ := json.Marshal(kyc2)
	
	mockIterator := &MockIterator{results: [][]byte{kyc1JSON, kyc2JSON}}
	
	ctx.stub.On("GetStateByRange", "", "").Return(mockIterator, nil)
	
	kycRecords, err := c.GetAllKYC(ctx)
	assert.NoError(t, err)
	assert.Len(t, kycRecords, 2)
	assert.Equal(t, "alice", kycRecords[0].Address)
	assert.Equal(t, "bob", kycRecords[1].Address)
}

func TestCompliance_GetAllAMLChecks(t *testing.T) {
	c := &Compliance{}
	ctx := &MockContext{stub: &MockStub{state: make(map[string][]byte)}}
	
	// Create mock iterator with AML check results
	aml1 := AMLCheck{Address: "alice", CheckType: "SANCTIONS"}
	aml2 := AMLCheck{Address: "alice", CheckType: "PEP"}
	
	aml1JSON, _ := json.Marshal(aml1)
	aml2JSON, _ := json.Marshal(aml2)
	
	mockIterator := &MockIterator{results: [][]byte{aml1JSON, aml2JSON}}
	
	ctx.stub.On("GetStateByRange", "alice_", "alice_\x00").Return(mockIterator, nil)
	
	amlChecks, err := c.GetAllAMLChecks(ctx, "alice")
	assert.NoError(t, err)
	assert.Len(t, amlChecks, 2)
	assert.Equal(t, "alice", amlChecks[0].Address)
	assert.Equal(t, "alice", amlChecks[1].Address)
}

