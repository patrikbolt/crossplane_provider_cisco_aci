package main

import (
    "flag"
    "fmt"
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

    client := clients.NewClient(*baseURL, *username, *password, *skipSSLVerify)
    if err := client.Authenticate(); err != nil {
        log.Fatalf("Authentifizierung fehlgeschlagen: %v", err)
    }

    epgClient := epg.NewEPGClient(client)

    if err := epgClient.DeleteEPG(*tenant, *appProfile, *epgName); err != nil {
        fmt.Printf("Fehler beim Löschen der EPG: %v\n", err)
    } else {
        fmt.Println("EPG erfolgreich gelöscht!")
    }
}

