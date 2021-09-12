/*
Copyright 2021.

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
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	webappv1 "github.com/cmmarslender/web-operator/api/v1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SimpleAppReconciler reconciles a SimpleApp object
type SimpleAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.k8s.cmm.io,resources=simpleapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.k8s.cmm.io,resources=simpleapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.k8s.cmm.io,resources=simpleapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SimpleApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SimpleAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info(fmt.Sprintf("SimpleApp name is %s", req.NamespacedName))

	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(r.Log))

	var app webappv1.SimpleApp
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Create a dummy config map with the value from the CRD, just as a test
	configMapObject := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: req.Namespace,
		},
		Data: map[string]string{
			"foo": app.Spec.Foo,
		},
	}

	err := ctrl.SetControllerReference(&app, configMapObject, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}
	result, err := resourceReconciler.ReconcileResource(configMapObject, reconciler.StatePresent)
	if err != nil {
		return ctrl.Result{}, err
	}

	if result != nil {
		return *result, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimpleAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Log = ctrl.Log.WithName("controllers").WithName("SimpleApp")

	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.SimpleApp{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
