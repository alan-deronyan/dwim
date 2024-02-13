package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type FlureeClient struct {
	endpoint string
}

func NewFlureeClient(endpoint string) *FlureeClient {
	client := &FlureeClient{
		endpoint: endpoint,
	}
	return client
}

func formatArray(arrayStr []string) string {
	if len(arrayStr) == 1 {
		return arrayStr[0]
	} else {
		return "[" + strings.Join(arrayStr, ",") + "]"
	}
}

func formatMap(mapStr map[string]string) string {
	if len(mapStr) == 0 {
		return "{}"
	}
	mapFmt := "{"
	for k, v := range mapStr {
		if len(mapFmt) > 1 {
			mapFmt += ","
		}
		mapFmt += fmt.Sprintf("\"%s\": \"%s\"", k, v)
	}
	mapFmt += "}"
	return mapFmt
}

func (client *FlureeClient) executePostRequest(method string, bodyStr string) (string, error) {
	fmt.Printf("POST %s...", method)

	body := strings.NewReader(bodyStr)
	resp, err := http.Post(client.endpoint+"/fluree/"+method, "application/json", body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %e", err)
		return "", err
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, bodyBytes, "", " ")
	if err != nil {
		fmt.Printf("Error: %e", err)
		return "", err
	}
	return prettyJSON.String(), err
}

func (client *FlureeClient) assembleRequestBody(ledger string, contexts map[string]string, inserts []string, deletes []string) (string, error) {
	if len(inserts) == 0 && len(deletes) == 0 {
		return "", errors.New("inserts and deletes cannot both be empty")
	}

	// build the parts array with the ledger, contexts, inserts, and deletes
	parts := []string{fmt.Sprintf(`"ledger": "%s"`, ledger)}

	if len(contexts) > 0 {
		parts = append(parts, fmt.Sprintf(`"@context": %s`, formatMap(contexts)))
	}

	if len(inserts) > 0 {
		parts = append(parts, fmt.Sprintf(`"insert": %s`, formatArray(inserts)))
	}

	if len(deletes) > 0 {
		parts = append(parts, fmt.Sprintf(`"delete": %s`, formatArray(deletes)))
	}

	return "{" + strings.Join(parts, ",") + "}", nil
}

func (client *FlureeClient) Create(ledger string, contexts map[string]string, inserts []string) (string, error) {
	body, err := client.assembleRequestBody(ledger, contexts, inserts, nil)
	if err != nil {
		return "", err
	}
	return client.executePostRequest("create", body)
}

func (client *FlureeClient) Transact(ledger string, contexts map[string]string, inserts []string, deletes []string) (string, error) {
	body, err := client.assembleRequestBody(ledger, contexts, inserts, deletes)
	if err != nil {
		return "", err
	}
	return client.executePostRequest("transact", body)
}

func (client *FlureeClient) Query(ledger string, query string) (string, error) {
	// build the parts array with the ledger and query
	parts := []string{
		fmt.Sprintf(`"from": "%s"`, ledger),
		query,
	}

	body := "{" + strings.Join(parts, ",") + "}"

	return client.executePostRequest("query", body)
}
