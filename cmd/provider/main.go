package main

import (
    "context"
    "flag"
    "fmt"
    "os"

    "github.com/crossplane/crossplane-runtime/pkg/logging"
    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients/tenant/epg"
    "go.uber.org/zap/zapcore"
    ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
    var (
        baseURL  string
        username string
        password string
    )

    // Flags für die ACI API Konfiguration
    flag.StringVar(&baseURL, "base-url", "", "Basis-URL der ACI API")
    flag.StringVar(&username, "username", "", "Benutzername für die ACI API")
    flag.StringVar(&password, "password", "", "Passwort für die ACI API")
    flag.Parse()

    // Überprüfe, ob alle notwendigen Parameter vorhanden sind
    if baseURL == "" || username == "" || password == "" {
        fmt.Println("Bitte geben Sie --base-url, --username und --password an")
        os.Exit(1)
    }

    // Logger initialisieren
    logger := logging.NewLogrLogger(ctrl.Log.WithName("provider-aci"))
    ctrl.SetLogger(logger.WithValues("provider", "aci"))

    // Erstelle einen neuen ACI API Client
    client := clients.NewClient(baseURL, username, password)
    err := client.Authenticate()
    if err != nil {
        logger.Error(err, "Authentifizierung bei der ACI API fehlgeschlagen")
        os.Exit(1)
    }

    // Erstelle einen neuen EPGClient
    epgClient := epg.NewEPGClient(client)

    // Beispiel: Erstelle ein neues EPG (dieser Teil sollte eigentlich in deinem Controller passieren)
    tenantName := "example-tenant"
    appProfileName := "example-app-profile"
    epgName := "example-epg"

    err = epgClient.CreateEPG(tenantName, appProfileName, epgName)
    if err != nil {
        logger.Error(err, "Fehler beim Erstellen des EPG")
        os.Exit(1)
    }

    logger.Info("EPG erfolgreich erstellt")

    // Hier würdest du den Controller-Manager initialisieren und deine Controller hinzufügen
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme: nil, // Füge hier dein Schema hinzu
    })
    if err != nil {
        logger.Error(err, "Fehler beim Erstellen des Managers")
        os.Exit(1)
    }

    // Starte den Manager (dieser Blockiert den Hauptthread)
    if err := mgr.Start(context.Background()); err != nil {
        logger.Error(err, "Manager beendet")
        os.Exit(1)
    }
}

