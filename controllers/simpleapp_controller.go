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

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	webappv1 "github.com/cmmarslender/web-operator/api/v1"
	util "github.com/cmmarslender/web-operator/pkg"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// Annotations
	lastAppliedAnnotationKey = "webapp.k8s.cmm.io/last-applied"

	// Labels
	typeLabelKey = "webapp.k8s.cmm.io/type" // SimpleApp, etc
	nameLabelKey = "webapp.k8s.cmm.io/name"
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

	var app webappv1.SimpleApp
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// @TODO move these to a desired state generator
	// Deployment
	objectMeta := r.getObjectMeta(app)
	deploymentObject := &appsv1.Deployment{
		ObjectMeta: objectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Selector: &metav1.LabelSelector{ // @TODO could make this a helper - takes obj meta, returns label selector
				MatchLabels: objectMeta.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{ // @TODO could make this a helper - takes obj meta, returns simple obj meta for template
					Labels: objectMeta.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            app.Name,
							Image:           app.Spec.Image,
							ImagePullPolicy: app.Spec.ImagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: app.Spec.ContainerPort,
								},
							},
						},
					},
					ImagePullSecrets: r.namesToLocalObjectRefs(app.Spec.ImagePullSecrets),
				},
			},
		},
	}

	// Service
	serviceObject := &corev1.Service{
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       app.Spec.ServicePort,
					TargetPort: intstr.IntOrString{IntVal: app.Spec.ContainerPort},
				},
			},
			Selector: objectMeta.Labels,
		},
	}

	// Ingress
	ingressObject := &networkingv1.Ingress{
		ObjectMeta: objectMeta,
		Spec: networkingv1.IngressSpec{
			// @TODO support TLS
			Rules: []networkingv1.IngressRule{
				{
					Host: app.Spec.Hostname,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: r.ingressPathsHelper(app),
						},
					},
				},
			},
		},
	}

	result, err := r.ReconcileResource(app, deploymentObject, reconciler.StatePresent)
	if result != nil || err != nil {
		return util.ReconcileReturnHelper(result, err)
	}

	// @TODO service should always be enabled if ingress is enabled
	result, err = r.ReconcileResource(app, serviceObject, util.ReconcilerStateHelper(app.Spec.ServiceEnabled))
	if result != nil || err != nil {
		return util.ReconcileReturnHelper(result, err)
	}

	result, err = r.ReconcileResource(app, ingressObject, util.ReconcilerStateHelper(app.Spec.IngressEnabled))
	if result != nil || err != nil {
		return util.ReconcileReturnHelper(result, err)
	}

	return reconcile.Result{}, nil
}

// ReconcileResource Sets ownership of the resource and then ensures the resource is in the correct state in the cluster
func (r *SimpleAppReconciler) ReconcileResource(app webappv1.SimpleApp, obj client.Object, state reconciler.DesiredState) (*reconcile.Result, error) {
	// @TODO this (along with the app) should probably live in some sort of parent reconciler struct
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(r.Log))

	err := ctrl.SetControllerReference(&app, obj, r.Scheme)
	if err != nil {
		return nil, err
	}
	return resourceReconciler.ReconcileResource(obj, state)
}

// getObjectMeta returns the object meta for resources owned by the SimpleApp
func (r *SimpleAppReconciler) getObjectMeta(app webappv1.SimpleApp) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: app.Namespace,
		Name:      app.Name,
		Labels:    r.labels(app),
	}
}

// labels returns labels for the provided app
func (r *SimpleAppReconciler) labels(app webappv1.SimpleApp) map[string]string {
	return map[string]string{
		typeLabelKey: app.Kind,
		nameLabelKey: app.Name,
	}
}

// namesToLocalObjectRefs returns []LocalObjectReference from []string
func (r *SimpleAppReconciler) namesToLocalObjectRefs(names []string) []corev1.LocalObjectReference {
	var refs []corev1.LocalObjectReference

	for _, name := range names {
		refs = append(refs, corev1.LocalObjectReference{Name: name})
	}

	return refs
}

// ingressPathsHelper returns generated ingress paths for the app
func (r *SimpleAppReconciler) ingressPathsHelper(app webappv1.SimpleApp) []networkingv1.HTTPIngressPath {
	var paths []networkingv1.HTTPIngressPath

	prefixType := networkingv1.PathTypePrefix

	for _, path := range app.Spec.IngressPaths {
		paths = append(paths, networkingv1.HTTPIngressPath{
			Path:     path,
			PathType: &prefixType,
			Backend: networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: app.Name,
					Port: networkingv1.ServiceBackendPort{
						Number: app.Spec.ContainerPort,
					},
				},
			},
		})
	}

	return paths
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimpleAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Set our own default annotator so we can control the key used for last-applied
	patch.DefaultAnnotator = patch.NewAnnotator(lastAppliedAnnotationKey)

	r.Log = ctrl.Log.WithName("controllers").WithName("SimpleApp")

	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.SimpleApp{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}
