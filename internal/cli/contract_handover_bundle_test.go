package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"cn.qfei/contract-cli/internal/cli"
	"cn.qfei/contract-cli/internal/config"
	contractskills "cn.qfei/contract-cli/skills"
)

func TestHandoverSkillInstalledReferencesAreSelfContained(t *testing.T) {
	target := t.TempDir()
	app := cli.New(cli.Options{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}, Store: config.NewStore(t.TempDir()), SkillsFS: contractskills.FS})
	if err := app.Run(context.Background(), []string{"skills", "install", "--target", target}); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(target, "contract-handover")
	for _, name := range []string{"identity.md", "contracts.md", "handover.md", "assemblies.md", "transport.md"} {
		if _, err := os.Stat(filepath.Join(root, "references", name)); err != nil {
			t.Error(err)
		}
	}
	assertSearchSkillLinksResolve(t, root)
	if _, err := os.Stat(filepath.Join(root, "references", "api-contract.yaml")); !os.IsNotExist(err) {
		t.Errorf("development contract must not ship as runtime instructions: %v", err)
	}
}

func TestHandoverDocumentedWriteRequestPreservesIdentityAndResourceTypes(t *testing.T) {
	raw, err := contractskills.FS.ReadFile("contract-handover/references/handover.md")
	if err != nil {
		t.Fatal(err)
	}
	blocks := regexp.MustCompile("(?s)```json\\s*\\n(.*?)\\n```").FindAllSubmatch(raw, -1)
	if len(blocks) != 1 {
		t.Fatalf("want one complete write example, got %d", len(blocks))
	}
	var request struct {
		Manual    bool   `json:"manual"`
		Reason    string `json:"reason"`
		Employee  string `json:"targetEmployeeId"`
		Tenant    string `json:"targetTenantId"`
		Resources []struct {
			ID     string `json:"businessId"`
			Type   int    `json:"businessTypeCode"`
			Parent string `json:"parentId"`
		} `json:"resourceList"`
	}
	if err := json.Unmarshal(blocks[0][1], &request); err != nil {
		t.Fatal(err)
	}
	if !request.Manual || request.Reason != "" || request.Employee == "" || request.Tenant == "" || len(request.Resources) != 1 {
		t.Fatalf("incomplete handover request: %+v", request)
	}
	if request.Resources[0].ID != "9223372036854775806" || request.Resources[0].Parent != "9223372036854775805" {
		t.Fatal("Long resource IDs must retain all decimal digits as strings")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(blocks[0][1], &fields); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"handoverByAdmin", "onlyOwner", "onlySubmitter", "operatorId", "token"} {
		if _, ok := fields[forbidden]; ok {
			t.Errorf("write example must not invent %s", forbidden)
		}
	}
}
