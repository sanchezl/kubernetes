package v1helpers

import (
	"testing"
	"time"

	apiv1 "github.com/openshift/api/apiserver/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRotateNodeRequestLogs(t *testing.T) {

	status := deprecatedApiRequestStatus(
		withActiveNodeRequestLog("test", testTime(0, 3, 28)),
	)

	expected := deprecatedApiRequestStatus(
		withArchiveNodeRequestLog("test", testTime(0, 3, 28)),
	)

	rotateNodeRequestLogs(status, "test", currentTime())

	assert.Equal(t, *expected, *status)

}

func deprecatedApiRequestStatus(options ...func(status *apiv1.DeprecatedAPIRequestStatus)) *apiv1.DeprecatedAPIRequestStatus {
	status := &apiv1.DeprecatedAPIRequestStatus{RequestsLastHour: apiv1.RequestLog{Nodes: []apiv1.NodeRequestLog{}}}
	for _, option := range options {
		option(status)
	}
	return status
}

func withActiveNodeRequestLog(node string, lastUpdate metav1.Time, options ...func(log *apiv1.NodeRequestLog)) func(status *apiv1.DeprecatedAPIRequestStatus) {
	return func(status *apiv1.DeprecatedAPIRequestStatus) {
		n := &apiv1.NodeRequestLog{
			NodeName:   node,
			LastUpdate: lastUpdate,
		}
		for _, option := range options {
			option(n)
		}
		status.RequestsLastHour.Nodes = append(status.RequestsLastHour.Nodes, *n)
	}
}

func withArchiveNodeRequestLog(node string, lastUpdate metav1.Time, options ...func(log *apiv1.NodeRequestLog)) func(status *apiv1.DeprecatedAPIRequestStatus) {
	return func(status *apiv1.DeprecatedAPIRequestStatus) {
		n := &apiv1.NodeRequestLog{
			NodeName:   node,
			LastUpdate: lastUpdate,
		}
		for _, option := range options {
			option(n)
		}
		for len(status.RequestsLast24h) <= lastUpdate.Hour() {
			status.RequestsLast24h = append(status.RequestsLast24h, apiv1.RequestLog{})
		}
		status.RequestsLast24h[lastUpdate.Hour()].Nodes = append(status.RequestsLast24h[lastUpdate.Hour()].Nodes, *n)
	}
}

func testTime(day, hour, min int) metav1.Time {
	return metav1.NewTime(time.Date(1974, 8, 18+day, hour, min, 0, 0, time.UTC))
}

func currentTime() time.Time {
	return testTime(0, 14, 0).Time
}
