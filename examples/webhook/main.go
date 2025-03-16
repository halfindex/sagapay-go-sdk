package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/halfindex/sagapay-go-sdk"
)

func main() {
	// Create a webhook handler
	webhookHandler := sagapay.NewWebhookHandler("your-api-secret")

	// Set up a handler for webhook notifications
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Process the webhook
		payload, err := webhookHandler.HandleRequest(r)
		if err != nil {
			log.Printf("Error processing webhook: %v", err)
			sagapay.SendErrorResponse(w, err)
			return
		}

		// Log the webhook
		log.Printf("Received webhook: ID=%s, Type=%s, Status=%s", payload.ID, payload.Type, payload.Status)

		// Handle different transaction statuses
		switch payload.Status {
		case sagapay.TransactionStatusCompleted:
			// Handle completed transaction
			log.Printf("Transaction %s completed: %s %s", payload.ID, payload.Amount, payload.Type)
			
			// Your business logic here...
			// For example, update order status in your database
			// if payload.Type == sagapay.TransactionTypeDeposit {
			//     updateOrderStatus(payload.UDF, "paid")
			// } else {
			//     updateWithdrawalStatus(payload.UDF, "completed")
			// }

		case sagapay.TransactionStatusFailed:
			// Handle failed transaction
			log.Printf("Transaction %s failed: %s %s", payload.ID, payload.Amount, payload.Type)
			
			// Your business logic here...
			// updateTransactionStatus(payload.UDF, "failed")

		case sagapay.TransactionStatusProcessing, sagapay.TransactionStatusPending:
			// Handle pending/processing transaction
			log.Printf("Transaction %s is %s: %s %s", payload.ID, payload.Status, payload.Amount, payload.Type)
			
			// Your business logic here...
			// updateTransactionStatus(payload.UDF, string(payload.Status))

		case sagapay.TransactionStatusCancelled:
			// Handle cancelled transaction
			log.Printf("Transaction %s cancelled: %s %s", payload.ID, payload.Amount, payload.Type)
			
			// Your business logic here...
			// updateTransactionStatus(payload.UDF, "cancelled")
		}

		// Send a success response
		sagapay.SendSuccessResponse(w)
	})

	// Start the server
	fmt.Println("Starting webhook server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}