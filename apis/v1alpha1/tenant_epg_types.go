package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantEPGSpec definiert den gewünschten Zustand eines TenantEPG.
type TenantEPGSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TenantEPGParameters `json:"forProvider"`
}

// TenantEPGParameters sind die konfigurierbaren Felder eines TenantEPG.
type TenantEPGParameters struct {
	Name       string `json:"name"`
	Tenant     string `json:"tenant"`
	AppProfile string `json:"appProfile"`
	Desc       string `json:"desc"`
	Bd         string `json:"bd"`
}

// TenantEPGStatus definiert den beobachteten Zustand eines TenantEPG.
type TenantEPGStatus struct {
	xpv1.ResourceStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TenantEPG ist das Schema für die TenantEPG API.
type TenantEPG struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TenantEPGSpec   `json:"spec"`
	Status            TenantEPGStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantEPGList enthält eine Liste von TenantEPGs.
type TenantEPGList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantEPG `json:"items"`
}

