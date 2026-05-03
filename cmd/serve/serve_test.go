package serve

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormatMac(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{"AA-BB-CC-DD-EE-FF", "AA:BB:CC:DD:EE:FF"},
		{"aa-bb-cc-dd-ee-ff", "aa:bb:cc:dd:ee:ff"},
		{"AA:BB:CC:DD:EE:FF", "AA:BB:CC:DD:EE:FF"}, // already colons
		{"", ""},
	}

	for _, tt := range tests {
		result := formatMac(tt.input)
		if result != tt.expected {
			t.Errorf("formatMac(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParse_MissingDeviceIP(t *testing.T) {
	t.Parallel()
	cmd := NewServeMetricsSubCommand()
	err := cmd.Parse([]string{"--password=secret"})
	if err == nil {
		t.Fatal("expected error for missing device-ip, got nil")
	}
}

func TestParse_MissingPassword(t *testing.T) {
	t.Parallel()
	cmd := NewServeMetricsSubCommand()
	err := cmd.Parse([]string{"--device-ip=192.168.1.254"})
	if err == nil {
		t.Fatal("expected error for missing password, got nil")
	}
}

func TestParse_ValidArgs(t *testing.T) {
	t.Parallel()
	cmd := NewServeMetricsSubCommand()
	err := cmd.Parse([]string{"--device-ip=192.168.1.254", "--password=secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_AllFlags(t *testing.T) {
	t.Parallel()
	cmd := NewServeMetricsSubCommand()
	err := cmd.Parse([]string{
		"--device-ip=192.168.1.254",
		"--password=secret",
		"--username=myadmin",
		"--host=127.0.0.1",
		"--port=9090",
		"--interval=60s",
		"--request-delay=2s",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHealthzEndpoint(t *testing.T) {
	t.Parallel()
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", rec.Body.String())
	}
}
