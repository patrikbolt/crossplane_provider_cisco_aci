module github.com/patrikbolt/crossplane_provider_cisco_aci

go 1.23

require (
    github.com/crossplane/crossplane-runtime v1.16.0
    github.com/crossplane/crossplane-tools v0.0.0-20230925130601-628280f8bf79
    github.com/google/go-cmp v0.6.0
    github.com/pkg/errors v0.9.1
    gopkg.in/alecthomas/kingpin.v2 v2.2.6
    k8s.io/apimachinery v0.29.2
    k8s.io/client-go v0.29.2
    sigs.k8s.io/controller-runtime v0.17.2
    sigs.k8s.io/controller-tools v0.14.0
)

replace github.com/patrikbolt/crossplane_provider_cisco_aci/apis/v1alpha1 => ./apis/v1alpha1
replace github.com/patrikbolt/crossplane_provider_cisco_aci/apis/tenant/epg => ./apis/tenant/epg

