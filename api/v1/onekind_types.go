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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OnekindSpec defines the desired state of Onekind
type OnekindSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Onekind. Edit onekind_types.go to remove/update
	// Foo      string `json:"foo,omitempty"`

	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=4
	Replicas int32 `json:"replicas"`
	//+kubebuilder:validation:Required
	Message string `json:"message"`
}

type StatusPhase string

const (
	RunningStatusPhase StatusPhase = "RUNNING"
	PendingStauePhase  StatusPhase = "PENDING"
	ErrorStatusPhase   StatusPhase = "ERROR"
)

// OnekindStatus defines the observed state of Onekind
type OnekindStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase StatusPhase `json:"phase"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`

// Onekind is the Schema for the onekinds API
type Onekind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OnekindSpec   `json:"spec,omitempty"`
	Status OnekindStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OnekindList contains a list of Onekind
type OnekindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Onekind `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Onekind{}, &OnekindList{})
}
