package logic

import (
	v1 "crd-controller/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetLables(meta v1.HandpaySpec) map[string]string {
	lables := map[string]string{}
	lables["project"] = meta.Project
	lables["owner"] = meta.Owner
	lables["serviceName"] = meta.ServiceName
	return lables
}

func ServiceMetaLogic(meta v1.HandpaySpec) *appsv1.Deployment {
	//测试环境所有公共服务副本数固定是1
	var replicas int32 = 1
	// 自定义lable标签
	lables := GetLables(meta)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.ServiceName,
			Namespace: meta.Env,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: lables,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: lables,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  meta.ServiceName,
							Image: meta.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          meta.ServiceName,
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: meta.Port,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}
