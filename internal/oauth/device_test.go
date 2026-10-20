package oauth_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"testing"

	"cn.qfei/contract-cli/internal/oauth"
)

func TestStartDeviceAuthorizationSendsFrozenFormAndReturnsCompleteURI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("client_id") != "client-a" || r.Form.Get("scope") != "contract:full" || r.Form.Get("resource") != "resource-a" {
			t.Fatalf("unexpected form: %v", r.Form)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"device_code": "secret-device", "user_code": "user-a",
			"verification_uri":          "https://example.test/device",
			"verification_uri_complete": "https://example.test/device?user_code=user-a",
			"expires_in":                600,
		})
	}))
	defer server.Close()

	result, err := oauth.StartDeviceAuthorization(context.Background(), server.Client(), oauth.DeviceAuthorizationRequest{
		Endpoint: server.URL, ClientID: "client-a", Scope: "contract:full", Resource: "resource-a",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.DeviceCode != "secret-device" || result.VerificationURIComplete != "https://example.test/device?user_code=user-a" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestStartDeviceAuthorizationClassifiesOnlyUnwrittenDialFailuresAsNotSent(t *testing.T) {
	tests := []struct {
		name               string
		markWritten        bool
		wantRequestNotSent bool
	}{
		{name: "tcp dial failed before write", wantRequestNotSent: true},
		{name: "request was already written", markWritten: true, wantRequestNotSent: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				if tt.markWritten {
					trace := httptrace.ContextClientTrace(request.Context())
					if trace != nil && trace.WroteRequest != nil {
						trace.WroteRequest(httptrace.WroteRequestInfo{})
					}
				}
				return nil, &net.OpError{Op: "dial", Net: "tcp", Err: timeoutError{}}
			})}

			_, err := oauth.StartDeviceAuthorization(context.Background(), client, oauth.DeviceAuthorizationRequest{
				Endpoint: "https://auth.example/device", ClientID: "client-a", Scope: "contract:full", Resource: "resource-a",
			})
			if err == nil {
				t.Fatal("error = nil, want transport failure")
			}
			if got := oauth.IsDeviceGrantRequestNotSent(err); got != tt.wantRequestNotSent {
				t.Fatalf("IsDeviceGrantRequestNotSent() = %v, want %v", got, tt.wantRequestNotSent)
			}
		})
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "i/o timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestCompleteDeviceAuthorizationClassifiesPendingWithoutPolling(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "authorization_pending"})
	}))
	defer server.Close()

	_, err := oauth.CompleteDeviceAuthorization(context.Background(), server.Client(), oauth.DeviceTokenRequest{
		Endpoint: server.URL, ClientID: "client-a", DeviceCode: "device-a",
	})
	if !oauth.IsDeviceGrantError(err, "authorization_pending") {
		t.Fatalf("error = %v", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestRefreshDeviceTokenRotatesToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "refresh_token" || r.Form.Get("refresh_token") != "refresh-old" {
			t.Fatalf("unexpected form: %v", r.Form)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "access-new", "refresh_token": "refresh-new", "token_type": "Bearer", "scope": "contract:full", "expires_in": 3600,
		})
	}))
	defer server.Close()

	token, err := oauth.RefreshDeviceToken(context.Background(), server.Client(), oauth.DeviceRefreshRequest{
		Endpoint: server.URL, ClientID: "client-a", RefreshToken: "refresh-old",
	})
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "access-new" || token.RefreshToken != "refresh-new" {
		t.Fatalf("unexpected token: %+v", token)
	}
}

func TestRevokeDeviceTokenSendsRefreshTokenOnce(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("client_id") != "client-a" || r.Form.Get("token") != "refresh-a" || r.Form.Get("token_type_hint") != "refresh_token" {
			t.Fatalf("unexpected form: %v", r.Form)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	if err := oauth.RevokeDeviceToken(context.Background(), server.Client(), oauth.DeviceRevokeRequest{
		Endpoint: server.URL, ClientID: "client-a", RefreshToken: "refresh-a",
	}); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}
