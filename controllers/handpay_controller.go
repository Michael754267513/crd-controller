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
	var err error
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
		// 添加OwnerReferences
		if meta.ObjectMeta.OwnerReferences == nil {
			log.Info("添加OwnerReferences")
			meta.ObjectMeta.OwnerReferences = getOwnerReferences(meta)
			if err = r.Update(ctx, meta); err != nil {
				log.Info("更新OwnerReferences错误")
				goto ERROR
			}
		}
		// 添加Finalizer
		if !containsString(meta.ObjectMeta.Finalizers, myFinalizerName) {
			meta.ObjectMeta.Finalizers = append(meta.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), meta); err != nil {
				goto ERROR
			}
		}
	} else {
		if containsString(meta.ObjectMeta.Finalizers, myFinalizerName) {
			if err := r.deleteExternalResources(meta); err != nil {
				goto ERROR
			}
			meta.ObjectMeta.Finalizers = removeString(meta.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(context.Background(), meta); err != nil {
				goto ERROR
			}
		}
	}
	// 判断是否是新建或者更新
	if meta.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Kind 新建")
		// 创建或者更新 deployment
		log.Info("获取deployment")
		deployment := logic.ServiceMetaLogic(meta.Spec)
		if deployment.ObjectMeta.OwnerReferences == nil {
			//关联OwnerReferences
			log.Info("关联OwnerReferences 当kind删除的时候 关联的资源也会自动删除")
			if err = ctrl.SetControllerReference(meta, deployment, r.Scheme); err != nil {
				log.Info("关联错误")
				goto ERROR
			}
		}
		//r.Update(ctx,meta)
		// 创建deployment
		log.Info("新建Deployment")
		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, deployment, func() error {
			return nil
		}); err != nil {
			goto ERROR
		}
		return ctrl.Result{}, nil
	}

	// 判断kind是否删除
	if meta.ObjectMeta.DeletionTimestamp != nil {
		log.Info("删除kind")
	}
ERROR:
	return ctrl.Result{}, err
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

func getOwnerReferences(meta *appsv1.Handpay) []metav1.OwnerReference {
	ownerRefs := []metav1.OwnerReference{}
	ownerRef := metav1.OwnerReference{}
	var k8sGC bool = true
	ownerRef.APIVersion = meta.APIVersion
	ownerRef.Name = meta.Name
	ownerRef.Kind = meta.Kind
	ownerRef.UID = meta.UID
	ownerRef.Controller = &k8sGC
	ownerRef.BlockOwnerDeletion = &k8sGC
	return append(ownerRefs, ownerRef)

}
