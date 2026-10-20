package cli_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"cn.qfei/contract-cli/internal/cli"
	"cn.qfei/contract-cli/internal/config"
)

func TestPaymentCommandsUseExpectedEndpoints(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityApp), true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}

	testCases := []struct {
		name         string
		args         []string
		wantMethod   string
		wantPath     string
		wantQuery    map[string]string
		wantBody     string
		responseBody string
	}{
		{
			name:         "create payment",
			args:         []string{"payment", "create", "--contract", "contract-1", "--profile", "contract", "--data", `{"apply_amount":100}`, "--user-id", "ou_123"},
			wantMethod:   http.MethodPost,
			wantPath:     "/open-apis/contract/v1/contracts/contract-1/payments",
			wantQuery:    map[string]string{"user_id_type": "user_id", "user_id": "ou_123"},
			wantBody:     `{"apply_amount":100}`,
			responseBody: `{"code":0,"data":{"payment":{"payment_id":"payment-1"}}}`,
		},
		{
			name:         "update payment",
			args:         []string{"payment", "update", "payment-1", "--contract", "contract-1", "--profile", "contract", "--data", `{"payment_status_code":9}`},
			wantMethod:   http.MethodPatch,
			wantPath:     "/open-apis/contract/v1/contracts/contract-1/payments/payment-1",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			wantBody:     `{"payment_status_code":9}`,
			responseBody: `{"code":0,"data":{"payment_id":"payment-1"}}`,
		},
		{
			name:         "get payment",
			args:         []string{"payment", "get", "payment-1", "--contract", "contract-1", "--profile", "contract"},
			wantMethod:   http.MethodGet,
			wantPath:     "/open-apis/contract/v1/contracts/contract-1/payments/payment-1",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			responseBody: `{"code":0,"data":{"payment":{"payment_id":"payment-1"}}}`,
		},
		{
			name:         "list payments",
			args:         []string{"payment", "list", "--contract", "contract-1", "--profile", "contract", "--page-size", "10", "--page-token", "next-1"},
			wantMethod:   http.MethodGet,
			wantPath:     "/open-apis/contract/v1/contracts/contract-1/payments",
			wantQuery:    map[string]string{"user_id_type": "user_id", "page_size": "10", "page_token": "next-1"},
			responseBody: `{"code":0,"data":{"items":[{"payment_id":"payment-1"}]}}`,
		},
		{
			name:         "notify payment plan",
			args:         []string{"payment", "plan", "notify", "--profile", "contract", "--data", `{"finance_number":"FN-1"}`},
			wantMethod:   http.MethodPost,
			wantPath:     "/open-apis/contract/v1/payment/notify",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			wantBody:     `{"finance_number":"FN-1"}`,
			responseBody: `{"code":0,"data":{"payment_record_resource_vo":{}}}`,
		},
		{
			name:         "search payment plans",
			args:         []string{"payment", "plan", "search", "--profile", "contract", "--data", `{"page_size":10}`},
			wantMethod:   http.MethodPost,
			wantPath:     "/open-apis/contract/v1/payments/search",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			wantBody:     `{"page_size":10}`,
			responseBody: `{"code":0,"data":{"items":[{"payment_plan_uuid":"plan-1"}]}}`,
		},
		{
			name:         "create payment record",
			args:         []string{"payment", "record", "create", "--contract", "contract-1", "--payment", "payment-1", "--profile", "contract", "--data", `{"transaction_amount":100}`},
			wantMethod:   http.MethodPost,
			wantPath:     "/open-apis/contract/v1/contracts/contract-1/payments/payment-1/payment_records",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			wantBody:     `{"transaction_amount":100}`,
			responseBody: `{"code":0,"data":{"payment_record":{"payment_record_id":"record-1"}}}`,
		},
		{
			name:         "update payment record",
			args:         []string{"payment", "record", "update", "record-1", "--contract", "contract-1", "--payment", "payment-1", "--profile", "contract", "--data", `{"transaction_amount":90}`},
			wantMethod:   http.MethodPatch,
			wantPath:     "/open-apis/contract/v1/contracts/contract-1/payments/payment-1/payment_records/record-1",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			wantBody:     `{"transaction_amount":90}`,
			responseBody: `{"code":0,"data":{"payment_record":{"payment_record_id":"record-1"}}}`,
		},
		{
			name:         "get payment record",
			args:         []string{"payment", "record", "get", "record-1", "--contract", "contract-1", "--payment", "payment-1", "--profile", "contract"},
			wantMethod:   http.MethodGet,
			wantPath:     "/open-apis/contract/v1/contracts/contract-1/payments/payment-1/payment_records/record-1",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			responseBody: `{"code":0,"data":{"payment_record":{"payment_record_id":"record-1"}}}`,
		},
		{
			name:         "list payment records by plan",
			args:         []string{"payment", "record", "list", "--plan", "plan-1", "--profile", "contract"},
			wantMethod:   http.MethodGet,
			wantPath:     "/open-apis/contract/v1/contracts/payments/plan-1/payment_records",
			wantQuery:    map[string]string{"user_id_type": "user_id"},
			responseBody: `{"code":0,"data":{"items":[{"payment_record_id":"record-1"}]}}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stdout := &bytes.Buffer{}
			app := cli.New(cli.Options{
				Stdout: stdout,
				Stderr: &bytes.Buffer{},
				Store:  store,
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						if req.Method != tc.wantMethod {
							t.Fatalf("method = %s, want %s", req.Method, tc.wantMethod)
						}
						if req.URL.Path != tc.wantPath {
							t.Fatalf("path = %s, want %s", req.URL.Path, tc.wantPath)
						}
						for key, want := range tc.wantQuery {
							if got := req.URL.Query().Get(key); got != want {
								t.Fatalf("query %s = %q, want %q", key, got, want)
							}
						}
						if req.Header.Get("Authorization") != "Bearer app-token" {
							t.Fatalf("authorization = %q", req.Header.Get("Authorization"))
						}
						body, err := io.ReadAll(req.Body)
						if err != nil {
							t.Fatalf("ReadAll() error = %v", err)
						}
						if string(body) != tc.wantBody {
							t.Fatalf("body = %q, want %q", string(body), tc.wantBody)
						}
						return jsonResponse(tc.responseBody), nil
					}),
				},
			})

			if err := app.Run(context.Background(), tc.args); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if !strings.Contains(stdout.String(), `"code": 0`) {
				t.Fatalf("unexpected output: %s", stdout.String())
			}
		})
	}
}

func TestPaymentCommandsRejectUserIdentityBeforeHTTP(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityUser), true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}

	testCases := [][]string{
		{"payment", "create", "--contract", "contract-1", "--profile", "contract", "--as", "user", "--data", `{"apply_amount":100}`},
		{"payment", "update", "payment-1", "--contract", "contract-1", "--profile", "contract", "--as", "user", "--data", `{"payment_status_code":9}`},
		{"payment", "get", "payment-1", "--contract", "contract-1", "--profile", "contract", "--as", "user"},
		{"payment", "list", "--contract", "contract-1", "--profile", "contract", "--as", "user"},
		{"payment", "plan", "notify", "--profile", "contract", "--as", "user", "--data", `{"finance_number":"FN-1"}`},
		{"payment", "plan", "search", "--profile", "contract", "--as", "user", "--data", `{"page_size":10}`},
		{"payment", "record", "create", "--contract", "contract-1", "--payment", "payment-1", "--profile", "contract", "--as", "user", "--data", `{"transaction_amount":100}`},
		{"payment", "record", "update", "record-1", "--contract", "contract-1", "--payment", "payment-1", "--profile", "contract", "--as", "user", "--data", `{"transaction_amount":90}`},
		{"payment", "record", "get", "record-1", "--contract", "contract-1", "--payment", "payment-1", "--profile", "contract", "--as", "user"},
		{"payment", "record", "list", "--plan", "plan-1", "--profile", "contract", "--as", "user"},
	}

	for _, args := range testCases {
		args := args
		t.Run(strings.Join(args[:min(4, len(args))], " "), func(t *testing.T) {
			t.Parallel()

			requests := 0
			app := cli.New(cli.Options{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
				Store:  store,
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						requests++
						return jsonResponse(`{"code":0}`), nil
					}),
				},
			})

			err := app.Run(context.Background(), args)
			if err == nil || !strings.Contains(err.Error(), "only supports --as app") {
				t.Fatalf("unexpected user error: %v", err)
			}
			if requests != 0 {
				t.Fatalf("user rejection should not send HTTP, got %d requests", requests)
			}
		})
	}
}

func TestPaymentCommandValidationErrors(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityApp), true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}

	testCases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "create missing contract",
			args:    []string{"payment", "create", "--profile", "contract", "--data", `{"apply_amount":100}`},
			wantErr: "--contract is required",
		},
		{
			name:    "create missing body",
			args:    []string{"payment", "create", "--contract", "contract-1", "--profile", "contract"},
			wantErr: "--input-file or --data is required",
		},
		{
			name:    "update missing id",
			args:    []string{"payment", "update", "--contract", "contract-1", "--profile", "contract", "--data", `{"payment_status_code":9}`},
			wantErr: "usage: contract-cli payment update <payment-id> --contract <contract-id> --input-file <path>|--data <json> [flags]",
		},
		{
			name:    "update missing contract",
			args:    []string{"payment", "update", "payment-1", "--profile", "contract", "--data", `{"payment_status_code":9}`},
			wantErr: "--contract is required",
		},
		{
			name:    "list rejects body",
			args:    []string{"payment", "list", "--contract", "contract-1", "--profile", "contract", "--data", `{"ignored":true}`},
			wantErr: "payment list does not accept --input-file or --data",
		},
		{
			name:    "plan missing subcommand",
			args:    []string{"payment", "plan"},
			wantErr: "missing payment plan subcommand",
		},
		{
			name:    "record create missing payment",
			args:    []string{"payment", "record", "create", "--contract", "contract-1", "--profile", "contract", "--data", `{"transaction_amount":100}`},
			wantErr: "--payment is required",
		},
		{
			name:    "record list missing plan",
			args:    []string{"payment", "record", "list", "--profile", "contract"},
			wantErr: "--plan is required",
		},
		{
			name:    "record get missing id",
			args:    []string{"payment", "record", "get", "--contract", "contract-1", "--payment", "payment-1", "--profile", "contract"},
			wantErr: "usage: contract-cli payment record get <payment-record-id> --contract <contract-id> --payment <payment-id> [flags]",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := cli.New(cli.Options{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
				Store:  store,
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						t.Fatalf("validation error should not send HTTP")
						return nil, nil
					}),
				},
			})

			err := app.Run(context.Background(), tc.args)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("unexpected error: %v, want %q", err, tc.wantErr)
			}
		})
	}
}
