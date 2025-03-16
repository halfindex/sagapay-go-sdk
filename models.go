package sagapay

import (
	"errors"
	"fmt"
	"time"
)

// NetworkType represents the blockchain network type
type NetworkType string

// Network types
const (
	NetworkTypeERC20   NetworkType = "ERC20"
	NetworkTypeBEP20   NetworkType = "BEP20"
	NetworkTypeTRC20   NetworkType = "TRC20"
	NetworkTypePOLYGON NetworkType = "POLYGON"
	NetworkTypeSOLANA  NetworkType = "SOLANA"
)

// TransactionType represents the type of transaction
type TransactionType string

// Transaction types
const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

// Transaction statuses
const (
	TransactionStatusPending    TransactionStatus = "PENDING"
	TransactionStatusProcessing TransactionStatus = "PROCESSING"
	TransactionStatusCompleted  TransactionStatus = "COMPLETED"
	TransactionStatusFailed     TransactionStatus = "FAILED"
	TransactionStatusCancelled  TransactionStatus = "CANCELLED"
)

// AddressType represents the type of address
type AddressType string

// Address types
const (
	AddressTypeTemporary AddressType = "TEMPORARY"
	AddressTypePermanent AddressType = "PERMANENT"
)

// CreateDepositParams represents the parameters for creating a deposit
type CreateDepositParams struct {
	NetworkType    NetworkType `json:"networkType"`
	ContractAddress string      `json:"contractAddress"`
	Amount         string      `json:"amount"`
	IPNUrl         string      `json:"ipnUrl"`
	UDF            string      `json:"udf,omitempty"`
	Type           AddressType `json:"type,omitempty"`
}

// Validate validates the create deposit parameters
func (p *CreateDepositParams) Validate() error {
	if p.NetworkType == "" {
		return errors.New("networkType is required")
	}
	if p.ContractAddress == "" {
		return errors.New("contractAddress is required")
	}
	if p.Amount == "" {
		return errors.New("amount is required")
	}
	if p.IPNUrl == "" {
		return errors.New("ipnUrl is required")
	}
	return nil
}

// CreateWithdrawalParams represents the parameters for creating a withdrawal
type CreateWithdrawalParams struct {
	NetworkType     NetworkType `json:"networkType"`
	ContractAddress string      `json:"contractAddress"`
	Address         string      `json:"address"`
	Amount          string      `json:"amount"`
	IPNUrl          string      `json:"ipnUrl"`
	UDF             string      `json:"udf,omitempty"`
}

// Validate validates the create withdrawal parameters
func (p *CreateWithdrawalParams) Validate() error {
	if p.NetworkType == "" {
		return errors.New("networkType is required")
	}
	if p.ContractAddress == "" {
		return errors.New("contractAddress is required")
	}
	if p.Address == "" {
		return errors.New("address is required")
	}
	if p.Amount == "" {
		return errors.New("amount is required")
	}
	if p.IPNUrl == "" {
		return errors.New("ipnUrl is required")
	}
	return nil
}

// DepositResponse represents the response from creating a deposit
type DepositResponse struct {
	ID        string            `json:"id"`
	Address   string            `json:"address"`
	ExpiresAt time.Time         `json:"expiresAt"`
	Amount    string            `json:"amount"`
	Status    TransactionStatus `json:"status"`
}

// WithdrawalResponse represents the response from creating a withdrawal
type WithdrawalResponse struct {
	ID     string            `json:"id"`
	Status TransactionStatus `json:"status"`
	Fee    string            `json:"fee"`
}

// Token represents a cryptocurrency token
type Token struct {
	NetworkType     NetworkType `json:"networkType"`
	ContractAddress string      `json:"contractAddress"`
	Symbol          string      `json:"symbol"`
	Name            string      `json:"name"`
	Decimals        int         `json:"decimals"`
}

// Transaction represents a cryptocurrency transaction
type Transaction struct {
	ID              string            `json:"id"`
	TransactionType TransactionType   `json:"transactionType"`
	Status          TransactionStatus `json:"status"`
	Amount          string            `json:"amount"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	TxHash          string            `json:"txHash,omitempty"`
	NetworkType     NetworkType       `json:"networkType"`
	ContractAddress string            `json:"contractAddress"`
	Address         string            `json:"address"`
	Token           Token             `json:"token"`
}

// TransactionStatusResponse represents the response from checking transaction status
type TransactionStatusResponse struct {
	Address         string        `json:"address"`
	TransactionType TransactionType `json:"transactionType"`
	Count           int           `json:"count"`
	Transactions    []Transaction `json:"transactions"`
}

// Balance represents a wallet balance
type Balance struct {
	Raw       string `json:"raw"`
	Formatted string `json:"formatted"`
}

// WalletBalanceResponse represents the response from fetching wallet balance
type WalletBalanceResponse struct {
	Address         string      `json:"address"`
	NetworkType     NetworkType `json:"networkType"`
	ContractAddress string      `json:"contractAddress"`
	Token           Token       `json:"token"`
	Balance         Balance     `json:"balance"`
}

// WebhookPayload represents the payload sent in webhook notifications
type WebhookPayload struct {
	ID              string            `json:"id"`
	Type            TransactionType   `json:"type"`
	Status          TransactionStatus `json:"status"`
	Address         string            `json:"address"`
	NetworkType     NetworkType       `json:"networkType"`
	Amount          string            `json:"amount"`
	UDF             string            `json:"udf,omitempty"`
	TxHash          string            `json:"txHash,omitempty"`
	Timestamp       time.Time         `json:"timestamp"`
}

// APIError represents an error response from the API
type APIError struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s - %s", e.Error, e.Message)
}