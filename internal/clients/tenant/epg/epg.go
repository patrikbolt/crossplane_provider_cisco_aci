package epg

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
)

type EPG struct {
    Name       string `json:"name"`
    Tenant     string `json:"tenant"`
    AppProfile string `json:"appProfile"`
}

// EPGClient wraps the ACIClient for EPG-specific operations
type EPGClient struct {
    client *clients.ACIClient
}

// NewEPGClient creates a new EPGClient
func NewEPGClient(client *clients.ACIClient) *EPGClient {
    return &EPGClient{client: client}
}

// CreateEPG creates a new EPG in Cisco ACI
func (c *EPGClient) CreateEPG(epg EPG) error {
    url := fmt.Sprintf("%s/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", c.client.baseURL, epg.Tenant, epg.AppProfile, epg.Name)
    data := map[string]interface{}{
        "fvAEPg": map[string]interface{}{
            "attributes": map[string]string{
                "name": epg.Name,
            },
        },
    }
    body, err := json.Marshal(data)
    if err != nil {
        return err
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    resp, err := c.client.doRequest(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to create EPG: %s", resp.Status)
    }
    return nil
}

// ObserveEPG retrieves the current state of the EPG from Cisco ACI
func (c *EPGClient) ObserveEPG(epg EPG) (bool, error) {
    url := fmt.Sprintf("%s/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", c.client.baseURL, epg.Tenant, epg.AppProfile, epg.Name)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false, err
    }
    resp, err := c.client.doRequest(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return false, nil // EPG does not exist
    } else if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("failed to observe EPG: %s", resp.Status)
    }

    // Parse response to verify the EPG exists and retrieve details
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return false, err
    }
    return true, nil // EPG exists
}

// UpdateEPG updates an existing EPG in Cisco ACI
func (c *EPGClient) UpdateEPG(epg EPG) error {
    url := fmt.Sprintf("%s/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", c.client.baseURL, epg.Tenant, epg.AppProfile, epg.Name)
    data := map[string]interface{}{
        "fvAEPg": map[string]interface{}{
            "attributes": map[string]string{
                "name": epg.Name,
            },
        },
    }
    body, err := json.Marshal(data)
    if err != nil {
        return err
    }
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    resp, err := c.client.doRequest(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to update EPG: %s", resp.Status)
    }
    return nil
}

// DeleteEPG deletes an existing EPG from Cisco ACI
func (c *EPGClient) DeleteEPG(epg EPG) error {
    url := fmt.Sprintf("%s/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", c.client.baseURL, epg.Tenant, epg.AppProfile, epg.Name)
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }
    resp, err := c.client.doRequest(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to delete EPG: %s", resp.Status)
    }
    return nil
}

