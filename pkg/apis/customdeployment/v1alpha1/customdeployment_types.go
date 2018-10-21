package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CustomDeploymentSpec defines the desired state of CustomDeployment
type CustomDeploymentSpec struct {
	Image    string `json:"image"`
	Version  string `json:"version"`
	Replicas int32  `json:"replicas"`
	// CpuRequest    string `json:"cpuRequest"`
	// MemoryRequest string `json:"memoryRequest"`
}

// CustomDeploymentStatus defines the observed state of CustomDeployment
type CustomDeploymentStatus struct {
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomDeployment is the Schema for the customdeployments API
// +k8s:openapi-gen=true
type CustomDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomDeploymentSpec   `json:"spec,omitempty"`
	Status CustomDeploymentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomDeploymentList contains a list of CustomDeployment
type CustomDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomDeployment{}, &CustomDeploymentList{})
}
