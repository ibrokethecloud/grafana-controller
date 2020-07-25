package grafana

import (
	"github.com/grafana-tools/sdk"
	"github.com/ibrokethecloud/grafana-controller/api/v1alpha1"
	json "github.com/json-iterator/go"
)

//set a dashboard and return details
func (a APIClient) SetDashboard(dSpec v1alpha1.DashboardSpec) (dStatus v1alpha1.DashboardStatus,
	err error) {
	folderid := 0
	board := sdk.Board{}
	err = json.Unmarshal([]byte(dSpec.Content), &board)
	if err != nil {
		dStatus.Message = err.Error()
		return dStatus, err
	}
	if len(dSpec.Folder) != 0 {
		// If we find no folder we will default to the main folder //
		folderid, _, _, _ = a.GetFolderByName(dSpec.Folder)
	}
	params := sdk.SetDashboardParams{
		FolderID:  folderid,
		Overwrite: true,
	}
	statusMessage, err := a.Client.SetDashboard(a.Context, board, params)
	if err != nil {
		dStatus.Message = err.Error()
	}

	return processStatusMessage(statusMessage), err
}

// delete the specified dashboard when the object is deleted
func (a APIClient) DeleteDashboard(slug string) (dStatus v1alpha1.DashboardStatus,
	err error) {
	statusMessage, err := a.Client.DeleteDashboard(a.Context, slug)
	if err != nil {
		dStatus.Message = err.Error()
	}
	return processStatusMessage(statusMessage), err
}

func processStatusMessage(status sdk.StatusMessage) (dStatus v1alpha1.DashboardStatus) {
	if status.ID != nil {
		dStatus.ID = *status.ID
	}
	if status.UID != nil {
		dStatus.UID = *status.UID
	}
	if status.Message != nil {
		dStatus.Message = *status.Message
	}
	if status.Slug != nil {
		dStatus.Slug = *status.Slug
	}
	if status.URL != nil {
		dStatus.URL = *status.URL
	}
	return
}
