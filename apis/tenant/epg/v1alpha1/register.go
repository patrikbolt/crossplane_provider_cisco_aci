// In apis/tenant/epg/v1alpha1/register.go

package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion is the group version used to register these objects.
var (
    GroupVersion = schema.GroupVersion{Group: "tenant.crossplane.io", Version: "v1alpha1"}
    SchemeBuilder = &metav1.SchemeBuilder{GroupVersion}
    AddToScheme = SchemeBuilder.AddToScheme
)

