package clients

import (
	"encoding/json"
	"fmt"
	"log"
)

// TenantEPGClient verwaltet Operationen für End Point Groups (EPGs) in Cisco ACI
type TenantEPGClient struct {
	client *Client
}

// NewTenantEPGClient initialisiert einen neuen TenantEPG-Client
func NewTenantEPGClient(client *Client) *TenantEPGClient {
	return &TenantEPGClient{
		client: client,
	}
}

// CreateTenantEPG erstellt eine neue End Point Group (EPG) in Cisco ACI
func (c *TenantEPGClient) CreateTenantEPG(tenant, appProfile, epgName, bd, desc string) error {
	url := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenant, appProfile, epgName)

	// Definiere die Payload-Struktur für die API-Anfrage
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

	// Logge Anfrage-Details zum Debuggen
	log.Printf("Sende POST-Anfrage an %s mit Daten: %v\n", url, data)

	// Führe die Anfrage mit der DoRequest-Funktion des Clients aus
	respBody, err := c.client.DoRequest("POST", url, data)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der TenantEPG: %v", err)
	}

	// Logge die Antwort zum Debuggen
	log.Printf("Antwort: %s\n", string(respBody))

	// Parse die Antwort, um nach Fehlern zu suchen
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("Fehler beim Unmarshalen der Antwort: %v", err)
	}

	// Überprüfe auf API-Fehler in der Antwort
	if imdata, ok := result["imdata"].([]interface{}); ok && len(imdata) > 0 {
		if errorInfo, found := imdata[0].(map[string]interface{})["error"]; found {
			attrs := errorInfo.(map[string]interface{})["attributes"].(map[string]interface{})
			return fmt.Errorf("Fehler beim Erstellen der TenantEPG: code=%s, text=%s", attrs["code"], attrs["text"])
		}
	}

	log.Println("TenantEPG erfolgreich erstellt!")
	return nil
}

// UpdateTenantEPG aktualisiert eine bestehende End Point Group (EPG) in Cisco ACI
func (c *TenantEPGClient) UpdateTenantEPG(tenant, appProfile, epgName, bd, desc string) error {
	url := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenant, appProfile, epgName)

	// Definiere die Payload-Struktur für die Update-Anfrage
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

	log.Printf("Sende POST-Anfrage an %s mit Daten: %v\n", url, data)

	respBody, err := c.client.DoRequest("POST", url, data)
	if err != nil {
		return fmt.Errorf("Fehler beim Aktualisieren der TenantEPG: %v", err)
	}

	log.Printf("Antwort: %s\n", string(respBody))

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("Fehler beim Unmarshalen der Antwort: %v", err)
	}

	if imdata, ok := result["imdata"].([]interface{}); ok && len(imdata) > 0 {
		if errorInfo, found := imdata[0].(map[string]interface{})["error"]; found {
			attrs := errorInfo.(map[string]interface{})["attributes"].(map[string]interface{})
			return fmt.Errorf("Fehler beim Aktualisieren der TenantEPG: code=%s, text=%s", attrs["code"], attrs["text"])
		}
	}

	log.Println("TenantEPG erfolgreich aktualisiert!")
	return nil
}

// DeleteTenantEPG löscht eine bestehende End Point Group (EPG) in Cisco ACI
func (c *TenantEPGClient) DeleteTenantEPG(tenant, appProfile, epgName string) error {
	url := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenant, appProfile, epgName)

	data := map[string]interface{}{
		"fvAEPg": map[string]interface{}{
			"attributes": map[string]string{
				"status": "deleted",
			},
		},
	}

	log.Printf("Sende POST-Anfrage an %s mit Daten: %v\n", url, data)

	respBody, err := c.client.DoRequest("POST", url, data)
	if err != nil {
		return fmt.Errorf("Fehler beim Löschen der TenantEPG: %v", err)
	}

	log.Printf("Antwort: %s\n", string(respBody))

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("Fehler beim Unmarshalen der Antwort: %v", err)
	}

	if imdata, ok := result["imdata"].([]interface{}); ok && len(imdata) > 0 {
		if errorInfo, found := imdata[0].(map[string]interface{})["error"]; found {
			attrs := errorInfo.(map[string]interface{})["attributes"].(map[string]interface{})
			return fmt.Errorf("Fehler beim Löschen der TenantEPG: code=%s, text=%s", attrs["code"], attrs["text"])
		}
	}

	log.Println("TenantEPG erfolgreich gelöscht!")
	return nil
}

// ObserveTenantEPG überprüft, ob eine spezifische TenantEPG existiert und gibt ihren Status zurück
func (c *TenantEPGClient) ObserveTenantEPG(tenantName, appProfileName, epgName string) (bool, error) {
	endpoint := fmt.Sprintf("/api/node/mo/uni/tn-%s/ap-%s/epg-%s.json", tenantName, appProfileName, epgName)
	response, err := c.client.DoRequest("GET", endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("Fehler beim Beobachten der TenantEPG: %w", err)
	}

	// Parse die Antwortdaten in eine Map, um den Status zu analysieren
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return false, fmt.Errorf("Fehler beim Parsen der Antwort: %w", err)
	}

	// Überprüfe, ob die TenantEPG in den Daten enthalten ist
	imdata, ok := result["imdata"].([]interface{})
	if !ok || len(imdata) == 0 {
		return false, fmt.Errorf("TenantEPG %s nicht gefunden in Tenant %s und Application Profile %s", epgName, tenantName, appProfileName)
	}

	return true, nil
}

