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
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tgikv1 "github.com/n3wscott/multiply/api/v1"
)

// SquareReconciler reconciles a Square object
type SquareReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tgik.tgik.io,resources=squares,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tgik.tgik.io,resources=squares/status,verbs=get;update;patch

func (r *SquareReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.WithValues("square", req.NamespacedName)

	obj := &tgikv1.Square{}

	if err := r.Client.Get(context.Background(), req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, err
	}

	obj.Status.Result = obj.Spec.Base * obj.Spec.Base
	obj.Status.Expression = fmt.Sprintf("%d^2", obj.Spec.Base)

	if err := r.Client.Status().Update(context.Background(), obj); err != nil {
		l.Info("unable to update status:", "err", err)
		return ctrl.Result{}, err
	}

	l.Info(fmt.Sprintf("--> updated status: %+v", obj))

	return ctrl.Result{}, nil
}

func (r *SquareReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tgikv1.Square{}).
		Complete(r)
}
