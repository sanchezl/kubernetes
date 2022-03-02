package admissionenablement

import (
	"encoding/json"

	configv1 "github.com/openshift/api/config/v1"
	kubecontrolplanev1 "github.com/openshift/api/kubecontrolplane/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/config/apis/webhookadmission/v2alpha1"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/mutating"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/validating"
)

func ExpandCardinalityOfAdmissionWebhookDurationMetrics(openshiftConfig *kubecontrolplanev1.KubeAPIServerConfig) error {

	webhookAdmissionConfig := &v2alpha1.WebhookAdmission{
		TypeMeta: metav1.TypeMeta{Kind: "WebhookAdmissionConfiguration", APIVersion: "apiserver.config.k8s.io/v2alpha1"},
	}
	webhookAdmissionConfig.Metrics.Duration.IncludeResourceLabelsFor = []v2alpha1.Rule{
		{Groups: []string{""}, Resources: []string{"endpoints", "events", "pods"}},
		{Groups: []string{"quota"}, Resources: []string{"resourcequotas"}},
		{Groups: []string{"apiserver.openshift.io"}, Resources: []string{"apirequestcounts"}},
		{Groups: []string{"discovery.k8s.io"}, Resources: []string{"endpointslices"}},
		{Groups: []string{"events.k8s.io"}, Resources: []string{"events"}},
		{Groups: []string{"quota.openshift.io"}, Resources: []string{"clusterresourcequotas"}},
	}
	raw, err := json.Marshal(webhookAdmissionConfig)
	if err != nil {
		return err
	}

	// TODO just overwriting for now, as upstream does not have default configs for admission webhooks
	admissionPluginConfig := configv1.AdmissionPluginConfig{Configuration: runtime.RawExtension{Raw: raw}}
	openshiftConfig.AdmissionConfig.PluginConfig[mutating.PluginName] = admissionPluginConfig
	openshiftConfig.AdmissionConfig.PluginConfig[validating.PluginName] = admissionPluginConfig

	return nil
}
