package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion definiert die API-Gruppen- und Versionsinformationen
var (
    GroupVersion   = schema.GroupVersion{Group: "aci.crossplane.io", Version: "v1alpha1"}
    SchemeBuilder  = &metav1.SchemeBuilder{GroupVersion} // SchemeBuilder definiert
    AddToScheme    = SchemeBuilder.AddToScheme
)

// init registriert die EPG- und ProviderConfig-Ressourcen im Scheme
func init() {
    SchemeBuilder.Register(
        &EPG{}, &EPGList{},                 // Registrierung der EPG-Ressource
        &ProviderConfig{}, &ProviderConfigList{}, // Registrierung der ProviderConfig-Ressource
    )
}

