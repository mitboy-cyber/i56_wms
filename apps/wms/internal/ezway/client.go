// Package ezway provides EZ Way (台湾关务署) real-name auth API integration.
package ezway

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EZWayClient is the HTTP client for EZ Way Customs API.
type EZWayClient struct {
	BaseURL    string
	APIKey     string
	APISecret  string
	httpClient *http.Client
}

// NewEZWayClient creates an EZ Way API client.
func NewEZWayClient(baseURL, apiKey, apiSecret string) *EZWayClient {
	return &EZWayClient{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		APISecret: apiSecret,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// ===================================================================
// Data types
// ===================================================================

// RealNameVerifyRequest is the request body for real-name verification.
type RealNameVerifyRequest struct {
	DeclarantName string `json:"declarant_name"`   // 申报人姓名
	DeclarantID   string `json:"declarant_id"`     // 身份证号/统一编号
	Phone         string `json:"phone"`            // 手机号
	Email         string `json:"email,omitempty"`  // 邮箱
	VerifyType    string `json:"verify_type"`      // "TW_ID" | "PASSPORT" | "ARC"
}

// RealNameVerifyResponse is the response from real-name verification.
type RealNameVerifyResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	VerifyID    string `json:"verify_id"`    // EZ Way verification ID
	MatchResult string `json:"match_result"` // "MATCH" | "MISMATCH" | "PENDING"
	VerifiedAt  string `json:"verified_at"`
}

// CustomsDeclaration is a customs declaration submission.
type CustomsDeclaration struct {
	DeclarantName string  `json:"declarant_name"`
	DeclarantID   string  `json:"declarant_id"`
	WaybillNo     string  `json:"waybill_no"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	DeclaredValue float64 `json:"declared_value"`
	Weight        float64 `json:"weight"`
	Country       string  `json:"country"`    // origin country
	VerifyID      string  `json:"verify_id"`  // from verify response
}

// DeclarationStatusResponse is the response from declaration status check.
type DeclarationStatusResponse struct {
	Success   bool   `json:"success"`
	Status    string `json:"status"` // "submitted" | "processing" | "cleared" | "rejected"
	Message   string `json:"message"`
	UpdatedAt string `json:"updated_at"`
}

// ===================================================================
// API Methods
// ===================================================================

// sign generates an HMAC-SHA256 signature for the request payload.
func (c *EZWayClient) sign(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(c.APISecret))
	mac.Write(payload)
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// doPost makes an authenticated POST request to the EZ Way API.
func (c *EZWayClient) doPost(ctx context.Context, path string, body, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("ezway: marshal request: %w", err)
	}

	sig := c.sign(jsonBody)
	url := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("ezway: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-Signature", sig)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ezway: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ezway: read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("ezway: API error %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("ezway: unmarshal response: %w", err)
		}
	}
	return nil
}

// VerifyRealName calls the EZ Way real-name verification endpoint.
func (c *EZWayClient) VerifyRealName(ctx context.Context, req *RealNameVerifyRequest) (*RealNameVerifyResponse, error) {
	var resp RealNameVerifyResponse
	if err := c.doPost(ctx, "/realname/verify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubmitDeclaration submits a customs declaration to EZ Way.
func (c *EZWayClient) SubmitDeclaration(ctx context.Context, decl *CustomsDeclaration) (string, error) {
	var resp struct {
		Success       bool   `json:"success"`
		DeclarationID string `json:"declaration_id"`
		Message       string `json:"message"`
	}
	if err := c.doPost(ctx, "/declaration/submit", decl, &resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("ezway: submission failed: %s", resp.Message)
	}
	return resp.DeclarationID, nil
}

// GetDeclarationStatus checks the status of a customs declaration.
func (c *EZWayClient) GetDeclarationStatus(ctx context.Context, declID string) (*DeclarationStatusResponse, error) {
	var resp DeclarationStatusResponse
	if err := c.doPost(ctx, "/declaration/status/"+declID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryReconciliation queries reconciliation records.
func (c *EZWayClient) QueryReconciliation(ctx context.Context, params map[string]string) ([]map[string]interface{}, error) {
	var resp struct {
		Records []map[string]interface{} `json:"records"`
	}
	if err := c.doPost(ctx, "/reconciliation/query", params, &resp); err != nil {
		return nil, err
	}
	return resp.Records, nil
}
