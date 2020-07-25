/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/ibrokethecloud/grafana-controller/grafana"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sv1alpha1 "github.com/ibrokethecloud/grafana-controller/api/v1alpha1"
)

// DatasourceReconciler reconciles a Datasource object
type DatasourceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	grafana.APIClient
}

// +kubebuilder:rbac:groups=k8s.grafana.io,resources=datasources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.grafana.io,resources=datasources/status,verbs=get;update;patch

func (r *DatasourceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	datasourceFinalizer := "datasource.k8s.grafana.io"
	ctx := context.Background()
	log := r.Log.WithValues("datasource", req.NamespacedName)

	var datasource k8sv1alpha1.Datasource
	if err := r.Get(ctx, req.NamespacedName, &datasource); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch datasource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if datasource.ObjectMeta.DeletionTimestamp.IsZero() {
		// Manage datasource updates //
		dStatus, err := r.APIClient.SetDatasource(datasource, r.Log)
		if err != nil {
			log.Error(err, "Unable to update datasource object")
			return ctrl.Result{}, err
		} else {
			datasource.Status = dStatus
			if !containsString(datasource.ObjectMeta.Finalizers, datasourceFinalizer) {
				datasource.ObjectMeta.Finalizers = append(datasource.ObjectMeta.Finalizers, datasourceFinalizer)
			}
			if err := r.Update(context.Background(), &datasource); err != nil {
				return ctrl.Result{}, err
			}

		}

	} else {
		// The object is being deleted
		if containsString(datasource.ObjectMeta.Finalizers, datasourceFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			if dStatus, err := r.deleteDatasourceResources(&datasource); err != nil {
				datasource.Status = dStatus
				if err := r.Update(ctx, &datasource); err != nil {
					return ctrl.Result{}, err
				}
			}

			datasource.ObjectMeta.Finalizers = removeString(datasource.ObjectMeta.Finalizers, datasourceFinalizer)
			if err := r.Update(context.Background(), &datasource); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *DatasourceReconciler) deleteDatasourceResources(datasource *k8sv1alpha1.Datasource) (dsStatus k8sv1alpha1.DatasourceStatus,
	err error) {
	dsStatus = *datasource.Status.DeepCopy()
	err = r.APIClient.DeleteDatasource(dsStatus.ID)
	if err != nil {
		dsStatus.Message = err.Error()
	}
	return dsStatus, err
}

func (r *DatasourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1alpha1.Datasource{}).
		Complete(r)
}
