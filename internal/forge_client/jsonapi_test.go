package forge_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractJsonApiSingle(t *testing.T) {
	body := []byte(`{
		"data": {
			"id": "42",
			"type": "servers",
			"attributes": {"name": "test-server", "ip_address": "1.2.3.4"}
		}
	}`)

	id, attrs, err := extractJsonApiSingle(body)
	if err != nil {
		t.Fatalf("extractJsonApiSingle failed: %v", err)
	}
	if id != 42 {
		t.Errorf("expected id 42, got %d", id)
	}

	var result struct {
		Name      string `json:"name"`
		IPAddress string `json:"ip_address"`
	}
	if err := json.Unmarshal(attrs, &result); err != nil {
		t.Fatalf("failed to unmarshal attributes: %v", err)
	}
	if result.Name != "test-server" {
		t.Errorf("expected name 'test-server', got '%s'", result.Name)
	}
	if result.IPAddress != "1.2.3.4" {
		t.Errorf("expected ip '1.2.3.4', got '%s'", result.IPAddress)
	}
}

func TestExtractJsonApiSingle_NonNumericID(t *testing.T) {
	body := []byte(`{
		"data": {
			"id": "abc",
			"type": "servers",
			"attributes": {"name": "test"}
		}
	}`)

	id, _, err := extractJsonApiSingle(body)
	if err != nil {
		t.Fatalf("extractJsonApiSingle failed: %v", err)
	}
	if id != 0 {
		t.Errorf("expected id 0 for non-numeric ID, got %d", id)
	}
}

func TestExtractJsonApiSingle_InvalidJSON(t *testing.T) {
	body := []byte(`not json`)
	_, _, err := extractJsonApiSingle(body)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestExtractJsonApiList(t *testing.T) {
	body := []byte(`{
		"data": [
			{"id": "1", "type": "servers", "attributes": {"name": "server-1"}},
			{"id": "2", "type": "servers", "attributes": {"name": "server-2"}},
			{"id": "3", "type": "servers", "attributes": {"name": "server-3"}}
		]
	}`)

	items, err := extractJsonApiList(body)
	if err != nil {
		t.Fatalf("extractJsonApiList failed: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0].ID != "1" {
		t.Errorf("expected first item ID '1', got '%s'", items[0].ID)
	}
	if items[2].ID != "3" {
		t.Errorf("expected third item ID '3', got '%s'", items[2].ID)
	}
}

func TestExtractJsonApiList_Empty(t *testing.T) {
	body := []byte(`{"data": []}`)

	items, err := extractJsonApiList(body)
	if err != nil {
		t.Fatalf("extractJsonApiList failed: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestUnmarshalList(t *testing.T) {
	items := []jsonApiResource{
		{ID: "10", Type: "servers", Attributes: json.RawMessage(`{"name": "server-a"}`)},
		{ID: "20", Type: "servers", Attributes: json.RawMessage(`{"name": "server-b"}`)},
	}

	type testStruct struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	result, err := unmarshalList(items, func(s *testStruct, id int64) { s.ID = id })
	if err != nil {
		t.Fatalf("unmarshalList failed: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	if result[0].ID != 10 {
		t.Errorf("expected first item ID 10, got %d", result[0].ID)
	}
	if result[0].Name != "server-a" {
		t.Errorf("expected first item name 'server-a', got '%s'", result[0].Name)
	}
	if result[1].ID != 20 {
		t.Errorf("expected second item ID 20, got %d", result[1].ID)
	}
}

func TestUnmarshalList_NilSetID(t *testing.T) {
	items := []jsonApiResource{
		{ID: "10", Type: "servers", Attributes: json.RawMessage(`{"name": "server-a"}`)},
	}

	type testStruct struct {
		Name string `json:"name"`
	}

	result, err := unmarshalList[testStruct](items, nil)
	if err != nil {
		t.Fatalf("unmarshalList failed: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if result[0].Name != "server-a" {
		t.Errorf("expected name 'server-a', got '%s'", result[0].Name)
	}
}

func TestOrgPath(t *testing.T) {
	client := NewClient("test-token")
	client.Organization = "my-org"

	path := client.orgPath("/servers/123")
	expected := "/orgs/my-org/servers/123"
	if path != expected {
		t.Errorf("expected '%s', got '%s'", expected, path)
	}
}

func TestGetJsonApi(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Bearer token, got %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/vnd.api+json")
		fmt.Fprint(w, `{
			"data": {
				"id": "99",
				"type": "servers",
				"attributes": {"name": "test-server", "status": "active"}
			}
		}`)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	var result struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	id, err := client.GetJsonApi(context.Background(), "/servers/99", &result)
	if err != nil {
		t.Fatalf("GetJsonApi failed: %v", err)
	}
	if id != 99 {
		t.Errorf("expected id 99, got %d", id)
	}
	if result.Name != "test-server" {
		t.Errorf("expected name 'test-server', got '%s'", result.Name)
	}
}

func TestPostJsonApi(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"data": {
				"id": "55",
				"type": "sites",
				"attributes": {"name": "new-site"}
			}
		}`)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	reqBody := map[string]string{"name": "new-site"}
	var result struct {
		Name string `json:"name"`
	}
	id, err := client.PostJsonApi(context.Background(), "/sites", reqBody, &result)
	if err != nil {
		t.Fatalf("PostJsonApi failed: %v", err)
	}
	if id != 55 {
		t.Errorf("expected id 55, got %d", id)
	}
	if result.Name != "new-site" {
		t.Errorf("expected name 'new-site', got '%s'", result.Name)
	}
}

func TestGetJsonApiList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		fmt.Fprint(w, `{
			"data": [
				{"id": "1", "type": "jobs", "attributes": {"command": "ls"}},
				{"id": "2", "type": "jobs", "attributes": {"command": "pwd"}}
			]
		}`)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	items, err := client.GetJsonApiList(context.Background(), "/jobs")
	if err != nil {
		t.Fatalf("GetJsonApiList failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestPutJsonApi(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/vnd.api+json")
		fmt.Fprint(w, `{
			"data": {
				"id": "77",
				"type": "servers",
				"attributes": {"name": "updated-server"}
			}
		}`)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	reqBody := map[string]string{"name": "updated-server"}
	var result struct {
		Name string `json:"name"`
	}
	id, err := client.PutJsonApi(context.Background(), "/servers/77", reqBody, &result)
	if err != nil {
		t.Fatalf("PutJsonApi failed: %v", err)
	}
	if id != 77 {
		t.Errorf("expected id 77, got %d", id)
	}
	if result.Name != "updated-server" {
		t.Errorf("expected name 'updated-server', got '%s'", result.Name)
	}
}

func TestGetJsonApi_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message": "internal error"}`)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	var result struct{}
	_, err := client.GetJsonApi(context.Background(), "/servers/1", &result)
	if err == nil {
		t.Error("expected error for 500 response")
	}
}
