# 自定义控制器
利用kubebuilder生成基础代码，用户只需要修改types定义参数文件和修改controller的逻辑代码实现块
来进行快速的逻辑控制

###默认环境变量

```cassandraql
("LANG", "en_US.UTF-8")
("POD_NAME", "name")
("POD_NAMESPACE", "namespace")
```
### 应用日志存放路径
```cassandraql
默认存放宿主机日志目录： /opt/logs/{POD_NAMESPACE}/{POD_DEPLOYMENT}/{POD_NAME}
容器应用日志目录：/logs
可以使用参数进行调整： LogDir容器日志目录 NodeLogDir宿主机日志存放目录(/opt/logs/)
```

### 使用 
```cassandraql
kubecte apply -f crd.yaml
go clone https://github.com/Michael754267513/crd-controller.git
cd crd-controller
go build main.go #进行测试debug 调试

```

###新增配置字段
```cassandraql
路径： api/v1/handpay_types.go
```
```
type HandpaySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Handpay. Edit Handpay_types.go to remove/update
	// 添加 所需要的yaml 变量
	Project     string            `json:"project,omitempty"`
	Owner       string            `json:"owner,omitempty"`
	ServiceName string            `json:"serviceName,omitempty"`
	Image       string            `json:"image"`
	Port        int32             `json:"port"`
	Hosts       []apiv1.HostAlias `json:"hosts"`
	Replicas    int32             `json:"replicas"`
	LogDir      string            `json:"logDir"`
	NodeLogDir  string            `json:"nodeLogDir"`
	PodEnv      []apiv1.EnvVar    `json:"podEnv"`
}
```
### 逻辑实现代码段
```cassandraql
路径: controllers/handpay_contriller.go
```
```cassandraql
func (r *HandpayReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {}
```
### 新建和更新逻辑
```cassandraql
判断 ObjectMeta.DeletionTimestamp.IsZero是否为空，当新建或者更新的时候该字段为空
判断新建方式：
    判断DeletionTimestamp为空的前提下
    r.GET 获取改资源是否存在当返回errors.IsNotFound 表示资源不存在
    r.CreateOrUpdate 进行资源的新建
资源更新方式：(当编辑yaml 发生变化的时候才会触发更新)
    判断DeletionTimestamp为空的前提下
    r.GET 获取改资源是否存在
    当资源存在的情况下使用!reflect.DeepEqual 来判断当前资源和更新的资源是否一致
    不一致进行r.Update资源更新

```
###资源删除
```cassandraql
利用Finalizer 延迟删除 可以记录删除资源
删除时候DeletionTimestamp不为空 
在资源新建的时候使用ctrl.SetControllerReference关联资源，利用k8s的垃圾回收机制对资源进行删除
```
### crd.yaml CRD文件
```cassandraql
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # metadata.name的内容是由"复数名.分组名"构成，如下，students是复数名，bolingcavalry.k8s.io是分组名
  name: handpay.apps.zenghao.com
spec:
  # 分组名，在REST API中也会用到的，格式是: /apis/分组名/CRD版本
  group: apps.zenghao.com
  # list of versions supported by this CustomResourceDefinition
  versions:
    - name: v1
      # 是否有效的开关.
      served: true
      # 只有一个版本能被标注为storage
      storage: true
  # 范围是属于namespace的
  scope: Namespaced
  names:
    # 复数名
    plural: handpay
    # 单数名
    singular: handpay
    # 类型名
    kind: Handpay
    # 简称，就像service的简称是svc
    shortNames:
    - handpay
```
### 测试文件 crd_test.yaml
```cassandraql
apiVersion: apps.zenghao.com/v1
kind: Handpay
metadata:
  name: test1
spec:
  project: risk-account
  image: nginx
  port: 80
  serviceName: nginx1
  owner: hzeng
  hosts: [{"hostnames": ["www.baidu.com","www.michael.com"],"ip": "1.1.1.1"},{"hostnames":["databases1"],"ip": "192.168.2.1"}]
```
### 默认参数配置
```cassandraql
	Project     string            `json:"project,omitempty"`
	Owner       string            `json:"owner,omitempty"`
	ServiceName string            `json:"serviceName,omitempty"`
	Image       string            `json:"image"`
	Port        int32             `json:"port"`
	Hosts       []apiv1.HostAlias `json:"hosts"`
	Replicas    int32             `json:"replicas"`
```