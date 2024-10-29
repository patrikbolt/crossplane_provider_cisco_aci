package epg

import (
    "fmt"

    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
)

// EPGClient bietet Funktionen zur Verwaltung von EPGs
type EPGClient struct {
    Client *clients.Client
}

// NewEPGClient erstellt einen neuen EPGClient
func NewEPGClient(client *clients.Client) *EPGClient {
    return &EPGClient{
        Client: client,
    }
}

// CreateEPG erstellt ein neues EPG
func (c *EPGClient) CreateEPG(tenantName, appProfileName, epgName string) error {
    endpoint := fmt.Sprintf("/api/mo/uni/tn-%s/ap-%s/epg-%s.json", tenantName, appProfileName, epgName)
    data := map[string]interface{}{
        "fvAEPg": map[string]interface{}{
            "attributes": map[string]string{
                "name": epgName,
            },
        },
    }
    _, err := c.Client.DoRequest("POST", endpoint, data)
    return err
}

// GetEPG ruft Informationen über ein EPG ab
func (c *EPGClient) GetEPG(tenantName, appProfileName, epgName string) ([]byte, error) {
    endpoint := fmt.Sprintf("/api/mo/uni/tn-%s/ap-%s/epg-%s.json", tenantName, appProfileName, epgName)
    return c.Client.DoRequest("GET", endpoint, nil)
}

// UpdateEPG aktualisiert ein bestehendes EPG
func (c *EPGClient) UpdateEPG(tenantName, appProfileName, epgName string, data map[string]interface{}) error {
    endpoint := fmt.Sprintf("/api/mo/uni/tn-%s/ap-%s/epg-%s.json", tenantName, appProfileName, epgName)
    _, err := c.Client.DoRequest("POST", endpoint, data)
    return err
}

// DeleteEPG löscht ein EPG
func (c *EPGClient) DeleteEPG(tenantName, appProfileName, epgName string) error {
    endpoint := fmt.Sprintf("/api/mo/uni/tn-%s/ap-%s/epg-%s.json", tenantName, appProfileName, epgName)
    data := map[string]interface{}{
        "fvAEPg": map[string]interface{}{
            "attributes": map[string]string{
                "status": "deleted",
            },
        },
    }
    _, err := c.Client.DoRequest("POST", endpoint, data)
    return err
}

