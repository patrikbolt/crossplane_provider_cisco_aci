package main

import (
    "os"
    "path/filepath"
    "time"

    "gopkg.in/alecthomas/kingpin.v2"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/cache"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"

    "github.com/crossplane/crossplane-runtime/pkg/controller"
    "github.com/crossplane/crossplane-runtime/pkg/logging"
    "github.com/crossplane/crossplane-runtime/pkg/ratelimiter"

    "github.com/patrikbolt/crossplane_provider_cisco_aci/apis/v1alpha1"
    epgcontroller "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/controller/tenant/epg"
)

func main() {
    var (
        app            = kingpin.New(filepath.Base(os.Args[0]), "Cisco ACI support for Crossplane.").DefaultEnvars()
        debug          = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
        leaderElection = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()

        syncInterval     = app.Flag("sync", "How often all resources will be double-checked for drift from the desired state.").Short('s').Default("1h").Duration()
        pollInterval     = app.Flag("poll", "How often individual resources will be checked for drift from the desired state").Default("1m").Duration()
        maxReconcileRate = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may checked for drift from the desired state.").Default("10").Int()
    )
    kingpin.MustParse(app.Parse(os.Args[1:]))

    zl := zap.New(zap.UseDevMode(*debug))
    log := logging.NewLogrLogger(zl.WithName("provider-template"))
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
        LeaderElectionID:           "crossplane-leader-election-provider-template",
    })
    kingpin.FatalIfError(err, "Cannot create controller manager")
    kingpin.FatalIfError(v1alpha1.AddToScheme(mgr.GetScheme()), "Cannot add Template APIs to scheme")

    o := controller.Options{
        Logger:                  log,
        MaxConcurrentReconciles: *maxReconcileRate,
        PollInterval:            *pollInterval,
    }

    kingpin.FatalIfError(epgcontroller.Setup(mgr, o), "Cannot setup EPG controller")
    kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

