package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// A map that stores Zuora REST endpoints
func getZuoraEndpoints() map[string]string {
	endpoints := map[string]string{
		"EUPROD": "https://rest.eu.zuora.com/",
		"USPROD": "https://rest.zuora.com/",
		"EUSBX":  "https://rest.eu.sandbox.zuora.com/",
		"EUTEST": "https://rest.test.eu.zuora.com/",
	}
	return endpoints
}

func getOAuthToken(clientID, clientSecret, grantType, urlEndpoint, tenantId string) (interface{}, error) {
	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", urlEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response %v", err)
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting OAuth token: %v", response)
	}

	return response, nil
}

func createJSONFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	return nil
}

func login(urlStr, username, password string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return resp, nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func createHiddenFile(fileName string, cookieName string, cookieValue string) error {
	htmlContent := createCookieSettingHTML(cookieName, cookieValue)
	err := os.WriteFile(fileName, []byte(htmlContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func createCookieSettingHTML(cookieName string, cookieValue string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Setting Cookie...</title>
		<script>
			document.cookie = "%s=%s; path=/;";
			console.log("Cookie set: %s=%s");
		</script>
	</head>
	<body>
		<p>Redirecting...</p>
	</body>
	</html>`, cookieName, cookieValue, cookieName, cookieValue)
}
