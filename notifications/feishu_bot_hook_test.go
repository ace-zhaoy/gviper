package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func customPayloadBuilder(configName string, err error) io.Reader {
	return bytes.NewBuffer([]byte(`{"msg_type":"text","content":{"text":"Custom payload"}}`))
}

func TestFeishuBotHook_Notify(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Expected method POST, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		expectedBody := `{"msg_type":"text","content":{"text":"Config test_config reload failed: test error"}}`
		if string(body) != expectedBody {
			t.Fatalf("Expected body %s, got %s", expectedBody, string(body))
		}

		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"code": 0,
			"msg":  "ok",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	hook := NewFeishuBotHook(ts.URL)

	hook.Notify("test_config", fmt.Errorf("test error"))
}

func TestFeishuBotHook_Notify_CustomPayload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		expectedBody := `{"msg_type":"text","content":{"text":"Custom payload"}}`
		if string(body) != expectedBody {
			t.Fatalf("Expected body %s, got %s", expectedBody, string(body))
		}

		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"code": 0,
			"msg":  "ok",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	hook := NewFeishuBotHook(ts.URL)
	hook.SetPayloadBuilder(customPayloadBuilder)

	hook.Notify("test_config", fmt.Errorf("test error"))
}

func TestFeishuBotHook_Notify_Non200Response(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	hook := NewFeishuBotHook(ts.URL)

	hook.Notify("test_config", fmt.Errorf("test error"))
}

func TestFeishuBotHook_Notify_ErrorResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"code": 1,
			"msg":  "error",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	hook := NewFeishuBotHook(ts.URL)

	hook.Notify("test_config", fmt.Errorf("test error"))
}
