package main

import (
	"flag"
	"fmt"

	"github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
)

func main() {
	baseURL := flag.String("base-url", "", "Cisco ACI base URL")
	username := flag.String("username", "", "Cisco ACI username")
	password := flag.String("password", "", "Cisco ACI password")
	tenant := flag.String("tenant", "", "Tenant name")
	appProfile := flag.String("app-profile", "", "Application Profile name")
	epgName := flag.String("epg-name", "", "EPG name")
	desc := flag.String("desc", "", "EPG description")
	bd := flag.String("bd", "", "Bridge domain name")
	skipSSLVerify := flag.Bool("skip-ssl-verify", false, "Skip SSL verification (insecure)")
	flag.Parse()

	// Create the ACI client
	client := clients.NewClient(*baseURL, *username, *password, *skipSSLVerify)

	// Authenticate the client
	if err := client.Authenticate(); err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		return
	}

	// Create the EPG client
	epgClient := clients.NewTenantEPGClient(client)

	// Update the EPG with the specified parameters
	err := epgClient.UpdateTenantEPG(*tenant, *appProfile, *epgName, *bd, *desc)
	if err != nil {
		fmt.Printf("Error updating EPG: %v\n", err)
	} else {
		fmt.Println("EPG updated successfully!")
	}
}

