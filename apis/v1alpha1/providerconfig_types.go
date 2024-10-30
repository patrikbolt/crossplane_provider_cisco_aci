package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProviderConfigSpec spezifiziert die Konfiguration für den ACI Provider
type ProviderConfigSpec struct {
	// URL der Cisco ACI API
	URL string `json:"url"`

	// CredentialsSecretRef verweist auf das Kubernetes Secret, das
	// die Anmeldeinformationen (Benutzername und Passwort) für die Cisco ACI API enthält
	CredentialsSecretRef xpv1.SecretKeySelector `json:"credentialsSecretRef"`

	// InsecureSkipVerify überspringt die SSL-Zertifikatüberprüfung, wenn auf true gesetzt
	InsecureSkipVerify bool `json:"insecureSkipVerify"`
}

// ProviderConfigStatus repräsentiert den Status der ProviderConfig
type ProviderConfigStatus struct {
	xpv1.ConditionedStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfig konfiguriert eine Verbindung zur Cisco ACI API
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList enthält eine Liste von ProviderConfig Objekten
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

