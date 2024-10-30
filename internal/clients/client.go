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

// Client represents the API client for communicating with Cisco ACI
type Client struct {
	BaseURL            string
	Username           string
	Password           string
	Token              string
	InsecureSkipVerify bool
}

// NewClient creates a new Client for the ACI API
func NewClient(baseURL, username, password string, insecureSkipVerify bool) *Client {
	return &Client{
		BaseURL:            baseURL,
		Username:           username,
		Password:           password,
		InsecureSkipVerify: insecureSkipVerify,
	}
}

// Authenticate authenticates the client and retrieves a token
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
	jsonData, _ := json.Marshal(authData)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	imdata, ok := result["imdata"].([]interface{})
	if !ok || len(imdata) == 0 {
		return fmt.Errorf("authentication failed: invalid response")
	}
	aaaLogin, ok := imdata[0].(map[string]interface{})["aaaLogin"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("authentication failed: missing aaaLogin data")
	}
	attributes, ok := aaaLogin["attributes"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("authentication failed: missing attributes")
	}
	c.Token, ok = attributes["token"].(string)
	if !ok {
		return fmt.Errorf("authentication failed: missing token")
	}
	log.Println("Authenticated successfully.")
	return nil
}

// DoRequest performs an HTTP request to the ACI API
func (c *Client) DoRequest(method, endpoint string, data interface{}) ([]byte, error) {
	if c.Token == "" {
		log.Println("No token found, authenticating...")
		if err := c.Authenticate(); err != nil {
			return nil, fmt.Errorf("failed to authenticate: %v", err)
		}
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Re-authenticate if 403 response and retry request
	if resp.StatusCode == 403 {
		log.Println("Token expired, re-authenticating...")
		if err := c.Authenticate(); err != nil {
			return nil, fmt.Errorf("failed to re-authenticate: %v", err)
		}
		req.Header.Set("Cookie", fmt.Sprintf("APIC-cookie=%s", c.Token))
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error from server: code=%d, status=%s", resp.StatusCode, resp.Status)
	}

	return body, nil
}
