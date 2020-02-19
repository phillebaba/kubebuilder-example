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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RemoteContentSpec defines the desired state of RemoteContent
type RemoteContentSpec struct {
	// Url to fetch content from.
	Url string `json:"url"`
}

// Indicates the status of the request
type RequestState string

var (
	Pending  RequestState = "Pending"
	Failed   RequestState = "Failed"
	Succeded RequestState = "Succeeded"
)

// RemoteContentStatus defines the observed state of RemoteContent
type RemoteContentStatus struct {
	State RequestState `json:"state"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// RemoteContent is the Schema for the remotecontents API
type RemoteContent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RemoteContentSpec   `json:"spec,omitempty"`
	Status RemoteContentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RemoteContentList contains a list of RemoteContent
type RemoteContentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RemoteContent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RemoteContent{}, &RemoteContentList{})
}
