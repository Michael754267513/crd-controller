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

	appsv1 "crd-controller/api/v1"
	"crd-controller/controllers/logic"
	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	indexOwnerKey = ".metadata.controller"
	apiGVStr      = appsv1.GroupVersion.String()
)

// HandpayReconciler reconciles a Handpay object
type HandpayReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.zenghao.com,resources=handpays,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.zenghao.com,resources=handpays/status,verbs=get;update;patch
// +kubebuilder:docs-gen:collapse=Imports

func (r *HandpayReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	/*
		控制器逻辑处理段，根据r.Get r.List r.Create r.Update r.Delete 对资源的增删改查（传送进来的资源进行逻辑处理）
	*/
	// 逻辑处理段
	ctx := context.Background()
	_ = r.Log.WithValues("handpay", req.NamespacedName)
	meta := &appsv1.Handpay{}
	// 判断是否存在
	if err := r.Get(ctx, req.NamespacedName, meta); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return ctrl.Result{}, nil
	}
	// Finalizer 意异步删除数据
	myFinalizerName := "zenghao.handpay.com.cn"
	if meta.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(meta.ObjectMeta.Finalizers, myFinalizerName) {
			meta.ObjectMeta.Finalizers = append(meta.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), meta); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(meta.ObjectMeta.Finalizers, myFinalizerName) {
			if err := r.deleteExternalResources(meta); err != nil {
				return ctrl.Result{}, err
			}
			meta.ObjectMeta.Finalizers = removeString(meta.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), meta); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	// 判断是否是新建或者更新
	if meta.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("创建deployment")
		// 创建或者更新 deployment
		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, logic.ServiceMetaLogic(meta.Spec), func() error {
			return nil
		}); err != nil {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	// 判断是否删除
	if meta.ObjectMeta.DeletionTimestamp != nil {
		log.Info("删除deployment")
		if err := r.Delete(ctx, logic.ServiceMetaLogic(meta.Spec)); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *HandpayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&appsv1.Handpay{}, indexOwnerKey, func(rawObj runtime.Object) []string {
		deployment := rawObj.(*appsv1.Handpay)
		owner := metav1.GetControllerOf(deployment)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "Handpay" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Handpay{}).
		Owns(&appsv1.Handpay{}).
		Complete(r)
}

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
func (r *HandpayReconciler) deleteExternalResources(META *appsv1.Handpay) error {
	//
	// delete any external resources associated with the cronJob
	//
	// Ensure that delete implementation is idempotent and safe to invoke
	// multiple types for same object.
	return nil
}
