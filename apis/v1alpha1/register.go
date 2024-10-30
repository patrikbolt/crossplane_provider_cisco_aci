package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
        // GroupVersion defines the Group and Version for the API
	GroupVersion = schema.GroupVersion{Group: "ciscoaci.crossplane.io", Version: "v1alpha1"}

        // SchemeBuilder registers teh API-Types with the Runtime-Scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

        // AddToScheme adds the API-Types to the given Scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(
		&ProviderConfig{},
		&ProviderConfigList{},
		&TenantEPG{},
		&TenantEPGList{},
		// add additional types here as example: Tenant_BD, Tenant_BDList, etc.
	)
}

