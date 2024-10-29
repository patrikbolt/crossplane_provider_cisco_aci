package epg

import (
    "k8s.io/apimachinery/pkg/runtime/schema"
    "sigs.k8s.io/controller-runtime/pkg/scheme"
)

// SchemeGroupVersion definiert die Gruppe und Version für die EPG-API
var SchemeGroupVersion = schema.GroupVersion{Group: "tenant.acme.org", Version: "v1alpha1"}

// SchemeBuilder registriert die EPG-Typen mit dem Runtime-Schema
var (
    SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
    AddToScheme   = SchemeBuilder.AddToScheme
)

