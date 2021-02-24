package deprecatedapirequest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	v1 "github.com/openshift/api/apiserver/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/mergepatch"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
)

type deprecateApiRequestHandler struct {
}

// TODO logs levels

func WithDeprecatedApiRequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		klog.V(1).Infof("[#####] Request: %-7s %v", req.Method, req.RequestURI)

		var debugType string
		debugLogged := &struct {
			b bool
		}{}
		for k := range deprecatedApiRemovedRelease {
			if strings.Contains(req.RequestURI, k.Resource) &&
				strings.Contains(req.RequestURI, k.Version) &&
				strings.Contains(req.RequestURI, k.Group) {
				debugType = k.Resource + "." + k.Version + "." + k.Group
			}
		}
		if len(debugType) > 0 {
			klog.V(1).Infof("[#####] BEGIN %s", debugType)
			defer func(b *struct{ b bool }) { klog.V(1).Infof("[#####] END   %s logged==%v", debugType, b.b) }(debugLogged)
		}

		audit := request.AuditEventFrom(req.Context())
		if audit == nil {
			klog.V(2).Info("No audit entry found in context.")
			next.ServeHTTP(w, req)
			return
		}

		klog.V(1).Infof("[#####] Audit: %v", audit)

		if audit.ObjectRef == nil {
			klog.V(2).Info("No objectRef found in audit context.")
			next.ServeHTTP(w, req)
			return
		}

		gvr := schema.GroupVersionResource{
			Group:    audit.ObjectRef.APIGroup,
			Version:  audit.ObjectRef.APIVersion,
			Resource: audit.ObjectRef.Resource,
		}

		verb := audit.Verb

		username := audit.User.Username
		if len(username) == 0 {
			if user, ok := request.UserFrom(req.Context()); ok {
				username = user.GetName()
			}
		}

		received := audit.RequestReceivedTimestamp

		removedRelease, deprecated := deprecatedApiRemovedRelease[gvr]
		if !deprecated {
			next.ServeHTTP(w, req)
			return
		}

		if err := todoLogDeprecatedAPIRequest(received.Time, gvr, verb, username, removedRelease); err != nil {
			runtime.HandleError(fmt.Errorf("error logging deprecated API request: %v", err))
		}
		klog.V(1).Infof("[#####] LOG   %s", debugType)
		debugLogged.b = true

		// continue to next handler
		next.ServeHTTP(w, req)

	})
}

func todoLogDeprecatedAPIRequest(timestamp time.Time, gvr schema.GroupVersionResource, verb string, username string, release string) error {

	r := &v1.DeprecatedAPIRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: gvr.Resource + "." + gvr.Version + "." + gvr.Group,
		},
		Spec: v1.DeprecatedAPIRequestSpec{
			RemovedRelease: release,
		},
		Status: v1.DeprecatedAPIRequestStatus{
			Conditions: []metav1.Condition{
				{
					Type:    string(v1.UsedInPastDay),
					Status:  metav1.ConditionTrue,
					Reason:  "Request",
					Message: "Yeah",
				},
			},
			RequestsLastHour: v1.RequestLog{
				Nodes: []v1.NodeRequestLog{
					{NodeName: "nodeanme", LastUpdate: metav1.Now(), Users: []v1.RequestUser{
						{
							UserName: username,
							Count:    1,
							Requests: []v1.RequestCount{
								{Verb: verb, Count: 1},
							},
						},
					}},
				},
			},
			RequestsLast24h: nil,
		},
	}

	klog.V(1).Infof("[#####] DeprecatedAPIRequest: %v", mergepatch.ToYAMLOrError(r))

	return nil
}
