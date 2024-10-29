package epg

import (
    "context"

    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "github.com/crossplane/crossplane-runtime/pkg/resource"
    "github.com/crossplane/crossplane-runtime/pkg/logging"
    "github.com/pkg/errors"

    "github.com/patrikbolt/crossplane_provider_cisco_aci/apis/tenant/epg/v1alpha1"
    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
)

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
    newClientFn func(config *clients.ACIConfig) (*clients.ACIClient, error)
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

    cfg := clients.NewConfig(pc)
    client, err := c.newClientFn(cfg)
    if err != nil {
        return nil, errors.Wrap(err, "cannot create ACI client")
    }
    return &external{client: client}, nil
}

type external struct {
    client *clients.ACIClient
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (resource.ExternalObservation, error) {
    cr, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return resource.ExternalObservation{}, errors.New("managed resource is not an EPG")
    }

    observation, err := c.client.ObserveEPG(ctx, cr.Spec.ForProvider)
    if err != nil {
        return resource.ExternalObservation{}, err
    }
    return observation, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (resource.ExternalCreation, error) {
    cr, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return resource.ExternalCreation{}, errors.New("managed resource is not an EPG")
    }

    creation, err := c.client.CreateEPG(ctx, cr.Spec.ForProvider)
    if err != nil {
        return resource.ExternalCreation{}, err
    }
    return creation, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (resource.ExternalUpdate, error) {
    cr, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return resource.ExternalUpdate{}, errors.New("managed resource is not an EPG")
    }

    update, err := c.client.UpdateEPG(ctx, cr.Spec.ForProvider)
    if err != nil {
        return resource.ExternalUpdate{}, err
    }
    return update, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
    cr, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return errors.New("managed resource is not an EPG")
    }

    return c.client.DeleteEPG(ctx, cr.Spec.ForProvider)
}

