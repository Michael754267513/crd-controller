package logic

import (
	v1 "crd-controller/api/v1"
	"github.com/prometheus/common/log"
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

func AddEnvString(name string, value string) apiv1.EnvVar {
	var env apiv1.EnvVar
	env.Name = name
	env.Value = value
	return env
}
func AddPodNameEnv() apiv1.EnvVar {
	var env apiv1.EnvVar
	fieldref := apiv1.ObjectFieldSelector{}
	fieldref.APIVersion = "v1"
	fieldref.FieldPath = "metadata.name"
	field := apiv1.EnvVarSource{FieldRef: &fieldref}
	env.Name = "POD_NAME"
	env.ValueFrom = &field
	return env
}
func ServiceMetaLogic(meta v1.HandpaySpec, namespace string) *appsv1.Deployment {
	//测试环境所有公共服务副本数固定是1
	var replicas int32 = 1
	var env []apiv1.EnvVar
	if meta.Replicas != 0 {
		replicas = meta.Replicas
	}
	log.Info("默认值：", meta.Replicas)
	// 自定义lable标签
	lables := GetLables(meta)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.ServiceName,
			Namespace: namespace,
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
	// 添加hosts解析
	deployment.Spec.Template.Spec.HostAliases = meta.Hosts
	// 添加容器环境变量
	env = append(env, AddEnvString("LANG", "en_US.UTF-8"))
	env = append(env, AddPodNameEnv())
	deployment.Spec.Template.Spec.Containers[0].Env = env
	return deployment
}
