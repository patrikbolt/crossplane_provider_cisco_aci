package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
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
		log.Fatalf("Authentication failed: %v", err)
	}

	epgClient := clients.NewTenantEPGClient(client)

	// Beobachten des EPG-Status
	exists, err := epgClient.ObserveTenantEPG(*tenant, *appProfile, *epgName)
	if err != nil {
		log.Fatalf("Error observing EPG: %v", err)
	}

	if exists {
		fmt.Printf("EPG %s observed successfully.\n", *epgName)
	} else {
		fmt.Printf("EPG %s does not exist.\n", *epgName)
	}
}
