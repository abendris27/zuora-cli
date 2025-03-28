package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	// "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

// Enum for the flag values
type OrgType string

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

const (
	EUPROD OrgType = "EUPROD"
	USPROD OrgType = "USPROD"
	EUSBX  OrgType = "EUSBX"
	EUTEST OrgType = "EUTEST"
)

const path = "./Auths/"

var allowedOrgTypes = []OrgType{EUPROD, USPROD, EUSBX, EUTEST}

func (o *OrgType) Set(val string) error {
	// Check if the value is a valid OrgType
	switch OrgType(val) {
	case EUPROD, USPROD, EUSBX, EUTEST:
		*o = OrgType(val)
		return nil
	default:
		return fmt.Errorf("invalid OrgType: %s. Allowed values are: %v", val, allowedOrgTypes)
	}
}

func (o *OrgType) Type() string {
	return "OrgType"
}

func (o *OrgType) String() string {
	return string(*o)
}

// Global variables for flags
var clientId, clientSecret, tenantId, alias string
var target OrgType

func getOAuthToken(clientID, clientSecret, grantType, urlEndpoint, tenantId string) (interface{}, error) {
	// Prepare the URL and form data
	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	// Create the POST request with headers and form data
	req, err := http.NewRequest("POST", urlEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body using io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, fmt.Errorf("Error reading response %v", err)
	}

	// Print the response
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	// Unmarshal the response body into an interface{}
	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting OAuth token: %v", response)
	}

	// Return the JSON response as interface{}
	return response, nil
}

// Function to create a JSON file from any data
func createJSONFile(filePath string, data interface{}) error {
	// Open or create the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create a JSON encoder and set it to pretty-print
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// Encode the data and write it to the file
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Successfully created JSON file
	return nil
}

// Function to print selected columns of JSON files as tables
func printJSONFilesAsTable(directory string, columnsToDisplay []string) error {
	// Open the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Initialize table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(append([]string{"File Name"}, columnsToDisplay...))

	// Loop over files in the directory
	for _, file := range files {
		// Only process JSON files
		if filepath.Ext(file.Name()) == ".json" {
			// Open the JSON file
			filePath := filepath.Join(directory, file.Name())
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			// Unmarshal JSON into a generic map
			var data map[string]interface{}
			if err := json.Unmarshal(fileContent, &data); err != nil {
				fmt.Printf("Error unmarshalling file %s: %v\n", file.Name(), err)
				continue
			}

			// Create a slice for row data with only selected columns
			row := []string{file.Name()}

			// Add selected columns to the row
			for _, col := range columnsToDisplay {
				if value, exists := data[col]; exists {
					row = append(row, fmt.Sprintf("%v", value))
				} else {
					row = append(row, "N/A") // If column doesn't exist in the data, add "N/A"
				}
			}

			// Add the row to the table
			table.Append(row)
		}
	}

	// Render the table
	table.Render()

	return nil
}

// Function to login
func login(urlStr, username, password string) (*http.Response, error) {
	// Initialize cookie jar to store session cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	// Initialize HTTP client with cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Create login payload
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	// Make the POST request to the login page
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set content-type header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	// Return the response
	return resp, nil
}

// Function to open the URL in the default web browser
func openBrowser(url string, cookieValue string) error {
	var cmd *exec.Cmd

	// Determine OS and open browser with cookie injection
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	// Set the environment variable for the browser to use the cookie
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZUNISESSION=%s", cookieValue))

	return cmd.Start()
}

// Start a local server to inject the cookie
func startServer(url string, name string, value string, domain string, path string, httpOnly bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// expiration := time.Now().Add(24 * time.Hour)
		fmt.Println("Setting cookie:", name, "=", value)

		// Set the session cookie in the HTTP response
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    value,
			Domain:   domain, // Make sure it's the right domain
			Path:     path,
			HttpOnly: httpOnly,
		})

		// Redirect to the target URL
		fmt.Println("Redirecting to:", url)
		http.Redirect(w, r, url, http.StatusFound)
	})

	fmt.Println("Starting local server on http://localhost:8080...")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Server error:", err)
		}
	}()
}

// ProxyHandler intercepts requests and injects the session cookie
func ProxyHandler(target *url.URL, name string, value string, domain string, path string, httpOnly bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔹 Incoming Request:", r.Method, r.URL)
		for name, values := range r.Header {
			log.Printf("  %s: %s\n", name, values)
		}
		// Inject the session cookie
		cookie := &http.Cookie{
			Name:   name,
			Value:  value,
			Domain: domain,
			Path:   path,
			// Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: httpOnly,
		}
		r.AddCookie(cookie)

		// Proxy the request to Zuora
		proxy := httputil.NewSingleHostReverseProxy(target)
		r.Host = target.Host
		proxy.ServeHTTP(w, r)
	})
}

