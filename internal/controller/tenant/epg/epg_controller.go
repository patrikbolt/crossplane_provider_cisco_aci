package epg

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/patrikbolt/crossplane_provider_cisco_aci/apis/tenant/epg/v1alpha1"
	"github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
	"github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients/tenant/epg"
)

// Setup sets up the EPG controller
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := "epg-controller"

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.EPG{}).
		Complete(resource.NewManagedReconciler(mgr,
			resource.ManagedKind(v1alpha1.EPGGroupVersionKind),
			resource.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: clients.NewClient}),
			o))
}

type connector struct {
	kube        client.Client
	newClientFn func(baseURL, username, password string, insecureSkipVerify bool) *clients.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (resource.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.EPG)
	if !ok {
		return nil, errors.New("managed resource is not an EPG custom resource")
	}

	pc := &v1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, client.ObjectKey{Name: cr.Spec.ProviderConfigRef.Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get ProviderConfig")
	}

	cfg := clients.NewClient(pc.Spec.BaseURL, pc.Spec.Username, pc.Spec.Password, pc.Spec.InsecureSkipVerify)
	return &external{client: epg.NewEPGClient(cfg)}, nil
}

type external struct {
	client *epg.EPGClient
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (resource.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.EPG)
	if !ok {
		return resource.ExternalObservation{}, errors.New("managed resource is not an EPG")
	}

	exists, err := c.client.ObserveEPG(epg.EPG{
		Name:       cr.Spec.ForProvider.Name,
		Tenant:     cr.Spec.ForProvider.Tenant,
		AppProfile: cr.Spec.ForProvider.AppProfile,
	})
	if err != nil {
		return resource.ExternalObservation{}, err
	}

	return resource.ExternalObservation{
		ResourceExists: exists,
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (resource.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.EPG)
	if !ok {
		return resource.ExternalCreation{}, errors.New("managed resource is not an EPG")
	}

	err := c.client.CreateEPG(epg.EPG{
		Name:       cr.Spec.ForProvider.Name,
		Tenant:     cr.Spec.ForProvider.Tenant,
		AppProfile: cr.Spec.ForProvider.AppProfile,
		Desc:       cr.Spec.ForProvider.Desc,
		Bd:         cr.Spec.ForProvider.Bd,
	})
	if err != nil {
		return resource.ExternalCreation{}, err
	}

	return resource.ExternalCreation{}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (resource.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.EPG)
	if !ok {
		return resource.ExternalUpdate{}, errors.New("managed resource is not an EPG")
	}

	err := c.client.UpdateEPG(epg.EPG{
		Name:       cr.Spec.ForProvider.Name,
		Tenant:     cr.Spec.ForProvider.Tenant,
		AppProfile: cr.Spec.ForProvider.AppProfile,
		Desc:       cr.Spec.ForProvider.Desc,
		Bd:         cr.Spec.ForProvider.Bd,
	})
	if err != nil {
		return resource.ExternalUpdate{}, err
	}

	return resource.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.EPG)
	if !ok {
		return errors.New("managed resource is not an EPG")
	}

	return c.client.DeleteEPG(epg.EPG{
		Name:       cr.Spec.ForProvider.Name,
		Tenant:     cr.Spec.ForProvider.Tenant,
		AppProfile: cr.Spec.ForProvider.AppProfile,
	})
}
