package webhookmetrics

import (
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/initializer"
)

func NewInitializer() *localInitializer {
	return &localInitializer{}
}

type localInitializer struct{}

func (i *localInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(initializer.WantsMutatingWebhookInvokerDecorator); ok {
		wants.AppendWebhookInvokerDecorator(ObserveMutatingWebhookOpenShiftMetrics)
	}
	if wants, ok := plugin.(initializer.WantsValidatingWebhookInvokerDecorator); ok {
		wants.AppendWebhookInvokerDecorator(ObserveValidatingWebhookOpenShiftMetrics)
	}
}
