// Package client provides the WMS API client for the device gateway.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// InboundTask represents an inbound task retrieved from the WMS.
type InboundTask struct {
	ID             int64   `json:"id"`
	WaybillNo      string  `json:"waybill_no"`
	TrackingNumber string  `json:"tracking_number"`
	SKUCode        string  `json:"sku_code"`
	ProductName    string  `json:"product_name"`
	PlannedQty     int     `json:"planned_qty"`
	DeclaredWeight float64 `json:"declared_weight"`
	Status         int     `json:"status"` // 0:待执行 1:执行中 2:已完成 3:异常
	LocationCode   string  `json:"location_code"`
}

// WeightRecord represents a weight record sent to the WMS.
type WeightRecord struct {
	ID             int64   `json:"id"`
	WaybillNo      string  `json:"waybill_no"`
	ParcelID       int64   `json:"parcel_id,omitempty"`
	GrossWeight    float64 `json:"gross_weight"`
	TareWeight     float64 `json:"tare_weight"`
	NetWeight      float64 `json:"net_weight"`
	DeclaredWeight float64 `json:"declared_weight"`
	WeightDiff     float64 `json:"weight_diff"`
	ScaleID        string  `json:"scale_id"`
	Status         int     `json:"status"`
}

// WMSClient is the HTTP client for the WMS API.
type WMSClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewWMSClient creates a new WMS API client.
func NewWMSClient(baseURL, apiKey string) *WMSClient {
	return &WMSClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetInboundTaskByBarcode queries the WMS for an inbound task by tracking number/barcode.
func (c *WMSClient) GetInboundTaskByBarcode(barcode string) (*InboundTask, error) {
	url := fmt.Sprintf("%s/api/device/inbound-task?barcode=%s", c.baseURL, barcode)
	resp, err := c.doGET(url)
	if err != nil {
		return nil, fmt.Errorf("wms: get inbound task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("wms: no inbound task found for barcode %q", barcode)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wms: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var task InboundTask
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("wms: decode task: %w", err)
	}
	return &task, nil
}

// RecordWeight sends a weight record to the WMS.
func (c *WMSClient) RecordWeight(waybillNo string, weight float64, deviceID string) (*WeightRecord, error) {
	rec := &WeightRecord{
		WaybillNo:   waybillNo,
		GrossWeight: weight,
		ScaleID:     deviceID,
		Status:      0,
	}
	rec.NetWeight = weight // simplified: assume tare is 0

	body, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("wms: marshal weight record: %w", err)
	}

	url := fmt.Sprintf("%s/api/device/weight-record", c.baseURL)
	resp, err := c.doPOST(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("wms: record weight: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wms: record weight failed (status %d): %s", resp.StatusCode, string(errBody))
	}

	var result WeightRecord
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Return the request record even if response parsing fails
		return rec, nil
	}
	return &result, nil
}

// ConfirmInbound confirms an inbound task completion with a location code.
func (c *WMSClient) ConfirmInbound(waybillNo string, locationCode string) error {
	payload := map[string]string{
		"waybill_no":     waybillNo,
		"location_code":  locationCode,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("wms: marshal confirm: %w", err)
	}

	url := fmt.Sprintf("%s/api/device/inbound-confirm", c.baseURL)
	resp, err := c.doPOST(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("wms: confirm inbound: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("wms: confirm inbound failed (status %d): %s", resp.StatusCode, string(errBody))
	}
	return nil
}

// Heartbeat sends a periodic heartbeat to the WMS for a device.
func (c *WMSClient) Heartbeat(deviceID string) error {
	payload := map[string]string{
		"device_id": deviceID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("wms: marshal heartbeat: %w", err)
	}

	url := fmt.Sprintf("%s/api/device/heartbeat", c.baseURL)
	resp, err := c.doPOST(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("wms: heartbeat: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *WMSClient) doGET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	return c.httpClient.Do(req)
}

func (c *WMSClient) doPOST(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", contentType)
	return c.httpClient.Do(req)
}

func (c *WMSClient) setHeaders(req *http.Request) {
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	req.Header.Set("X-Device-Gateway", "i56-device-gateway/1.0")
}
