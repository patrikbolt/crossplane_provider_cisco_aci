package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// SchemeGroupVersion ist die Gruppen- und Versionsdefinition
var SchemeGroupVersion = schema.GroupVersion{Group: "ciscoaci.crossplane.io", Version: "v1alpha1"}

// SchemeBuilder registriert die CRD-Typen mit dem Runtime-Schema
var (
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
	AddToScheme   = SchemeBuilder.AddToScheme
)
