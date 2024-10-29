package main

import (
    "flag"
    "fmt"

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

    // Erstellen Sie den ACI-Client mit den übergebenen Parametern
    client := clients.NewClient(*baseURL, *username, *password, *skipSSLVerify)

    // EPG-Client erstellen
    epgClient := epg.NewEPGClient(client)

    // Erstelle EPG basierend auf den übergebenen Parametern
    err := epgClient.CreateEPG(*tenant, *appProfile, *epgName)
    if err != nil {
        fmt.Printf("Fehler beim Erstellen der EPG: %v\n", err)
    } else {
        fmt.Println("EPG erfolgreich erstellt!")
    }
}

