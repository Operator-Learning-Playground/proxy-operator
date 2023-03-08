package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Proxy
type Proxy struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProxySpec `json:"spec,omitempty"`
}

type Rules struct {
	Path Path `json:"path"`
}

type Path struct {
	Backend Backend `json:"backend"`
}

type Backend struct {
	Prefix string `json:"prefix"`
	Url    string `json:"url"`
}

type Server struct {
	Ip   string `json:"ip"`
	Port int 	`json:"port"`
}


type ProxySpec struct {
	Rules  []Rules `json:"rules"`
	Server Server  `json:"server"`
}



// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxyList
type ProxyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Proxy `json:"items"`
}


