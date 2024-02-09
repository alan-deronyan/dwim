package clients

import (
	"github.com/deronyan-llc/fluree-client/fluree"
)

type FlureeClient struct {
	Config *fluree.Configuration
	*fluree.APIClient
}

func NewFlureeClient() *FlureeClient {
	config := fluree.NewConfiguration()
	config.BasePath = "http://localhost:58090"
	config.AddDefaultHeader("Content-Type", "application/json")
	client := &FlureeClient{
		APIClient: fluree.NewAPIClient(config),
		Config:    config,
	}

	if client.APIClient == nil {
		panic("Failed to create fluree client")
	}

	return client
}
