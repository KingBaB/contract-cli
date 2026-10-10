package cli_test

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"cn.qfei/contract-cli/internal/cli"
	"cn.qfei/contract-cli/internal/config"
)

func TestContractSummaryQueryUsesUserIdentityAndContractNumber(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityUser), true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}

	app := cli.New(cli.Options{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Store:  store,
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodPost || req.URL.Path != "/open-apis/contract/v1/summary/query" {
					t.Fatalf("request = %s %s", req.Method, req.URL.Path)
				}
				if got := req.Header.Get("Authorization"); got != "Bearer user-token" {
					t.Fatalf("authorization = %q", got)
				}
				if got := req.URL.Query().Get("user_id_type"); got != "user_id" {
					t.Fatalf("user_id_type = %q", got)
				}
				assertJSONNumberRequestBody(t, req, `{"contractNumber":"CN-001"}`)
				return jsonResponse(`{"code":0,"data":{"contractNumber":"CN-001","status":2,"overview":"摘要"},"msg":"success"}`), nil
			}),
		},
	})

	err := app.Run(context.Background(), []string{
		"contract", "summary", "query",
		"--profile", "contract",
		"--as", "user",
		"--contract-number", "CN-001",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestContractSummaryQuerySupportsExactContractName(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityUser), true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}

	app := cli.New(cli.Options{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Store:  store,
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				assertJSONNumberRequestBody(t, req, `{"contractName":"采购合同"}`)
				return jsonResponse(`{"code":0,"data":{"contractName":"采购合同","status":2},"msg":"success"}`), nil
			}),
		},
	})

	err := app.Run(context.Background(), []string{
		"contract", "summary", "query",
		"--profile", "contract",
		"--as", "user",
		"--contract-name", "采购合同",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestContractSummaryQueryRejectsAmbiguousSelectorsBeforeHTTP(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityUser), true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}
	app := cli.New(cli.Options{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Store:  store,
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
				return nil, nil
			}),
		},
	})

	err := app.Run(context.Background(), []string{
		"contract", "summary", "query",
		"--profile", "contract",
		"--as", "user",
		"--contract-number", "CN-001",
		"--contract-name", "采购合同",
	})
	if err == nil || !strings.Contains(err.Error(), "exactly one of --contract-number or --contract-name") {
		t.Fatalf("Run() error = %v, want selector validation", err)
	}
}

func TestContractSummaryQueryRejectsAppIdentityBeforeHTTP(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityApp), true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}
	app := cli.New(cli.Options{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		Store:  store,
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
				return nil, nil
			}),
		},
	})

	err := app.Run(context.Background(), []string{
		"contract", "summary", "query",
		"--profile", "contract",
		"--as", "app",
		"--contract-number", "CN-001",
	})
	if err == nil || !strings.Contains(err.Error(), "only supports --as user") {
		t.Fatalf("Run() error = %v, want user-only identity rejection", err)
	}
}
