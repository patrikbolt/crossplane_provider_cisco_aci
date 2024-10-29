package clients

import (
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
)

type ACIClient struct {
    baseURL    string
    httpClient *http.Client
    authToken  string
}

// NewACIClient initializes a new ACI client with authentication
func NewACIClient(baseURL, username, password string) (*ACIClient, error) {
    client := &ACIClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
    }
    if err := client.authenticate(username, password); err != nil {
        return nil, err
    }
    return client, nil
}

// authenticate establishes a session and stores the auth token
func (c *ACIClient) authenticate(username, password string) error {
    authURL := fmt.Sprintf("%s/api/aaaLogin.json", c.baseURL)
    authData := map[string]interface{}{
        "aaaUser": map[string]string{
            "name": username,
            "pwd":  password,
        },
    }
    body, err := json.Marshal(authData)
    if err != nil {
        return err
    }
    resp, err := c.httpClient.Post(authURL, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("authentication failed with status: %s", resp.Status)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }

    // Extract the token
    token, ok := result["imdata"].([]interface{})[0].(map[string]interface{})["aaaLogin"].(map[string]interface{})["token"].(string)
    if !ok {
        return fmt.Errorf("failed to parse authentication token")
    }
    c.authToken = token
    return nil
}

// doRequest handles authenticated requests
func (c *ACIClient) doRequest(req *http.Request) (*http.Response, error) {
    req.Header.Set("Cookie", "APIC-cookie="+c.authToken)
    req.Header.Set("Content-Type", "application/json")
    return c.httpClient.Do(req)
}

