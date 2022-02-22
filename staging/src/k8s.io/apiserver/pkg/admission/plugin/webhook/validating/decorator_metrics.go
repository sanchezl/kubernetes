package validating

import (
	"context"
	"time"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionmetrics "k8s.io/apiserver/pkg/admission/metrics"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/generic"
	webhookutil "k8s.io/apiserver/pkg/util/webhook"
)

// MetricsDecorator wraps webhook invocations with metric observations.
var MetricsDecorator = WebhookInvokerDecoratorFunc(func(invoker WebhookInvoker) WebhookInvoker {
	return func(ctx context.Context, hook *admissionregistrationv1.ValidatingWebhook, invocation *generic.WebhookInvocation, attr *generic.VersionedAttributes) error {
		t := time.Now()
		err := invoker(ctx, hook, invocation, attr)
		ignoreClientCallFailures := hook.FailurePolicy != nil && *hook.FailurePolicy == admissionregistrationv1.Ignore
		rejected := false
		if err != nil {
			switch err := err.(type) {
			case *webhookutil.ErrCallingWebhook:
				if !ignoreClientCallFailures {
					rejected = true
					admissionmetrics.Metrics.ObserveWebhookRejection(ctx, hook.Name, "validating", string(attr.Attributes.GetOperation()), admissionmetrics.WebhookRejectionCallingWebhookError, int(err.Status.ErrStatus.Code))
				}
				admissionmetrics.Metrics.ObserveWebhook(ctx, hook.Name, time.Since(t), rejected, attr.Attributes, "validating", int(err.Status.ErrStatus.Code))
			case *webhookutil.ErrWebhookRejection:
				rejected = true
				admissionmetrics.Metrics.ObserveWebhookRejection(ctx, hook.Name, "validating", string(attr.Attributes.GetOperation()), admissionmetrics.WebhookRejectionNoError, int(err.Status.ErrStatus.Code))
				admissionmetrics.Metrics.ObserveWebhook(ctx, hook.Name, time.Since(t), rejected, attr.Attributes, "validating", int(err.Status.ErrStatus.Code))
			default:
				rejected = true
				admissionmetrics.Metrics.ObserveWebhookRejection(ctx, hook.Name, "validating", string(attr.Attributes.GetOperation()), admissionmetrics.WebhookRejectionAPIServerInternalError, 0)
				admissionmetrics.Metrics.ObserveWebhook(ctx, hook.Name, time.Since(t), rejected, attr.Attributes, "validating", 0)
			}
		} else {
			admissionmetrics.Metrics.ObserveWebhook(ctx, hook.Name, time.Since(t), rejected, attr.Attributes, "validating", 200)
		}
		return err
	}
})
