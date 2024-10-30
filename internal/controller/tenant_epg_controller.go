package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// Import global API types
	v1alpha1 "github.com/patrikbolt/crossplane_provider_cisco_aci/apis/v1alpha1"

	"github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"

	corev1 "k8s.io/api/core/v1"
)

// Options definiert die Konfigurationsoptionen für den TenantEPG-Controller
type Options struct {
	Logger                  logging.Logger
	MaxConcurrentReconciles int
	PollInterval            time.Duration
}

// SetupTenantEPGController richtet den TenantEPG-Controller mit dem Manager ein.
func SetupTenantEPGController(mgr ctrl.Manager, o Options) error {
	name := managed.ControllerName(v1alpha1.TenantEPGGroupKind)

	// Definieren des Controllers mit den gewünschten Optionen
	_, err := ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.TenantEPG{}).
		WithOptions(ctrl.Options{
			MaxConcurrentReconciles: o.MaxConcurrentReconciles,
		}).
		Complete(managed.NewReconciler(mgr,
			managed.WithExternalConnecter(&connector{
				kube:        mgr.GetClient(),
				newClientFn: clients.NewClient,
			}),
			managed.WithLogger(o.Logger.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		))
	if err != nil {
		return errors.Wrap(err, "cannot create TenantEPG controller")
	}

	return nil
}

type connector struct {
	kube        client.Client
	newClientFn func(baseURL, username, password string, insecureSkipVerify bool) *clients.Client
}

// Credentials hält den Benutzernamen und das Passwort, die aus dem Secret extrahiert wurden
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *connector) Connect(ctx context.Context, mg managed.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.TenantEPG)
	if !ok {
		return nil, errors.New("managed resource is not a TenantEPG custom resource")
	}

	// ProviderConfig abrufen
	pc := &v1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, client.ObjectKey{Name: cr.Spec.ProviderConfigReference.Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get ProviderConfig")
	}

	// Credentials aus dem Secret abrufen
	credsSecret := &corev1.Secret{}
	if err := c.kube.Get(ctx, client.ObjectKey{Namespace: pc.Spec.CredentialsSecretRef.Namespace, Name: pc.Spec.CredentialsSecretRef.Name}, credsSecret); err != nil {
		return nil, errors.Wrap(err, "cannot get credentials secret")
	}

	// Benutzername und Passwort aus dem Secret extrahieren
	credsBytes, ok := credsSecret.Data[pc.Spec.CredentialsSecretRef.Key]
	if !ok {
		return nil, errors.New("credentials secret does not contain the specified key")
	}

	var creds Credentials
	if err := json.Unmarshal(credsBytes, &creds); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal credentials secret")
	}

	// insecureSkipVerify aus ProviderConfig abrufen
	insecureSkipVerify := pc.Spec.InsecureSkipVerify

	// Client erstellen
	apiClient := c.newClientFn(pc.Spec.URL, creds.Username, creds.Password, insecureSkipVerify)

	return &external{client: clients.NewTenantEPGClient(apiClient)}, nil
}

type external struct {
	client *clients.TenantEPGClient
}

func (c *external) Observe(ctx context.Context, mg managed.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.TenantEPG)
	if !ok {
		return managed.ExternalObservation{}, errors.New("managed resource is not a TenantEPG")
	}

	// ObserveTenantEPG mit tenant, appProfile, epgName aufrufen
	exists, err := c.client.ObserveTenantEPG(cr.Spec.ForProvider.Tenant, cr.Spec.ForProvider.AppProfile, cr.Spec.ForProvider.Name)
	if err != nil {
		// Wenn das TenantEPG nicht gefunden wird, setzen wir ResourceExists auf false
		if !exists {
			return managed.ExternalObservation{
				ResourceExists: false,
			}, nil
		}
		return managed.ExternalObservation{}, err
	}

	// Wenn epgData vorhanden ist, setzen wir ResourceExists auf true
	return managed.ExternalObservation{
		ResourceExists: true,
	}, nil
}

func (c *external) Create(ctx context.Context, mg managed.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.TenantEPG)
	if !ok {
		return managed.ExternalCreation{}, errors.New("managed resource is not a TenantEPG")
	}

	// CreateTenantEPG mit tenant, appProfile, epgName, bd, desc aufrufen
	err := c.client.CreateTenantEPG(cr.Spec.ForProvider.Tenant, cr.Spec.ForProvider.AppProfile, cr.Spec.ForProvider.Name, cr.Spec.ForProvider.Bd, cr.Spec.ForProvider.Desc)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	return managed.ExternalCreation{}, nil
}

func (c *external) Update(ctx context.Context, mg managed.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.TenantEPG)
	if !ok {
		return managed.ExternalUpdate{}, errors.New("managed resource is not a TenantEPG")
	}

	// UpdateTenantEPG mit tenant, appProfile, epgName, bd, desc aufrufen
	err := c.client.UpdateTenantEPG(cr.Spec.ForProvider.Tenant, cr.Spec.ForProvider.AppProfile, cr.Spec.ForProvider.Name, cr.Spec.ForProvider.Bd, cr.Spec.ForProvider.Desc)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	return managed.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg managed.Managed) error {
	cr, ok := mg.(*v1alpha1.TenantEPG)
	if !ok {
		return errors.New("managed resource is not a TenantEPG")
	}

	// DeleteTenantEPG mit tenant, appProfile, epgName aufrufen
	err := c.client.DeleteTenantEPG(cr.Spec.ForProvider.Tenant, cr.Spec.ForProvider.AppProfile, cr.Spec.ForProvider.Name)
	if err != nil {
		return err
	}

	return nil
}

