/*
Copyright 2018 The Knative Authors

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
	"github.com/knative/pkg/apis/duck"
	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Important: Run "make generate" to regenerate code after modifying this file
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GcpPubSubSource is the Schema for the gcppubsubsources API.
// +k8s:openapi-gen=true
// +kubebuilder:categories=all,knative,eventing,sources
type GcpPubSubSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GcpPubSubSourceSpec   `json:"spec,omitempty"`
	Status GcpPubSubSourceStatus `json:"status,omitempty"`
}

// Check that GcpPubSubSource can be validated and can be defaulted.
var _ runtime.Object = (*GcpPubSubSource)(nil)

// Check that GcpPubSubSource implements the Conditions duck type.
var _ = duck.VerifyType(&GcpPubSubSource{}, &duckv1alpha1.Conditions{})

// GcpPubSubSourceSpec defines the desired state of the GcpPubSubSource.
type GcpPubSubSourceSpec struct {
	// GcpCredsSecret is the name of a K8s Secret in the same namespace as this GcpPubSubSource that
	// contains the credentials used to authenticate with the GCP PubSub service. The credentials
	// must be in the entry named 'key.json' and be in the JSON Service Account format.
	GcpCredsSecret string `json:"gcpCredsSecret,omitempty`

	// GoogleCloudProject is the ID of the Google Cloud Project that the PubSub Topic exists in.
	GoogleCloudProject string `json:"googleCloudProject,omitempty"`

	// Topic is the ID of the GCP PubSub Topic to Subscribe to. It must be in the form of the
	// unique identifier within the project, not the entire name. E.g. it must be 'laconia', not
	// 'projects/my-gcp-project/topics/laconia'.
	Topic string `json:"topic,omitempty"`

	// Sink is a reference to an object that will resolve to a domain name to use as the sink.
	// +optional
	Sink *corev1.ObjectReference `json:"sink,omitempty"`
}

const (
	// GcpPubSubSourceConditionReady has status True when the GcpPubSubSource is ready to send events.
	GcpPubSubSourceConditionReady = duckv1alpha1.ConditionReady

	// GcpPubSubSourceConditionSinkProvided has status True when the GcpPubSubSource has been configured with a sink target.
	GcpPubSubSourceConditionSinkProvided duckv1alpha1.ConditionType = "SinkProvided"

	// GcpPubSubSourceConditionDeployed has status True when the GcpPubSubSource has had it's receive adapter deployment created.
	GcpPubSubSourceConditionDeployed duckv1alpha1.ConditionType = "Deployed"

	// GcpPubSubSourceConditionSubscribed has status True when a GCP PubSub Subscription has been created pointing at the created receive adapter deployment.
	GcpPubSubSourceConditionSubscribed duckv1alpha1.ConditionType = "Subscribed"
)

var gcpPubSubSourceCondSet = duckv1alpha1.NewLivingConditionSet(
	GcpPubSubSourceConditionSinkProvided,
	GcpPubSubSourceConditionDeployed,
	GcpPubSubSourceConditionSubscribed)

// GcpPubSubSourceStatus defines the observed state of GcpPubSubSource.
type GcpPubSubSourceStatus struct {
	// Conditions holds the state of a source at a point in time.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions duckv1alpha1.Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// SinkURI is the current active sink URI that has been configured for the GcpPubSubSource.
	// +optional
	SinkURI string `json:"sinkUri,omitempty"`
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *GcpPubSubSourceStatus) GetCondition(t duckv1alpha1.ConditionType) *duckv1alpha1.Condition {
	return gcpPubSubSourceCondSet.Manage(s).GetCondition(t)
}

// IsReady returns true if the resource is ready overall.
func (s *GcpPubSubSourceStatus) IsReady() bool {
	return gcpPubSubSourceCondSet.Manage(s).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *GcpPubSubSourceStatus) InitializeConditions() {
	gcpPubSubSourceCondSet.Manage(s).InitializeConditions()
}

// MarkSink sets the condition that the source has a sink configured.
func (s *GcpPubSubSourceStatus) MarkSink(uri string) {
	s.SinkURI = uri
	if len(uri) > 0 {
		gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubSourceConditionSinkProvided)
	} else {
		gcpPubSubSourceCondSet.Manage(s).MarkUnknown(GcpPubSubSourceConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *GcpPubSubSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkFalse(GcpPubSubSourceConditionSinkProvided, reason, messageFormat, messageA...)
}

// MarkDeployed sets the condition that the source has been deployed.
func (s *GcpPubSubSourceStatus) MarkDeployed() {
	gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubSourceConditionDeployed)
}

// MarkDeploying sets the condition that the source is deploying.
func (s *GcpPubSubSourceStatus) MarkDeploying(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkUnknown(GcpPubSubSourceConditionDeployed, reason, messageFormat, messageA...)
}

// MarkNotDeployed sets the condition that the source has not been deployed.
func (s *GcpPubSubSourceStatus) MarkNotDeployed(reason, messageFormat string, messageA ...interface{}) {
	gcpPubSubSourceCondSet.Manage(s).MarkFalse(GcpPubSubSourceConditionDeployed, reason, messageFormat, messageA...)
}

func (s *GcpPubSubSourceStatus) MarkSubscribed() {
	gcpPubSubSourceCondSet.Manage(s).MarkTrue(GcpPubSubSourceConditionSubscribed)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GcpPubSubSourceList contains a list of GcpPubSubSources.
type GcpPubSubSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GcpPubSubSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GcpPubSubSource{}, &GcpPubSubSourceList{})
}
