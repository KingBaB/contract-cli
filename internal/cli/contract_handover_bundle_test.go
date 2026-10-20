package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func TestHandoverParameterReferencesHaveTypedRequiredness(t *testing.T) {
	for _, name := range []string{"identity.md", "contracts.md", "handover.md", "assemblies.md"} {
		raw, err := contractskills.FS.ReadFile("contract-handover/references/" + name)
		if err != nil {
			t.Fatal(err)
		}
		text := string(raw)
		if strings.Contains(text, "| 参数 | 位置 | 必填 |") {
			t.Errorf("%s conflates API requiredness with recipe policy and omits types", name)
		}
		if !strings.Contains(text, "| 参数 | 位置 | 类型 | 本流程要求 | 值、来源与限制 |") {
			t.Errorf("%s has no complete parameter definition", name)
		}
	}
	raw, err := contractskills.FS.ReadFile("contract-handover/references/resource-types.md")
	if err != nil {
		t.Fatal(err)
	}
	rows := regexp.MustCompile(`(?m)^\| ([0-9]+) \|`).FindAllSubmatch(raw, -1)
	codes := map[string]bool{}
	for _, row := range rows {
		codes[string(row[1])] = true
	}
	for _, code := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "99"} {
		if !codes[code] {
			t.Errorf("missing business enum %s", code)
		}
	}
}

func TestHandoverSearchExampleUsesExplicitTypedConditions(t *testing.T) {
	raw, err := contractskills.FS.ReadFile("contract-handover/references/contracts.md")
	if err != nil {
		t.Fatal(err)
	}
	blocks := regexp.MustCompile("(?s)```json\\s*\\n(.*?)\\n```").FindAllSubmatch(raw, -1)
	if len(blocks) == 0 {
		t.Fatal("missing complete search request")
	}
	for _, block := range blocks {
		var request struct {
			Page       int `json:"pageIndex"`
			Size       int `json:"pageSize"`
			Tab        int `json:"searchTabCode"`
			Conditions []struct {
				Field string `json:"searchField"`
				Value string `json:"searchValue"`
				Union string `json:"unionType"`
			} `json:"conditionUnits"`
			All        bool   `json:"allConditionUnits"`
			Sort       string `json:"sortType"`
			Order      string `json:"order"`
			Statistics bool   `json:"withStatistics"`
		}
		decoder := json.NewDecoder(bytes.NewReader(block[1]))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.Page != 0 || request.Size != 100 || (request.Tab != 0 && request.Tab != 1) || request.All || request.Statistics || request.Sort != "SUBMITTED_TIME" || request.Order != "DESC" || len(request.Conditions) != 1 {
			t.Fatalf("unsafe or incomplete search example: %+v", request)
		}
		unit := request.Conditions[0]
		if (unit.Field != "CONTRACT_NUMBER_FIELD" && unit.Field != "CONTRACT_NAME_FIELD") || unit.Value == "" || unit.Union != "SHOULD" {
			t.Fatalf("invalid search condition: %+v", unit)
		}
	}
}
