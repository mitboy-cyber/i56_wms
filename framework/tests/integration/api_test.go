package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APIResponse struct {
	Data  interface{} `json:"data"`
	Meta  *APIMeta    `json:"meta,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

type APIMeta struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestHealthEndpoint(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()

	resp := doGet(t, ts.URL+"/api/health")
	assertStatus(t, resp, 200)

	var r APIResponse
	decode(t, resp, &r)
	data := r.Data.(map[string]interface{})
	if data["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", data["status"])
	}
}

func TestLoginAPI(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()

	resp := doPost(t, ts.URL+"/api/v1/auth/login", map[string]string{
		"username": "admin", "password": "admin",
	})
	assertStatus(t, resp, 200)

	var r APIResponse
	decode(t, resp, &r)
	data := r.Data.(map[string]interface{})
	if data["access_token"] == "" {
		t.Error("expected non-empty access_token")
	}
}

func TestAuthRequired(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()

	resp := doGet(t, ts.URL+"/api/v1/me")
	if resp.StatusCode < 400 {
		t.Errorf("expected error without token, got %d", resp.StatusCode)
	}
}

func TestClients(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/clients", token)
	assertStatus(t, resp, 200)
	var r APIResponse
	decode(t, resp, &r)
	if len(r.Data.([]interface{})) == 0 {
		t.Error("expected at least 1 client")
	}
}

func TestParcels(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/parcels", token)
	assertStatus(t, resp, 200)
	var r APIResponse
	decode(t, resp, &r)
	if len(r.Data.([]interface{})) != 10 {
		t.Errorf("expected 10 parcels, got %d", len(r.Data.([]interface{})))
	}
}

func TestOrders(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/orders", token)
	assertStatus(t, resp, 200)
	var r APIResponse
	decode(t, resp, &r)
	if len(r.Data.([]interface{})) != 3 {
		t.Errorf("expected 3 orders, got %d", len(r.Data.([]interface{})))
	}
}

func TestServiceTypes(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/services/types", token)
	assertStatus(t, resp, 200)
	var r APIResponse
	decode(t, resp, &r)
	if len(r.Data.([]interface{})) != 13 {
		t.Errorf("expected 13 service types, got %d", len(r.Data.([]interface{})))
	}
}

func TestWorkOrders(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/workorders", token)
	assertStatus(t, resp, 200)
	var r APIResponse
	decode(t, resp, &r)
	if len(r.Data.([]interface{})) != 2 {
		t.Errorf("expected 2 workorders, got %d", len(r.Data.([]interface{})))
	}
}

func TestReports(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	for _, ep := range []string{"/api/v1/reports/orders", "/api/v1/reports/clients", "/api/v1/reports/routes"} {
		resp := doGetWithAuth(t, ts.URL+ep, token)
		assertStatus(t, resp, 200)
	}
}

func TestWebhooksCRUD(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/webhooks", token)
	assertStatus(t, resp, 200)

	resp = doPostWithAuth(t, ts.URL+"/api/v1/webhooks", token, map[string]interface{}{
		"url": "https://example.com/hook", "event": "order.created", "secret": "whsec_test",
	})
	assertStatus(t, resp, 201)

	resp = doGetWithAuth(t, ts.URL+"/api/v1/webhooks/logs", token)
	assertStatus(t, resp, 200)
}

func TestDashboard(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/dashboard", token)
	assertStatus(t, resp, 200)
}

func TestPrints(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/prints", token)
	assertStatus(t, resp, 200)
	var r APIResponse
	decode(t, resp, &r)
	if len(r.Data.([]interface{})) != 3 {
		t.Errorf("expected 3 print templates, got %d", len(r.Data.([]interface{})))
	}
}

func TestCouriers(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
	defer ts.Close()
	token := loginAndGetToken(t, ts.URL)

	resp := doGetWithAuth(t, ts.URL+"/api/v1/couriers", token)
	assertStatus(t, resp, 200)
	var r APIResponse
	decode(t, resp, &r)
	if len(r.Data.([]interface{})) != 9 {
		t.Errorf("expected 9 couriers, got %d", len(r.Data.([]interface{})))
	}
}

// Helpers

func loginAndGetToken(t *testing.T, baseURL string) string {
	t.Helper()
	resp := doPost(t, baseURL+"/api/v1/auth/login", map[string]string{"username": "admin", "password": "admin"})
	var r APIResponse
	decode(t, resp, &r)
	return r.Data.(map[string]interface{})["access_token"].(string)
}

func doGet(t *testing.T, url string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("GET %s: %v", url, err) }
	return resp
}

func doGetWithAuth(t *testing.T, url, token string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("GET %s: %v", url, err) }
	return resp
}

func doPost(t *testing.T, url string, body interface{}) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("POST %s: %v", url, err) }
	return resp
}

func doPostWithAuth(t *testing.T, url, token string, body interface{}) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("POST %s: %v", url, err) }
	return resp
}

func assertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected %d, got %d: %s", expected, resp.StatusCode, string(body))
	}
}

func decode(t *testing.T, resp *http.Response, v interface{}) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode: %v", err)
	}
}
