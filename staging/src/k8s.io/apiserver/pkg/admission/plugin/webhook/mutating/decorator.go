package mutating

// A WebhookInvokerDecorator decorates a WebHookInvoker.
type WebhookInvokerDecorator interface {
	// Decorate returns a WebhookInvoker wrapped by a decorator.
	Decorate(invoker WebhookInvoker) WebhookInvoker
}

// The WebhookInvokerDecoratorFunc type adapts an ordinary function as a WebhookInvokerDecorator.
type WebhookInvokerDecoratorFunc func(invoker WebhookInvoker) WebhookInvoker

// Decorate returns a WebhookInvoker wrapped by a decorator.
func (d WebhookInvokerDecoratorFunc) Decorate(invoker WebhookInvoker) WebhookInvoker {
	return d(invoker)
}

// WebhookInvokerDecorators presents multiple webhook decorators as one.
type WebhookInvokerDecorators []WebhookInvokerDecorator

// Decorate returns WebhookInvoker wrapped by multiple webhook decorators.
func (d WebhookInvokerDecorators) Decorate(invoker WebhookInvoker) WebhookInvoker {
	decorated := invoker
	for _, decorator := range d {
		decorated = decorator.Decorate(decorated)
	}
	return decorated
}
