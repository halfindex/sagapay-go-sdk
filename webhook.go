package sagapay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// WebhookHandler handles SagaPay webhook notifications
type WebhookHandler struct {
	apiSecret string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(apiSecret string) *WebhookHandler {
	return &WebhookHandler{
		apiSecret: apiSecret,
	}
}

// HandleRequest processes a webhook notification from an HTTP request
func (h *WebhookHandler) HandleRequest(r *http.Request) (*WebhookPayload, error) {
	// Get the signature from the headers
	signature := r.Header.Get("x-sagapay-signature")
	if signature == "" {
		return nil, errors.New("missing SagaPay signature in headers")
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	// Verify the signature
	if !h.VerifySignature(body, signature) {
		return nil, errors.New("invalid webhook signature")
	}

	// Parse the webhook payload
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	return &payload, nil
}

// ProcessWebhook processes a webhook notification from raw body and signature
func (h *WebhookHandler) ProcessWebhook(body []byte, signature string) (*WebhookPayload, error) {
	// Verify the signature
	if !h.VerifySignature(body, signature) {
		return nil, errors.New("invalid webhook signature")
	}

	// Parse the webhook payload
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	return &payload, nil
}

// VerifySignature verifies the HMAC signature of a webhook payload
func (h *WebhookHandler) VerifySignature(payload []byte, signature string) bool {
	// Calculate the HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(h.apiSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare with the provided signature
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// SendSuccessResponse sends a success response for a webhook
func SendSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"received": true})
}

// SendErrorResponse sends an error response for a webhook
func SendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	// Still return 200 OK to prevent retries
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"received": false,
		"error":    err.Error(),
	})
}