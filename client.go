// Package sagapay provides a Go client for the SagaPay blockchain payment gateway API.
//
// SagaPay is the world's first free, non-custodial blockchain payment gateway service provider.
// This package enables Go developers to integrate cryptocurrency payments without holding customer funds.
package sagapay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the SagaPay API
	DefaultBaseURL = "https://api.sagapay.net"
	
	// DefaultTimeout is the default timeout for API requests
	DefaultTimeout = 30 * time.Second
)

// Client is the SagaPay API client
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Base URL for API requests
	baseURL *url.URL

	// API credentials
	apiKey    string
	apiSecret string
}

// Config contains the configuration options for the SagaPay client
type Config struct {
	// BaseURL is the base URL for the SagaPay API
	BaseURL string

	// APIKey is your SagaPay API key
	APIKey string

	// APISecret is your SagaPay API secret
	APISecret string

	// Timeout is the timeout for API requests
	Timeout time.Duration

	// HTTPClient is the HTTP client to use for API requests
	HTTPClient *http.Client
}

// NewClient creates a new SagaPay API client
func NewClient(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if config.APISecret == "" {
		return nil, fmt.Errorf("API secret is required")
	}

	baseURL := DefaultBaseURL
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		timeout := DefaultTimeout
		if config.Timeout > 0 {
			timeout = config.Timeout
		}
		httpClient = &http.Client{
			Timeout: timeout,
		}
	}

	return &Client{
		client:    httpClient,
		baseURL:   parsedURL,
		apiKey:    config.APIKey,
		apiSecret: config.APISecret,
	}, nil
}

// CreateDeposit creates a new deposit address for receiving cryptocurrency
func (c *Client) CreateDeposit(ctx context.Context, params CreateDepositParams) (*DepositResponse, error) {
	endpoint := "/create-deposit"
	
	// Validate params
	if err := params.Validate(); err != nil {
		return nil, err
	}

	var response DepositResponse
	err := c.sendRequest(ctx, http.MethodPost, endpoint, params, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateWithdrawal creates a cryptocurrency withdrawal request
func (c *Client) CreateWithdrawal(ctx context.Context, params CreateWithdrawalParams) (*WithdrawalResponse, error) {
	endpoint := "/create-withdrawal"
	
	// Validate params
	if err := params.Validate(); err != nil {
		return nil, err
	}

	var response WithdrawalResponse
	err := c.sendRequest(ctx, http.MethodPost, endpoint, params, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CheckTransactionStatus gets the status of transactions for a specific blockchain address
func (c *Client) CheckTransactionStatus(ctx context.Context, address string, transactionType TransactionType) (*TransactionStatusResponse, error) {
	endpoint := "/check-transaction-status"

	// Validate params
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Add("address", address)
	queryParams.Add("type", string(transactionType))

	var response TransactionStatusResponse
	err := c.sendRequestWithQuery(ctx, http.MethodGet, endpoint, queryParams, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// FetchWalletBalance gets the balance of a specific wallet address for a token or native currency
func (c *Client) FetchWalletBalance(ctx context.Context, address string, networkType NetworkType, contractAddress string) (*WalletBalanceResponse, error) {
	endpoint := "/fetch-wallet-balance"

	// Validate params
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Add("address", address)
	queryParams.Add("networkType", string(networkType))
	if contractAddress != "" {
		queryParams.Add("contractAddress", contractAddress)
	}

	var response WalletBalanceResponse
	err := c.sendRequestWithQuery(ctx, http.MethodGet, endpoint, queryParams, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// sendRequest sends an API request and parses the response
func (c *Client) sendRequest(ctx context.Context, method, path string, body interface{}, v interface{}) error {
	return c.sendRequestWithQuery(ctx, method, path, nil, body, v)
}

// sendRequestWithQuery sends an API request with query parameters and parses the response
func (c *Client) sendRequestWithQuery(ctx context.Context, method, path string, query url.Values, body interface{}, v interface{}) error {
	// Create the request URL
	u, err := url.Parse(path)
	if err != nil {
		return err
	}

	u = c.baseURL.ResolveReference(u)

	// Add query parameters if any
	if query != nil {
		u.RawQuery = query.Encode()
	}

	// Create the request body if any
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return err
		}
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("x-api-secret", c.apiSecret)

	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the response
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("HTTP error: %d - failed to parse error response", resp.StatusCode)
		}
		return &apiErr
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}