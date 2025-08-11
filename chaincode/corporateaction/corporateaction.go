package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CorporateAction represents the corporate action contract
type CorporateAction struct {
	contractapi.Contract
}

// CouponPayment represents a coupon payment
type CouponPayment struct {
	ID          string    `json:"id"`
	BondID      string    `json:"bondId"`
	PaymentDate time.Time `json:"paymentDate"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // "PENDING", "PAID", "FAILED"
	PaidAt      time.Time `json:"paidAt"`
	TxID        string    `json:"txId"`
	Metadata    map[string]string `json:"metadata"`
}

// Redemption represents a bond redemption
type Redemption struct {
	ID          string    `json:"id"`
	BondID      string    `json:"bondId"`
	RedemptionDate time.Time `json:"redemptionDate"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // "PENDING", "COMPLETED", "FAILED"
	CompletedAt time.Time `json:"completedAt"`
	TxID        string    `json:"txId"`
	Metadata    map[string]string `json:"metadata"`
}

// CorporateActionEvent represents a corporate action event
type CorporateActionEvent struct {
	Type      string    `json:"type"`
	BondID    string    `json:"bondId"`
	Details   string    `json:"details"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	TxID      string    `json:"txId"`
}

// Init initializes the contract
func (ca *CorporateAction) Init(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("CorporateAction contract initialized")
	return nil
}

// CreateCouponPayment creates a new coupon payment
func (ca *CorporateAction) CreateCouponPayment(ctx contractapi.TransactionContextInterface, bondID, paymentDateStr string, amount float64) error {
	// Generate unique ID for coupon payment
	couponID := fmt.Sprintf("COUPON_%s_%s", bondID, time.Now().Format("20060102"))
	
	// Parse payment date
	paymentDate, err := time.Parse("2006-01-02", paymentDateStr)
	if err != nil {
		return fmt.Errorf("invalid payment date format: %v", err)
	}

	// Create new coupon payment
	couponPayment := CouponPayment{
		ID:          couponID,
		BondID:      bondID,
		PaymentDate: paymentDate,
		Amount:      amount,
		Status:      "PENDING",
		Metadata:    make(map[string]string),
	}

	// Store coupon payment
	couponJSON, err := json.Marshal(couponPayment)
	if err != nil {
		return fmt.Errorf("failed to marshal coupon payment: %v", err)
	}

	err = ctx.GetStub().PutState(couponID, couponJSON)
	if err != nil {
		return fmt.Errorf("failed to store coupon payment: %v", err)
	}

	// Emit event
	event := CorporateActionEvent{
		Type:      "COUPON_PAYMENT_CREATED",
		BondID:    bondID,
		Details:   fmt.Sprintf("Coupon payment created for bond %s", bondID),
		Amount:    amount,
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("CorporateActionEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// ProcessCouponPayment processes a coupon payment
func (ca *CorporateAction) ProcessCouponPayment(ctx contractapi.TransactionContextInterface, couponID string) error {
	couponPayment, err := ca.GetCouponPayment(ctx, couponID)
	if err != nil {
		return fmt.Errorf("failed to get coupon payment: %v", err)
	}

	if couponPayment.Status != "PENDING" {
		return fmt.Errorf("coupon payment %s is not pending", couponID)
	}

	// Update status to paid
	couponPayment.Status = "PAID"
	couponPayment.PaidAt = time.Now()
	couponPayment.TxID = ctx.GetStub().GetTxID()

	// Store updated coupon payment
	couponJSON, err := json.Marshal(couponPayment)
	if err != nil {
		return fmt.Errorf("failed to marshal coupon payment: %v", err)
	}

	err = ctx.GetStub().PutState(couponID, couponJSON)
	if err != nil {
		return fmt.Errorf("failed to update coupon payment: %v", err)
	}

	// Emit event
	event := CorporateActionEvent{
		Type:      "COUPON_PAYMENT_PROCESSED",
		BondID:    couponPayment.BondID,
		Details:   fmt.Sprintf("Coupon payment %s processed", couponID),
		Amount:    couponPayment.Amount,
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("CorporateActionEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// CreateRedemption creates a new bond redemption
func (ca *CorporateAction) CreateRedemption(ctx contractapi.TransactionContextInterface, bondID, redemptionDateStr string, amount float64) error {
	// Generate unique ID for redemption
	redemptionID := fmt.Sprintf("REDEMPTION_%s_%s", bondID, time.Now().Format("20060102"))
	
	// Parse redemption date
	redemptionDate, err := time.Parse("2006-01-02", redemptionDateStr)
	if err != nil {
		return fmt.Errorf("invalid redemption date format: %v", err)
	}

	// Create new redemption
	redemption := Redemption{
		ID:             redemptionID,
		BondID:         bondID,
		RedemptionDate: redemptionDate,
		Amount:         amount,
		Status:         "PENDING",
		Metadata:       make(map[string]string),
	}

	// Store redemption
	redemptionJSON, err := json.Marshal(redemption)
	if err != nil {
		return fmt.Errorf("failed to marshal redemption: %v", err)
	}

	err = ctx.GetStub().PutState(redemptionID, redemptionJSON)
	if err != nil {
		return fmt.Errorf("failed to store redemption: %v", err)
	}

	// Emit event
	event := CorporateActionEvent{
		Type:      "REDEMPTION_CREATED",
		BondID:    bondID,
		Details:   fmt.Sprintf("Redemption created for bond %s", bondID),
		Amount:    amount,
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("CorporateActionEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// ProcessRedemption processes a bond redemption
func (ca *CorporateAction) ProcessRedemption(ctx contractapi.TransactionContextInterface, redemptionID string) error {
	redemption, err := ca.GetRedemption(ctx, redemptionID)
	if err != nil {
		return fmt.Errorf("failed to get redemption: %v", err)
	}

	if redemption.Status != "PENDING" {
		return fmt.Errorf("redemption %s is not pending", redemptionID)
	}

	// Update status to completed
	redemption.Status = "COMPLETED"
	redemption.CompletedAt = time.Now()
	redemption.TxID = ctx.GetStub().GetTxID()

	// Store updated redemption
	redemptionJSON, err := json.Marshal(redemption)
	if err != nil {
		return fmt.Errorf("failed to marshal redemption: %v", err)
	}

	err = ctx.GetStub().PutState(redemptionID, redemptionJSON)
	if err != nil {
		return fmt.Errorf("failed to update redemption: %v", err)
	}

	// Emit event
	event := CorporateActionEvent{
		Type:      "REDEMPTION_PROCESSED",
		BondID:    redemption.BondID,
		Details:   fmt.Sprintf("Redemption %s processed", redemptionID),
		Amount:    redemption.Amount,
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("CorporateActionEvent", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// GetCouponPayment retrieves a coupon payment
func (ca *CorporateAction) GetCouponPayment(ctx contractapi.TransactionContextInterface, couponID string) (*CouponPayment, error) {
	couponJSON, err := ctx.GetStub().GetState(couponID)
	if err != nil {
		return nil, fmt.Errorf("failed to read coupon payment: %v", err)
	}
	if couponJSON == nil {
		return nil, fmt.Errorf("coupon payment %s does not exist", couponID)
	}

	var couponPayment CouponPayment
	err = json.Unmarshal(couponJSON, &couponPayment)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal coupon payment: %v", err)
	}

	return &couponPayment, nil
}

// GetRedemption retrieves a redemption
func (ca *CorporateAction) GetRedemption(ctx contractapi.TransactionContextInterface, redemptionID string) (*Redemption, error) {
	redemptionJSON, err := ctx.GetStub().GetState(redemptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to read redemption: %v", err)
	}
	if redemptionJSON == nil {
		return nil, fmt.Errorf("redemption %s does not exist", redemptionID)
	}

	var redemption Redemption
	err = json.Unmarshal(redemptionJSON, &redemption)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal redemption: %v", err)
	}

	return &redemption, nil
}

// GetCouponPaymentsByBond returns all coupon payments for a specific bond
func (ca *CorporateAction) GetCouponPaymentsByBond(ctx contractapi.TransactionContextInterface, bondID string) ([]*CouponPayment, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var couponPayments []*CouponPayment
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		// Check if this is a coupon payment for the specific bond
		if len(queryResult.Key) > 8 && queryResult.Key[:8] == "COUPON_" && contains(queryResult.Key, bondID) {
			var couponPayment CouponPayment
			err = json.Unmarshal(queryResult.Value, &couponPayment)
			if err == nil && couponPayment.BondID == bondID {
				couponPayments = append(couponPayments, &couponPayment)
			}
		}
	}

	return couponPayments, nil
}

// GetRedemptionsByBond returns all redemptions for a specific bond
func (ca *CorporateAction) GetRedemptionsByBond(ctx contractapi.TransactionContextInterface, bondID string) ([]*Redemption, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var redemptions []*Redemption
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		// Check if this is a redemption for the specific bond
		if len(queryResult.Key) > 11 && queryResult.Key[:11] == "REDEMPTION_" && contains(queryResult.Key, bondID) {
			var redemption Redemption
			err = json.Unmarshal(queryResult.Value, &redemption)
			if err == nil && redemption.BondID == bondID {
				redemptions = append(redemptions, &redemption)
			}
		}
	}

	return redemptions, nil
}

// GetPendingCouponPayments returns all pending coupon payments
func (ca *CorporateAction) GetPendingCouponPayments(ctx contractapi.TransactionContextInterface) ([]*CouponPayment, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var pendingPayments []*CouponPayment
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		// Check if this is a pending coupon payment
		if len(queryResult.Key) > 8 && queryResult.Key[:8] == "COUPON_" {
			var couponPayment CouponPayment
			err = json.Unmarshal(queryResult.Value, &couponPayment)
			if err == nil && couponPayment.Status == "PENDING" {
				pendingPayments = append(pendingPayments, &couponPayment)
			}
		}
	}

	return pendingPayments, nil
}

// GetPendingRedemptions returns all pending redemptions
func (ca *CorporateAction) GetPendingRedemptions(ctx contractapi.TransactionContextInterface) ([]*Redemption, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var pendingRedemptions []*Redemption
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		// Check if this is a pending redemption
		if len(queryResult.Key) > 11 && queryResult.Key[:11] == "REDEMPTION_" {
			var redemption Redemption
			err = json.Unmarshal(queryResult.Value, &redemption)
			if err == nil && redemption.Status == "PENDING" {
				pendingRedemptions = append(pendingRedemptions, &redemption)
			}
		}
	}

	return pendingRedemptions, nil
}

// CalculateCouponAmount calculates the coupon amount for a bond
func (ca *CorporateAction) CalculateCouponAmount(ctx contractapi.TransactionContextInterface, bondID string, faceValue float64, couponRate float64) (float64, error) {
	// Simple calculation: (Face Value * Coupon Rate) / 100
	couponAmount := (faceValue * couponRate) / 100
	return couponAmount, nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && s[:len(substr)] == substr))
}

func main() {
	chaincode, err := contractapi.NewChaincode(&CorporateAction{})
	if err != nil {
		fmt.Printf("Error creating CorporateAction chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting CorporateAction chaincode: %s", err.Error())
	}
}

