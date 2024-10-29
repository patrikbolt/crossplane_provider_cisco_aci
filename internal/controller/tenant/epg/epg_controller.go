package epg

import (
    "context"
    "fmt"

    "github.com/patrikbolt/crossplane_provider_cisco_aci/apis/tenant/epg/v1alpha1" // Importiere die EPG-API
    aciv1alpha1 "github.com/patrikbolt/crossplane_provider_cisco_aci/apis/v1alpha1" // Importiere die ProviderConfig API
    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"         // Importiere den ACI-Client
    "github.com/crossplane/crossplane-runtime/pkg/controller"
    "github.com/crossplane/crossplane-runtime/pkg/logging"
    "github.com/crossplane/crossplane-runtime/pkg/event"
    "github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
    "github.com/crossplane/crossplane-runtime/pkg/resource"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    corev1 "k8s.io/api/core/v1"
)

type connector struct {
    kube client.Client
}

type external struct {
    client *clients.ACIClient
}

// Setup adds a controller that reconciles EPG managed resources
func Setup(mgr ctrl.Manager, o controller.Options) error {
    name := managed.ControllerName(v1alpha1.EPGGroupKind)

    return ctrl.NewControllerManagedBy(mgr).
        Named(name).
        For(&v1alpha1.EPG{}).
        Complete(managed.NewReconciler(mgr,
            resource.ManagedKind(v1alpha1.GroupVersion.WithKind("EPG")),
            managed.WithExternalConnecter(&connector{kube: mgr.GetClient()}),
            managed.WithLogger(o.Logger.WithValues("controller", name)),
            managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
        ))
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
    // Konvertiere die Managed Resource zu einem EPG
    epg, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return nil, fmt.Errorf("managed resource is not an EPG")
    }

    // Lade die ProviderConfig
    pc := &aciv1alpha1.ProviderConfig{}
    if err := c.kube.Get(ctx, client.ObjectKey{Name: epg.GetProviderConfigReference().Name}, pc); err != nil {
        return nil, err
    }

    // Lade das Kubernetes-Secret für die Zugangsdaten
    secret := &corev1.Secret{}
    if err := c.kube.Get(ctx, client.ObjectKey{
        Namespace: pc.Spec.CredentialsSecretRef.Namespace,
        Name:      pc.Spec.CredentialsSecretRef.Name,
    }, secret); err != nil {
        return nil, fmt.Errorf("cannot get secret %s: %w", pc.Spec.CredentialsSecretRef.Name, err)
    }

    // Extrahiere die Zugangsdaten
    username := string(secret.Data["username"])
    password := string(secret.Data["password"])

    // Initialisiere den ACIClient mit URL und Zugangsdaten
    aciClient, err := clients.NewACIClient(pc.Spec.URL, username, password)
    if err != nil {
        return nil, err
    }
    return &external{client: aciClient}, nil
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
    epg, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return managed.ExternalObservation{}, fmt.Errorf("managed resource is not an EPG")
    }

    exists, err := c.client.ObserveEPG(clients.EPG{
        Name:       epg.Spec.ForProvider.Name,
        Tenant:     epg.Spec.ForProvider.Tenant,
        AppProfile: epg.Spec.ForProvider.AppProfile,
    })
    if err != nil {
        return managed.ExternalObservation{}, err
    }

    return managed.ExternalObservation{
        ResourceExists:    exists,
        ResourceUpToDate:  true, // Hier könntest du eine tiefere Prüfung für Aktualität einbauen
    }, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
    epg, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return managed.ExternalCreation{}, fmt.Errorf("managed resource is not an EPG")
    }

    err := c.client.CreateEPG(clients.EPG{
        Name:       epg.Spec.ForProvider.Name,
        Tenant:     epg.Spec.ForProvider.Tenant,
        AppProfile: epg.Spec.ForProvider.AppProfile,
    })
    return managed.ExternalCreation{}, err
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
    epg, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return managed.ExternalUpdate{}, fmt.Errorf("managed resource is not an EPG")
    }

    err := c.client.UpdateEPG(clients.EPG{
        Name:       epg.Spec.ForProvider.Name,
        Tenant:     epg.Spec.ForProvider.Tenant,
        AppProfile: epg.Spec.ForProvider.AppProfile,
    })
    return managed.ExternalUpdate{}, err
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
    epg, ok := mg.(*v1alpha1.EPG)
    if !ok {
        return fmt.Errorf("managed resource is not an EPG")
    }

    return c.client.DeleteEPG(clients.EPG{
        Name:       epg.Spec.ForProvider.Name,
        Tenant:     epg.Spec.ForProvider.Tenant,
        AppProfile: epg.Spec.ForProvider.AppProfile,
    })
}

