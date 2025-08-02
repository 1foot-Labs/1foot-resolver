package relayer_communication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload structure
type FulfillOrderRequest struct {
	OrderID        string `json:"orderId"`
	TakerAddress   string `json:"takerAddress"`
	EthHTLCAddress string `json:"ethHTLCAddress"`
	BtcHTLCAddress string `json:"btcHTLCAddress"`
}

// FulfillOrder sends a POST request to the fixed API endpoint
func FulfillOrder(reqData FulfillOrderRequest) (*http.Response, error) {
	const url = "http://localhost:3002/api/fulfill-order"

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}