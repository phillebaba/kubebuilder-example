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
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webv1alpha1 "github.com/phillebaba/kubebuilder-example/api/v1alpha1"
)

// RemoteContentReconciler reconciles a RemoteContent object
type RemoteContentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=web.phillebaba.io,resources=remotecontents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=web.phillebaba.io,resources=remotecontents/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *RemoteContentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("remotecontent", req.NamespacedName)

	set_state := func(rc *webv1alpha1.RemoteContent, state webv1alpha1.RequestState) {
		rc.Status.State = state
		if err := r.Status().Update(ctx, rc); err != nil {
			log.Error(err, "unable to update RemoteContent status")
		}
	}

	var rc webv1alpha1.RemoteContent
	if err := r.Get(ctx, req.NamespacedName, &rc); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if rc.Status.State == "" {
		set_state(&rc, webv1alpha1.Pending)
	}

	resp, err := http.Get(rc.Spec.Url)
	if err != nil {
		log.Error(err, "unable to make request", "url", rc.Spec.Url)
		set_state(&rc, webv1alpha1.Failed)
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, "could not read response body")
		set_state(&rc, webv1alpha1.Failed)
		return ctrl.Result{}, err
	}
	content := string(body)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rc.Name,
			Namespace: rc.Namespace,
		},
	}

	if _, err := ctrl.CreateOrUpdate(ctx, r.Client, cm, func() error {
		cm.Data = map[string]string{"content": content}
		return ctrl.SetControllerReference(&rc, cm, r.Scheme)
	}); err != nil {
		set_state(&rc, webv1alpha1.Failed)
		return ctrl.Result{}, err
	}

	set_state(&rc, webv1alpha1.Succeded)
	return ctrl.Result{}, nil
}

func (r *RemoteContentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1alpha1.RemoteContent{}).
		Complete(r)
}
