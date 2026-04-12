package paymentclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: "http://localhost:8081",
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *Client) Pay(orderID string, amount int64) (string, error) {
	body := map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	}

	jsonData, _ := json.Marshal(body)

	resp, err := c.client.Post(c.baseURL+"/payments", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	status, _ := result["status"].(string)

	return status, nil
}
