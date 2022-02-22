package webhookmetrics

import (
	"context"
	"time"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/generic"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/mutating"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/validating"
	"k8s.io/apiserver/pkg/util/webhook"
	"k8s.io/klog/v2"
)

var ObserveMutatingWebhookOpenShiftMetrics = mutating.WebhookInvokerDecoratorFunc(func(invoker mutating.WebhookInvoker) mutating.WebhookInvoker {
	return func(ctx context.Context,
		hook *admissionregistrationv1.MutatingWebhook,
		invocation *generic.WebhookInvocation,
		attr *generic.VersionedAttributes,
		annotator mutating.WebhookAnnotator,
		interfaces admission.ObjectInterfaces,
		round, idx int) (bool, error) {

		t := time.Now()
		changed, err := invoker(ctx, hook, invocation, attr, annotator, interfaces, round, idx)
		latency := time.Since(t)

		rejected := err != nil
		if _, ok := err.(*webhook.ErrCallingWebhook); ok && !(hook.FailurePolicy != nil && *hook.FailurePolicy == admissionregistrationv1.Ignore) {
			rejected = true
		}

		// labels for metric
		labels := []interface{}{
			"name", hook.Name,
			"operation", attr.GetOperation(),
			"type", "mutating",
			"rejected", rejected,
			"group", invocation.Resource.Group,
			"version", invocation.Resource.Version,
			"resource", invocation.Resource.Resource,
			"subresource", invocation.Subresource,
		}

		// TODO observe
		klog.V(1).InfoS("##### MutatingWebhookInvoked", append(labels, "latency", float64(latency)/float64(time.Second))...)

		return changed, err
	}
})

func observeMetrics() {

}

var ObserveValidatingWebhookOpenShiftMetrics = validating.WebhookInvokerDecoratorFunc(func(invoker validating.WebhookInvoker) validating.WebhookInvoker {
	return func(ctx context.Context,
		hook *admissionregistrationv1.ValidatingWebhook,
		invocation *generic.WebhookInvocation,
		attr *generic.VersionedAttributes,
	) error {

		t := time.Now()
		err := invoker(ctx, hook, invocation, attr)
		latency := time.Since(t)

		rejected := err != nil
		if _, ok := err.(*webhook.ErrCallingWebhook); ok && !(hook.FailurePolicy != nil && *hook.FailurePolicy == admissionregistrationv1.Ignore) {
			rejected = true
		}

		// labels for metric
		labels := []interface{}{
			"name", hook.Name,
			"operation", attr.GetOperation(),
			"type", "validating",
			"rejected", rejected,
			"group", invocation.Resource.Group,
			"version", invocation.Resource.Version,
			"resource", invocation.Resource.Resource,
			"subresource", invocation.Subresource,
		}

		// TODO observe
		klog.V(1).InfoS("##### ValidatingWebhookInvoked", append(labels, "latency", float64(latency)/float64(time.Second))...)

		return err
	}
})
