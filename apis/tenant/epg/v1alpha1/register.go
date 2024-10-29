package epg

import (
    "k8s.io/apimachinery/pkg/runtime/schema"
    "sigs.k8s.io/controller-runtime/pkg/scheme"
)

// SchemeGroupVersion defines the group and version for the EPG API
var SchemeGroupVersion = schema.GroupVersion{Group: "tenant.acme.org", Version: "v1alpha1"}

// SchemeBuilder registers the EPG types with the runtime scheme
var (
    SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
    AddToScheme   = SchemeBuilder.AddToScheme
)

