package main

import (
    "flag"
    "log"

    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients/tenant/epg"
)

func main() {
    baseURL := flag.String("base-url", "", "Cisco ACI base URL")
    username := flag.String("username", "", "Cisco ACI username")
    password := flag.String("password", "", "Cisco ACI password")
    tenant := flag.String("tenant", "", "Tenant name")
    appProfile := flag.String("app-profile", "", "Application Profile name")
    epgName := flag.String("epg-name", "", "EPG name")
    skipSSLVerify := flag.Bool("skip-ssl-verify", false, "Skip SSL verification (insecure)")
    flag.Parse()

    // Überprüfe die Eingabeparameter
    if *baseURL == "" || *username == "" || *password == "" || *tenant == "" || *appProfile == "" || *epgName == "" {
        log.Fatal("Bitte geben Sie alle erforderlichen Parameter an: --base-url, --username, --password, --tenant, --app-profile, --epg-name")
    }

    // Erstelle den ACI-Client mit den übergebenen Parametern
    client := clients.NewClient(*baseURL, *username, *password, *skipSSLVerify)

    // Authentifizierung bei der ACI API
    if err := client.Authenticate(); err != nil {
        log.Fatalf("Fehler bei der Authentifizierung: %v", err)
    }
    log.Println("Authentifizierung erfolgreich!")

    // EPG-Client erstellen
    epgClient := epg.NewEPGClient(client)

    // Erstelle EPG basierend auf den übergebenen Parametern
    log.Printf("Erstelle EPG: Tenant=%s, AppProfile=%s, EPGName=%s", *tenant, *appProfile, *epgName)
    if err := epgClient.CreateEPG(*tenant, *appProfile, *epgName); err != nil {
        log.Fatalf("Fehler beim Erstellen der EPG: %v", err)
    }
    log.Println("EPG erfolgreich erstellt!")
}

