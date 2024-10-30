module github.com/patrikbolt/crossplane_provider_cisco_aci

go 1.23

require (
	github.com/crossplane/crossplane-runtime v1.16.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.29.2
	k8s.io/apimachinery v0.29.2
	sigs.k8s.io/controller-runtime v0.17.2
)

replace github.com/patrikbolt/crossplane_provider_cisco_aci/apis/v1alpha1 => ./apis/v1alpha1
