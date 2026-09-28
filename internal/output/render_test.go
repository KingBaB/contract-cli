package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"cn.qfei/contract-cli/internal/output"
)

func TestRendererRenderJSON(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout)

	if err := renderer.Render(output.FormatJSON, json.RawMessage(`{"code":0,"data":{"vendorId":"123"}}`)); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, `"vendorId": "123"`) {
		t.Fatalf("unexpected json output: %s", got)
	}
}

func TestRendererRenderJSONDoesNotEscapeHTML(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout)

	if err := renderer.Render(output.FormatJSON, map[string]any{"message": "0.1.0 -> 1.2.0"}); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := stdout.String()
	if strings.Contains(got, `\u003e`) || !strings.Contains(got, `"message": "0.1.0 -> 1.2.0"`) {
		t.Fatalf("json should not html-escape message: %s", got)
	}
}

func TestRendererRenderJSONPreservesLargeIntegerPrecision(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout)

	if err := renderer.Render(output.FormatJSON, json.RawMessage(`{"code":0,"data":{"contract_id":1140439949925941626,"tenant_id":1118543858007802222}}`)); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := stdout.String()
	for _, want := range []string{
		`"contract_id": 1140439949925941626`,
		`"tenant_id": 1118543858007802222`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("json output should preserve %q, got: %s", want, got)
		}
	}
	if strings.Contains(got, "1140439949925941600") {
		t.Fatalf("json output lost integer precision: %s", got)
	}
}

func TestRendererRenderRaw(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout)

	if err := renderer.RenderRaw([]byte(`{"ok":true}`)); err != nil {
		t.Fatalf("RenderRaw() error = %v", err)
	}
	if stdout.String() != `{"ok":true}` {
		t.Fatalf("unexpected raw output: %s", stdout.String())
	}
}

func TestRendererRenderTable(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout)

	value := []map[string]any{
		{"vendorId": "123", "name": "Acme"},
	}
	if err := renderer.Render(output.FormatTable, value); err != nil {
		t.Fatalf("Render(table) error = %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "vendorId") || !strings.Contains(got, "Acme") {
		t.Fatalf("unexpected table output: %s", got)
	}
}

func TestRendererRenderYAML(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout)

	if err := renderer.Render(output.FormatYAML, map[string]any{
		"code": 0,
		"data": map[string]any{"vendorId": "123"},
	}); err != nil {
		t.Fatalf("Render(yaml) error = %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "code: 0") || !strings.Contains(got, "vendorId: \"123\"") {
		t.Fatalf("unexpected yaml output: %s", got)
	}
}

func TestRendererInjectsNoticeIntoJSONObject(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout).WithNotice(map[string]any{
		"update": map[string]any{
			"current": "0.1.0-beta.1",
			"latest":  "0.1.0-beta.2",
		},
	})

	if err := renderer.Render(output.FormatJSON, json.RawMessage(`{"code":0,"data":{"vendorId":"123"}}`)); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	notice, ok := got["_notice"].(map[string]any)
	if !ok {
		t.Fatalf("missing _notice in output: %+v", got)
	}
	update, ok := notice["update"].(map[string]any)
	if !ok || update["latest"] != "0.1.0-beta.2" {
		t.Fatalf("unexpected _notice.update: %+v", notice["update"])
	}
}

func TestRendererMergesExistingNotice(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	renderer := output.NewRenderer(&stdout).WithNotice(map[string]any{
		"update": map[string]any{"latest": "0.1.0-beta.2"},
	})

	if err := renderer.Render(output.FormatJSON, json.RawMessage(`{"ok":true,"_notice":{"skills":{"message":"keep"}}}`)); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	notice := got["_notice"].(map[string]any)
	if _, ok := notice["skills"].(map[string]any); !ok {
		t.Fatalf("existing notice was not preserved: %+v", notice)
	}
	if _, ok := notice["update"].(map[string]any); !ok {
		t.Fatalf("update notice was not injected: %+v", notice)
	}
}

func TestRendererSkipsNoticeForNonObjectAndNonJSON(t *testing.T) {
	t.Parallel()

	notice := map[string]any{"update": map[string]any{"latest": "0.1.0-beta.2"}}

	var jsonOut bytes.Buffer
	if err := output.NewRenderer(&jsonOut).WithNotice(notice).Render(output.FormatJSON, json.RawMessage(`[{"code":0}]`)); err != nil {
		t.Fatalf("Render(json array) error = %v", err)
	}
	if strings.Contains(jsonOut.String(), "_notice") {
		t.Fatalf("json array should not receive notice: %s", jsonOut.String())
	}

	var yamlOut bytes.Buffer
	if err := output.NewRenderer(&yamlOut).WithNotice(notice).Render(output.FormatYAML, map[string]any{"code": 0}); err != nil {
		t.Fatalf("Render(yaml) error = %v", err)
	}
	if strings.Contains(yamlOut.String(), "_notice") {
		t.Fatalf("yaml output should not receive notice: %s", yamlOut.String())
	}

	var tableOut bytes.Buffer
	if err := output.NewRenderer(&tableOut).WithNotice(notice).Render(output.FormatTable, map[string]any{"code": 0}); err != nil {
		t.Fatalf("Render(table) error = %v", err)
	}
	if strings.Contains(tableOut.String(), "_notice") {
		t.Fatalf("table output should not receive notice: %s", tableOut.String())
	}
}
