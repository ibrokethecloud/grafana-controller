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

// DashboardReconciler reconciles a Dashboard object
type DashboardReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	grafana.APIClient
}

// +kubebuilder:rbac:groups=k8s.grafana.io,resources=dashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.grafana.io,resources=dashboards/status,verbs=get;update;patch

func (r *DashboardReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	dashboardFinalizer := "dashboard.k8s.grafana.io"
	ctx := context.Background()
	log := r.Log.WithValues("dashboard", req.NamespacedName)

	var dashboard k8sv1alpha1.Dashboard
	if err := r.Get(ctx, req.NamespacedName, &dashboard); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch dashboard")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if dashboard.ObjectMeta.DeletionTimestamp.IsZero() {
		// Manage dashboard updates //
		dStatus, err := r.APIClient.SetDashboard(dashboard.Spec)
		if err != nil {
			log.Error(err, "Unable to update dashboard object")
			return ctrl.Result{}, err
		} else {
			dashboard.Status = dStatus
			if !containsString(dashboard.ObjectMeta.Finalizers, dashboardFinalizer) {
				dashboard.ObjectMeta.Finalizers = append(dashboard.ObjectMeta.Finalizers, dashboardFinalizer)
			}
			if err := r.Update(context.Background(), &dashboard); err != nil {
				return ctrl.Result{}, err
			}

		}

	} else {
		// The object is being deleted
		if containsString(dashboard.ObjectMeta.Finalizers, dashboardFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			if dStatus, err := r.deleteDashboardResources(&dashboard); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				dashboard.Status = dStatus
				if err := r.Update(ctx, &dashboard); err != nil {
					return ctrl.Result{}, err
				}
			}

			dashboard.ObjectMeta.Finalizers = removeString(dashboard.ObjectMeta.Finalizers, dashboardFinalizer)
			if err := r.Update(context.Background(), &dashboard); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *DashboardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1alpha1.Dashboard{}).
		Complete(r)
}

func (r *DashboardReconciler) deleteDashboardResources(dashboard *k8sv1alpha1.Dashboard) (k8sv1alpha1.DashboardStatus,
	error) {
	slug := dashboard.Status.Slug
	dStatus, err := r.APIClient.DeleteDashboard(slug)
	return dStatus, err
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
