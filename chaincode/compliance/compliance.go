package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Compliance represents the compliance contract
type Compliance struct {
	contractapi.Contract
}

// KYCRecord represents a KYC record
type KYCRecord struct {
	Address       string    `json:"address"`
	FullName      string    `json:"fullName"`
	DateOfBirth   string    `json:"dateOfBirth"`
	Nationality   string    `json:"nationality"`
	IDType        string    `json:"idType"`
	IDNumber      string    `json:"idNumber"`
	Status        string    `json:"status"` // "PENDING", "APPROVED", "REJECTED"
	RiskLevel     string    `json:"riskLevel"` // "LOW", "MEDIUM", "HIGH"
	ApprovedBy    string    `json:"approvedBy"`
	ApprovedAt    time.Time `json:"approvedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Metadata      map[string]string `json:"metadata"`
}

// AMLCheck represents an AML check
type AMLCheck struct {
	Address       string    `json:"address"`
	CheckType     string    `json:"checkType"` // "SANCTIONS", "PEP", "ADVERSE_MEDIA"
	Status        string    `json:"status"` // "PASSED", "FAILED", "PENDING"
	RiskScore     int       `json:"riskScore"`
	CheckDate     time.Time `json:"checkDate"`
	ExpiryDate    time.Time `json:"expiryDate"`
	Details       string    `json:"details"`
	CheckedBy     string    `json:"checkedBy"`
}

// ComplianceRule represents a compliance rule
type ComplianceRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "KYC", "AML", "TRADING_LIMIT"
	Status      string    `json:"status"` // "ACTIVE", "INACTIVE"
	Parameters  map[string]interface{} `json:"parameters"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ComplianceEvent represents a compliance event
type ComplianceEvent struct {
	Type      string    `json:"type"`
	Address   string    `json:"address"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
	TxID      string    `json:"txId"`
}

// Init initializes the contract
func (c *Compliance) Init(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("Compliance contract initialized")
	return nil
}

// CreateKYC creates a new KYC record
func (c *Compliance) CreateKYC(ctx contractapi.TransactionContextInterface, address, fullName, dateOfBirth, nationality, idType, idNumber string) error {
	// Check if KYC already exists
	exists, err := c.KYCExists(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to check KYC existence: %v", err)
	}
	if exists {
		return fmt.Errorf("KYC for address %s already exists", address)
	}

	// Create new KYC record
	kyc := KYCRecord{
		Address:     address,
		FullName:    fullName,
		DateOfBirth: dateOfBirth,
		Nationality: nationality,
		IDType:      idType,
		IDNumber:    idNumber,
		Status:      "PENDING",
		RiskLevel:   "MEDIUM",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]string),
	}

	// Store KYC record
	kycJSON, err := json.Marshal(kyc)
	if err != nil {
		return fmt.Errorf("failed to marshal KYC: %v", err)
	}

	err = ctx.GetStub().PutState(address, kycJSON)
	if err != nil {
		return fmt.Errorf("failed to store KYC: %v", err)
	}

	// Emit event
	event := ComplianceEvent{
		Type:      "KYC_CREATED",
		Address:   address,
		Details:   fmt.Sprintf("KYC created for %s", fullName),
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("KYCEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// ApproveKYC approves a KYC record
func (c *Compliance) ApproveKYC(ctx contractapi.TransactionContextInterface, address, approvedBy, riskLevel string) error {
	kyc, err := c.GetKYC(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get KYC: %v", err)
	}

	kyc.Status = "APPROVED"
	kyc.RiskLevel = riskLevel
	kyc.ApprovedBy = approvedBy
	kyc.ApprovedAt = time.Now()
	kyc.UpdatedAt = time.Now()

	kycJSON, err := json.Marshal(kyc)
	if err != nil {
		return fmt.Errorf("failed to marshal KYC: %v", err)
	}

	err = ctx.GetStub().PutState(address, kycJSON)
	if err != nil {
		return fmt.Errorf("failed to update KYC: %v", err)
	}

	// Emit event
	event := ComplianceEvent{
		Type:      "KYC_APPROVED",
		Address:   address,
		Details:   fmt.Sprintf("KYC approved by %s", approvedBy),
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("KYCEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// RejectKYC rejects a KYC record
func (c *Compliance) RejectKYC(ctx contractapi.TransactionContextInterface, address, rejectedBy, reason string) error {
	kyc, err := c.GetKYC(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get KYC: %v", err)
	}

	kyc.Status = "REJECTED"
	kyc.UpdatedAt = time.Now()
	kyc.Metadata["rejection_reason"] = reason
	kyc.Metadata["rejected_by"] = rejectedBy

	kycJSON, err := json.Marshal(kyc)
	if err != nil {
		return fmt.Errorf("failed to marshal KYC: %v", err)
	}

	err = ctx.GetStub().PutState(address, kycJSON)
	if err != nil {
		return fmt.Errorf("failed to update KYC: %v", err)
	}

	// Emit event
	event := ComplianceEvent{
		Type:      "KYC_REJECTED",
		Address:   address,
		Details:   fmt.Sprintf("KYC rejected by %s: %s", rejectedBy, reason),
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("KYCEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// CreateAMLCheck creates a new AML check
func (c *Compliance) CreateAMLCheck(ctx contractapi.TransactionContextInterface, address, checkType string, riskScore int, details string) error {
	checkKey := fmt.Sprintf("%s_%s", address, checkType)
	
	// Create new AML check
	amlCheck := AMLCheck{
		Address:    address,
		CheckType:  checkType,
		Status:     "PENDING",
		RiskScore:  riskScore,
		CheckDate:  time.Now(),
		ExpiryDate: time.Now().AddDate(0, 6, 0), // 6 months validity
		Details:    details,
		CheckedBy:  "SYSTEM",
	}

	// Store AML check
	checkJSON, err := json.Marshal(amlCheck)
	if err != nil {
		return fmt.Errorf("failed to marshal AML check: %v", err)
	}

	err = ctx.GetStub().PutState(checkKey, checkJSON)
	if err != nil {
		return fmt.Errorf("failed to store AML check: %v", err)
	}

	// Emit event
	event := ComplianceEvent{
		Type:      "AML_CHECK_CREATED",
		Address:   address,
		Details:   fmt.Sprintf("AML check created for %s: %s", address, checkType),
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("AMLEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// UpdateAMLCheck updates an AML check
func (c *Compliance) UpdateAMLCheck(ctx contractapi.TransactionContextInterface, address, checkType, status string, riskScore int, details string) error {
	checkKey := fmt.Sprintf("%s_%s", address, checkType)
	
	checkJSON, err := ctx.GetStub().GetState(checkKey)
	if err != nil {
		return fmt.Errorf("failed to read AML check: %v", err)
	}
	if checkJSON == nil {
		return fmt.Errorf("AML check %s does not exist", checkKey)
	}

	var amlCheck AMLCheck
	err = json.Unmarshal(checkJSON, &amlCheck)
	if err != nil {
		return fmt.Errorf("failed to unmarshal AML check: %v", err)
	}

	amlCheck.Status = status
	amlCheck.RiskScore = riskScore
	amlCheck.Details = details
	amlCheck.CheckDate = time.Now()

	// Store updated AML check
	updatedCheckJSON, err := json.Marshal(amlCheck)
	if err != nil {
		return fmt.Errorf("failed to marshal updated AML check: %v", err)
	}

	err = ctx.GetStub().PutState(checkKey, updatedCheckJSON)
	if err != nil {
		return fmt.Errorf("failed to update AML check: %v", err)
	}

	// Emit event
	event := ComplianceEvent{
		Type:      "AML_CHECK_UPDATED",
		Address:   address,
		Details:   fmt.Sprintf("AML check updated for %s: %s - %s", address, checkType, status),
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("AMLEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// CheckCompliance checks if an address is compliant
func (c *Compliance) CheckCompliance(ctx contractapi.TransactionContextInterface, address string) (bool, string, error) {
	// Check KYC status
	kyc, err := c.GetKYC(ctx, address)
	if err != nil {
		return false, "KYC record not found", nil
	}

	if kyc.Status != "APPROVED" {
		return false, fmt.Sprintf("KYC status: %s", kyc.Status), nil
	}

	// Check AML status
	sanctionsKey := fmt.Sprintf("%s_SANCTIONS", address)
	pepKey := fmt.Sprintf("%s_PEP", address)
	
	sanctionsCheck, err := c.GetAMLCheck(ctx, sanctionsKey)
	if err == nil && sanctionsCheck.Status == "FAILED" {
		return false, "Sanctions check failed", nil
	}

	pepCheck, err := c.GetAMLCheck(ctx, pepKey)
	if err == nil && pepCheck.Status == "FAILED" {
		return false, "PEP check failed", nil
	}

	return true, "Compliant", nil
}

// GetKYC retrieves a KYC record
func (c *Compliance) GetKYC(ctx contractapi.TransactionContextInterface, address string) (*KYCRecord, error) {
	kycJSON, err := ctx.GetStub().GetState(address)
	if err != nil {
		return nil, fmt.Errorf("failed to read KYC: %v", err)
	}
	if kycJSON == nil {
		return nil, fmt.Errorf("KYC for address %s does not exist", address)
	}

	var kyc KYCRecord
	err = json.Unmarshal(kycJSON, &kyc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal KYC: %v", err)
	}

	return &kyc, nil
}

// GetAMLCheck retrieves an AML check
func (c *Compliance) GetAMLCheck(ctx contractapi.TransactionContextInterface, checkKey string) (*AMLCheck, error) {
	checkJSON, err := ctx.GetStub().GetState(checkKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read AML check: %v", err)
	}
	if checkJSON == nil {
		return nil, fmt.Errorf("AML check %s does not exist", checkKey)
	}

	var amlCheck AMLCheck
	err = json.Unmarshal(checkJSON, &amlCheck)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal AML check: %v", err)
	}

	return &amlCheck, nil
}

// KYCExists checks if a KYC record exists
func (c *Compliance) KYCExists(ctx contractapi.TransactionContextInterface, address string) (bool, error) {
	kycJSON, err := ctx.GetStub().GetState(address)
	if err != nil {
		return false, fmt.Errorf("failed to read KYC: %v", err)
	}
	return kycJSON != nil, nil
}

// GetAllKYC returns all KYC records
func (c *Compliance) GetAllKYC(ctx contractapi.TransactionContextInterface) ([]*KYCRecord, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var kycRecords []*KYCRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		// Check if this is a KYC record (not an AML check)
		if !contains(queryResult.Key, "_") {
			var kyc KYCRecord
			err = json.Unmarshal(queryResult.Value, &kyc)
			if err == nil && kyc.Address != "" {
				kycRecords = append(kycRecords, &kyc)
			}
		}
	}

	return kycRecords, nil
}

// GetAllAMLChecks returns all AML checks for an address
func (c *Compliance) GetAllAMLChecks(ctx contractapi.TransactionContextInterface, address string) ([]*AMLCheck, error) {
	startKey := fmt.Sprintf("%s_", address)
	endKey := fmt.Sprintf("%s_%c", address, 0)

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var amlChecks []*AMLCheck
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		var amlCheck AMLCheck
		err = json.Unmarshal(queryResult.Value, &amlCheck)
		if err == nil && amlCheck.Address == address {
			amlChecks = append(amlChecks, &amlCheck)
		}
	}

	return amlChecks, nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && s[:len(substr)] == substr))
}

func main() {
	chaincode, err := contractapi.NewChaincode(&Compliance{})
	if err != nil {
		fmt.Printf("Error creating Compliance chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Compliance chaincode: %s", err.Error())
	}
}

