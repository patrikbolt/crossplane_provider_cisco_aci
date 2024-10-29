// In apis/tenant/epg/v1alpha1/epg_types.go

package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// EPGSpec defines the desired state of an EPG.
type EPGSpec struct {
    xpv1.ResourceSpec `json:",inline"`
    ForProvider       EPGParameters `json:"forProvider"`
}

// EPGParameters are the configurable fields of an EPG.
type EPGParameters struct {
    Name       string `json:"name"`
    Tenant     string `json:"tenant"`
    AppProfile string `json:"appProfile"`
}

// EPGStatus represents the observed state of an EPG.
type EPGStatus struct {
    xpv1.ResourceStatus `json:",inline"`
    AtProvider          EPGObservation `json:"atProvider,omitempty"`
}

// EPGObservation are the observable fields of an EPG.
type EPGObservation struct {
    ID string `json:"id,omitempty"`
}

// +kubebuilder:object:root=true

// EPG is the Schema for the EPGs API
type EPG struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   EPGSpec   `json:"spec"`
    Status EPGStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EPGList contains a list of EPGs.
type EPGList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []EPG `json:"items"`
}

