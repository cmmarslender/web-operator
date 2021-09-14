/*
Copyright 2021.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SimpleAppSpec defines the desired state of SimpleApp
type SimpleAppSpec struct {
	// Image is the container image to deploy
	Image string `json:"image,omitempty"`

	// ImagePullPolicy describes a policy for if/when to pull a container image
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets names of the secrets with image pull credentials
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// @TODO LivenessProbe
	// @TODO ReadinessProbe

	// ContainerPort is the port the container is set to listen on
	// +kubebuilder:default:=80
	ContainerPort int32 `json:"containerPort,omitempty"`

	// Replicas how many replicas in the deployment
	// +kubebuilder:default:=1
	Replicas *int32 `json:"replicas,omitempty"`

	// ServiceEnabled sets whether an ingress should be enabled
	// +kubebuilder:default:=true
	ServiceEnabled bool `json:"serviceEnabled,omitempty"`

	// ServicePort is the port the service will listen on
	// traffic will be forwarded at the service to the ContainerPort
	// +kubebuilder:default:=80
	ServicePort int32 `json:"servicePort,omitempty"`

	// IngressEnabled sets whether an ingress should be enabled
	// +kubebuilder:default:=true
	IngressEnabled bool `json:"ingressEnabled,omitempty"`

	// Hostname is the hostname to use for the Ingress
	Hostname string `json:"hostname,omitempty"`

	// IngressPaths are the paths the ingress will serve traffic on
	// The default below looks like an object, but it's actually an array in kubebuilder syntax
	// +kubebuilder:default:={"/"}
	IngressPaths []string `json:"ingressPaths,omitempty"`

	// IngressAnnotations map of annotations that should be added to an ingress
	// If a key is present in this, it will override the global ingress annotation with the same key
	IngressAnnotations map[string]string `json:"ingressAnnotations,omitempty"`
}

// SimpleAppStatus defines the observed state of SimpleApp
type SimpleAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SimpleApp is the Schema for the simpleapps API
type SimpleApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SimpleAppSpec   `json:"spec,omitempty"`
	Status SimpleAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SimpleAppList contains a list of SimpleApp
type SimpleAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SimpleApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SimpleApp{}, &SimpleAppList{})
}
