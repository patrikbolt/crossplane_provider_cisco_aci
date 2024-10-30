package clients

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client repräsentiert den API-Client für die Kommunikation mit Cisco ACI
type Client struct {
	BaseURL            string
	Username           string
	Password           string
	Token              string
	InsecureSkipVerify bool
}

// NewClient erstellt einen neuen Client für die ACI API
func NewClient(baseURL, username, password string, insecureSkipVerify bool) *Client {
	return &Client{
		BaseURL:            baseURL,
		Username:           username,
		Password:           password,
		InsecureSkipVerify: insecureSkipVerify,
	}
}

// Authenticate authentifiziert den Client und ruft ein Token ab
func (c *Client) Authenticate() error {
	url := fmt.Sprintf("%s/api/aaaLogin.json", c.BaseURL)
	authData := map[string]interface{}{
		"aaaUser": map[string]interface{}{
			"attributes": map[string]string{
				"name": c.Username,
				"pwd":  c.Password,
			},
		},
	}
	jsonData, err := json.Marshal(authData)
	if err != nil {
		return fmt.Errorf("Fehler beim Marshalen der Authentifizierungsdaten: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der Authentifizierungsanfrage: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Fehler bei der Authentifizierungsanfrage: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Fehler beim Lesen der Authentifizierungsantwort: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("Fehler beim Unmarshalen der Authentifizierungsantwort: %v", err)
	}
	imdata, ok := result["imdata"].([]interface{})
	if !ok || len(imdata) == 0 {
		return fmt.Errorf("Authentifizierung fehlgeschlagen: ungültige Antwort")
	}
	aaaLogin, ok := imdata[0].(map[string]interface{})["aaaLogin"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Authentifizierung fehlgeschlagen: fehlende aaaLogin-Daten")
	}
	attributes, ok := aaaLogin["attributes"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Authentifizierung fehlgeschlagen: fehlende Attribute")
	}
	c.Token, ok = attributes["token"].(string)
	if !ok {
		return fmt.Errorf("Authentifizierung fehlgeschlagen: fehlendes Token")
	}
	log.Println("Erfolgreich authentifiziert.")
	return nil
}

// DoRequest führt eine HTTP-Anfrage an die ACI API durch
func (c *Client) DoRequest(method, endpoint string, data interface{}) ([]byte, error) {
	if c.Token == "" {
		log.Println("Kein Token gefunden, authentifiziere...")
		if err := c.Authenticate(); err != nil {
			return nil, fmt.Errorf("Authentifizierung fehlgeschlagen: %v", err)
		}
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	var reqBody *bytes.Buffer
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Marshalen der Anfrage-Daten: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Erstellen der Anfrage: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("APIC-cookie=%s", c.Token))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Fehler bei der Anfrage: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Lesen der Antwort: %v", err)
	}

	// Re-authentifiziere, wenn eine 403-Antwort empfangen wird, und versuche die Anfrage erneut
	if resp.StatusCode == 403 {
		log.Println("Token abgelaufen, re-authentifiziere...")
		if err := c.Authenticate(); err != nil {
			return nil, fmt.Errorf("Re-Authentifizierung fehlgeschlagen: %v", err)
		}
		req.Header.Set("Cookie", fmt.Sprintf("APIC-cookie=%s", c.Token))
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Fehler bei der erneuten Anfrage nach Re-Authentifizierung: %v", err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Lesen der erneuten Antwort: %v", err)
		}
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Fehler vom Server: code=%d, status=%s", resp.StatusCode, resp.Status)
	}

	return body, nil
}

