package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// BondToken represents a bond token on the blockchain
type BondToken struct {
	contractapi.Contract
}

// Bond represents a corporate bond
type Bond struct {
	ID              string    `json:"id"`
	IssuerID        string    `json:"issuerId"`
	IssuerName      string    `json:"issuerName"`
	FaceValue       float64   `json:"faceValue"`
	CouponRate      float64   `json:"couponRate"`
	MaturityDate    time.Time `json:"maturityDate"`
	IssueDate       time.Time `json:"issueDate"`
	TotalSupply     int64     `json:"totalSupply"`
	AvailableSupply int64     `json:"availableSupply"`
	Status          string    `json:"status"` // "ACTIVE", "MATURED", "DEFAULTED"
	Currency        string    `json:"currency"`
	ISIN            string    `json:"isin"`
	Rating          string    `json:"rating"`
	Collateral      string    `json:"collateral"`
}

// TokenHolder represents a token holder
type TokenHolder struct {
	Address     string            `json:"address"`
	BondID      string            `json:"bondId"`
	Quantity    int64             `json:"quantity"`
	LastUpdated time.Time         `json:"lastUpdated"`
	Metadata    map[string]string `json:"metadata"`
}

// TransferEvent represents a token transfer event
type TransferEvent struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	BondID    string    `json:"bondId"`
	Quantity  int64     `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
	TxID      string    `json:"txId"`
}

// Init initializes the contract
func (bt *BondToken) Init(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("BondToken contract initialized")
	return nil
}

// IssueBond issues a new bond
func (bt *BondToken) IssueBond(ctx contractapi.TransactionContextInterface, bondID, issuerID, issuerName, currency, isin, rating, collateral string, faceValue float64, couponRate float64, totalSupply int64, maturityDateStr string) error {
	// Check if bond already exists
	exists, err := bt.BondExists(ctx, bondID)
	if err != nil {
		return fmt.Errorf("failed to check bond existence: %v", err)
	}
	if exists {
		return fmt.Errorf("bond %s already exists", bondID)
	}

	// Parse maturity date
	maturityDate, err := time.Parse("2006-01-02", maturityDateStr)
	if err != nil {
		return fmt.Errorf("invalid maturity date format: %v", err)
	}

	// Create new bond
	bond := Bond{
		ID:              bondID,
		IssuerID:        issuerID,
		IssuerName:      issuerName,
		FaceValue:       faceValue,
		CouponRate:      couponRate,
		MaturityDate:    maturityDate,
		IssueDate:       time.Now(),
		TotalSupply:     totalSupply,
		AvailableSupply: totalSupply,
		Status:          "ACTIVE",
		Currency:        currency,
		ISIN:            isin,
		Rating:          rating,
		Collateral:      collateral,
	}

	// Store bond
	bondJSON, err := json.Marshal(bond)
	if err != nil {
		return fmt.Errorf("failed to marshal bond: %v", err)
	}

	err = ctx.GetStub().PutState(bondID, bondJSON)
	if err != nil {
		return fmt.Errorf("failed to store bond: %v", err)
	}

	// Emit event
	event := TransferEvent{
		From:      "SYSTEM",
		To:        issuerID,
		BondID:    bondID,
		Quantity:  totalSupply,
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("BondIssued", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// Transfer transfers tokens from one address to another
func (bt *BondToken) Transfer(ctx contractapi.TransactionContextInterface, from, to, bondID string, quantity int64) error {
	// Check if bond exists
	exists, err := bt.BondExists(ctx, bondID)
	if err != nil {
		return fmt.Errorf("failed to check bond existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("bond %s does not exist", bondID)
	}

	// Check if quantity is positive
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Get sender's balance
	senderKey := fmt.Sprintf("%s_%s", from, bondID)
	senderHolder, err := bt.GetTokenHolder(ctx, senderKey)
	if err != nil {
		return fmt.Errorf("failed to get sender holder: %v", err)
	}

	if senderHolder.Quantity < quantity {
		return fmt.Errorf("insufficient balance: %d < %d", senderHolder.Quantity, quantity)
	}

	// Get recipient's balance
	recipientKey := fmt.Sprintf("%s_%s", to, bondID)
	recipientHolder, err := bt.GetTokenHolder(ctx, recipientKey)
	if err != nil {
		// Create new holder if doesn't exist
		recipientHolder = &TokenHolder{
			Address:     to,
			BondID:      bondID,
			Quantity:    0,
			LastUpdated: time.Now(),
			Metadata:    make(map[string]string),
		}
	}

	// Update balances
	senderHolder.Quantity -= quantity
	senderHolder.LastUpdated = time.Now()

	recipientHolder.Quantity += quantity
	recipientHolder.LastUpdated = time.Now()

	// Store updated holders
	senderJSON, err := json.Marshal(senderHolder)
	if err != nil {
		return fmt.Errorf("failed to marshal sender holder: %v", err)
	}

	recipientJSON, err := json.Marshal(recipientHolder)
	if err != nil {
		return fmt.Errorf("failed to marshal recipient holder: %v", err)
	}

	err = ctx.GetStub().PutState(senderKey, senderJSON)
	if err != nil {
		return fmt.Errorf("failed to store sender holder: %v", err)
	}

	err = ctx.GetStub().PutState(recipientKey, recipientJSON)
	if err != nil {
		return fmt.Errorf("failed to store recipient holder: %v", err)
	}

	// Emit transfer event
	event := TransferEvent{
		From:      from,
		To:        to,
		BondID:    bondID,
		Quantity:  quantity,
		Timestamp: time.Now(),
		TxID:      ctx.GetStub().GetTxID(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().SetEvent("TokensTransferred", eventJSON)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	return nil
}

// GetBond retrieves a bond by ID
func (bt *BondToken) GetBond(ctx contractapi.TransactionContextInterface, bondID string) (*Bond, error) {
	bondJSON, err := ctx.GetStub().GetState(bondID)
	if err != nil {
		return nil, fmt.Errorf("failed to read bond: %v", err)
	}
	if bondJSON == nil {
		return nil, fmt.Errorf("bond %s does not exist", bondID)
	}

	var bond Bond
	err = json.Unmarshal(bondJSON, &bond)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bond: %v", err)
	}

	return &bond, nil
}

// GetTokenHolder retrieves a token holder
func (bt *BondToken) GetTokenHolder(ctx contractapi.TransactionContextInterface, holderKey string) (*TokenHolder, error) {
	holderJSON, err := ctx.GetStub().GetState(holderKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read holder: %v", err)
	}
	if holderJSON == nil {
		return nil, fmt.Errorf("holder %s does not exist", holderKey)
	}

	var holder TokenHolder
	err = json.Unmarshal(holderJSON, &holder)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal holder: %v", err)
	}

	return &holder, nil
}

// BondExists checks if a bond exists
func (bt *BondToken) BondExists(ctx contractapi.TransactionContextInterface, bondID string) (bool, error) {
	bondJSON, err := ctx.GetStub().GetState(bondID)
	if err != nil {
		return false, fmt.Errorf("failed to read bond: %v", err)
	}
	return bondJSON != nil, nil
}

// GetBalance returns the balance of a specific bond for a specific address
func (bt *BondToken) GetBalance(ctx contractapi.TransactionContextInterface, address, bondID string) (int64, error) {
	holderKey := fmt.Sprintf("%s_%s", address, bondID)
	holder, err := bt.GetTokenHolder(ctx, holderKey)
	if err != nil {
		// Return 0 if holder doesn't exist
		return 0, nil
	}
	return holder.Quantity, nil
}

// GetAllBonds returns all bonds
func (bt *BondToken) GetAllBonds(ctx contractapi.TransactionContextInterface) ([]*Bond, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var bonds []*Bond
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		// Check if this is a bond (not a holder)
		if len(queryResult.Key) < 20 { // Bonds have shorter keys than holders
			var bond Bond
			err = json.Unmarshal(queryResult.Value, &bond)
			if err == nil && bond.ID != "" {
				bonds = append(bonds, &bond)
			}
		}
	}

	return bonds, nil
}

// UpdateBondStatus updates the status of a bond
func (bt *BondToken) UpdateBondStatus(ctx contractapi.TransactionContextInterface, bondID, newStatus string) error {
	bond, err := bt.GetBond(ctx, bondID)
	if err != nil {
		return fmt.Errorf("failed to get bond: %v", err)
	}

	bond.Status = newStatus
	bondJSON, err := json.Marshal(bond)
	if err != nil {
		return fmt.Errorf("failed to marshal bond: %v", err)
	}

	err = ctx.GetStub().PutState(bondID, bondJSON)
	if err != nil {
		return fmt.Errorf("failed to update bond: %v", err)
	}

	return nil
}

// GetBondHolders returns all holders of a specific bond
func (bt *BondToken) GetBondHolders(ctx contractapi.TransactionContextInterface, bondID string) ([]*TokenHolder, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by range: %v", err)
	}
	defer resultsIterator.Close()

	var holders []*TokenHolder
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		// Check if this is a holder for the specific bond
		if len(queryResult.Key) > 20 && queryResult.Key[len(queryResult.Key)-len(bondID)-1:] == "_"+bondID {
			var holder TokenHolder
			err = json.Unmarshal(queryResult.Value, &holder)
			if err == nil && holder.BondID == bondID {
				holders = append(holders, &holder)
			}
		}
	}

	return holders, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&BondToken{})
	if err != nil {
		fmt.Printf("Error creating BondToken chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting BondToken chaincode: %s", err.Error())
	}
}

