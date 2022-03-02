package metrics

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/config/apis/webhookadmission"
)

func resourceLabels(rules []webhookadmission.Rule, attr admission.Attributes) []string {
	if len(rules) == 0 {
		return []string{"", "", "", ""}
	}
	gvr := attr.GetResource()
	subresource := attr.GetSubresource()
	ns := attr.GetNamespace()
	if rulesMatch(rules, gvr, ns) {
		return []string{gvr.Group, gvr.Version, gvr.Resource, subresource}
	}
	return []string{"", "", "", ""}
}

func rulesMatch(filters []webhookadmission.Rule, gvr schema.GroupVersionResource, ns string) bool {
	for _, rule := range filters {
		if ruleMatches(rule, gvr, ns) {
			return true
		}
	}
	return false
}

func ruleMatches(rule webhookadmission.Rule, gvr schema.GroupVersionResource, ns string) bool {
	// namespace
	namespace := len(rule.Namespaces) == 0
	for _, n := range rule.Namespaces {
		if n == "*" || n == ns {
			namespace = true
			break
		}
	}
	// group
	group := len(rule.Groups) == 0
	for _, g := range rule.Groups {
		if g == "*" || g == gvr.Group {
			group = true
			break
		}
	}
	// version
	version := len(rule.Versions) == 0
	for _, v := range rule.Versions {
		if v == "*" || v == gvr.Version {
			version = true
			break
		}
	}
	// resource
	resource := len(rule.Resources) == 0
	for _, rr := range rule.Resources {
		// ignore sub-resource, match only resource
		segments := strings.SplitN(rr, "/", 2)
		r := segments[0]
		if r == "*" || r == gvr.Resource {
			resource = true
		}
	}
	return namespace && group && version && resource
}
