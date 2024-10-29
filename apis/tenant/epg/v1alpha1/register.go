package epg

import (
    "sigs.k8s.io/controller-runtime/pkg/scheme"
    "k8s.io/apimachinery/pkg/runtime"
)

// SchemeBuilder registers the EPG type with the runtime scheme.
var (
    SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
    AddToScheme   = SchemeBuilder.AddToScheme
)

// SchemeGroupVersion defines the group and version for the EPG API.
var SchemeGroupVersion = runtime.NewSchemeGroupVersion("tenant.acme.org", "v1alpha1")

