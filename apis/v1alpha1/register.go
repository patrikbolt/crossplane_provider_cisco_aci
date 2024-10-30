package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion definiert die Gruppe und Version für die API
	GroupVersion = schema.GroupVersion{Group: "ciscoaci.crossplane.io", Version: "v1alpha1"}

	// SchemeBuilder registriert die API-Typen mit dem Runtime-Scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme fügt die API-Typen zum gegebenen Scheme hinzu
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

