package libre

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"

	"librelink-up-tg/config"
)

const USER_AGENT = "Mozilla/5.0 (iPhone; CPU OS 17_4.1 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/17.4.1 Mobile/10A5355d Safari/8536.25"
const LIBRE_LINK_UP_VERSION = "4.12.0"
const LIBRE_LINK_UP_PRODUCT = "llu.ios"

type Client struct {
	httpClient *http.Client
	config     *config.Config
	authTicket AuthTicket
	userID     string
	endpoint   string
}

func NewClient(config *config.Config) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	endpoint := LLU_API_ENDPOINTS[CountryCode(config.LinkUpRegion)]

	return &Client{
		httpClient: &http.Client{Jar: jar},
		config:     config,
		endpoint:   endpoint,
	}, nil
}

func (c *Client) Login() error {
	loginURL := fmt.Sprintf("https://%s/llu/auth/login", c.endpoint)

	data := map[string]string{
		"email":    c.config.LinkUpUsername,
		"password": c.config.LinkUpPassword,
	}

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	c.setCommonHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return err
	}

	c.authTicket = loginResp.Data.AuthTicket
	c.userID = loginResp.Data.User.ID
	return nil
}

func (c *Client) getConnections() ([]Connection, error) {
	connectionsURL := fmt.Sprintf("https://%s/llu/connections", c.endpoint)

	req, _ := http.NewRequest("GET", connectionsURL, nil)
	c.setAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var connResp ConnectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&connResp); err != nil {
		return nil, err
	}

	return connResp.Data, nil
}

func (c *Client) getPatientID() (patientID string, err error) {
	connections, err := c.getConnections()
	if err != nil {
		return "", fmt.Errorf("Failed to get connections: %w", err)
	}

	if len(connections) == 0 {
		return "", fmt.Errorf("No connections found")
	}

	if c.config.LinkUpConnection != "" {
		for _, conn := range connections {
			if conn.PatientID == c.config.LinkUpConnection {
				patientID = conn.PatientID
				break
			}
		}
		if patientID == "" {
			log.Fatal("Specified connection not found")
		}
	} else {
		patientID = connections[0].PatientID
		log.Printf("Using first connection: %s %s (%s)",
			connections[0].FirstName,
			connections[0].LastName,
			patientID)
	}

	return patientID, nil
}

func (c *Client) GetGlucoseData() (*GraphData, error) {
	patientID, err := c.getPatientID()
	if err != nil {
		return nil, err
	}

	graphURL := fmt.Sprintf(
		"https://%s/llu/connections/%s/graph",
		c.endpoint,
		patientID,
	)

	req, _ := http.NewRequest("GET", graphURL, nil)
	c.setAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var graphResp GraphResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphResp); err != nil {
		return nil, err
	}

	return &graphResp.Data, nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("version", LIBRE_LINK_UP_VERSION)
	req.Header.Set("product", LIBRE_LINK_UP_PRODUCT)
	req.Header.Set("Content-Type", "application/json")
}

func (c *Client) setAuthHeaders(req *http.Request) {
	c.setCommonHeaders(req)
	req.Header.Set("Authorization", "Bearer "+c.authTicket.Token)

	if c.userID != "" {
		hash := sha256.Sum256([]byte(c.userID))
		req.Header.Set("account-id", fmt.Sprintf("%x", hash))
	}
}