// StartProxy starts the local proxy server
func StartProxy(name string, value string, domain string, path string, httpOnly bool) {
	target, _ := url.Parse("https://one.zuora.com")
	http.Handle("/", ProxyHandler(target, name, value, domain, path, httpOnly))

	fmt.Println("🚀 Local proxy running at http://localhost" + ":8080" + "...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetCookie(name string, value string, domain string, path string, httpOnly bool, sameSite string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		err := network.SetCookie(name, value).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSameSite(network.CookieSameSite(sameSite)).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func openBrowserWithChromeDp(url, name string, value string, domain string, path string, httpOnly bool, sameSite string) error {
	fmt.Println("Opening browser with chromedp...", name)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.Flag("headless", false))
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Run browser automation
	err := chromedp.Run(ctx,
		SetCookie(name, value, domain, path, httpOnly, sameSite),
		chromedp.Navigate(url), // Open the page first
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the cookie using JavaScript execution
			script := fmt.Sprintf(`document.cookie = "%s=%s";`, name, value)
			return chromedp.Evaluate(script, nil).Do(ctx)

		}),
		chromedp.WaitVisible("body"), // Wait for body to be visible to confirm page load
		chromedp.Reload(),            // Reload the page to apply the cookie
	)

	return err
}

func OpenBrowserWithProxy() error {
	proxyURL := "http://localhost:8080"
	targetURL := "https://one.zuora.com/one-id/home"

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "chrome", "--proxy-server="+proxyURL, targetURL)
	case "darwin": // macOS
		cmd = exec.Command("open", "-a", "Google Chrome", "--args", "--proxy-server="+proxyURL, targetURL)
	case "linux":
		cmd = exec.Command("google-chrome", "--proxy-server="+proxyURL, targetURL)
	default:
		return fmt.Errorf("❌ Unsupported platform")
	}
	return cmd.Start()
}

