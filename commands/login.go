package commands

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
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
				fmt.Println(resp)
				fmt.Println("Login successful!")

				var cookie = resp.Header.Get("Set-Cookie")
				parts := strings.SplitN(cookie, ";", 2)
				cookieParts := strings.SplitN(parts[0], "=", 2)
				fmt.Println(cookieParts)

				cookiename := cookieParts[0]
				cookieValue := cookieParts[1]
				createHiddenFile("cookie.html", cookiename, cookieValue)
				openBrowser("./cookie.html")
			} else {
				fmt.Printf("Login failed with status code: %d\n", resp.StatusCode)
			}
		},
	}

	cmd.Flags().StringP("userName", "u", "", "UserName for zuora one Id")
	cmd.Flags().StringP("password", "p", "", "Password")

	cmd.MarkFlagRequired("userName")
	cmd.MarkFlagRequired("password")

	return cmd
}
