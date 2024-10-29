/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
    "context"
    "os"
    "path/filepath"
    "time"

    "gopkg.in/alecthomas/kingpin.v2"
    kerrors "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/leaderelection/resourcelock"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/cache"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"

    xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
    "github.com/crossplane/crossplane-runtime/pkg/controller"
    "github.com/crossplane/crossplane-runtime/pkg/feature"
    "github.com/crossplane/crossplane-runtime/pkg/logging"
    "github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
    "github.com/crossplane/crossplane-runtime/pkg/resource"

    "github.com/patrikbolt/crossplane_provider_cisco_aci/apis" // Passe diesen Pfad an
    "github.com/patrikbolt/crossplane_provider_cisco_aci/apis/tenant/epg/v1alpha1" // Passe diesen Pfad an
    epgcontroller "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/controller/tenant/epg" // Passe diesen Pfad an
)

func main() {
    var (
        app            = kingpin.New(filepath.Base(os.Args[0]), "Cisco ACI support for Crossplane.").DefaultEnvars()
        debug          = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
        leaderElection = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()

        syncInterval     = app.Flag("sync", "How often all resources will be double-checked for drift from the desired state.").Short('s').Default("1h").Duration()
        pollInterval     = app.Flag("poll", "How often individual resources will be checked for drift from the desired state").Default("1m").Duration()
        maxReconcileRate = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may checked for drift from the desired state.").Default("10").Int()

        namespace                  = app.Flag("namespace", "Namespace used to set as default scope in default secret store config.").Default("crossplane-system").Envar("POD_NAMESPACE").String()
        enableExternalSecretStores = app.Flag("enable-external-secret-stores", "Enable support for ExternalSecretStores.").Default("false").Envar("ENABLE_EXTERNAL_SECRET_STORES").Bool()
        enableManagementPolicies   = app.Flag("enable-management-policies", "Enable support for Management Policies.").Default("false").Envar("ENABLE_MANAGEMENT_POLICIES").Bool()
    )
    kingpin.MustParse(app.Parse(os.Args[1:]))

    zl := zap.New(zap.UseDevMode(*debug))
    log := logging.NewLogrLogger(zl.WithName("provider-cisco-aci"))
    if *debug {
        ctrl.SetLogger(zl)
    }

    cfg, err := ctrl.GetConfig()
    kingpin.FatalIfError(err, "Cannot get API server rest config")

    mgr, err := ctrl.NewManager(ratelimiter.LimitRESTConfig(cfg, *maxReconcileRate), ctrl.Options{
        Cache: cache.Options{
            SyncPeriod: syncInterval,
        },
        LeaderElection:             *leaderElection,
        LeaderElectionID:           "crossplane-leader-election-provider-cisco-aci",
        LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
        LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
        RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
    })
    kingpin.FatalIfError(err, "Cannot create controller manager")
    kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add ACI APIs to scheme")

    o := controller.Options{
        Logger:                  log,
        MaxConcurrentReconciles: *maxReconcileRate,
        PollInterval:            *pollInterval,
        GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
        Features:                &feature.Flags{},
    }

    if *enableExternalSecretStores {
        o.Features.Enable(feature.Flag("EnableExternalSecretStores"))
        log.Info("Alpha feature enabled", "flag", "EnableExternalSecretStores")

        kingpin.FatalIfError(resource.Ignore(kerrors.IsAlreadyExists, mgr.GetClient().Create(context.Background(), &v1alpha1.StoreConfig{
            ObjectMeta: metav1.ObjectMeta{
                Name: "default",
            },
            Spec: v1alpha1.StoreConfigSpec{
                SecretStoreConfig: xpv1.SecretStoreConfig{
                    DefaultScope: *namespace,
                },
            },
        })), "cannot create default store config")
    }

    if *enableManagementPolicies {
        o.Features.Enable(feature.Flag("EnableManagementPolicies"))
        log.Info("Alpha feature enabled", "flag", "EnableManagementPolicies")
    }

    // Setup EPG Controller
    kingpin.FatalIfError(epgcontroller.Setup(mgr, o), "Cannot setup EPG controller")

    kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

