/*
Copyright 2020 The Kube Diagnoser Authors.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	diagnosisv1 "netease.com/k8s/kube-diagnoser/api/v1"
)

// AbnormalReconciler reconciles a Abnormal object
type AbnormalReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	abnormalSourceCh chan diagnosisv1.Abnormal
}

func NewAbnormalReconciler(
	cli client.Client,
	log logr.Logger,
	scheme *runtime.Scheme,
	abnormalSourceCh chan diagnosisv1.Abnormal,
) *AbnormalReconciler {
	return &AbnormalReconciler{
		Client:           cli,
		Log:              log,
		Scheme:           scheme,
		abnormalSourceCh: abnormalSourceCh,
	}
}

// +kubebuilder:rbac:groups=diagnosis.netease.com,resources=abnormals,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=diagnosis.netease.com,resources=abnormals/status,verbs=get;update;patch

func (r *AbnormalReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("abnormal", req.NamespacedName)

	var abnormal diagnosisv1.Abnormal
	if err := r.Get(ctx, req.NamespacedName, &abnormal); err != nil {
		log.Error(err, "unable to fetch Abnormal")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.abnormalSourceCh <- abnormal

	return ctrl.Result{}, nil
}

func (r *AbnormalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&diagnosisv1.Abnormal{}).
		Complete(r)
}