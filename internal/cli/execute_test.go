package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cn.qfei/contract-cli/internal/cli"
	"cn.qfei/contract-cli/internal/config"
	"cn.qfei/contract-cli/internal/invocation"
	updatecheck "cn.qfei/contract-cli/internal/update"
)

func TestExecuteCancelledJSONIncludesTraceAndCachedNoticeWithoutExtraNetwork(t *testing.T) {
	t.Parallel()
	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(uploadProfile(config.IdentityApp), true); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	cache := updatecheck.Cache{CheckedAt: now, Channel: "latest", CurrentVersion: "1.0.0", LatestVersion: "1.0.1", UpdateAvailable: true}
	if err := updatecheck.SaveCache(filepath.Join(filepath.Dir(store.Path()), "update-check.json"), cache); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	calls, trace := 0, ""
	app := cli.New(cli.Options{Store: store, Stdout: stdout, Stderr: stderr,
		UpdateCurrentVersion: "1.0.0", Now: func() time.Time { return now },
		LookupEnv:          func(string) (string, bool) { return "", false },
		InspectEnvironment: func(context.Context, int) invocation.Result { return invocation.Result{} },
		HTTPClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if strings.Contains(req.URL.Hostname(), "registry") {
				t.Fatal("error output started a new update request")
			}
			trace = req.Header.Get("X-Log-Id")
			cancel()
			return nil, req.Context().Err()
		})},
	})
	args := []string{"contract", "get", "1", "--profile", "contract", "--as", "app", "--output", "json"}
	if code := app.Execute(ctx, args); code != 130 || calls != 1 || stdout.Len() != 0 {
		t.Fatalf("cancel code/calls/output: code=%d calls=%d stdout=%s stderr=%s", code, calls, stdout, stderr)
	}
	var envelope struct {
		OK    bool `json:"ok"`
		Error struct {
			Type    string `json:"type"`
			TraceID string `json:"trace_id"`
		} `json:"error"`
		Notice map[string]map[string]any `json:"_notice"`
	}
	if err := json.Unmarshal(stderr.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.OK || envelope.Error.Type != "cancelled" || envelope.Error.TraceID != trace || envelope.Notice["update"]["latest"] != "1.0.1" {
		t.Fatalf("invalid cancellation envelope: %s", stderr)
	}
	stderr.Reset()
	if code := app.Execute(context.Background(), []string{"--debug=bad", "version", "--output", "json"}); code != 1 || !json.Valid(stderr.Bytes()) || bytes.Contains(stderr.Bytes(), []byte(`"_notice"`)) {
		t.Fatalf("bad flags reused stale notice or lost JSON: code=%d stderr=%s", code, stderr)
	}
}

func TestExecuteOnlyExplicitJSONChangesErrorPresentation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		args []string
		json bool
	}{
		{[]string{"contract", "get"}, false},
		{[]string{"contract", "get", "--output", "json"}, true},
		{[]string{"contract", "get", "--output=json", "--debug"}, true},
		{[]string{"contract", "get", "--output", "table"}, false},
		{[]string{"unknown", "--json"}, false},
		{[]string{"unknown", "--data", "--output=json"}, false},
		{[]string{"environment", "inspect", "--output", "json", "--depth", "not-an-int"}, true},
		{[]string{"auth", "init", "--output", "json", "--unknown"}, true},
		{[]string{"environment", "inspect", "--depth", "not-an-int"}, false},
	} {
		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		app := cli.New(cli.Options{Store: config.NewStore(t.TempDir()), Stdout: stdout, Stderr: stderr,
			LookupEnv: func(key string) (string, bool) { return "1", key == "CONTRACT_CLI_NO_UPDATE_CHECK" },
		})
		if code := app.Execute(context.Background(), tc.args); code != 1 || stdout.Len() != 0 {
			t.Fatalf("ordinary error code/output: %d %s", code, stdout)
		}
		if tc.json {
			for _, line := range bytes.Split(bytes.TrimSpace(stderr.Bytes()), []byte("\n")) {
				if !json.Valid(line) {
					t.Fatalf("stderr mixed JSON and text: %s", stderr)
				}
			}
		} else if !strings.HasPrefix(stderr.String(), "error: ") {
			t.Fatalf("default/text error changed to JSON: %s", stderr)
		}
	}
}

func TestExecuteDistinguishesHTTPDeadlineAndUncertainCancellation(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"http", "timeout", "uncertain_write"} {
		t.Run(kind, func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			if err := store.UpsertProfile(uploadProfile(config.IdentityApp), true); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if kind == "timeout" {
				var stop context.CancelFunc
				ctx, stop = context.WithDeadline(ctx, time.Now().Add(-time.Second))
				defer stop()
			}
			stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
			calls := 0
			app := cli.New(cli.Options{Store: store, Stdout: stdout, Stderr: stderr,
				LookupEnv:          func(key string) (string, bool) { return "1", key == "CONTRACT_CLI_NO_UPDATE_CHECK" },
				InspectEnvironment: func(context.Context, int) invocation.Result { return invocation.Result{} },
				HTTPClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					calls++
					if kind == "uncertain_write" {
						cancel()
						return nil, req.Context().Err()
					}
					response := jsonResponse(`{"code":50001}`)
					response.StatusCode = http.StatusInternalServerError
					return response, nil
				})},
			})
			args := []string{"contract", "get", "1", "--profile", "contract", "--as", "app", "--output", "json"}
			code := 1
			if kind == "uncertain_write" {
				args[1] = "submit"
				code = 130
			}
			if got := app.Execute(ctx, args); got != code || calls > 1 {
				t.Fatalf("error classification changed exit/retry behavior: code=%d calls=%d", got, calls)
			}
			var envelope struct {
				Error struct {
					Type       string `json:"type"`
					HTTPStatus int    `json:"http_status"`
					Cancelled  bool   `json:"cancelled"`
				} `json:"error"`
			}
			if err := json.Unmarshal(stderr.Bytes(), &envelope); err != nil || envelope.Error.Type != kind || envelope.Error.Cancelled != (kind == "uncertain_write") {
				t.Fatalf("error envelope = %s (%v)", stderr, err)
			}
			if kind == "http" && envelope.Error.HTTPStatus != 500 {
				t.Fatalf("HTTP status lost: %s", stderr)
			}
		})
	}
}
