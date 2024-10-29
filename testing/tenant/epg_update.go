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
    desc := flag.String("desc", "", "New description")
    skipSSLVerify := flag.Bool("skip-ssl-verify", false, "Skip SSL verification (insecure)")
    flag.Parse()

    client := clients.NewClient(*baseURL, *username, *password, *skipSSLVerify)
    if err := client.Authenticate(); err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }

    epgClient := epg.NewEPGClient(client)
    err := epgClient.UpdateEPG(*tenant, *appProfile, *epgName, *desc)
    if err != nil {
        log.Fatalf("Failed to update EPG: %v", err)
    }
    fmt.Println("EPG updated successfully!")
}

