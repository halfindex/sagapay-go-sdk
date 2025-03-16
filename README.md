# SagaPay Go SDK

Go SDK for [SagaPay](https://sagapay.net) - the world's first free, non-custodial blockchain payment gateway service provider. This SDK enables Go developers to seamlessly integrate cryptocurrency payments without holding customer funds. With enterprise-grade security and zero transaction fees, SagaPay empowers merchants to accept crypto payments across multiple blockchains while maintaining full control of their digital assets.

## Installation

```bash
go get github.com/sagapay/sagapay-go-sdk
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sagapay/sagapay-go-sdk"
)

func main() {
	// Initialize the SagaPay client
	client, err := sagapay.NewClient(sagapay.Config{
		APIKey:    "your-api-key",
		APISecret: "your-api-secret",
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a deposit address
	depositResponse, err := client.CreateDeposit(context.Background(), sagapay.CreateDepositParams{
		NetworkType:     sagapay.NetworkTypeBEP20,
		ContractAddress: "0", // Use '0' for native tokens (BNB)
		Amount:          "1.5",
		IPNUrl:          "https://yourwebsite.com/webhook",
		UDF:             "order-123",
		Type:            sagapay.AddressTypeTemporary,
	})
	if err != nil {
		log.Fatalf("Failed to create deposit: %v", err)
	}

	fmt.Printf("Deposit address created: %s\n", depositResponse.Address)
}
```

## Features

- Deposit address generation
- Withdrawal processing
- Transaction status checking
- Wallet balance fetching
- Multi-chain support (ERC20, BEP20, TRC20, POLYGON, SOLANA)
- Webhook notifications (IPN)
- Custom UDF field support
- Non-custodial architecture
- Context support for proper cancellation handling

## API Reference

### Create Deposit

```go
depositResponse, err := client.CreateDeposit(ctx, sagapay.CreateDepositParams{
    NetworkType:     sagapay.NetworkTypeBEP20,     // Required: Blockchain network type
    ContractAddress: "0",                          // Required: Contract address or '0' for native coins
    Amount:          "1.5",                        // Required: Expected deposit amount
    IPNUrl:          "https://example.com/webhook", // Required: URL for notifications
    UDF:             "order-123",                  // Optional: User-defined field
    Type:            sagapay.AddressTypeTemporary, // Optional: TEMPORARY or PERMANENT
})
```

### Create Withdrawal

```go
withdrawalResponse, err := client.CreateWithdrawal(ctx, sagapay.CreateWithdrawalParams{
    NetworkType:     sagapay.NetworkTypeERC20,
    ContractAddress: "0xdAC17F958D2ee523a2206206994597C13D831ec7", // USDT on Ethereum
    Address:         "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
    Amount:          "10.5",
    IPNUrl:          "https://example.com/webhook",
    UDF:             "withdrawal-456",
})
```

### Check Transaction Status

```go
statusResponse, err := client.CheckTransactionStatus(
    ctx,
    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", // Address
    sagapay.TransactionTypeDeposit,                // Transaction type
)
```

### Fetch Wallet Balance

```go
balanceResponse, err := client.FetchWalletBalance(
    ctx,
    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", // Address
    sagapay.NetworkTypeERC20,                      // Network type
    "0xdAC17F958D2ee523a2206206994597C13D831ec7", // Contract address
)
```

## Handling Webhooks (IPN)

SagaPay sends webhook notifications to your specified `ipnUrl` when transaction statuses change. Use the `WebhookHandler` to process these notifications:

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/sagapay/sagapay-go-sdk"
)

func main() {
    // Create a webhook handler
    webhookHandler := sagapay.NewWebhookHandler("your-api-secret")

    // Set up a handler for webhook notifications
    http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
        // Process the webhook
        payload, err := webhookHandler.HandleRequest(r)
        if err != nil {
            log.Printf("Error processing webhook: %v", err)
            sagapay.SendErrorResponse(w, err)
            return
        }

        // Handle the webhook data
        switch payload.Status {
        case sagapay.TransactionStatusCompleted:
            // Payment successful, update your database
            log.Printf("Transaction %s completed", payload.ID)
        case sagapay.TransactionStatusFailed:
            // Handle failed transaction
            log.Printf("Transaction %s failed", payload.ID)
        }

        // Send a success response
        sagapay.SendSuccessResponse(w)
    })

    // Start the server
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Webhook Payload Format

When SagaPay sends a webhook to your endpoint, it will include the following payload:

```json
{
  "id": "transaction-uuid",
  "type": "deposit|withdrawal",
  "status": "PENDING|PROCESSING|COMPLETED|FAILED|CANCELLED",
  "address": "0x123abc...",
  "networkType": "ERC20|BEP20|TRC20|POLYGON|SOLANA",
  "amount": "10.5",
  "udf": "your-optional-user-defined-field",
  "txHash": "0xabc123...",
  "timestamp": "2025-03-16T14:30:00Z"
}
```

## Error Handling

The SDK includes comprehensive error handling:

```go
depositResponse, err := client.CreateDeposit(ctx, params)
if err != nil {
    // Check if it's an API error
    if apiErr, ok := err.(*sagapay.APIError); ok {
        fmt.Printf("API Error: %s - %s\n", apiErr.Error, apiErr.Message)
        return
    }
    
    // Handle other errors
    fmt.Printf("Error: %v\n", err)
    return
}
```

## License

This SDK is released under the MIT License.

## Support

For questions or support, please contact support@sagapay.net or visit [https://sagapay.net](https://sagapay.net).