package v1alpha1

import (
    xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProviderConfigSpec specifies the configuration for the ACI Provider.
type ProviderConfigSpec struct {
    // URL of the Cisco ACI API.
    URL string `json:"url"`

    // CredentialsSecretRef refers to the Kubernetes Secret containing
    // the credentials (username and password) for the Cisco ACI API.
    CredentialsSecretRef xpv1.SecretKeySelector `json:"credentialsSecretRef"`

    // InsecureSkipVerify skips SSL certificate verification when set to true.
    InsecureSkipVerify bool `json:"insecureSkipVerify"`
}

// ProviderConfigStatus represents the status of the ProviderConfig.
type ProviderConfigStatus struct {
    xpv1.ResourceStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfig configures a connection to the Cisco ACI API.
type ProviderConfig struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   ProviderConfigSpec   `json:"spec"`
    Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig objects.
type ProviderConfigList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []ProviderConfig `json:"items"`
}

