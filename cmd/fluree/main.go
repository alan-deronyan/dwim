package main

import (
	"context"

	"github.com/deronyan-llc/dwim/internal/clients"
	"github.com/deronyan-llc/fluree-client/fluree"
)

func main() {
	client := clients.NewFlureeClient()

	requestBody := fluree.Body1{
		FlureeServerComponentsHttpquery: `{
			"from": "cookbook/base",
			"where": {
				"@id": "?s",
				"schema:name": "?name"
			},
			"select": {"?s": ["*"]}
		}`,
	}

	rtn, resp, err := client.APIClient.DefaultApi.FlureeQueryPost(context.Background(), requestBody)
	if err != nil {
		panic(err)
	}
	println(resp)
	println(rtn)
}
