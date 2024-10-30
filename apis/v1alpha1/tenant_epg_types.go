package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantEPGSpec defines the wished State of TenantEPG
type TenantEPGSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TenantEPGParameters `json:"forProvider"`
}

// TenantEPGParameters are the configurable Fields of a TenantEPG
type TenantEPGParameters struct {
	Name       string `json:"name"`
	Tenant     string `json:"tenant"`
	AppProfile string `json:"appProfile"`
	Desc       string `json:"desc"`
	Bd         string `json:"bd"`
}

// TenantEPGStatus defines the observed State of a TenantEPG
type TenantEPGStatus struct {
	xpv1.ResourceStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TenantEPG is the Schema for the TenantEPG API
type TenantEPG struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TenantEPGSpec   `json:"spec"`
	Status            TenantEPGStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantEPGList contains a List of TenantEPGs
type TenantEPGList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantEPG `json:"items"`
}