func main() {

	// ANSI escape code for green text
	green := "\033[32m"
	// ANSI escape code to reset text color
	reset := "\033[0m"
	// ANSI escape code for green text
	red := "\033[31m"
	zuoraEndpoints := getZuoraEndpoints()

	var target OrgType

	var rootCmd = &cobra.Command{
		Use:   "zuo",
		Short: "Zuora CLI Tool",
	}

	var orgCmd = &cobra.Command{
		Use:   "org",
		Short: "Manage Zuora organizations",
	}

	// Define the "authorize" subcommand
	var authorizeCmd = &cobra.Command{
		Use:   "authorize",
		Short: "Authorize a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			clientId, _ := cmd.Flags().GetString("clientId")
			clientSecret, _ := cmd.Flags().GetString("clientSecret")
			tenantId, _ := cmd.Flags().GetString("tenantId")
			alias, _ := cmd.Flags().GetString("alias")

			fmt.Printf("Authorizing Zuora organization with target: %s, clientId: %s, clientSecret: %s, tenantId: %s, alias: %s\n",
				target, clientId, clientSecret, tenantId, alias)

			// Get the Zuora REST endpoint
			var endpoint = zuoraEndpoints[string(target)] + "oauth/token"
			var tokeres, err = getOAuthToken(clientId, clientSecret, "client_credentials", endpoint, tenantId)
			if err != nil {
				fmt.Println(red+"Error getting OAuth token:"+reset, err)
				return
			}
			tokeres.(map[string]interface{})["alias"] = alias
			tokeres.(map[string]interface{})["clientId"] = clientId
			tokeres.(map[string]interface{})["clientSecret"] = clientSecret
			tokeres.(map[string]interface{})["tenantId"] = tenantId
			tokeres.(map[string]interface{})["expireAt"] = time.Now().Add(time.Duration(tokeres.(map[string]interface{})["expires_in"].(float64)) * time.Second)
			tokeres.(map[string]interface{})["firstGenerated"] = time.Now()

			var errFile = createJSONFile(path+alias+"_token.json", tokeres)
			if errFile != nil {
				fmt.Println(red+"Error creating JSON file:"+reset, errFile)
				return
			}

			// Print a message in green
			fmt.Println(green + "Susscefully authenticated to Zuora Well done!" + reset)
			// Add the actual authorization logic here
		},
	}

	// Define flags for the "authorize" command
	authorizeCmd.Flags().VarP(&target, "target", "t", "Target environment (e.g., PROD)")
	authorizeCmd.Flags().StringP("clientId", "i", "", "Zuora client ID")
	authorizeCmd.Flags().StringP("clientSecret", "s", "", "Zuora client secret")
	authorizeCmd.Flags().StringP("tenantId", "d", "", "Zuora tenant ID (optional)")
	authorizeCmd.Flags().StringP("alias", "a", "", "Alias for the Zuora organization")

	// Mark the greeting flag as required
	authorizeCmd.MarkFlagRequired("target")
	authorizeCmd.MarkFlagRequired("clientId")
	authorizeCmd.MarkFlagRequired("clientSecret")
	authorizeCmd.MarkFlagRequired("alias")

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login to oneId Zuora",
		Run: func(cmd *cobra.Command, args []string) {
			userName, _ := cmd.Flags().GetString("userName")
			password, _ := cmd.Flags().GetString("password")

			fmt.Printf("login to one.zuora.com for : %s\n", userName)

			resp, err := login("https://one.zuora.com/api/login", userName, password)
			if err != nil {
				log.Fatalf("Login failed: %v", err)
			}
			// Check if login is successful by inspecting response status
			if resp.StatusCode == http.StatusOK {
				// fmt.Println(resp)
				fmt.Println("Login successful!")

				var cookieStr = resp.Header.Get("Set-Cookie")
				parts := strings.SplitN(cookieStr, ";", 2)
				cookieParts := strings.SplitN(parts[0], "=", 2)
				fmt.Println(cookieParts)

				cookieName := cookieParts[0]
				cookieValue := cookieParts[1]

				// go StartProxy(cookieName, cookieValue, "one.zuora.com", "/", true)
				// time.Sleep(5 * time.Second)
				// err := OpenBrowserWithProxy()

				// 	fmt.Println("Starting local server to set cookie...")
				// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// 		html := fmt.Sprintf(`
				// 		<!DOCTYPE html>
				// 		<html>
				// 		<head>
				// 			<title>Zuora Login</title>
				// 			<script>
				// 				document.cookie = "%s=%s; path=/";
				// 				window.location.href = "https://one.zuora.com/one-id/home";
				// 			</script>
				// 		</head>
				// 		<body>
				// 			<p>Logging in to Zuora...</p>
				// 		</body>
				// 		</html>
				// 		`, cookieName, cookieValue)

				// 		fmt.Println("Setting cookie:", html)
				// 		w.Header().Set("Content-Type", "text/html")
				// 		w.Write([]byte(html))
				// 	})

				// 	server := &http.Server{Addr: ":8085"}
				// 	go func() {
				// 		fmt.Println("Server started on http://localhost:8085")
				// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				// 			fmt.Printf("Server error: %v\n", err)
				// 		}
				// 	}()
				// 	// Open the HTML file in the default browser
				// 	var cmd *exec.Cmd
				// 	switch runtime.GOOS {
				// 	case "windows":
				// 		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:8085")
				// 	case "darwin":
				// 		cmd = exec.Command("open", "http://localhost:8085")
				// 	case "linux":
				// 		cmd = exec.Command("xdg-open", "http://localhost:8085")
				// 	default:
				// 		fmt.Println("Unsupported platform")
				// 		return
				// 	}

				// 	if err := cmd.Start(); err != nil {
				// 		fmt.Printf("Failed to open browser: %v\n", err)
				// 		return
				// 	}

				// 	fmt.Println("Browser opened successfully!")

				// 	// Wait for user to press enter before shutting down server
				// 	fmt.Println("Press Enter to exit...")
				// 	fmt.Scanln()
				// 	server.Shutdown(context.Background())
				// } else {
				// 	fmt.Printf("Login failed with status code: %d\n", resp.StatusCode)
				// }

				// After successful login, open the browser to the desired page
				err = openBrowserWithChromeDp("https://one.zuora.com/one-id/home", cookieName, cookieValue, ".one.zuora.com", "/", true, "Lax")
				if err != nil {
					fmt.Printf("Failed to open browser: %v\n", err)
				} else {
					fmt.Println("Browser opened successfully!")
				}
			} else {
				fmt.Printf("Login failed with status code: %d\n", resp.StatusCode)
			}
		},
	}

	// Define flags for the "login" command
	loginCmd.Flags().StringP("userName", "u", "", "UserName for zuora one Id")
	loginCmd.Flags().StringP("password", "p", "", "Password")

	// Mark the greeting flag as required
	loginCmd.MarkFlagRequired("userName")
	loginCmd.MarkFlagRequired("password")

	// Define the "delete" subcommand
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle deletion
			fmt.Printf("Deleting Zuora organization with alias: %s\n", alias)
			// Add the actual deletion logic here
		},
	}

	// Define flags for the "list" command
	deleteCmd.Flags().StringVar(&alias, "alias", "", "Alias for the Zuora organization")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "list the stored Zuora organizations",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle deletion
			fmt.Printf("listing the stored Zuora organizations\n")
			// Add the actual deletion logic here
			columnsToDisplay := []string{"alias", "clientId", "tenantId", "firstGenerated"}
			err := printJSONFilesAsTable(path, columnsToDisplay)
			if err != nil {
				fmt.Println(red+"Error listing JSON files:"+reset, err)
				return
			}
		},
	}

	// Define the "delete" subcommand
	var viewOrgCmd = &cobra.Command{
		Use:   "view",
		Short: "view a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle deletion
			fmt.Printf("viewing Zuora organization with alias: %s\n", alias)
			// Add the actual deletion logic here
		},
	}

	// Define flags for the "list" command
	viewOrgCmd.Flags().StringVar(&alias, "alias", "a", "Alias for the Zuora organization")

	// Add the commands to the org command
	orgCmd.AddCommand(authorizeCmd)
	orgCmd.AddCommand(deleteCmd)
	orgCmd.AddCommand(listCmd)
	orgCmd.AddCommand(viewOrgCmd)
	orgCmd.AddCommand(loginCmd)

	// Add org command to the root command
	rootCmd.AddCommand(orgCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
