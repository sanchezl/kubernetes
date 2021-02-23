// Package v1 is an api version in the apiserver.openshift.io group
package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

// DeprecatedAPIRequest tracts requests made to a deprecated API. The instance name should
// be of the form `resource.version.group`, matching the deprecated resource.
type DeprecatedAPIRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the characteristics of the resource.
	// +kubebuilder:validation:Required
	// +required
	Spec DeprecatedAPIRequestSpec `json:"spec"`

	// Status contains the observed state of the resource.
	Status DeprecatedAPIRequestStatus `json:"status,omitempty"`
}

type DeprecatedAPIRequestSpec struct {
	// RemovedRelease is when the API will be removed.
	// +kubebuilder:validation:MaxLength=15
	// +required
	RemovedRelease string `json:"removedRelease"`
}

// +k8s:deepcopy-gen=true
type DeprecatedAPIRequestStatus struct {

	// Conditions contains details of the current status of this API Resource.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions"`

	// RequestsLastHour contains request history for the current hour. This is porcelain to make the API
	// easier to read by humans seeing if they addressed a problem.
	RequestsLastHour RequestLog `json:"requestsLastHour"`

	// RequestsLast24h contains request history for the last 24 hours, indexed by the hour, so
	// 12:00AM-12:59 is in index 0, 6am-6:59am is index 6, etc..
	RequestsLast24h []RequestLog `json:"requestsLast24h"`
}

// RequestLog logs request for various nodes.
type RequestLog struct {

	// Nodes contains logs of requests per node.
	Nodes []NodeRequestLog `json:"nodes"`
}

// NodeRequestLog contains logs of requests to a certain node.
type NodeRequestLog struct {

	// NodeName where the request are being handled.
	NodeName string `json:"nodeName"`

	// LastUpdate should *always* being within the hour this is for.  This is a time indicating
	// the last moment the server is recording for, not the actual update time.
	LastUpdate metav1.Time `json:"lastUpdate"`

	// Users contains request details by top 10 users.
	Users []RequestUser `json:"users"`
}

type DeprecatedAPIRequestConditionType string

const (
	// UsedInPastDay condition indicates a request has been made against the deprecated api in the last 24h.
	UsedInPastDay DeprecatedAPIRequestConditionType = "UsedInPastDay"
)

// RequestUser contains logs of a user's requests.
type RequestUser struct {

	// UserName that made the request.
	UserName string `json:"username"`

	// Count of requests.
	Count int `json:"count"`

	// Requests details by verb.
	Requests []RequestCount `json:"requests"`
}

// RequestCount counts requests by API request verb.
type RequestCount struct {

	// Verb of API request (get, list, create, etc...)
	Verb string `json:"verb"`

	// Count of requests for verb.
	Count int `json:"count"`
}

// DeprecatedAPIRequestList is a list of DeprecatedAPIRequest resources.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DeprecatedAPIRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []DeprecatedAPIRequest `json:"items"`
}
