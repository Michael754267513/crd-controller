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
	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HandpayReconciler reconciles a Handpay object
type HandpayReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.zenghao.com,resources=handpays,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.zenghao.com,resources=handpays/status,verbs=get;update;patch

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
		log.Error(err)
	} else {
		//获取参数
		log.Info("名称：", req.Name, "镜像名称: ", meta.Spec.Image, "副本数: ", meta.Spec.Replicas)
	}
	// 遍历筛选自定义kind
	handpayList := &appsv1.HandpayList{}
	//  根据req.Namespace 遍历，client.MatchingFields{} 根据lable字段筛选； 可参考kubebuilder 添加相关筛选字段
	if err := r.List(ctx, handpayList, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err)
	} else {
		for _, v := range handpayList.Items {
			log.Info(v)
		}
	}
	//  r.Delete()  r.Update()  r.Create() 可以对资源进行增 删 改，根据自定crd的参数进行操作

	return ctrl.Result{}, nil
}

func (r *HandpayReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Handpay{}).
		Complete(r)
}
