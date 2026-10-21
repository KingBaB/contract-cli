package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"cn.qfei/contract-cli/internal/cli"
	"cn.qfei/contract-cli/internal/config"
	contractskills "cn.qfei/contract-cli/skills"
)

type handoverLookupRequest struct {
	PageIndex         int    `json:"pageIndex"`
	PageSize          int    `json:"pageSize"`
	SearchTabCode     int    `json:"searchTabCode"`
	AllConditionUnits bool   `json:"allConditionUnits"`
	SortType          string `json:"sortType"`
	Order             string `json:"order"`
	WithStatistics    bool   `json:"withStatistics"`
	Conditions        []struct {
		Field string `json:"searchField"`
		Value string `json:"searchValue"`
		Union string `json:"unionType"`
	} `json:"conditionUnits"`
}

func TestHandoverAssemblyLookupExampleUsesSelectedCandidate(t *testing.T) {
	body := handoverLookupExample(t)
	var request handoverLookupRequest
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		t.Fatal(err)
	}
	if request.PageIndex != 0 || request.PageSize != 100 || request.SearchTabCode != 0 || request.AllConditionUnits || request.WithStatistics || request.SortType != "SUBMITTED_TIME" || request.Order != "DESC" || len(request.Conditions) != 1 {
		t.Fatalf("lookup must retain one bounded candidate condition: %+v", request)
	}
	condition := request.Conditions[0]
	if condition.Field != "CONTRACT_NUMBER_FIELD" || condition.Value == "" || condition.Union != "SHOULD" {
		t.Fatalf("lookup must use the selected candidate's number: %+v", condition)
	}
}

func TestHandoverAssemblyLookupExamplePreservesMemberOnlyMapping(t *testing.T) {
	body := handoverLookupExample(t)
	responseBody := `{"code":200,"success":true,"data":{"content":[{"groupId":"90001","searchContractListItems":[{"contractId":"9223372036854775806","contractNumber":"CG-EXAMPLE-1017","businessTypeCode":11,"assemblyId":"9223372036854775805"}]}]}}`
	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(apiWebTestProfile(), true); err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(t.TempDir(), "lookup.json")
	if err := os.WriteFile(input, body, 0600); err != nil {
		t.Fatal(err)
	}
	stdout := &bytes.Buffer{}
	requests := 0
	transport := handoverLookupTransport(t, body, responseBody, &requests)
	app := cli.New(cli.Options{Store: store, Stdout: stdout, Stderr: &bytes.Buffer{}, HTTPClient: &http.Client{Transport: transport}})
	if err := app.Run(context.Background(), []string{"api", "call", "--channel", "web", "--profile", "contract", "--as", "user", "--input-file", input, "POST", "/open-web/web/clm/es/search/contractList"}); err != nil {
		t.Fatal(err)
	}
	assertHandoverLookupMember(t, stdout.Bytes(), requests)
}

func handoverLookupTransport(t *testing.T, body []byte, responseBody string, requests *int) roundTripFunc {
	t.Helper()
	return func(request *http.Request) (*http.Response, error) {
		*requests++
		if request.Method != http.MethodPost || request.URL.Path != "/open-web/web/clm/es/search/contractList" || request.URL.RawQuery != "" {
			t.Fatalf("unexpected lookup request: %s %s", request.Method, request.URL)
		}
		actual, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(actual, body) {
			t.Fatalf("documented lookup body changed in transit: %s", actual)
		}
		return jsonResponse(responseBody), nil
	}
}

func assertHandoverLookupMember(t *testing.T, raw []byte, requests int) {
	t.Helper()
	var response struct {
		Data struct {
			Content []struct {
				Members []struct {
					ContractID string `json:"contractId"`
					AssemblyID string `json:"assemblyId"`
				} `json:"searchContractListItems"`
			} `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatal(err)
	}
	if requests != 1 || len(response.Data.Content) != 1 || len(response.Data.Content[0].Members) != 1 {
		t.Fatalf("lookup must retain the member without requiring a parent search row: %s", raw)
	}
	member := response.Data.Content[0].Members[0]
	if member.ContractID != "9223372036854775806" || member.AssemblyID != "9223372036854775805" {
		t.Fatalf("member and assembly identities must remain distinct, exact strings: %+v", member)
	}
}

func handoverLookupExample(t *testing.T) []byte {
	t.Helper()
	raw, err := contractskills.FS.ReadFile("contract-handover/references/assembly-lookup.md")
	if err != nil {
		t.Fatal(err)
	}
	blocks := regexp.MustCompile("(?s)```json\\s*\\n(.*?)\\n```").FindAllSubmatch(raw, -1)
	if len(blocks) != 1 {
		t.Fatalf("want one complete lookup request example, got %d", len(blocks))
	}
	return blocks[0][1]
}
