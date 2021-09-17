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

	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/util/intstr"
	msaberdevv1 "custom-k8s-operator/api/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	// "sigs.k8s.io/controller-runtime/pkg/source"
)

// OnekindReconciler reconciles a Onekind object
type OnekindReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// func (r *OnekindReconciler) updateResourceState(ctx context.Context, kind msaberdevv1.Onekind, phase msaberdevv1.StatusPhase) (ctrl.Result, error) {
// 	kind.Status.Phase = msaberdevv1.ErrorStatusPhase

// 	err := r.Status().Update(ctx, &kind)
// 	if err != nil {
// 		return ctrl.Result{}, err
// 	}
// 	return ctrl.Result{}, nil
// }

//+kubebuilder:rbac:groups=msaber.dev,resources=onekinds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=msaber.dev,resources=onekinds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=msaber.dev,resources=onekinds/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Onekind object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *OnekindReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log := log.Log.WithValues("onekinds", req.NamespacedName)
	log.Info("reconciling guestbook")

	var kind msaberdevv1.Onekind

	// Check object before creation
	if err := r.Get(ctx, req.NamespacedName, &kind); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	kind.Status.Phase = msaberdevv1.PendingStauePhase

	err := r.Status().Update(ctx, &kind)
	if err != nil {
		return ctrl.Result{}, err
	}

	cm, err := r.desiredConfigMap(kind)
	if err != nil {
		return ctrl.Result{}, err
	}

	deployment, err := r.desiredDeployment(kind, *cm)
	if err != nil {
		return ctrl.Result{}, err
	}

	svc, err := r.desiredService(kind)
	if err != nil {
		return ctrl.Result{}, err
	}

	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("kind-controller")}

	err = r.Patch(ctx, cm, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Patch(ctx, deployment, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Patch(ctx, svc, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	kind.Status.Phase = msaberdevv1.RunningStatusPhase

	err = r.Status().Update(ctx, &kind)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled kind")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnekindReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&msaberdevv1.Onekind{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&apps.Deployment{}).
		// Watches(
		// 	&source.Kind{Type: &msaberdevv1.Onekind{}},
		// 	&handler.EnqueueRequestForObject{},
		// ).
		Complete(r)
}
