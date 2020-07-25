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

// FolderReconciler reconciles a Folder object
type FolderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	grafana.APIClient
}

// +kubebuilder:rbac:groups=k8s.grafana.io,resources=folders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.grafana.io,resources=folders/status,verbs=get;update;patch

func (r *FolderReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	folderFinalizer := "folder.k8s.grafana.io"
	ctx := context.Background()
	log := r.Log.WithValues("folder", req.NamespacedName)

	var folder k8sv1alpha1.Folder
	if err := r.Get(ctx, req.NamespacedName, &folder); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch folder")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if folder.ObjectMeta.DeletionTimestamp.IsZero() {
		// Manage folder updates //
		fStatus, err := r.APIClient.SetFolder(folder.Spec)
		if err != nil {
			log.Error(err, "Unable to update folder object")
			return ctrl.Result{}, err
		} else {
			folder.Status = fStatus
			if !containsString(folder.ObjectMeta.Finalizers, folderFinalizer) {
				folder.ObjectMeta.Finalizers = append(folder.ObjectMeta.Finalizers, folderFinalizer)
			}
			if err := r.Update(context.Background(), &folder); err != nil {
				return ctrl.Result{}, err
			}

		}

	} else {
		// The object is being deleted
		if containsString(folder.ObjectMeta.Finalizers, folderFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			if fStatus, err := r.deleteFolderResources(&folder); err != nil {
				folder.Status = fStatus
				if err := r.Update(ctx, &folder); err != nil {
					return ctrl.Result{}, err
				}
			}

			folder.ObjectMeta.Finalizers = removeString(folder.ObjectMeta.Finalizers, folderFinalizer)
			if err := r.Update(context.Background(), &folder); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *FolderReconciler) deleteFolderResources(folder *k8sv1alpha1.Folder) (k8sv1alpha1.FolderStatus,
	error) {
	uid := folder.Status.UID
	err := r.APIClient.DeleteFolder(uid)
	folderStatus := folder.Status.DeepCopy()
	if err != nil {
		folderStatus.Message = "Unable to remote the folder"
	}
	return *folderStatus, err
}

func (r *FolderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1alpha1.Folder{}).
		Complete(r)
}
