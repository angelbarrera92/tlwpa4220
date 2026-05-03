package tlwpa4220_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/angelbarrera92/tlwpa4220/pkg/tlwpa4220"
)

// newTestClient returns a Client pointed at the given test server URL.
func newTestClient(serverURL string) tlwpa4220.Client {
	// Strip "http://" — Client prepends it internally.
	host := strings.TrimPrefix(serverURL, "http://")
	return tlwpa4220.Client{
		Username: "admin",
		Password: "testpass",
		IP:       host,
	}
}

func TestWirelessStatistics_Success(t *testing.T) {
	t.Parallel()
	response := tlwpa4220.WirelessStatistics{
		Success: true,
		Timeout: false,
		Data: []tlwpa4220.WirelessStatisticsData{
			{
				Mac:        "AA-BB-CC-DD-EE-FF",
				Type:       "2.4GHz",
				Encryption: "wpa2-psk",
				Rxpkts:     "0",
				Txpkts:     "0",
				IP:         "192.168.1.100",
				DevName:    "mylaptop",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.String(), "wireless") {
			t.Errorf("unexpected path: %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	stats, err := client.WirelessStatistics()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stats.Success {
		t.Error("expected success=true")
	}
	if len(stats.Data) != 1 {
		t.Fatalf("expected 1 device, got %d", len(stats.Data))
	}
	if stats.Data[0].Mac != "AA-BB-CC-DD-EE-FF" {
		t.Errorf("unexpected MAC: %s", stats.Data[0].Mac)
	}
	if stats.Data[0].Type != "2.4GHz" {
		t.Errorf("unexpected type: %s", stats.Data[0].Type)
	}
	if stats.Data[0].Encryption != "wpa2-psk" {
		t.Errorf("unexpected encryption: %s", stats.Data[0].Encryption)
	}
}

func TestWirelessStatistics_NonOKStatus(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.WirelessStatistics()
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestWirelessStatistics_InvalidJSON(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.WirelessStatistics()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestWirelessStatistics_ConnectionRefused(t *testing.T) {
	t.Parallel()
	client := tlwpa4220.Client{
		Username: "admin",
		Password: "testpass",
		IP:       "127.0.0.1:19999", // nothing listening here
	}
	_, err := client.WirelessStatistics()
	if err == nil {
		t.Fatal("expected error for connection refused, got nil")
	}
}

func TestPowerLineStatistics_Success(t *testing.T) {
	t.Parallel()
	response := tlwpa4220.PowerLineStatistics{
		Success: true,
		Timeout: false,
		Data: []tlwpa4220.PowerLineStatisticsData{
			{
				DeviceMac: "11-22-33-44-55-66",
				RxRate:    "254",
				TxRate:    "139",
				Status:    "on",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	stats, err := client.PowerLineStatistics()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stats.Success {
		t.Error("expected success=true")
	}
	if len(stats.Data) != 1 {
		t.Fatalf("expected 1 device, got %d", len(stats.Data))
	}
	if stats.Data[0].DeviceMac != "11-22-33-44-55-66" {
		t.Errorf("unexpected MAC: %s", stats.Data[0].DeviceMac)
	}
	if stats.Data[0].TxRate != "139" {
		t.Errorf("unexpected tx_rate: %s", stats.Data[0].TxRate)
	}
	if stats.Data[0].RxRate != "254" {
		t.Errorf("unexpected rx_rate: %s", stats.Data[0].RxRate)
	}
	if stats.Data[0].Status != "on" {
		t.Errorf("unexpected status: %s", stats.Data[0].Status)
	}
}

func TestPowerLineStatistics_NonOKStatus(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.PowerLineStatistics()
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestReboot_Success(t *testing.T) {
	t.Parallel()
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	err := client.Reboot()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected server to be called")
	}
}

func TestReboot_NonOKStatus(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	err := client.Reboot()
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestAuthCookieFormat(t *testing.T) {
	t.Parallel()
	var receivedCookie string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCookie = r.Header.Get("Cookie")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_ = client.Reboot()

	if !strings.HasPrefix(receivedCookie, "Authorization=Basic admin:") {
		t.Errorf("unexpected cookie format: %s", receivedCookie)
	}
	// MD5 of "testpass" is 179ad45c6ce2cb97cf1029e212046e81
	if !strings.Contains(receivedCookie, "179ad45c6ce2cb97cf1029e212046e81") {
		t.Errorf("unexpected MD5 in cookie: %s", receivedCookie)
	}
}

func TestRequestHeaders(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "Mozilla/5.0" {
			t.Errorf("unexpected User-Agent: %s", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("unexpected Content-Type: %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Referer") == "" {
			t.Error("expected Referer header to be set")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_ = client.Reboot()
}
