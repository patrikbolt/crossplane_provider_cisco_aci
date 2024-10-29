package epg

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/patrikbolt/crossplane_provider_cisco_aci/internal/clients"
)

// EPGClient manages operations for End Point Groups (EPGs) in Cisco ACI
type EPGClient struct {
    client *clients.Client
}

// NewEPGClient initializes a new EPG client
func NewEPGClient(client *clients.Client) *EPGClient {
    return &EPGClient{
        client: client,
    }
}

// CreateEPG creates a new End Point Group (EPG) in Cisco ACI
func (c *EPGClient) CreateEPG(tenant, appProfile, epgName, bd, desc string) error {
    url := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenant, appProfile, epgName)

    // Define the payload structure for the API request
    data := map[string]interface{}{
        "fvAEPg": map[string]interface{}{
            "attributes": map[string]string{
                "dn":     fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, appProfile, epgName),
                "prio":   "level3",
                "name":   epgName,
                "descr":  desc,
                "rn":     fmt.Sprintf("epg-%s", epgName),
                "status": "created",
            },
            "children": []interface{}{
                map[string]interface{}{
                    "fvRsBd": map[string]interface{}{
                        "attributes": map[string]string{
                            "tnFvBDName": bd,
                            "status":     "created,modified",
                        },
                        "children": []interface{}{},
                    },
                },
            },
        },
    }

    // Log request details for debugging
    log.Printf("Sending POST request to %s with data: %v\n", url, data)

    // Execute the request using the DoRequest function from the client
    respBody, err := c.client.DoRequest("POST", url, data)
    if err != nil {
        return err
    }

    // Log the response for debugging
    log.Printf("Response: %s\n", string(respBody))

    // Parse the response to check for any errors
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        return err
    }

    // Check for API errors in the response
    if imdata, ok := result["imdata"].([]interface{}); ok && len(imdata) > 0 {
        if errorInfo, found := imdata[0].(map[string]interface{})["error"]; found {
            attrs := errorInfo.(map[string]interface{})["attributes"].(map[string]interface{})
            return fmt.Errorf("error creating EPG: code=%s, text=%s", attrs["code"], attrs["text"])
        }
    }

    log.Println("EPG created successfully!")
    return nil
}

// UpdateEPG updates an existing End Point Group (EPG) in Cisco ACI
func (c *EPGClient) UpdateEPG(tenant, appProfile, epgName, bd, desc string) error {
    url := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenant, appProfile, epgName)

    // Define the payload structure for the update request
    data := map[string]interface{}{
        "fvAEPg": map[string]interface{}{
            "attributes": map[string]string{
                "descr":  desc,
                "name":   epgName,
                "status": "modified",
            },
            "children": []interface{}{
                map[string]interface{}{
                    "fvRsBd": map[string]interface{}{
                        "attributes": map[string]string{
                            "tnFvBDName": bd,
                            "status":     "created,modified",
                        },
                        "children": []interface{}{},
                    },
                },
            },
        },
    }

    log.Printf("Sending POST request to %s with data: %v\n", url, data)

    respBody, err := c.client.DoRequest("POST", url, data)
    if err != nil {
        return err
    }

    log.Printf("Response: %s\n", string(respBody))

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        return err
    }

    if imdata, ok := result["imdata"].([]interface{}); ok && len(imdata) > 0 {
        if errorInfo, found := imdata[0].(map[string]interface{})["error"]; found {
            attrs := errorInfo.(map[string]interface{})["attributes"].(map[string]interface{})
            return fmt.Errorf("error updating EPG: code=%s, text=%s", attrs["code"], attrs["text"])
        }
    }

    log.Println("EPG updated successfully!")
    return nil
}

// DeleteEPG deletes an existing End Point Group (EPG) in Cisco ACI
func (c *EPGClient) DeleteEPG(tenant, appProfile, epgName string) error {
    url := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenant, appProfile, epgName)

    data := map[string]interface{}{
        "fvAEPg": map[string]interface{}{
            "attributes": map[string]string{
                "status": "deleted",
            },
        },
    }

    log.Printf("Sending POST request to %s with data: %v\n", url, data)

    respBody, err := c.client.DoRequest("POST", url, data)
    if err != nil {
        return err
    }

    log.Printf("Response: %s\n", string(respBody))

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        return err
    }

    if imdata, ok := result["imdata"].([]interface{}); ok && len(imdata) > 0 {
        if errorInfo, found := imdata[0].(map[string]interface{})["error"]; found {
            attrs := errorInfo.(map[string]interface{})["attributes"].(map[string]interface{})
            return fmt.Errorf("error deleting EPG: code=%s, text=%s", attrs["code"], attrs["text"])
        }
    }

    log.Println("EPG deleted successfully!")
    return nil
}

// ObserveEPG prüft, ob ein bestimmtes EPG existiert, und gibt den Status zurück
func (c *EPGClient) ObserveEPG(tenantName, appProfileName, epgName string) (map[string]interface{}, error) {
    endpoint := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenantName, appProfileName, epgName)
    response, err := c.client.DoRequest("GET", endpoint, nil)
    if err != nil {
        return nil, fmt.Errorf("error observing EPG: %w", err)
    }

    // Die Antwortdaten in eine Map unmarshallen, um den Status zu analysieren
    var result map[string]interface{}
    if err := json.Unmarshal(response, &result); err != nil {
        return nil, fmt.Errorf("error parsing response: %w", err)
    }

    // Prüfen, ob das EPG in den Daten enthalten ist
    imdata, ok := result["imdata"].([]interface{})
    if !ok || len(imdata) == 0 {
        return nil, fmt.Errorf("EPG %s not found in tenant %s and application profile %s", epgName, tenantName, appProfileName)
    }

    epgData := imdata[0].(map[string]interface{})
    return epgData, nil
}
