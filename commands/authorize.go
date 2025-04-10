package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

type OrgType string

const (
	EUPROD OrgType = "EUPROD"
	USPROD OrgType = "USPROD"
	EUSBX  OrgType = "EUSBX"
	EUTEST OrgType = "EUTEST"
)

const path = "./Auths/"

var allowedOrgTypes = []OrgType{EUPROD, USPROD, EUSBX, EUTEST}

func (o *OrgType) Set(val string) error {
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

func NewAuthorizeCommand() *cobra.Command {
	var target OrgType
	var clientId, clientSecret, tenantId, alias string

	cmd := &cobra.Command{
		Use:   "authorize",
		Short: "Authorize a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Authorizing Zuora organization with target: %s, clientId: %s, clientSecret: %s, tenantId: %s, alias: %s\n",
				target, clientId, clientSecret, tenantId, alias)

			// Get the Zuora REST endpoint
			var endpoint = getZuoraEndpoints()[string(target)] + "oauth/token"
			var tokeres, err = getOAuthToken(clientId, clientSecret, "client_credentials", endpoint, tenantId)
			if err != nil {
				fmt.Println("Error getting OAuth token:", err)
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
				fmt.Println("Error creating JSON file:", errFile)
				return
			}

			fmt.Println("Successfully authenticated to Zuora!")
		},
	}

	cmd.Flags().VarP(&target, "target", "t", "Target environment (e.g., PROD)")
	cmd.Flags().StringP("clientId", "i", "", "Zuora client ID")
	cmd.Flags().StringP("clientSecret", "s", "", "Zuora client secret")
	cmd.Flags().StringP("tenantId", "d", "", "Zuora tenant ID (optional)")
	cmd.Flags().StringP("alias", "a", "", "Alias for the Zuora organization")

	cmd.MarkFlagRequired("target")
	cmd.MarkFlagRequired("clientId")
	cmd.MarkFlagRequired("clientSecret")
	cmd.MarkFlagRequired("alias")

	return cmd
} 