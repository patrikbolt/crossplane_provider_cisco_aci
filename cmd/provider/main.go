package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	// Import global API types
	v1alpha1 "github.com/patrikbolt/crossplane_provider_cisco_aci/apis/v1alpha1"

	// Import TenantEPG-Controller
	epgcontroller "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/controller/tenant_epg"
)

func main() {
	var (
		debug        = flag.Bool("debug", false, "Enable debug logging.")
		syncInterval = flag.Duration("sync-interval", time.Hour, "Sync interval for all resources.")
		maxReconcile = flag.Int("max-reconcile", 10, "Maximum reconcile rate per second.")
		pollInterval = flag.Duration("poll-interval", time.Minute, "Poll interval for drift checks.")
	)
	flag.Parse()

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-aci"))
	ctrl.SetLogger(zl)

	cfg, err := ctrl.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		SyncPeriod: syncInterval,
	})
	if err != nil {
		log.Error(err, "Error creating controller manager")
		os.Exit(1)
	}

	// Register API schema
	if err := v1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "Error adding API schema")
		os.Exit(1)
	}

	o := epgcontroller.Options{
		Logger:                  log,
		MaxConcurrentReconciles: *maxReconcile,
		PollInterval:            *pollInterval,
	}

	// Setup TenantEPG controller
	if err := epgcontroller.SetupTenantEPGController(mgr, o); err != nil {
		log.Error(err, "Error setting up TenantEPG controller")
		os.Exit(1)
	}

	// Start the manager
	log.Info("Starting controller manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "Error running manager")
		os.Exit(1)
	}
}

