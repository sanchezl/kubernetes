package deprecatedapirequest

import "k8s.io/apimachinery/pkg/runtime/schema"

var deprecatedApiRemovedRelease = map[schema.GroupVersionResource]string{
	schema.GroupVersionResource{Group: "flowcontrol.apiserver.k8s.io", Version: "v1alpha1", Resource: "flowschemas"}:                    "1.21",
	schema.GroupVersionResource{Group: "flowcontrol.apiserver.k8s.io", Version: "v1alpha1", Resource: "prioritylevelconfigurations"}:    "1.21",
	schema.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "ingresses"}:                                         "1.22",
	schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1beta1", Resource: "validatingwebhookconfigurations"}: "1.22",
	schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1beta1", Resource: "customresourcedefinitions"}:               "1.22",
	schema.GroupVersionResource{Group: "admissionregistration.k8s.io", Version: "v1beta1", Resource: "mutatingwebhookconfigurations"}:   "1.22",
	schema.GroupVersionResource{Group: "certificates.k8s.io", Version: "v1beta1", Resource: "certificatesigningrequests"}:               "1.22",
	schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1beta1", Resource: "ingresses"}:                                  "1.22",
	schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1beta1", Resource: "clusterrolebindings"}:                "1.22",
	schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1beta1", Resource: "rolebindings"}:                       "1.22",
	schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1beta1", Resource: "roles"}:                              "1.22",
}
