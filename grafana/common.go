package grafana

import (
	"net/http"

	"github.com/grafana-tools/sdk"
)

type APIClient struct {
	*sdk.Client
}

func NewAPIClient(grafanaEndpoint string, grafanaToken string) APIClient {
	return APIClient{
		sdk.NewClient(grafanaEndpoint, grafanaToken, &http.Client{}),
	}
}
