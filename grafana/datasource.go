package grafana

import (
	"github.com/go-logr/logr"
	"github.com/grafana-tools/sdk"
	"github.com/ibrokethecloud/grafana-controller/api/v1alpha1"
	json "github.com/json-iterator/go"
)

func (a APIClient) SetDatasource(datasource v1alpha1.Datasource, log logr.Logger) (dsStatus v1alpha1.DatasourceStatus,
	err error) {
	log.Info("About to create the datasource")
	ds := sdk.Datasource{}
	if err = json.Unmarshal([]byte(datasource.Spec.Content), &ds); err != nil {
		dsStatus.Message = err.Error()
		return dsStatus, err
	}

	if datasource.Status.ID != 0 {
		// Datasource already exists.. lets just updated it //
		ds.ID = datasource.Status.ID
		statusMessage, err := a.Client.UpdateDatasource(a.Context, ds)
		if err != nil {
			return dsStatus, err
		}
		dsStatus.ID = *statusMessage.ID
	} else {
		// Looks like upstream DS creation function doesnt return
		// a status message. This causes issues.
		_, err = a.Client.CreateDatasource(a.Context, ds)
		if err != nil {
			dsStatus.Message = err.Error()
			return dsStatus, err
		}
		fetchDS, err := a.Client.GetDatasourceByName(a.Context, ds.Name)
		if err != nil {
			dsStatus.Message = err.Error()
			return dsStatus, err
		}
		dsStatus.ID = fetchDS.ID
		dsStatus.Message = "Datasource created by controller"
	}

	return dsStatus, err
}

func (a APIClient) DeleteDatasource(id uint) (err error) {
	_, err = a.Client.DeleteDatasource(a.Context, id)
	return err
}
