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

package v1

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HandpaySpec defines the desired state of Handpay
// handpay kind Spec 定义参数 列表
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
}

// HandpayStatus defines the observed state of Handpay
type HandpayStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Handpay is the Schema for the handpays API
type Handpay struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HandpaySpec   `json:"spec,omitempty"`
	Status HandpayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HandpayList contains a list of Handpay
type HandpayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Handpay `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Handpay{}, &HandpayList{})
}
