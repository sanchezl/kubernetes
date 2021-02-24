package v1helpers

import (
	"context"
	"time"

	apiv1 "github.com/openshift/api/apiserver/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
)

type DeprecatedAPIRequestClient interface {
	Get(name string) (*apiv1.DeprecatedAPIRequest, error)
	UpdateStatus(ctx context.Context, podNetworkConnectivityCheck *apiv1.DeprecatedAPIRequest, opts metav1.UpdateOptions) (*apiv1.DeprecatedAPIRequest, error)
}

type UpdateStatusFunc func(status *apiv1.DeprecatedAPIRequestStatus)

func UpdateStatus(ctx context.Context, client DeprecatedAPIRequestClient, name string, updateFuncs ...UpdateStatusFunc) (*apiv1.DeprecatedAPIRequestStatus, bool, error) {
	updated := false
	var updatedStatus *apiv1.DeprecatedAPIRequestStatus
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		check, err := client.Get(name)
		if err != nil {
			return err
		}
		oldStatus := check.Status
		newStatus := oldStatus.DeepCopy()
		for _, update := range updateFuncs {
			update(newStatus)
		}
		if equality.Semantic.DeepEqual(oldStatus, newStatus) {
			updatedStatus = newStatus
			return nil
		}
		check, err = client.Get(name)
		if err != nil {
			return err
		}
		check.Status = *newStatus
		updatedCheck, err := client.UpdateStatus(ctx, check, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		updatedStatus = &updatedCheck.Status
		updated = true
		return err
	})
	return updatedStatus, updated, err
}

func AppendRequestLog(node string, timestamp time.Time, gvr schema.GroupVersionResource, verb, username, release string) UpdateStatusFunc {
	return func(status *apiv1.DeprecatedAPIRequestStatus) {
		if timestamp.Before(time.Now().AddDate(0, 0, -1)) {
			klog.V(6).Infof("Ignoring request older that 24 hours")
			return
		}
		rotateNodeRequestLogs(status, node, time.Now())
		nodeRequestLog := findNodeRequestLog(status, node, timestamp)
		addRequestToNodeLog(nodeRequestLog, timestamp, gvr, verb, username)
	}
}

func addRequestToNodeLog(log *apiv1.NodeRequestLog, timestamp time.Time, gvr schema.GroupVersionResource, verb string, username string) {
	// TODO find or add user
	// TODO find or add verb
	// TODO request.Count++
	// TODO user.Count++
	if log.LastUpdate.Time.Before(timestamp) {
		log.LastUpdate = metav1.NewTime(timestamp)
	}
}

func findNodeRequestLog(status *apiv1.DeprecatedAPIRequestStatus, node string, timestamp time.Time) *apiv1.NodeRequestLog {
	_, log := nodeRequestLogForNode(status.RequestsLastHour.Nodes, node)
	if log == nil {
		_, log = nodeRequestLogForNode(status.RequestsLast24h[timestamp.Hour()].Nodes, node)
	}
	return log
}

func rotateNodeRequestLogs(status *apiv1.DeprecatedAPIRequestStatus, nodeName string, currentTime time.Time) {

	currentHour := currentTime.Hour()

	if currentHour > 23 || currentHour < 0 {
		klog.V(6).Infof("0 <= currentHour <= 23: %v", currentHour)
		return
	}

	activeLogs := &status.RequestsLastHour.Nodes
	activeLogIndex, activeLog := nodeRequestLogForNode(*activeLogs, nodeName)
	if activeLog == nil {
		klog.V(6).Infof("requestsLastHour does not contain request logs for node %s", nodeName)
		return
	}

	// archive the requestsLastHour node log to the requestsLast24h node log
	activeLogHour := activeLog.LastUpdate.Hour()
	if activeLogHour != currentHour {

		// ensure the hour index exists
		for len(status.RequestsLast24h) <= activeLogHour {
			status.RequestsLast24h = append(status.RequestsLast24h, apiv1.RequestLog{})
		}

		// find or create the node log for the node
		archivedLogs := &status.RequestsLast24h[activeLogHour].Nodes
		archivedLogIndex, _ := nodeRequestLogForNode(*archivedLogs, nodeName)
		switch archivedLogIndex {
		case -1:
			*archivedLogs = append(*archivedLogs, *activeLog)
		default:
			(*archivedLogs)[archivedLogIndex] = *activeLog
		}

		// delete from active log
		*activeLogs = append((*activeLogs)[:activeLogIndex], (*activeLogs)[activeLogIndex+1:]...)
	}
}

func nodeRequestLogForNode(logs []apiv1.NodeRequestLog, node string) (int, *apiv1.NodeRequestLog) {
	for i, log := range logs {
		if log.NodeName == node {
			return i, &log
		}
	}
	return -1, nil
}
