package homeassistant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func CheckConfig(haUrl string, haToken string) error {
	url, err := url.JoinPath(haUrl, "/api/config/core/check_config")
	if err != nil {
		return fmt.Errorf("failed to create homeassistant check_config URL: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create homeassistant check_config request: %v", err)
	}

	var bearer = "Bearer " + haToken
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send homeassistant check_config request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read homeassistant check_config response: %v", err)
	}

	var r map[string]string
	err = json.Unmarshal(body, &r)
	if err != nil {
		return fmt.Errorf("failed to unmarshal homeassistant check_config response: %v", err)
	}

	if r["result"] != "valid" || r["errors"] != "" || r["warnings"] != "" {
		return fmt.Errorf("homeassistant check config failed: %s", body)
	}

	return nil
}
