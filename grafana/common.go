package grafana

import (
	"context"
	"net/http"

	"github.com/grafana-tools/sdk"
)

type APIClient struct {
	*sdk.Client
	context.Context
}

func NewAPIClient(grafanaEndpoint string, grafanaToken string) APIClient {
	return APIClient{
		sdk.NewClient(grafanaEndpoint, grafanaToken, &http.Client{}),
		context.Background(),
	}
}
