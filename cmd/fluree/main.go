package main

import (
	"fmt"

	"github.com/deronyan-llc/dwim/internal/clients"
)

func assertResponse(resp string, err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}

func main() {
	flureeClient := clients.NewFlureeClient("http://localhost:58090")

	assertResponse(flureeClient.Create("foo/bar", nil, []string{`
		{
			"name": "foo"
		}
	`}))

	assertResponse(flureeClient.Query("foo/bar", `
		"where": {
			"@id": "?s"
		},
		"select": {"?s": ["*"]}
	`))
}
