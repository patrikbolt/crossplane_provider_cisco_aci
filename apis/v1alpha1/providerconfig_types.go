package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProviderConfigSpec specifies the Configuration for the ACI Provider
type ProviderConfigSpec struct {
	// URL of the Cisco ACI API
	URL string `json:"url"`

	// CredentialsSecretRef refers to the Kubernetes Secret, where
	// the Login Credentials  (Username and Password) for the Cisco ACI API is stored
	CredentialsSecretRef xpv1.SecretKeySelector `json:"credentialsSecretRef"`

	// InsecureSkipVerify when set to true skips the SSL Certificate Verification
	InsecureSkipVerify bool `json:"insecureSkipVerify"`
}

// ProviderConfigStatus represents the Status of ProviderConfig
type ProviderConfigStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfig configures the Connection to the Cisco ACI API
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a List of ProviderConfig Objects
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}
