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

// Client ist der API-Client für die Kommunikation mit Cisco ACI
type Client struct {
    BaseURL           string
    Username          string
    Password          string
    Token             string
    InsecureSkipVerify bool
}

// NewClient erstellt einen neuen Client für die ACI API
func NewClient(baseURL, username, password string, insecureSkipVerify bool) *Client {
    return &Client{
        BaseURL:           baseURL,
        Username:          username,
        Password:          password,
        InsecureSkipVerify: insecureSkipVerify,
    }
}

// Authenticate authentifiziert den Client und erhält ein Token
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

    log.Printf("Sending authentication request to %s", url)
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("authentication request failed: %v", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    log.Printf("Authentication response: %s", body)

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
    log.Println("Authentication successful, token acquired")
    return nil
}

// DoRequest führt eine HTTP-Anfrage an die ACI API aus
func (c *Client) DoRequest(method, endpoint string, data interface{}) ([]byte, error) {
    url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
    var jsonData []byte
    var err error
    if data != nil {
        jsonData, err = json.Marshal(data)
        if err != nil {
            return nil, fmt.Errorf("error marshalling data: %v", err)
        }
    }
    req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    if c.Token != "" {
        req.Header.Set("Cookie", fmt.Sprintf("APIC-cookie=%s", c.Token))
    }

    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
        },
    }

    log.Printf("Sending %s request to %s with data: %s", method, url, jsonData)
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %v", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    log.Printf("Response: %s", body)
    return body, nil
}

